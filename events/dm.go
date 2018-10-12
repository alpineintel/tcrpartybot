package events

import (
	"fmt"
	goTwitter "github.com/dghubble/go-twitter/twitter"
	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
	"log"
	"strconv"
	"time"
)

// RegistrationEventData collects the required data for keeping track of
// the user registration flow into a struct
type RegistrationEventData struct {
	Event     *Event
	Challenge *models.RegistrationChallengeRegistrationQuestion
	Account   *models.Account
}

// ListenForTwitterDM is a blocking function which polls the Twitter API for
// new direct messages and sends them off to the eventChan for further
// processing as they are received.
func ListenForTwitterDM(handle string, eventChan chan<- *Event, errChan chan<- error) {
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

func processDM(event *Event, errChan chan<- error) {
	log.Printf("Received DM from %s: %s", event.SourceHandle, event.Message)

	// If they don't have an acccount, do nothing.
	account, err := models.FindAccountByHandle(event.SourceHandle)
	if account == nil || err != nil {
		return
	}

	// Are they just doing general stuff?
	if account.PassedRegistrationChallengeAt != nil {
		// They're already registered, trying to send some kind of command to
		// the bot

		var msg string
		switch event.Message {
		case "balance":
			balance, err := contracts.GetTokenBalance(account.ETHAddress)
			if err != nil {
				errChan <- err
				return
			}

			msg = fmt.Sprintf("Your balance is %d TCRP", balance)
		default:
			msg = "🎉 You're registered to party 🎉. Hang tight while we prepare to distribute our token."
		}

		err := twitter.SendDM(account.TwitterID, msg)
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
		response := fmt.Sprintf("🙅‍♀️ That's not right! %s", data.Challenge.Question)
		err := twitter.SendDM(data.Account.TwitterID, response)
		if err != nil {
			errChan <- err
		}
		return
	}

	// They got it! Let's mark the challenge as completed and give them another
	// question if they have any remaining
	err := models.MarkChallengeCompleted(data.Challenge.RegistrationChallenge.ID)
	if err != nil {
		errChan <- err
		return
	}

	// Are they completely done?
	completedChallenges, err := models.AccountHasCompletedChallenges(data.Account.ID)
	if err != nil {
		errChan <- err
		return
	}

	// Yes! Let's let them know that they're good to go.
	if completedChallenges {
		err := models.MarkAccountRegistered(data.Account.ID)
		if err != nil {
			errChan <- err
			return
		}

		err = twitter.SendDM(data.Account.TwitterID, "🎉 Awesome! You've been registered for the party. We'll reach out once we're ready to distribute TCRP tokens 🎈.")
		if err != nil {
			errChan <- err
		}
		return
	}

	// Nope, so let's send them the next challenge question
	activeChallenge, err := models.FindUnsentChallenge(data.Account.ID)
	if err != nil {
		errChan <- err
		return
	}

	err = models.MarkRegistrationChallengeSent(activeChallenge.RegistrationChallenge.ID)
	if err != nil {
		errChan <- err
		return
	}

	response := fmt.Sprintf("Nice, that's it! Here's another one for you: %s", activeChallenge.Question)
	err = twitter.SendDM(data.Account.TwitterID, response)
	if err != nil {
		errChan <- err
	}
}
