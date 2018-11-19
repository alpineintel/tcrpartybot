package events

import (
	"context"
	"errors"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// StartETHListener begins listening for relevant events on the ETH blockchain
func StartETHListener(eventChan chan<- *ETHEvent, errChan chan<- error) {
	client, err := contracts.GetClientSession()
	if err != nil {
		errChan <- err
		return
	}

	// Prepare the ABI parsers we'll need
	walletFactoryABI, err := abi.JSON(strings.NewReader(string(contracts.MultiSigWalletFactoryABI)))
	if err != nil {
		errChan <- err
		return
	}

	// Get the contract address we'll be watching
	walletFactoryAddress := common.HexToAddress(os.Getenv("WALLET_FACTORY_ADDRESS"))

	// Begin the watching loop
	for {
		time.Sleep(1 * time.Second)

		// Get the previously synced block number
		val, err := models.GetKey(models.LatestSyncedBlockKey)
		if err != nil {
			errChan <- err
			continue
		} else if val == "" {
			val = os.Getenv("START_BLOCK")
		}

		blockCursor := new(big.Int)
		blockCursor, ok := blockCursor.SetString(val, 10)
		if !ok {
			errChan <- errors.New("Could not parse previous block cursor")
			continue
		}

		// Get the latest block number on the chain
		latestBlock, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			errChan <- err
		}

		// Do we need to sync?
		if latestBlock.Number.Cmp(blockCursor) == 0 {
			continue
		}

		// Finally, query the ETH node for any interesting events
		query := ethereum.FilterQuery{
			FromBlock: blockCursor,
			ToBlock:   latestBlock.Number,
			Addresses: []common.Address{walletFactoryAddress},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			errChan <- err
		}

		for _, log := range logs {
			event := contracts.MultiSigWalletFactoryContractInstantiation{}
			err := walletFactoryABI.Unpack(&event, "ContractInstantiation", log.Data)
			if err != nil {
				errChan <- err
			}

			account, err := models.FindAccountByMultisigFactoryIdentifier(event.Identifier.Int64())
			if err != nil {
				errChan <- err
				continue
			} else if account == nil {
				continue
			}

			log.Printf("Wallet %d detected at %s\n", event.Identifier, event.Instantiation.Hex())

			err = account.SetMultisigAddress(event.Instantiation.Hex())
			if err != nil {
				errChan <- err
			}
		}

		err = models.SetKey(models.LatestSyncedBlockKey, latestBlock.Number.String())
		if err != nil {
			errChan <- err
		}
	}
}
