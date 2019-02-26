package events

import (
	"context"
	"fmt"
	"log"
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
	"github.com/ethereum/go-ethereum/core/types"
)

type watchedContract struct {
	ABI          abi.ABI
	Address      common.Address
	TopicStructs map[common.Hash]interface{}
}

// WatchedContractGenerator is a function that should be passed into
// ETHListener to generate a list of contracts that need to be watched.
type WatchedContractGenerator func() ([]*watchedContract, error)

func generateBotContracts() ([]*watchedContract, error) {
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

func generateAllContracts() ([]*watchedContract, error) {
	walletFactoryABI, err := abi.JSON(strings.NewReader(string(contracts.MultiSigWalletFactoryABI)))
	if err != nil {
		return nil, errors.Wrap(err)
	}

	registryABI, err := abi.JSON(strings.NewReader(string(contracts.RegistryABI)))
	if err != nil {
		return nil, errors.Wrap(err)
	}

	tokenABI, err := abi.JSON(strings.NewReader(string(contracts.TCRPartyPointsABI)))
	if err != nil {
		return nil, errors.Wrap(err)
	}

	plcrABI, err := abi.JSON(strings.NewReader(string(contracts.PLCRVotingABI)))
	if err != nil {
		return nil, errors.Wrap(err)
	}

	plcrAddress, err := contracts.GetPLCRContractAddress()
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
		&watchedContract{
			ABI:     tokenABI,
			Address: common.HexToAddress(os.Getenv("TOKEN_ADDRESS")),
		},
		&watchedContract{
			ABI:     plcrABI,
			Address: plcrAddress,
		},
	}, nil
}

// StartBotListener starts ETHListener, passing in the set of contracts needed
// for the bot to operate
func StartBotListener(ethEvents chan<- *ETHEvent, errChan chan<- error) {
	ETHListener(ethEvents, errChan, generateBotContracts, models.LatestSyncedBlockKey)
}

// StartScannerListener starts ETH Listener, passing in all contracts that make
// up the TCR for logging purposes.
func StartScannerListener(ethEvents chan<- *ETHEvent, errChan chan<- error) {
	ETHListener(ethEvents, errChan, generateAllContracts, models.LatestLoggedBlockKey)
}

// ETHListener begins listening for relevant events on the ETH blockchain
func ETHListener(ethEvents chan<- *ETHEvent, errChan chan<- error, generateContracts WatchedContractGenerator, lastSyncedKey string) {
	generatedContracts, err := generateContracts()
	if err != nil {
		errChan <- errors.Wrap(err)
		return
	}

	// Create some structures that will make our lives easier below
	watchedAddresses := []common.Address{}
	watchedTopics := map[common.Hash]string{}

	for _, contract := range generatedContracts {
		watchedAddresses = append(watchedAddresses, contract.Address)

		for _, event := range contract.ABI.Events {
			watchedTopics[event.Id()] = event.Name
		}
	}

	query := ethereum.FilterQuery{
		Addresses: watchedAddresses,
	}

	syncing := false

	// Begin the watching loop
	for {
		if !syncing {
			time.Sleep(3 * time.Second)
		} else {
			syncing = false
		}

		client, err := contracts.GetClientSession()
		if err != nil {
			errChan <- errors.Wrap(err)
			return
		}

		// Get the previously synced block number
		val, err := models.GetKey(lastSyncedKey)
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
		latestBlockHeader, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			errChan <- errors.Wrap(err)
			continue
		} else if latestBlockHeader == nil || latestBlockHeader.Number == nil {
			errChan <- errors.Errorf("Could not fetch latest block number, skipping watch loop")
			continue
		} else if latestBlockHeader.Number.Cmp(blockCursor) == -1 {
			errChan <- errors.Errorf("ethereum node reported block number lower than last seen. Eth reported: %s, I last saw: %s", latestBlockHeader.Number.String(), blockCursor.String())
			continue
		}

		diff := new(big.Int)
		diff.Sub(latestBlockHeader.Number, blockCursor)
		toBlock := new(big.Int)
		toBlock.Set(latestBlockHeader.Number)

		// Batch sync into chunks of 10k blocks
		if diff.Cmp(big.NewInt(10000)) == 1 {
			toBlock.Add(blockCursor, big.NewInt(10000))
			log.Printf("Syncing from block %s to %s", blockCursor.String(), toBlock.String())
			syncing = true
		}

		// The filter is inclusive, therefore we should add 1 to the last seen block
		query.FromBlock = blockCursor.Add(blockCursor, big.NewInt(1))
		query.ToBlock = toBlock

		// Query the ETH node for any interesting events
		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			errChan <- errors.Wrap(err)
			continue
		}

		// Create a cache for block header information, allowing us to share
		// block timestamp data between multiple events within the same block
		blockInfo := map[uint64]*types.Header{}
		for _, ethLog := range logs {
			// Get the block timestamp
			var blockTime *time.Time
			if header := blockInfo[ethLog.BlockNumber]; header != nil {
				time := time.Unix(header.Time.Int64(), 0)
				blockTime = &time
			} else {
				block := big.NewInt(int64(ethLog.BlockNumber))
				if block.Cmp(blockCursor) == -1 {
					fmt.Printf("Found old event skipping")
					continue
				}

				header, err := client.HeaderByNumber(context.Background(), block)
				if err != nil {
					errChan <- errors.Wrap(err)
				}

				if header != nil {
					time := time.Unix(header.Time.Int64(), 0)
					blockTime = &time
					blockInfo[ethLog.BlockNumber] = header
				} else {
					// If INFURA's API fails us let's just mark the block time
					// as nil and be on our merry way
					blockTime = nil
				}
			}

			// Look for a topic that we're interested in
			for _, topicHash := range ethLog.Topics {
				eventName := watchedTopics[topicHash]
				if eventName == "" {
					continue
				}

				// We've found one! Let's submit it to our event channel
				event := &ETHEvent{
					EventType:   eventName,
					CreatedAt:   blockTime,
					BlockNumber: ethLog.BlockNumber,
					TxHash:      ethLog.TxHash.Hex(),
					TxIndex:     ethLog.TxIndex,
					LogIndex:    ethLog.Index,
					Data:        ethLog.Data,
					Topics:      ethLog.Topics,
				}

				ethEvents <- event
			}
		}

		err = models.SetKey(lastSyncedKey, toBlock.String())
		if err != nil {
			errChan <- errors.Wrap(err)
		}
	}
}
