package events

import (
	"fmt"
	"log"
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
	nominateErrorMsg             = "Whoops, looks like you forgot something. Try again with something like 'nominate [twitter handle]'. Eg: 'apply weratedogs'"
	nominateInsufficientFundsMsg = "Drat, looks like you don't have enough TCRP to nominate to the party"
	nominateSubmissionErrorMsg   = "There was an error trying to submit your nomination. The admins have been notified!"
	nominateSuccessMsg           = "We've submitted your nomination to the registry (tx: %s) Keep an eye on @TCRPartyVIP for updates."

	balanceMsg                  = "Your balance is %d TCRP"
	awaitingPartyBeginMsg       = "🎉 You're registered to party 🎉. Hang tight while we prepare to distribute our token."
	invalidChallengeResponseMsg = "🙅‍♀️ That's not right! %s"
	nextChallengeMsg            = "Nice, that's it! Here's another one for you: %s"
	preregistrationSuccessMsg   = "🎉 Awesome! You've been registered for the party. We'll reach out once we're ready to distribute TCRP tokens 🎈."
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
			if account.MultisigAddress == nil || !account.MultisigAddress.Valid {
				return
			}

			balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
			humanBalance := balance / contracts.TokenDecimals
			if err != nil {
				errChan <- err
				return
			}
			sendDM(fmt.Sprintf(balanceMsg, humanBalance))
		case "nominate":
			if len(argv) != 2 {
				sendDM(nominateErrorMsg)
				return
			}

			if account.MultisigAddress == nil || !account.MultisigAddress.Valid {
				return
			}

			balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
			humanBalance := balance / contracts.TokenDecimals
			if err != nil {
				errChan <- err
				return
			}

			if humanBalance < 500 {
				sendDM(nominateInsufficientFundsMsg)
				return
			}

			tx, err := contracts.Apply(500, argv[1])
			if err != nil {
				errChan <- err
				sendDM(nominateSubmissionErrorMsg)
				return
			}
			msg := fmt.Sprintf(nominateSuccessMsg, tx.Hash().Hex())
			sendDM(msg)
		default:
			sendDM(awaitingPartyBeginMsg)
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
