package events

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"

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
	blockCursor := big.NewInt(40)
	for {
		latestBlock, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			errChan <- err
		}

		if latestBlock.Number.Cmp(blockCursor) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		query := ethereum.FilterQuery{
			FromBlock: blockCursor,
			ToBlock:   latestBlock.Number,
			Addresses: []common.Address{walletFactoryAddress},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			errChan <- err
		}

		fmt.Println("\nLOGS:")
		for _, log := range logs {
			event := contracts.MultiSigWalletFactoryContractInstantiation{}
			err := walletFactoryABI.Unpack(&event, "ContractInstantiation", log.Data)
			if err != nil {
				errChan <- err
			}
			fmt.Printf("Wallet instantiated! It's at %s\n", event.Instantiation.Hex())
		}

		blockCursor = latestBlock.Number
		time.Sleep(1 * time.Second)
	}
}
