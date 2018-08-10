package events

import (
	"github.com/tokenfoundry/tcrpartybot/models"
	"log"
)

func processDM(event *Event, errChan chan<- error) {
	log.Printf("Received DM from %s: %s", event.SourceHandle, event.Message)

	account := models.FindAccountByHandle(event.SourceHandle)

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

	if activeChallenge != nil {
		// Check to see if they've responded with the correct answer
		return
	}
}
