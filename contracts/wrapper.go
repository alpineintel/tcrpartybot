package contracts

import (
	"context"
	"math/big"
	"math/rand"
	"os"
	"time"

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

// GetClientSession returns an ethereum client
func GetClientSession() (*ethclient.Client, error) {
	if session != nil {
		return session, nil
	}

	session, err := ethclient.Dial(os.Getenv("ETH_NODE_URI"))
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetTokenBalance returns the balance (in the smallest delimiter) of a given wallet address
func GetTokenBalance(address string) (int64, error) {
	client, err := GetClientSession()
	if err != nil {
		return 0, err
	}

	tokenAddress := common.HexToAddress(os.Getenv("TOKEN_ADDRESS"))
	token, err := NewTCRPartyPoints(tokenAddress, client)
	if err != nil {
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
	client, err := GetClientSession()
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

func submitTransaction(multisigAddress string, tx *proxiedTransaction) (*types.Transaction, error) {
	client, err := GetClientSession()
	if err != nil {
		return nil, err
	}

	txOpts, err := setupTransactionOpts(os.Getenv("MASTER_PRIVATE_KEY"), 500000)
	if err != nil {
		return nil, err
	}

	contractAddress := common.HexToAddress(multisigAddress)
	wallet, err := NewMultiSigWallet(contractAddress, client)
	if err != nil {
		return nil, err
	}

	submitTX, err := wallet.SubmitTransaction(txOpts, tx.To, tx.Value, tx.Data)
	return submitTX, err
}

// Apply creates a new listing application on the TCR for the given twitter
// handle
func Apply(multisigAddress string, amount int64, twitterHandle string) (*types.Transaction, error) {
	// Generate a listing hash
	listingHash := make([]byte, 32)
	_, err := rand.Read(listingHash)
	if err != nil {
		return nil, err
	}

	// Convert that hash into the type it needs to be
	var txListingHash [32]byte
	copy(txListingHash[:], listingHash[0:4])
	txAmount := big.NewInt(amount)

	// Generate a new proxied transaction to be submitted via the wallet
	contractAddress := common.HexToAddress(os.Getenv("TCR_ADDRESS"))
	proxiedTX, err := newProxiedTransaction(
		contractAddress,
		RegistryABI,
		"apply",
		txListingHash,
		txAmount,
		twitterHandle,
	)

	if err != nil {
		return nil, err
	}

	tx, err := submitTransaction(multisigAddress, proxiedTX)
	return tx, err
}

// DeployWallet creates a new instance of the multisig wallet and returns the
// resulting transaction and an identifier which will be broadcast in the
// ContractInstantiation event on the wallet factory contract
func DeployWallet() (*types.Transaction, int64, error) {
	client, err := GetClientSession()
	if err != nil {
		return nil, 0, err
	}

	txOpts, err := setupTransactionOpts(os.Getenv("MASTER_PRIVATE_KEY"), 5000000)
	if err != nil {
		return nil, 0, err
	}

	contractAddress := common.HexToAddress(os.Getenv("WALLET_FACTORY_ADDRESS"))
	factory, err := NewMultiSigWalletFactory(contractAddress, client)
	if err != nil {
		return nil, 0, err
	}

	botKey, err := getPublicAddress(os.Getenv("MASTER_PRIVATE_KEY"))
	if err != nil {
		return nil, 0, err
	}

	identifier := rand.Int63()
	owners := []common.Address{botKey}
	tx, err := factory.Create(txOpts, owners, big.NewInt(1), big.NewInt(identifier))

	return tx, identifier, err
}

// AwaitTransactionConfirmation will block on a transaction until it is
// confirmed onto the network. It will return with a transaction receipt or an
// error (ie in the case of a timeout)
func AwaitTransactionConfirmation(txHash common.Hash) (*types.Receipt, error) {
	client, err := GetClientSession()
	if err != nil {
		return nil, err
	}

	// This request should expire after 3 minutes
	emptyCtx := context.Background()
	ctx, cancel := context.WithDeadline(emptyCtx, time.Now().Add(time.Minute*3))
	defer cancel()

	return client.TransactionReceipt(ctx, txHash)
}
