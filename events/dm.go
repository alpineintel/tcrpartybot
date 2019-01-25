package events

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	goTwitter "github.com/stevenleeg/go-twitter/twitter"

	"github.com/dustin/go-humanize"
	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

const (
	noAccountMsg       = "Hmm, I haven't met you yet. If you want to join the TCR party send me a tweet that says \"let's party\""
	inactiveAccountMsg = "The party hasn't started just yet. Keep an eye on @TCRPartyBot for a launch announcement."

	votingArgErrorMsg          = "Whoops, looks like you forgot something. Try again with something like 'vote [twitter handle] yes' or 'vote [twitter handle] no'"
	votingListingNotFoundMsg   = "Hmm, I couldn't find a registry listing for that twitter handle. Are you sure they've been nominated to the registry?"
	votingChallengeNotFoundMsg = "Looks like that twitter handle doesn't have a challenge opened on it yet. If you'd like to challenge their place on the registry respond with 'challenge %s'."
	votingChallengeErrorMsg    = "There was an error committing your vote. The admins have been notified!"
	votedAlreadyMsg            = "Oops! Looks like you've already voted %s on this challenge."
	votingEndedMsg             = "Ack! Looks like the voting period has ended for this challenge. Hang tight, we'll announce the result on %s."
	votingSuccessMsg           = "Your vote has been committed! Hang tight, we'll announce the results on %s."

	challengeArgErrorMsg          = "Whoops, looks like you forgot something. Try again with something like 'challenge [twitter handle]'. Eg: 'challenge weratedogs'"
	challengeNotFoundMsg          = "Looks like nobody has tried creating a listing for this twitter handle yet."
	challengeAlreadyExistsMsg     = "Looks like somebody has already begun a challenge for this twitter handle. You can support this challenge by voting yes (respond with 'vote %s yes')."
	challengeInsufficientFundsMsg = "Drat, looks like you don't have enough TCRP to start a challenge. You'll need 500 available in your wallet."
	challengeSubmissionErrorMsg   = "There was an error trying to submit your challenge. The admins have been notified!"
	challengeSuccessMsg           = "We've submitted your challenge to the registry (tx: %s). Keep an eye on @TCRPartyVIP for updates."

	nominateArgErrorMsg          = "Whoops, looks like you forgot something. Try again with something like 'nominate [twitter handle]'. Eg: 'apply weratedogs'"
	nominateAlreadyAppliedMsg    = "Looks like that Twitter handle has already been submitted to the TCR. A twitter handle can only appear on the TCR once, so you'll need to wait for a successful challenge (or a delisting) in order to re-nominate them."
	nominateInsufficientFundsMsg = "Drat! Looks like you don't have enough TCRP to start a nomination. You'll need 500 available in your wallet."
	nominateSubmissionErrorMsg   = "There was an error trying to submit your nomination. The admins have been notified!"
	nominateSuccessMsg           = "We've submitted your nomination to the registry (tx: %s). Keep an eye on @TCRPartyVIP for updates."

	balanceMsg                      = "Your balance is %d TCRP"
	awaitingPartyBeginMsg           = "🎉 You're registered to party 🎉. Hang tight while we prepare to distribute our token."
	invalidChallengeResponseMsg     = "🙅‍♀️ That's not right! %s"
	nextChallengeMsg                = "Nice, that's it! Here's another one for you: %s"
	preregistrationSuccessMsg       = "🎉 Awesome! You've been registered for the party. We'll reach out once we're ready to distribute TCRP tokens 🎈."
	registrationSuccessMsg          = "🎉 Awesome! Now that you're registered I'll need a few minutes to build your wallet and give you some TCR Party Points to get started with. I'll send you a DM once I'm done."
	invalidCommandMsg               = "Whoops, I don't recognize that command. Try typing help to see what you can say to me."
	helpMsg                         = "Here are the commands I recognize:\n• balance - See your TCRP balance\n• nominate [handle] = Nominate the given Twitter handle to be on the TCR\n• challenge [handle] - Begin a challenge for a listing on the TCR\n• vote [handle] [kick/keep] - Vote on an existing listing's challenge."
	cannotHitFaucetMsg              = "Ack, I can only let you hit the faucet once per day. Try again %s."
	hitFaucetMsg                    = "You got it. %d TCRP headed your way. TX Hash: %s"
	errorMsg                        = "Yikes, we ran into an error: %s. Try tweeting at @stevenleeg for help."
	voteBalanceMsg                  = "You have %d tokens deposited to vote. This means you can vote with a maximum weight of %d."
	plcrDepositInsufficientFundsMsg = "Whoops, looks like you don't have enough tokens to deposit this amount to your maximum voting weight. Your current balance is %d"
	plcrDepositSuccessMsg           = "Your tokens have been deposited successfully! TX Hash: %s"
	plcrWithdrawInsufficientFunds   = "Whoops, you only have %d tokens locked up."
	plcrWithdrawSuccessMsg          = "Your tokens have been withdrawn successfully! TX hash: %s"

	depositAmount = 500
	voteAmount    = 50
	faucetAmount  = 50
)

