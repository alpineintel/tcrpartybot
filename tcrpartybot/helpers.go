package main

import (
	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"log"
)

func logErrors(errChan <-chan error) {
	for err := range errChan {
		log.Printf("\n%s", err)
	}
}

func deployWallet(errChan chan<- error) {
	tx, identifier, err := contracts.DeployWallet()
	if err != nil {
		errChan <- err
		return
	}

	log.Printf("TX: %s ID: %d", tx.Hash().String(), identifier)
}
