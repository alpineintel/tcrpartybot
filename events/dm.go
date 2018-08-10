package events

import (
	"fmt"
	"github.com/tokenfoundry/tcrpartybot/models"
	"github.com/tokenfoundry/tcrpartybot/twitter"
	"log"
)

type RegistrationEventData struct {
	Event     *Event
	Challenge *models.RegistrationChallengeRegistrationQuestion
	Account   *models.Account
}

func processDM(event *Event, errChan chan<- error) {
	log.Printf("Received DM from %s: %s", event.SourceHandle, event.Message)

	// If they don't have an acccount, do nothing.
	account, err := models.FindAccountByHandle(event.SourceHandle)
	if err != nil {
		return
	}

	// Are they just doing general stuff?
	fullyRegistered, err := models.AccountIsRegistered(account.ID)
	if err != nil {
		errChan <- err
		return
	}

	if fullyRegistered {
		// They're already registered, trying to send some kind of command to
		// the bot
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
		twitter.SendDM(data.Account.TwitterHandle, response)
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
	fullyRegistered, err := models.AccountIsRegistered(data.Account.ID)
	if err != nil {
		errChan <- err
		return
	}

	if fullyRegistered {
		twitter.SendDM(data.Account.TwitterHandle, "🎉 Awesome! You've been registered for the party. We'll reach out once we're ready to distribute TCRP tokens 🎈.")
		return
	}

	// Nope, so let's send them the next challenge question
	activeChallenge, err := models.FindUnsentChallenge(data.Account.ID)
	if err != nil {
		errChan <- err
		return
	}

	response := fmt.Sprintf("Nice, that's it! Here's another one for you: %s", activeChallenge.Question)
	twitter.SendDM(data.Account.TwitterHandle, response)
}
