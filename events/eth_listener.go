package events

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/errors"
	"gitlab.com/alpinefresh/tcrpartybot/models"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type watchedContract struct {
	ABI          abi.ABI
	Address      common.Address
	TopicStructs map[common.Hash]interface{}
}

type topicResource struct {
	Name     string
	Contract watchedContract
}

func generateContracts() ([]*watchedContract, error) {
	walletFactoryABI, err := abi.JSON(strings.NewReader(string(contracts.MultiSigWalletFactoryABI)))
	if err != nil {
		return nil, errors.Wrap(err)
	}

	registryABI, err := abi.JSON(strings.NewReader(string(contracts.RegistryABI)))
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return []*watchedContract{
		&watchedContract{
			ABI:     walletFactoryABI,
			Address: common.HexToAddress(os.Getenv("WALLET_FACTORY_ADDRESS")),
		},
		&watchedContract{
			ABI:     registryABI,
			Address: common.HexToAddress(os.Getenv("TCR_ADDRESS")),
		},
	}, nil
}

// StartETHListener begins listening for relevant events on the ETH blockchain
func StartETHListener(ethEvents chan<- *ETHEvent, errChan chan<- error) {
	client, err := contracts.GetClientSession()
	if err != nil {
		errChan <- errors.Wrap(err)
		return
	}

	contracts, err := generateContracts()
	if err != nil {
		errChan <- errors.Wrap(err)
		return
	}

	// Create some structures that will make our lives easier below
	watchedAddresses := []common.Address{}
	watchedTopics := map[common.Hash]string{}

	for _, contract := range contracts {
		watchedAddresses = append(watchedAddresses, contract.Address)

		for _, event := range contract.ABI.Events {
			watchedTopics[event.Id()] = event.Name
		}
	}

	query := ethereum.FilterQuery{
		Addresses: watchedAddresses,
	}

	// Begin the watching loop
	for {
		time.Sleep(3 * time.Second)

		// Get the previously synced block number
		val, err := models.GetKey(models.LatestSyncedBlockKey)
		if err != nil {
			errChan <- errors.Wrap(err)
			continue
		} else if val == "" {
			val = os.Getenv("START_BLOCK")
		}

		blockCursor := new(big.Int)
		blockCursor, ok := blockCursor.SetString(val, 10)
		if !ok {
			errChan <- errors.Errorf("Could not parse previous block cursor")
			continue
		}

		// Get the latest block number on the chain
		latestBlock, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			errChan <- errors.Wrap(err)
			continue
		} else if latestBlock == nil || latestBlock.Number == nil {
			errChan <- errors.Errorf("Could not fetch latest block number, skipping watch loop")
			continue
		} else if latestBlock.Number.Cmp(blockCursor) == -1 {
			errChan <- errors.Errorf("ethereum node reported block number lower than last seen. Eth reported: %s, I last saw: %s", latestBlock.Number.String(), blockCursor.String())
			continue
		}

		// The filter is inclusive, therefore we should add 1 to the last seen block
		query.FromBlock = blockCursor.Add(blockCursor, big.NewInt(1))

		// Query the ETH node for any interesting events
		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			errChan <- err
			continue
		}

		for _, ethLog := range logs {
			block := big.NewInt(int64(ethLog.BlockNumber))
			if block.Cmp(blockCursor) == -1 {
				fmt.Printf("Found old event skipping")
				continue
			}

			// Look for a topic that we're interested in
			for _, topicHash := range ethLog.Topics {
				eventName := watchedTopics[topicHash]
				if eventName == "" {
					continue
				}

				// We've found one! Let's submit it to our event channel
				event := &ETHEvent{
					EventType: eventName,
					Data:      ethLog.Data,
					Topics:    ethLog.Topics,
				}

				ethEvents <- event
				go processETHEvent(event, errChan)
				go scheduleUpdateForEvent(event, errChan)
			}
		}

		err = models.SetKey(models.LatestSyncedBlockKey, latestBlock.Number.String())
		if err != nil {
			errChan <- errors.Wrap(err)
		}
	}
}
