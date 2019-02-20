package main

import (
	"log"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/errors"
	"gitlab.com/alpinefresh/tcrpartybot/models"
)

func deployWallet(errChan chan<- error) {
	tx, identifier, err := contracts.DeployWallet()
	if err != nil {
		errChan <- errors.Wrap(err)
		return
	}

	log.Printf("TX: %s ID: %d", tx.Hash().String(), identifier)
}

func deleteAccount(twitterHandle string) error {
	account, err := models.FindAccountByHandle(twitterHandle)
	if err != nil {
		return err
	} else if account == nil {
		return errors.Errorf("Could not find account for handle %s", twitterHandle)
	}

	return account.Destroy()
}
