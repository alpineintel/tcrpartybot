package events

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tokenfoundry/tcrpartybot/models"
	"strings"
)

func processMention(event *Event, errChan chan<- error) {
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

	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey := hex.EncodeToString(key.D.Bytes())

	// Store the association between their handle and their wallet in our db
	newAccount := &models.Account{
		TwitterHandle: event.SourceHandle,
		ETHAddress:    address,
		ETHPrivateKey: privateKey,
	}
	models.CreateAccount(newAccount, errChan)
}