// RegistrationEventData collects the required data for keeping track of
// the user registration flow into a struct
type RegistrationEventData struct {
	Event     *TwitterEvent
	Challenge *models.RegistrationChallengeRegistrationQuestion
	Account   *models.Account
}

// ListenForTwitterDM is a blocking function which polls the Twitter API for
// new direct messages and sends them off to the eventChan for further
// processing as they are received.
func ListenForTwitterDM(handle string, eventChan chan<- *TwitterEvent, errChan chan<- error) {
	client, token, err := twitter.GetClientFromHandle(handle)
	if err != nil {
		log.Println("Could not establish client listening to DMs")
		errChan <- err
		return
	}

	for {
		latestID, err := models.GetKey("latestDirectMessageID")
		if err != nil {
			log.Println("Error fetching latest ID")
			errChan <- err
			return
		}

		var cursor string
		for {
			params := &goTwitter.DirectMessageEventsListParams{
				Cursor: cursor,
				Count:  20,
			}
			events, _, err := client.DirectMessages.EventsList(params)
			if err != nil {
				log.Println("Could not fetch event feed")
				errChan <- err
				time.Sleep(2 * time.Minute)
				break
			}

			// Store the latest cursor in our keyval store
			models.SetKey("latestDirectMessageID", events.Events[0].ID)

			for _, event := range events.Events {
				if event.Type != "message_create" {
					continue
				}

				// If this condition is true we've hit the most recently processed
				// event on the last pull and don't need to process the remainder
				// of the list.
				if event.ID == latestID {
					break
				}

				// If we are the sender we can safely ignore the value
				if event.Message.SenderID == strconv.FormatInt(token.TwitterID, 10) {
					continue
				}

				log.Printf("Received DM from %s: %s", event.Message.SenderID, event.Message.Data.Text)
			}

			time.Sleep(1 * time.Minute)
			if events.NextCursor == "" {
				break
			}
			cursor = events.NextCursor
		}
	}
}

func generateSendDM(account *models.Account, errChan chan<- error) func(message string) {
	return func(message string) {
		// Prevent us from spamming
		if account.LastDMAt != nil && time.Now().Sub(*account.LastDMAt) < 2*time.Second {
			return
		}

		go func() {
			// Wait a second in order to prevent users from spamming
			time.Sleep(1 * time.Second)
			if err := twitter.SendDM(account.TwitterID, message); err != nil {
				errChan <- err
				return
			}

			if err := account.UpdateLastDMAt(); err != nil {
				errChan <- err
				return
			}
		}()
	}

}

func processDM(event *TwitterEvent, errChan chan<- error) {
	log.Printf("Received DM from %s: %s", event.SourceHandle, event.Message)

	// If they don't have an acccount, do nothing.
	account, err := models.FindAccountByHandle(event.SourceHandle)
	if account == nil || err != nil {
		err := twitter.SendDM(event.SourceID, noAccountMsg)
		if err != nil {
			errChan <- err
			return
		}
		return
	}

	// Define a helper function that will be passed around below
	sendDM := generateSendDM(account, errChan)

	// Are they still in the registration challenge stage?
	if account.PassedRegistrationChallengeAt == nil {
		activeChallenge, err := models.FindIncompleteChallenge(account.ID)
		if err != nil {
			errChan <- err
			return
		}

		data := RegistrationEventData{
			Event:     event,
			Challenge: activeChallenge,
			Account:   account,
		}

		if activeChallenge != nil {
			verifyAnswer(data, errChan)
			return
		}

		log.Println("Account has yet to pass activation challenge but does not have any incomplete challenges")
		return
	}

	// They've registered but their account hasn't been activated yet (perhaps
	// because we're still in the preregistration phase?
	if account.ActivatedAt == nil {
		sendDM(inactiveAccountMsg)
		return
	}

	// Looks like they're good to go and are sending a command to the bot
	argv := strings.Split(event.Message, " ")

	// If they haven't been activated yet (ie pre-registration) then we'll
	// stop them here
	if account.ActivatedAt == nil {
		sendDM(preregistrationSuccessMsg)
		return
	}

	switch argv[0] {
	case "balance":
		err = handleBalance(account, argv, sendDM)
	case "nominate":
		err = handleNomination(account, argv, sendDM)
	case "challenge":
		err = handleChallenge(account, argv, sendDM)
	case "vote":
		err = handleVote(account, argv, sendDM)
	case "faucet":
		err = handleFaucet(account, argv, sendDM)
	case "vote-balance":
		err = handleVoteBalance(account, argv, sendDM)
	case "vote-deposit":
		err = handleVoteDeposit(account, argv, sendDM)
	case "vote-withdraw":
		err = handleVoteWithdraw(account, argv, sendDM)
	case "help":
		sendDM(helpMsg)
	default:
		sendDM(invalidCommandMsg)
	}

	if err != nil {
		errChan <- err
	}
}

