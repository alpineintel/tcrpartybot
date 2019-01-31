package events

import (
	"fmt"
	"log"
	"time"

	"github.com/dustin/go-humanize"
	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
)

const (
	cannotHitFaucetMsg = "Ack, I can only let you hit the faucet once per day. Try again %s."
	hitFaucetMsg       = "You got it. %d TCRP headed your way. TX Hash: %s"
	balanceMsg         = "Your balance is %d TCRP"
)

func handleBalance(account *models.Account, argv []string, sendDM func(string)) error {
	if account.MultisigAddress == nil || !account.MultisigAddress.Valid {
		return nil
	}

	balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}

	humanBalance := contracts.GetHumanTokenAmount(balance).Int64()
	sendDM(fmt.Sprintf(balanceMsg, humanBalance))

	return nil
}

func handleFaucet(account *models.Account, argv []string, sendDM func(string)) error {
	// Have they hit the faucet recently?
	lastHit, err := models.LatestFaucetHit(account.ID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if lastHit != nil && now.Sub(*lastHit.Timestamp) < 24*time.Hour {
		nextHit := lastHit.Timestamp.Add(24 * time.Hour)
		sendDM(fmt.Sprintf(cannotHitFaucetMsg, humanize.Time(nextHit)))
		return nil
	}

	// If they don't yet have a multisig wallet we'll have to stop here
	if !account.MultisigAddress.Valid {
		return fmt.Errorf("Could not faucet tokens to account w/o multisig address: %d", account.ID)
	}

	atomicAmount := contracts.GetAtomicTokenAmount(faucetAmount)
	tx, err := contracts.MintTokens(account.MultisigAddress.String, atomicAmount)
	if err != nil {
		return err
	}

	err = models.RecordFaucetHit(account.ID, atomicAmount)
	if err != nil {
		return err
	}

	log.Printf("Faucet hit: %d tokens to %s (%d). TX: %s", faucetAmount, account.TwitterHandle, account.ID, tx.Hash().Hex())
	sendDM(fmt.Sprintf(hitFaucetMsg, faucetAmount, tx.Hash().Hex()))
	return nil
}
