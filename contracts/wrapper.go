package contracts

import (
	"context"
	"errors"
	"log"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	applicationTopicHash = "0xa27f550c3c7a7c6d8369e5383fdc7a3b4850d8ce9e20066f9d496f6989f00864"
)

var session *ethclient.Client

// RegistryListing represents a listing on the TCR
type RegistryListing struct {
	ListingHash       [32]byte
	ApplicationExpiry *big.Int
	Whitelisted       bool
	Owner             common.Address
	UnstakedDeposit   *big.Int
	ChallengeID       *big.Int
	ExitTime          *big.Int
	ExitTimeExpiry    *big.Int
}

const (
	// TokenDecimals is the number you can multiply/divide by in order to
	// arrive at a human readable TCRP balance
	TokenDecimals = 15
)

// GetAtomicTokenAmount inputs an amount in human-readable tokens and outputs
// the same amount of TCRP in its smallest denomination
func GetAtomicTokenAmount(amount int64) *big.Int {
	tokens := big.NewInt(amount)
	multi := new(big.Int).Exp(big.NewInt(10), big.NewInt(TokenDecimals), nil)
	atomicAmount := new(big.Int).Mul(tokens, multi)

	return atomicAmount
}

// GetHumanTokenAmount takes an input amount in the smallest token denomination
// and returns a value in normal TCRP
func GetHumanTokenAmount(amount *big.Int) *big.Int {
	multi := new(big.Int).Exp(big.NewInt(10), big.NewInt(TokenDecimals), nil)
	humanAmount := new(big.Int).Div(amount, multi)

	return humanAmount
}

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
func GetTokenBalance(address string) (*big.Int, error) {
	client, err := GetClientSession()
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(os.Getenv("TOKEN_ADDRESS"))
	token, err := NewTCRPartyPoints(tokenAddress, client)
	if err != nil {
		return nil, err
	}

	balance, err := token.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// MintTokens assigns new tokens to the given ETH address
func MintTokens(address string, amount *big.Int) (*types.Transaction, error) {
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

	tx, err := token.Mint(txOpts, txAddress, amount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// TCRPApprove permits a given address spend to N TCRP on a wallet's behalf
func TCRPApprove(multisigAddress string, spenderAddress string, amount *big.Int) (*types.Transaction, error) {
	log.Printf("Approving %d TCRP for spender %s", GetHumanTokenAmount(amount).Int64(), spenderAddress)
	tokenAddress := common.HexToAddress(os.Getenv("TOKEN_ADDRESS"))
	proxiedTX, err := newProxiedTransaction(
		tokenAddress,
		TCRPartyPointsABI,
		"approve",
		common.HexToAddress(spenderAddress),
		amount,
	)

	if err != nil {
		return nil, err
	}

	tx, err := submitTransaction(multisigAddress, proxiedTX)
	return tx, err
}

// Apply creates a new listing application on the TCR for the given twitter
// handle
func Apply(multisigAddress string, amount *big.Int, twitterHandle string) (*types.Transaction, error) {
	// First let's approve `amount` tokens for spending by the TCR
	approvalTX, err := TCRPApprove(multisigAddress, os.Getenv("TCR_ADDRESS"), amount)
	if err != nil {
		return nil, err
	}

	_, err = AwaitTransactionConfirmation(approvalTX.Hash())
	if err != nil {
		return nil, err
	}

	// Generate a new proxied transaction to be submitted via the wallet
	listingHash := getListingHash(twitterHandle)
	contractAddress := common.HexToAddress(os.Getenv("TCR_ADDRESS"))
	proxiedTX, err := newProxiedTransaction(
		contractAddress,
		RegistryABI,
		"apply",
		listingHash,
		amount,
		twitterHandle,
	)

	if err != nil {
		return nil, err
	}

	tx, err := submitTransaction(multisigAddress, proxiedTX)
	return tx, err
}

// Challenge initiates a new challenge against the given twitter handle
func Challenge(multisigAddress string, amount *big.Int, twitterHandle string) (*types.Transaction, error) {
	// First let's approve `amount` tokens for spending by the TCR
	approvalTX, err := TCRPApprove(multisigAddress, os.Getenv("TCR_ADDRESS"), amount)
	if err != nil {
		return nil, err
	}

	_, err = AwaitTransactionConfirmation(approvalTX.Hash())
	if err != nil {
		return nil, err
	}

	// Generate a new proxied transaction to be submitted via the wallet
	listingHash := getListingHash(twitterHandle)
	contractAddress := common.HexToAddress(os.Getenv("TCR_ADDRESS"))
	proxiedTX, err := newProxiedTransaction(
		contractAddress,
		RegistryABI,
		"challenge",
		listingHash,
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

// ApplicationWasMade returns true or false depending on whether or not a
// twitter handle is already an application or listing on the registry
func ApplicationWasMade(twitterHandle string) (bool, error) {
	listingHash := getListingHash(twitterHandle)
	client, err := GetClientSession()
	if err != nil {
		return false, err
	}

	registryAddress := common.HexToAddress(os.Getenv("TCR_ADDRESS"))
	registry, err := NewRegistry(registryAddress, client)
	if err != nil {
		return false, err
	}

	result, err := registry.AppWasMade(nil, listingHash)
	if err != nil {
		return false, err
	}

	return result, nil
}

// GetListingFromHandle returns a listing from the TCR based on a twitter handle
func GetListingFromHandle(twitterHandle string) (*RegistryListing, error) {
	// Generate a listing hash from the handle's string value
	listingHash := getListingHash(twitterHandle)
	return GetListingFromHash(listingHash)
}

// GetListingFromHash returns a listing from the TCR based on a [32]byte hash
func GetListingFromHash(listingHash [32]byte) (*RegistryListing, error) {
	client, err := GetClientSession()
	if err != nil {
		return nil, err
	}

	registryAddress := common.HexToAddress(os.Getenv("TCR_ADDRESS"))
	registry, err := NewRegistry(registryAddress, client)
	if err != nil {
		return nil, err
	}

	result, err := registry.Listings(nil, listingHash)
	if err != nil {
		return nil, err
	}

	if result.Owner.Big().Cmp(big.NewInt(0)) == 0 {
		return nil, nil
	}

	listing := RegistryListing{
		ListingHash:       listingHash,
		ApplicationExpiry: result.ApplicationExpiry,
		Whitelisted:       result.Whitelisted,
		Owner:             result.Owner,
		UnstakedDeposit:   result.UnstakedDeposit,
		ChallengeID:       result.ChallengeID,
		ExitTime:          result.ExitTime,
		ExitTimeExpiry:    result.ExitTimeExpiry,
	}
	return &listing, nil
}

// GetUnwhitelistedListings returns a list of registry listings that are still
// in their approval stage
func GetUnwhitelistedListings() ([]*RegistryListing, error) {
	client, err := GetClientSession()
	if err != nil {
		return nil, err
	}

	blockCursor := new(big.Int)
	fromBlock, ok := blockCursor.SetString(os.Getenv("START_BLOCK"), 10)
	if !ok {
		return nil, errors.New("Could not set fromBlock while getting active registry listings")
	}

	// Filter all _Application events
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			common.HexToAddress(os.Getenv("TCR_ADDRESS")),
		},
		FromBlock: fromBlock,
		Topics: [][]common.Hash{
			{common.HexToHash(applicationTopicHash)},
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var listings []*RegistryListing
	for _, event := range logs {
		// Get the listing hash
		var listingHash [32]byte
		copy(listingHash[:], event.Topics[1].Bytes()[0:32])

		// Make sure this listing still exists
		listing, err := GetListingFromHash(listingHash)
		if err != nil {
			return nil, err
		} else if listing == nil || listing.Whitelisted {
			continue
		}

		listings = append(listings, listing)
	}

	return listings, nil
}
