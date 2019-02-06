package events

import (
	"log"
	"os"
	"strings"
	"time"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/errors"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

const (
	noAccountMsg       = "Hmm, I haven't met you yet. If you want to join the TCR party send me a tweet that says \"let's party\""
	inactiveAccountMsg = "We're still working on activating your account. Hang tight, we'll message you shortly!"

	awaitingPartyBeginMsg       = "🎉 You're registered to party 🎉. Hang tight while we prepare to distribute our token."
	invalidChallengeResponseMsg = "🙅‍♀️ That's not right! %s"
	nextChallengeMsg            = "Nice, that's it! Here's another one for you: %s"
	preregistrationSuccessMsg   = "🎉 Awesome! You've been registered for the party. We'll reach out once we're ready to distribute TCRP tokens 🎈."
	registrationSuccessMsg      = "🎉 Awesome! Now that you're registered I'll need a few minutes to build your wallet and give you some TCR Party Points to get started with. I'll send you a DM once I'm done."
	invalidCommandMsg           = "Whoops, I don't recognize that command. Try typing help to see what you can say to me."
	helpMsg                     = "Here are the commands I recognize:\n• balance - See your TCRP balance\n• nominate [handle] = Nominate the given Twitter handle to be on the TCR\n• challenge [handle] - Begin a challenge for a listing on the TCR\n• vote [handle] [kick/keep] - Vote on an existing listing's challenge.\n• faucet – Get 100 free tokens per day.\nThose are the basics, but you can check out https://www.tcr.party for more advanced commands."
	errorMsg                    = "Yikes, we ran into an error: %s. Try tweeting at @stevenleeg for help."
	activatingAccountMsg        = "Welcome back to the party! We unfortunately had to reset the TCR after killing Ropsten, but we're glad to see you back. Give me a minute while I rebuild your wallet 👷‍♀️..."

	depositAmount     = 500
	defaultVoteWeight = 50
	faucetAmount      = 100
)

// RegistrationEventData collects the required data for keeping track of
// the user registration flow into a struct
type RegistrationEventData struct {
	Event     *TwitterEvent
	Challenge *models.RegistrationChallengeRegistrationQuestion
	Account   *models.Account
}

func parseHandle(handle string) string {
	return strings.ToLower(strings.Replace(handle, "@", "", -1))
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
				errChan <- errors.Wrap(err)
				return
			}

			if err := account.UpdateLastDMAt(); err != nil {
				errChan <- errors.Wrap(err)
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
			errChan <- errors.Wrap(err)
			return
		}
		return
	}

	// Define a helper function that will be passed around below
	sendDM := generateSendDM(account, errChan)

	if msg := os.Getenv("MAINTENANCE_MESSAGE"); msg != "" && account.TwitterHandle != "stevenleeg" {
		sendDM(msg)
		return
	}

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

	// They've registered but their account hasn't been activated yet, let's
	// activate it for them.
	if account.ActivatedAt == nil && (account.MultisigFactoryIdentifier == nil || !account.MultisigFactoryIdentifier.Valid) {
		_, identifier, err := contracts.DeployWallet()
		if err != nil {
			errChan <- errors.Wrap(err)
			return
		}

		err = account.SetMultisigFactoryIdentifier(identifier)
		if err != nil {
			errChan <- errors.Wrap(err)
			return
		}

		sendDM(activatingAccountMsg)
		return
	} else if account.ActivatedAt == nil {
		sendDM(inactiveAccountMsg)
		return
	}

	// Looks like they're good to go and are sending a command to the bot
	argv := strings.FieldsFunc(event.Message, func(c rune) bool {
		return c == ' '
	})

	// If they haven't been activated yet (ie pre-registration) then we'll
	// stop them here
	if account.ActivatedAt == nil {
		sendDM(preregistrationSuccessMsg)
		return
	}

	switch strings.ToLower(argv[0]) {
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