func verifyAnswer(data RegistrationEventData, errChan chan<- error) {
	// Check to see if they've responded with the correct answer
	if strings.ToLower(data.Event.Message) != strings.ToLower(data.Challenge.Answer) {
		response := fmt.Sprintf(invalidChallengeResponseMsg, data.Challenge.Question)
		err := twitter.SendDM(data.Account.TwitterID, response)
		if err != nil {
			errChan <- err
		}
		return
	}

	// They got it! Let's mark the challenge as completed and give them another
	// question if they have any remaining
	err := data.Challenge.MarkCompleted()
	if err != nil {
		errChan <- err
		return
	}

	// Are they completely done?
	completedChallenges, err := data.Account.HasCompletedChallenges()
	if err != nil {
		errChan <- err
		return
	}

	// Yes! Let's let them know that they're good to go.
	if completedChallenges {
		err := data.Account.MarkRegistered()
		if err != nil {
			errChan <- err
			return
		}

		// Send them a DM letting them know they're good to go
		msg := registrationSuccessMsg
		if os.Getenv("PREREGISTRATION") == "true" {
			msg = preregistrationSuccessMsg
		} else {
			err := data.Account.MarkActivated()
			if err != nil {
				errChan <- err
				return
			}
		}
		err = twitter.SendDM(data.Account.TwitterID, msg)
		if err != nil {
			errChan <- err
		}

		// Let's also register a new multisig wallet for them
		go provisionWallet(data.Account, errChan)
		return
	}

	// Nope, so let's send them the next challenge question
	activeChallenge, err := models.FindUnsentChallenge(data.Account.ID)
	if err != nil {
		errChan <- err
		return
	}

	err = activeChallenge.MarkSent()
	if err != nil {
		errChan <- err
		return
	}

	response := fmt.Sprintf(nextChallengeMsg, activeChallenge.Question)
	err = twitter.SendDM(data.Account.TwitterID, response)
	if err != nil {
		errChan <- err
	}
}

func provisionWallet(account *models.Account, errChan chan<- error) {
	tx, identifier, err := contracts.DeployWallet()
	if err != nil {
		errChan <- err
		return
	}

	err = account.SetMultisigFactoryIdentifier(identifier)
	if err != nil {
		errChan <- err
		return
	}

	receipt, err := contracts.AwaitTransactionConfirmation(tx.Hash())
	if err != nil {
		errChan <- err
		return
	}

	// Make sure the wallet creation actually succeeded
	if receipt.Status == ethTypes.ReceiptStatusFailed {
		errChan <- fmt.Errorf("Could not create multisig wallet for account %d", account.ID)
		return
	}
}

func handleNomination(account *models.Account, argv []string, sendDM func(string)) error {
	if len(argv) != 2 {
		sendDM(nominateArgErrorMsg)
		return nil
	}

	if account.MultisigAddress == nil || !account.MultisigAddress.Valid {
		return errors.New("User attempted to nominate without a multisig address")
	}

	balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}

	if balance.Cmp(contracts.GetAtomicTokenAmount(depositAmount)) == -1 {
		sendDM(nominateInsufficientFundsMsg)
		return nil
	}

	alreadyApplied, err := contracts.ApplicationWasMade(argv[1])
	if err != nil {
		return err
	}

	if alreadyApplied {
		sendDM(nominateAlreadyAppliedMsg)
		return nil
	}

	// Calculate the atomic number of tokens needed to apply
	tokens := contracts.GetAtomicTokenAmount(depositAmount)
	tx, err := contracts.Apply(account.MultisigAddress.String, tokens, argv[1])
	if err != nil {
		sendDM(nominateSubmissionErrorMsg)
		return err
	}
	msg := fmt.Sprintf(nominateSuccessMsg, tx.Hash().Hex())
	sendDM(msg)

	return nil
}

