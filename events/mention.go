package events

import (
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tokenfoundry/tcrpartybot/models"
	"github.com/tokenfoundry/tcrpartybot/twitter"
	"log"
	"strings"
)

func processRegistration(event *Event, errChan chan<- error) {
	// If they already have an account we don't need to continue
	account := models.FindAccountByHandle(event.SourceHandle)
	if account != nil {
		return
	}

	log.Printf("Creating account for %s", event.SourceHandle)
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
		ETHAddress:    address,
		ETHPrivateKey: privateKey,
	}
	err = models.CreateAccount(account)

	if err != nil {
		errChan <- err
	}

	// Send them a direct message asking them for the answer to a challenge
	// question
	question := models.FetchRandomChallengeQuestion()
	if question == nil {
		errChan <- errors.New("Could not fetch challenge question from db")
		return
	}

	twitter.SendDirectMessage(account.TwitterHandle, question.Question)
}

func processMention(event *Event, errChan chan<- error) {
	log.Println("Received mention:", event.Message)
	// Filter based on let's party
	lower := strings.ToLower(event.Message)
	if strings.Contains(lower, "let's party") {
		processRegistration(event, errChan)
	}
}
