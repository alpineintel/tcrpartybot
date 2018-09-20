package events

import (
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
	"log"
	"strings"
)

func processMention(event *Event, errChan chan<- error) {
	log.Printf("\nReceived mention from %s [%d]: %s", event.SourceHandle, event.SourceID, event.Message)
	// Filter based on let's party
	lower := strings.ToLower(event.Message)
	if strings.Contains(lower, "party") {
		processRegistration(event, errChan)
	}
}

func processRegistration(event *Event, errChan chan<- error) {
	// If they already have an account we don't need to continue
	account, err := models.FindAccountByID(event.SourceID)
	if account != nil || err != nil {
		return
	}

	log.Printf("Creating account for %s", event.SourceID)
	// Let's create a wallet for them
	key, err := crypto.GenerateKey()
	if err != nil {
		errChan <- err
		return
	}

	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey := hex.EncodeToString(key.D.Bytes())

	// Store the association between their handle and their wallet in our db
	account = &models.Account{
		TwitterHandle: event.SourceHandle,
		TwitterID:     event.SourceID,
		ETHAddress:    address,
		ETHPrivateKey: privateKey,
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
		}
	}

	// Send them a direct message asking them for the answer to a challenge
	// question
	if len(questions) == 0 {
		errChan <- errors.New("Could not fetch registration question from db")
		return
	}

	firstChallenge := challenges[0]
	err = twitter.SendDM(account.TwitterID, questions[0].Question)
	if err != nil {
		errChan <- err
		return
	}

	err = models.MarkRegistrationChallengeSent(firstChallenge.ID)
	if err != nil {
		errChan <- err
	}
}