func handleBalance(account *models.Account, argv []string, sendDM func(string)) error {
	if account.MultisigAddress == nil || !account.MultisigAddress.Valid {
		return nil
	}

	balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}

	humanBalance := contracts.GetHumanTokenAmount(balance).Int64()
	sendDM(fmt.Sprintf(balanceMsg, humanBalance))

	return nil
}

func handleChallenge(account *models.Account, argv []string, sendDM func(string)) error {
	if len(argv) != 2 {
		sendDM(challengeArgErrorMsg)
		return nil
	}

	if account.MultisigAddress == nil || !account.MultisigAddress.Valid {
		return errors.New("User attempted to challenge without a multisig address")
	}

	balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}

	if balance.Cmp(contracts.GetAtomicTokenAmount(depositAmount)) == -1 {
		sendDM(challengeInsufficientFundsMsg)
		return nil
	}

	listing, err := contracts.GetListingFromHandle(argv[1])
	if err != nil {
		return err
	} else if listing == nil {
		sendDM(challengeNotFoundMsg)
		return nil
	} else if listing.ChallengeID.Cmp(big.NewInt(0)) != 0 {
		challenge, err := contracts.GetChallenge(listing.ChallengeID)
		if err != nil {
			return err
		}

		if !challenge.Resolved {
			sendDM(fmt.Sprintf(challengeAlreadyExistsMsg, argv[1]))
			return nil
		}
	}

	tokens := contracts.GetAtomicTokenAmount(depositAmount)
	tx, err := contracts.CreateChallenge(account.MultisigAddress.String, tokens, argv[1])
	if err != nil {
		sendDM(challengeSubmissionErrorMsg)
		return err
	}
	msg := fmt.Sprintf(challengeSuccessMsg, tx.Hash().Hex())
	sendDM(msg)

	return nil
}

func handleVote(account *models.Account, argv []string, sendDM func(string)) error {
	if len(argv) != 3 {
		sendDM(votingArgErrorMsg)
		return nil
	}

	if argv[2] != "keep" && argv[2] != "kick" {
		sendDM(votingArgErrorMsg)
		return nil
	}

	if account.MultisigAddress == nil || !account.MultisigAddress.Valid {
		return errors.New("User attempted to vote without a multisig address")
	}

	// TODO: Check token balance on the PLCR contract

	// Check to make sure there is an active poll for the given listing
	listing, err := contracts.GetListingFromHandle(argv[1])
	if err != nil {
	} else if listing == nil {
		sendDM(votingListingNotFoundMsg)
		return nil
	} else if listing.ChallengeID.Cmp(big.NewInt(0)) == 0 {
		sendDM(fmt.Sprintf(votingChallengeNotFoundMsg, argv[1]))
		return nil
	}

	// Fetch the poll they want to vote on
	poll, err := contracts.GetPoll(listing.ChallengeID)
	if err != nil {
		return err
	} else if poll == nil {
		log.Printf("Poll doesn't exist for listing: %x, challenge: %d", listing.ListingHash, listing.ChallengeID)
		sendDM(votingChallengeErrorMsg)
		return err
	}

	commitEndDate := time.Unix(poll.CommitEndDate.Int64(), 0)
	revealDate := time.Unix(poll.RevealEndDate.Int64(), 0)
	fmtRevealDate := revealDate.Format(time.RFC1123)

	// Make sure we're still in the commit phase
	if commitEndDate.Before(time.Now()) {
		sendDM(fmt.Sprintf(votingEndedMsg, fmtRevealDate))
		return nil
	}

	// Make sure they aren't voting twice
	vote, err := models.FindVote(listing.ChallengeID.Int64(), account.ID)
	if err != nil {
		return err
	} else if vote != nil {
		voteValue := "keep"
		if !vote.Vote {
			voteValue = "kick"
		}
		sendDM(fmt.Sprintf(votedAlreadyMsg, voteValue))
		return nil
	}

	voteValue := argv[2] == "keep"
	salt, _, err := contracts.PLCRCommitVote(account.MultisigAddress.String, listing.ChallengeID, big.NewInt(voteAmount), voteValue)
	if err != nil {
		return err
	}

	_, err = models.CreateVote(account, listing.ChallengeID.Int64(), salt, voteValue)
	if err != nil {
		return err
	}

	sendDM(fmt.Sprintf(votingSuccessMsg, fmtRevealDate))
	return nil
}

