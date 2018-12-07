package events

import (
	"errors"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
	"log"
	"os"
	"strings"
)

func processFollow(event *TwitterEvent, errChan chan<- error) {
	err := twitter.Follow(event.SourceID)
	if err != nil {
		errChan <- err
	}
}

func processMention(event *TwitterEvent, errChan chan<- error) {
	if strings.ToLower(event.SourceHandle) == os.Getenv("PARTY_BOT_HANDLE") {
		return
	} else if strings.ToLower(event.SourceHandle) == os.Getenv("VIP_BOT_HANDLE") {
		return
	}

	log.Printf("\nReceived mention from %s [%d]: %s", event.SourceHandle, event.SourceID, event.Message)
	// Filter based on let's party
	lower := strings.ToLower(event.Message)
	if strings.Contains(lower, "party") {
		processRegistration(event, errChan)
	}
}

func processRegistration(event *TwitterEvent, errChan chan<- error) {
	// If they already have an account we don't need to continue
	account, err := models.FindAccountByTwitterID(event.SourceID)
	if account != nil {
		return
	} else if err != nil {
		errChan <- err
		return
	}

	log.Printf("Creating account for %d", event.SourceID)

	account = &models.Account{
		TwitterHandle: event.SourceHandle,
		TwitterID:     event.SourceID,
	}
	err = models.CreateAccount(account)

	if err != nil {
		errChan <- err
	}

	// Generate three registration challenges for them
	questions, err := models.FetchRandomRegistrationQuestions(models.RegistrationChallengeCount)

	if err != nil {
		errChan <- err
		return
	}

	// Create a list of challenges for the new user to complete
	challenges := make([]*models.RegistrationChallenge, models.RegistrationChallengeCount)
	for i, question := range questions {
		challenges[i], err = models.CreateRegistrationChallenge(account, &question)

		if err != nil {
			errChan <- err
			return
		}
	}

	// Send them a direct message asking them for the answer to a challenge
	// question
	if len(questions) == 0 {
		errChan <- errors.New("Could not fetch registration question from db")
		return
	}

	firstChallenge := challenges[0]
	text := "Welcome to the party! Before we put you on the VIP list, we need to make sure you're a human. Here's an easy question for you: " + questions[0].Question
	err = twitter.SendDM(account.TwitterID, text)
	if err != nil {
		errChan <- err
		return
	}

	err = firstChallenge.MarkSent()
	if err != nil {
		errChan <- err
	}
}
