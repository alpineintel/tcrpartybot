package contracts

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var session *ethclient.Client

const (
	// TokenDecimals is the number you can multiply/divide by in order to
	// arrive at a human readable TCRP balance
	TokenDecimals = 10 ^ 15
)

func getClientSession() (*ethclient.Client, error) {
	if session != nil {
		return session, nil
	}

	session, err := ethclient.Dial(os.Getenv("ETH_NODE_URI"))
	if err != nil {
		return nil, err
	}

	return session, nil
}

func GetTokenBalance(address string) (int64, error) {
	client, err := getClientSession()
	if err != nil {
		return 0, err
	}

	tokenAddress := common.HexToAddress(os.Getenv("TOKEN_ADDRESS"))
	token, err := NewTCRPartyPoints(tokenAddress, client)
	if err != nil {
		fmt.Println("something went wrong init")
		return 0, err
	}

	balance, err := token.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		return 0, err
	}

	return balance.Int64(), nil
}

// MintTokens assigns new tokens to the given ETH address
func MintTokens(address string, amount int64) (*types.Transaction, error) {
	client, err := getClientSession()
	if err != nil {
		return nil, err
	}

	txOpts, err := setupTransactionOpts(os.Getenv("MASTER_PRIVATE_KEY"), 500000)
	if err != nil {
		return nil, err
	}

	contractAddress := common.HexToAddress(os.Getenv("TOKEN_ADDRESS"))
	token, err := NewTCRPartyPoints(contractAddress, client)
	if err != nil {
		return nil, err
	}

	txAddress := common.HexToAddress(address)
	txAmount := big.NewInt(amount)

	tx, err := token.Mint(txOpts, txAddress, txAmount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Apply creates a new listing application on the TCR for the given twitter
// handle
func Apply(privateKey string, amount int64, twitterHandle string) (*types.Transaction, error) {
	client, err := getClientSession()
	if err != nil {
		return nil, err
	}

	txOpts, err := setupTransactionOpts(privateKey, 500000)
	if err != nil {
		return nil, err
	}

	contractAddress := common.HexToAddress(os.Getenv("TCR_ADDRESS"))
	registry, err := NewRegistry(contractAddress, client)
	if err != nil {
		return nil, err
	}

	// Generate a listing hash
	listingHash := make([]byte, 32)
	_, err = rand.Read(listingHash)
	if err != nil {
		return nil, err
	}

	var txListingHash [32]byte
	copy(txListingHash[:], listingHash[0:4])
	txAmount := big.NewInt(amount)
	tx, err := registry.Apply(txOpts, txListingHash, txAmount, twitterHandle)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
