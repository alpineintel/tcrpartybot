package main

import (
	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"log"
)

const (
	initialDistributionAmount = 1500*10 ^ 15
)

func distributeTokens(errChan chan<- error) {
	accounts, err := models.AllAccounts()
	if err != nil {
		errChan <- err
		return
	}

	log.Println("Minting tokens to all holders on the list")
	for accounts.Next() {
		var account models.Account
		err = accounts.StructScan(&account)
		if err != nil {
			errChan <- err
			return
		}

		// Mint new tokens to the wallet's address
		tx, err := contracts.MintTokens(account.ETHAddress, initialDistributionAmount)
		if err != nil {
			errChan <- err
			return
		}

		log.Printf("\tMinted tokens to %s (%s)", account.TwitterHandle, account.ETHAddress)
		log.Printf("\ttx: %s", tx.Hash().Hex())
	}
}
