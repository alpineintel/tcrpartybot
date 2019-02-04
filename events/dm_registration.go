package events

import (
	"fmt"
	"os"
	"strings"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

func verifyAnswer(data RegistrationEventData, errChan chan<- error) {
	// Check to see if they've responded with the correct answer
	if strings.ToLower(data.Event.Message) != strings.ToLower(data.Challenge.Answer) {
		response := fmt.Sprintf(invalidChallengeResponseMsg, data.Challenge.Question)
		err := twitter.SendDM(data.Account.TwitterID, response)
		if err != nil {
			errChan <- err
		}
		return
	}

	// They got it! Let's mark the challenge as completed and give them another
	// question if they have any remaining
	err := data.Challenge.MarkCompleted()
	if err != nil {
		errChan <- err
		return
	}

	// Are they completely done?
	completedChallenges, err := data.Account.HasCompletedChallenges()
	if err != nil {
		errChan <- err
		return
	}

	// Yes! Let's let them know that they're good to go.
	if completedChallenges {
		err := data.Account.MarkRegistered()
		if err != nil {
			errChan <- err
			return
		}

		// Send them a DM letting them know they're good to go
		msg := registrationSuccessMsg
		if os.Getenv("PREREGISTRATION") == "true" {
			msg = preregistrationSuccessMsg
		}
		err = twitter.SendDM(data.Account.TwitterID, msg)
		if err != nil {
			errChan <- err
		}

		// Let's also register a new multisig wallet for them
		go provisionWallet(data.Account, errChan)
		return
	}

	// Nope, so let's send them the next challenge question
	activeChallenge, err := models.FindUnsentChallenge(data.Account.ID)
	if err != nil {
		errChan <- err
		return
	}

	err = activeChallenge.MarkSent()
	if err != nil {
		errChan <- err
		return
	}

	response := fmt.Sprintf(nextChallengeMsg, activeChallenge.Question)
	err = twitter.SendDM(data.Account.TwitterID, response)
	if err != nil {
		errChan <- err
	}
}

func provisionWallet(account *models.Account, errChan chan<- error) {
	tx, identifier, err := contracts.DeployWallet()
	if err != nil {
		errChan <- err
		return
	}

	err = account.SetMultisigFactoryIdentifier(identifier)
	if err != nil {
		errChan <- err
		return
	}

	receipt, err := contracts.AwaitTransactionConfirmation(tx.Hash())
	if err != nil {
		errChan <- err
		return
	}

	// Make sure the wallet creation actually succeeded
	if receipt.Status == ethTypes.ReceiptStatusFailed {
		errChan <- fmt.Errorf("Could not create multisig wallet for account %d", account.ID)
		return
	}
}
