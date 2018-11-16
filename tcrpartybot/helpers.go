package main

import (
	"fmt"
	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
	"log"
)

func logErrors(errChan <-chan error) {
	for err := range errChan {
		log.Printf("\n%s", err)
	}
}

func authenticateHandle(handle string, errChan chan<- error) {
	request := &twitter.OAuthRequest{
		Handle: handle,
	}

	url, err := request.GetOAuthURL()
	if err != nil {
		errChan <- err
		return
	}

	fmt.Printf("Go to this URL to generate an access token:\n%s", url)
	fmt.Print("\nEnter PIN: ")

	_, err = fmt.Scanf("%s", &request.PIN)
	if err != nil {
		log.Println("Error receiving PIN")
		errChan <- err
		return
	}

	err = request.ReceivePIN()
	if err != nil {
		log.Println("Error fetching token")
		errChan <- err
		return
	}

	log.Println("Access token saved!")
}

func createWebhook(errChan chan<- error) {
	webhookID, err := models.GetKey("webhookID")
	if err != nil {
		errChan <- err
		return
	}

	// If we don't already have a webhook ID we should create it
	if webhookID == "" {
		id, err := twitter.CreateWebhook()
		if err != nil {
			errChan <- err
			return
		}

		log.Printf("Webhook %s created successfully", id)
		models.SetKey("webhookID", id)
	}

	// And subscribe to TCRPartyVIP's DMs
	if err := twitter.CreateSubscription(); err != nil {
		errChan <- err
		return
	}
	log.Printf("Subscription created successfully")
}

func deployWallet(errChan chan<- error) {
	tx, identifier, err := contracts.DeployWallet()
	if err != nil {
		errChan <- err
		return
	}

	log.Printf("TX: %s ID: %d", tx.Hash().String(), identifier)
}
