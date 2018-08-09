package events

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tokenfoundry/tcrpartybot/models"
	"log"
	"strings"
)

func processMention(event *Event, errChan chan<- error) {
	log.Println("Received mention:", event.Message)
	// Filter based on let's party
	lower := strings.ToLower(event.Message)
	if !strings.Contains(lower, "let's party") {
		return
	}

	// Let's create a wallet for them
	key, err := crypto.GenerateKey()
	if err != nil {
		errChan <- err
		return
	}

	account := models.FindAccountByHandle(event.SourceHandle)
	if account == nil {
		log.Printf("Creating account for %s", event.SourceHandle)
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
	}
}
