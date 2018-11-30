package events

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	goTwitter "github.com/stevenleeg/go-twitter/twitter"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

const (
	votingArgErrorMsg          = "Whoops, looks like you forgot something. Try again with something like 'vote [twitter handle] yes' or 'vote [twitter handle] no'"
	votingListingNotFoundMsg   = "Hmm, I couldn't find a registry listing for that twitter handle. Are you sure they've been nominated to the registry?"
	votingChallengeNotFoundMsg = "Looks like that twitter handle doesn't have a challenge opened on it yet. If you'd like to challenge their place on the registry respond with 'challenge %s'."
	votingChallengeErrorMsg    = "There was an error committing your vote. The admins have been notified!"
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

	balanceMsg                  = "Your balance is %d TCRP"
	awaitingPartyBeginMsg       = "🎉 You're registered to party 🎉. Hang tight while we prepare to distribute our token."
	invalidChallengeResponseMsg = "🙅‍♀️ That's not right! %s"
	nextChallengeMsg            = "Nice, that's it! Here's another one for you: %s"
	preregistrationSuccessMsg   = "🎉 Awesome! You've been registered for the party. We'll reach out once we're ready to distribute TCRP tokens 🎈."

	depositAmount = 500
	voteAmount    = 50
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

func processDM(event *TwitterEvent, errChan chan<- error) {
	log.Printf("Received DM from %s: %s", event.SourceHandle, event.Message)

	// If they don't have an acccount, do nothing.
	account, err := models.FindAccountByHandle(event.SourceHandle)
	if account == nil || err != nil {
		return
	}

	sendDM := func(message string) {
		err := twitter.SendDM(account.TwitterID, message)
		if err != nil {
			errChan <- err
		}
	}

	// If they're already registered they're trying to send some kind of
	// command to the bot
	if account.PassedRegistrationChallengeAt != nil {
		// Split up the message into a command
		argv := strings.Split(event.Message, " ")

		switch argv[0] {
		case "balance":
			err = handleBalance(account, argv, sendDM)
		case "nominate":
			err = handleNomination(account, argv, sendDM)
		case "challenge":
			err = handleChallenge(account, argv, sendDM)
		case "vote":
			err = handleVote(account, argv, sendDM)
		default:
			sendDM(awaitingPartyBeginMsg)
		}

		if err != nil {
			errChan <- err
		}

		return
	}

	// Are they still in the dance (registration challenge) stage?
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
	}
}

func verifyAnswer(data RegistrationEventData, errChan chan<- error) {
	// Check to see if they've responded with the correct answer
	if data.Event.Message != data.Challenge.Answer {
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

		err = twitter.SendDM(data.Account.TwitterID, preregistrationSuccessMsg)
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
		sendDM(fmt.Sprintf(challengeAlreadyExistsMsg, argv[1]))
		return nil
	}

	tokens := contracts.GetAtomicTokenAmount(depositAmount)
	tx, err := contracts.Challenge(account.MultisigAddress.String, tokens, argv[1])
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

	if argv[2] != "yes" && argv[2] != "no" {
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

	voteValue := argv[2] == "yes"
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