func handleFaucet(account *models.Account, argv []string, sendDM func(string)) error {
	// Have they hit the faucet recently?
	lastHit, err := models.LatestFaucetHit(account.ID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if lastHit != nil && now.Sub(*lastHit.Timestamp) < 24*time.Hour {
		nextHit := lastHit.Timestamp.Add(24 * time.Hour)
		sendDM(fmt.Sprintf(cannotHitFaucetMsg, humanize.Time(nextHit)))
		return nil
	}

	// If they don't yet have a multisig wallet we'll have to stop here
	if !account.MultisigAddress.Valid {
		return fmt.Errorf("Could not faucet tokens to account w/o multisig address: %d", account.ID)
	}

	atomicAmount := contracts.GetAtomicTokenAmount(faucetAmount)
	tx, err := contracts.MintTokens(account.MultisigAddress.String, atomicAmount)
	if err != nil {
		return err
	}

	err = models.RecordFaucetHit(account.ID, atomicAmount)
	if err != nil {
		return err
	}

	log.Printf("Faucet hit: %d tokens to %s (%d). TX: %s", faucetAmount, account.TwitterHandle, account.ID, tx.Hash().Hex())
	sendDM(fmt.Sprintf(hitFaucetMsg, faucetAmount, tx.Hash().Hex()))
	return nil
}

func handleVoteBalance(account *models.Account, argv []string, sendDM func(string)) error {
	if !account.MultisigAddress.Valid {
		err := errors.New("User attempted to fetch PLCR balance without a multisig address")
		sendDM(fmt.Sprintf(errorMsg, err.Error()))
		return err
	}

	balance, err := contracts.PLCRFetchBalance(account.MultisigAddress.String)
	if err != nil {
		sendDM(fmt.Sprintf(errorMsg, err.Error()))
		return err
	}

	humanBalance := contracts.GetHumanTokenAmount(balance).Int64()
	sendDM(fmt.Sprintf(voteBalanceMsg, humanBalance, humanBalance))
	return nil
}

func handleVoteDeposit(account *models.Account, argv []string, sendDM func(string)) error {
	if !account.MultisigAddress.Valid {
		err := errors.New("User attempted to PLCRDeposit without a multisig address")
		sendDM(fmt.Sprintf(errorMsg, err.Error()))
		return err
	}
	amount, err := strconv.ParseInt(argv[1], 10, 64)
	toDeposit := contracts.GetAtomicTokenAmount(amount)

	balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}

	if balance.Cmp(toDeposit) == -1 {
		msg := fmt.Sprintf(plcrDepositInsufficientFundsMsg, contracts.GetHumanTokenAmount(balance).Int64())
		sendDM(msg)
		return nil
	}

	tx, err := contracts.PLCRDeposit(account.MultisigAddress.String, toDeposit)
	if err != nil {
		sendDM(fmt.Sprintf(errorMsg, err.Error()))
		return err
	}

	sendDM(fmt.Sprintf(plcrDepositSuccessMsg, tx.Hash().Hex()))
	return nil
}

func handleVoteWithdraw(account *models.Account, argv []string, sendDM func(string)) error {
	if !account.MultisigAddress.Valid {
		err := errors.New("User attempted to PLCRwithdraw without a multisig address")
		sendDM(fmt.Sprintf(errorMsg, err.Error()))
		return err
	}
	amount, err := strconv.ParseInt(argv[1], 10, 64)
	toWithdraw := contracts.GetAtomicTokenAmount(amount)

	balance, err := contracts.PLCRFetchBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}

	if balance.Cmp(toWithdraw) == -1 {
		msg := fmt.Sprintf(plcrWithdrawInsufficientFunds, contracts.GetHumanTokenAmount(balance).Int64())
		sendDM(msg)
		return nil
	}

	tx, err := contracts.PLCRWithdraw(account.MultisigAddress.String, toWithdraw)
	if err != nil {
		sendDM(fmt.Sprintf(errorMsg, err.Error()))
		return err
	}

	sendDM(fmt.Sprintf(plcrWithdrawSuccessMsg, tx.Hash().Hex()))
	return nil
}
