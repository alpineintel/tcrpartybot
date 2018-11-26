package events

import (
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

const (
	newApplicationWithHandleTweet    = "New #TCRParty listing! @%s has nominated @%s to be on the list for %s TCRP. Challenge this application by DMing 'challenge @%s'."
	newApplicationWithoutHandleTweet = "New #TCRParty listing! @%s has been nominated to be on the list for %s TCRP. Challenge this application by DMing 'challenge @%s'."
)

func processMultisigWalletCreation(event *ETHEvent) error {
	walletFactoryABI, err := abi.JSON(strings.NewReader(string(contracts.MultiSigWalletFactoryABI)))
	if err != nil {
		return err
	}

	data := contracts.MultiSigWalletFactoryContractInstantiation{}
	err = walletFactoryABI.Unpack(&data, "ContractInstantiation", event.Data)

	account, err := models.FindAccountByMultisigFactoryIdentifier(data.Identifier.Int64())
	if err != nil {
		return err
	} else if account == nil {
		return nil
	}

	err = account.SetMultisigAddress(data.Instantiation.Hex())
	if err != nil {
		return err
	}

	log.Printf("Wallet at %s linked to %s\n", data.Instantiation.Hex(), account.TwitterHandle)
	return nil
}

func processNewApplication(event *ETHEvent) error {
	registryABI, err := abi.JSON(strings.NewReader(string(contracts.RegistryABI)))
	if err != nil {
		return err
	}

	data := contracts.RegistryApplication{}
	err = registryABI.Unpack(&data, "_Application", event.Data)
	if err != nil {
		return err
	}

	// See if we can find an applicant in our database
	log.Printf("New application from %s for %s (hash: %s)", data.Applicant.Hex(), data.Data, data.ListingHash)
	account, err := models.FindAccountByMultisigAddress(data.Applicant.Hex())
	if err != nil {
		return err
	}

	tweet := ""
	depositAmount := contracts.GetHumanTokenAmount(data.Deposit).String()
	if account != nil {
		tweet = fmt.Sprintf(
			newApplicationWithHandleTweet,
			account.TwitterHandle,
			data.Data,
			depositAmount,
			data.Data,
		)
	} else {
		tweet = fmt.Sprintf(
			newApplicationWithoutHandleTweet,
			data.Data,
			depositAmount,
			data.Data,
		)
	}

	return twitter.SendTweet(twitter.VIPBotHandle, tweet)
}
