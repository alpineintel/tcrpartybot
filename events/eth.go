package events

import (
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
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
