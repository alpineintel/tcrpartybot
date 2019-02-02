package contracts

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

const (
	applicationTopicHash = "0xa27f550c3c7a7c6d8369e5383fdc7a3b4850d8ce9e20066f9d496f6989f00864"
	gasLimit             = 5000000
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

// Poll represents a poll conducted on the PLCR voting contract
type Poll struct {
	CommitEndDate *big.Int
	RevealEndDate *big.Int
	VoteQuorum    *big.Int
	VotesFor      *big.Int
	VotesAgainst  *big.Int
}

// Challenge represents the challenge struct on the Registry contract
type Challenge struct {
	RewardPool  *big.Int
	Challenger  common.Address
	Resolved    bool
	Stake       *big.Int
	TotalTokens *big.Int
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

	ethSession, err := ethclient.Dial(os.Getenv("ETH_NODE_URI"))
	session = ethSession
	if err != nil {
		session = nil
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

	log.Printf("Minting %d tokens to %s", GetHumanTokenAmount(amount).Int64(), address)

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
	_, err = AwaitTransactionConfirmation(tx.Hash())
	if err != nil {
		return nil, err
	}

	return tx, err
}

// Withdraw calls the withdraw method on the registry contract, taking unstaked
// tokens out of the registry and returning it to a listing's owner.
func Withdraw(twitterHandle string, amount *big.Int) (*types.Transaction, error) {
	listingHash := getListingHash(twitterHandle)

	listing, err := GetListingFromHash(listingHash)
	if err != nil {
		return nil, err
	} else if listing == nil {
		return nil, fmt.Errorf("no listing for %s", twitterHandle)
	}

	log.Printf("Withdrawing %d tokens from listing %s for %s", GetHumanTokenAmount(amount), twitterHandle, listing.Owner.Hash().Hex())
	contractAddress := common.HexToAddress(os.Getenv("TCR_ADDRESS"))
	proxiedTX, err := newProxiedTransaction(
		contractAddress,
		RegistryABI,
		"withdraw",
		listingHash,
		amount,
	)

	if err != nil {
		return nil, err
	}

	tx, err := submitTransaction(listing.Owner.Hash().Hex(), proxiedTX)
	return tx, err
}

// CreateChallenge initiates a new challenge against the given twitter handle
func CreateChallenge(multisigAddress string, amount *big.Int, twitterHandle string) (*types.Transaction, error) {
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

	txOpts, err := setupTransactionOpts(os.Getenv("MASTER_PRIVATE_KEY"), gasLimit)
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

	var retryDelay time.Duration = 16
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)

		if err != nil && retryDelay > 128 {
			return nil, err
		} else if err == ethereum.NotFound {
			time.Sleep(retryDelay * time.Second)
			retryDelay = retryDelay * 2
			continue
		}

		return receipt, err
	}
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

func getApplicationEventFromHash(listingHash [32]byte) (*RegistryApplication, error) {
	client, err := GetClientSession()
	if err != nil {
		return nil, err
	}

	blockCursor := new(big.Int)
	fromBlock, ok := blockCursor.SetString(os.Getenv("START_BLOCK"), 10)
	if !ok {
		return nil, errors.New("Could not set fromBlock while getting active registry listings")
	}

	// Since listing data is only broadcast once (in the
	// initial application event), we need to filter out the specific event we're
	// looking for in order to fetch the associated data field.
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			common.HexToAddress(os.Getenv("TCR_ADDRESS")),
		},
		FromBlock: fromBlock,
		Topics: [][]common.Hash{
			{common.HexToHash(applicationTopicHash)},
			{common.HexToHash(common.Bytes2Hex(listingHash[:]))},
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var latestEvent *RegistryApplication
	for _, event := range logs {
		event, err := DecodeApplicationEvent(event.Topics, event.Data)
		if err != nil {
			return nil, err
		}

		latestEvent = event
	}

	return latestEvent, nil
}

// GetListingDataFromHash returns the data field (ideally a Twitter handle)
// from a given listing hash.
func GetListingDataFromHash(listingHash [32]byte) (string, error) {
	applicationEvent, err := getApplicationEventFromHash(listingHash)
	if err != nil {
		return "", nil
	}

	return applicationEvent.Data, nil
}

// GetListingOwnerFromHash returns the applicant field from the most recent
// Application event for the given listing hash. This is mostly useful for
// retreiving information about a listing that has been removed from the list.
func GetListingOwnerFromHash(listingHash [32]byte) (*common.Address, error) {
	applicationEvent, err := getApplicationEventFromHash(listingHash)
	if err != nil {
		return nil, nil
	}

	return &applicationEvent.Applicant, nil
}

// GetPLCRContractAddress returns the address of the PLCR voting contract,
// fetched directly from the registry contract
func GetPLCRContractAddress() (common.Address, error) {
	client, err := GetClientSession()
	if err != nil {
		return common.Address{}, err
	}

	// Fetch the address of the PLCR voting contract
	registryAddress := common.HexToAddress(os.Getenv("TCR_ADDRESS"))
	registry, err := NewRegistry(registryAddress, client)
	if err != nil {
		return common.Address{}, err
	}

	return registry.Voting(nil)
}

// GetAllListings returns a list of registry listings
func GetAllListings() ([]*RegistryListing, error) {
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
	seenListings := map[[32]byte]bool{}
	for _, event := range logs {
		// Get the listing hash
		var listingHash [32]byte
		copy(listingHash[:], event.Topics[1].Bytes()[0:32])

		// If we've already seen this listing before, move on
		if seenListings[listingHash] {
			continue
		}

		// Make sure this listing still exists
		listing, err := GetListingFromHash(listingHash)
		if err != nil {
			return nil, err
		} else if listing == nil {
			continue
		}

		listings = append(listings, listing)
		seenListings[listingHash] = true
	}

	return listings, nil
}

// GetWhitelistedListings returns a list of all listings that are active
// members of the TCR.
func GetWhitelistedListings() ([]*RegistryListing, error) {
	listings, err := GetAllListings()
	if err != nil {
		return nil, err
	}

	var activeListings []*RegistryListing
	for _, listing := range listings {
		if listing.Whitelisted {
			activeListings = append(activeListings, listing)
		}
	}

	return activeListings, nil
}

// UpdateStatus calls the updateStatus method on the registry contract,
// allowing a listing past its application period to be promoted to a
// whitelisted listing.
func UpdateStatus(listingHash [32]byte) (*types.Transaction, error) {
	client, err := GetClientSession()
	if err != nil {
		return nil, err
	}

	registryAddress := common.HexToAddress(os.Getenv("TCR_ADDRESS"))
	registry, err := NewRegistry(registryAddress, client)
	if err != nil {
		return nil, err
	}

	txOpts, err := setupTransactionOpts(os.Getenv("MASTER_PRIVATE_KEY"), gasLimit)
	tx, err := registry.UpdateStatus(txOpts, listingHash)
	return tx, err
}

// GetPoll returns a poll from the PLCR voting contract given a challenge ID
func GetPoll(pollID *big.Int) (*Poll, error) {
	client, err := GetClientSession()
	if err != nil {
		return nil, err
	}

	plcrAddress, err := GetPLCRContractAddress()
	if err != nil {
		return nil, err
	}

	plcr, err := NewPLCRVoting(plcrAddress, client)
	if err != nil {
		return nil, err
	}

	pollData, err := plcr.PollMap(nil, pollID)
	if err != nil {
		return nil, err
	} else if pollData.CommitEndDate == nil {
		return nil, nil
	}

	poll := Poll(pollData)
	return &poll, nil
}

// GetChallenge returns a challenge struct given its ID
func GetChallenge(challengeID *big.Int) (*Challenge, error) {
	client, err := GetClientSession()
	if err != nil {
		return nil, err
	}

	registry, err := NewRegistry(common.HexToAddress(os.Getenv("TCR_ADDRESS")), client)
	if err != nil {
		return nil, err
	}

	challengeData, err := registry.Challenges(nil, challengeID)
	if err != nil {
		return nil, err
	}

	challenge := Challenge(challengeData)
	return &challenge, nil
}

// PLCRDeposit locks up a number of tokens in the TCR's PLCR voting contract
func PLCRDeposit(multisigAddress string, amount *big.Int) (*types.Transaction, error) {
	// Fetch the PLCR contract's address
	plcrAddress, err := GetPLCRContractAddress()
	if err != nil {
		return nil, err
	}

	approvalTX, err := TCRPApprove(multisigAddress, plcrAddress.Hex(), amount)
	if err != nil {
		return nil, err
	}

	// Wait a bit for the transaction to propagate
	_, err = AwaitTransactionConfirmation(approvalTX.Hash())
	if err != nil {
		return nil, err
	}

	log.Printf("Depositing %d into PLCR contract for %s", GetHumanTokenAmount(amount).Int64(), multisigAddress)
	// Send off a request for voting rights for the given amount
	proxiedTX, err := newProxiedTransaction(
		plcrAddress,
		PLCRVotingABI,
		"requestVotingRights",
		amount,
	)
	if err != nil {
		return nil, err
	}

	return submitTransaction(multisigAddress, proxiedTX)
}

// PLCRCommitVote calls the commitVote method on the PLCR voting contract
func PLCRCommitVote(multisigAddress string, pollID *big.Int, amount *big.Int, vote bool) (int64, *types.Transaction, error) {
	plcrAddress, err := GetPLCRContractAddress()
	if err != nil {
		return 0, nil, err
	}

	// Generate a salt
	salt := rand.Int63()

	voteOption := int64(0)
	if vote {
		voteOption = 1
	}

	// Generating out secret voting hash is a bit of a process...
	uintType, err := abi.NewType("uint256")
	if err != nil {
		return 0, nil, err
	}

	arguments := abi.Arguments{
		{Type: uintType},
		{Type: uintType},
	}

	hashBytes, err := arguments.Pack(big.NewInt(voteOption), big.NewInt(salt))
	if err != nil {
		return 0, nil, err
	}

	// Hash the packed arguments
	var secretHashBuf []byte
	hashing := sha3.NewLegacyKeccak256()
	hashing.Write(hashBytes)
	secretHashBuf = hashing.Sum(secretHashBuf)

	// Convert from []byte to [32]byte
	var secretHash [32]byte
	copy(secretHash[:], secretHashBuf[0:32])

	proxiedTX, err := newProxiedTransaction(
		plcrAddress,
		PLCRVotingABI,
		"commitVote",
		pollID,
		secretHash,
		amount,
		big.NewInt(0),
	)
	if err != nil {
		return 0, nil, err
	}

	tx, err := submitTransaction(multisigAddress, proxiedTX)
	log.Printf("Committing %d vote for %s on poll %s (tx: %s)", voteOption, multisigAddress, pollID, tx.Hash().Hex())

	return salt, tx, err
}

// PLCRRevealVote reveals a previously committed vote on the PLCR voting
// contract
func PLCRRevealVote(multisigAddress string, pollID *big.Int, vote bool, salt int64) (*types.Transaction, error) {
	plcrAddress, err := GetPLCRContractAddress()
	if err != nil {
		return nil, err
	}

	intVote := int64(0)
	if vote {
		intVote = 1
	}

	proxiedTX, err := newProxiedTransaction(
		plcrAddress,
		PLCRVotingABI,
		"revealVote",
		pollID,
		big.NewInt(intVote),
		big.NewInt(salt),
	)

	return submitTransaction(multisigAddress, proxiedTX)
}

// PLCRFetchBalance returns the amount of tokens deposited into the voting
// contract
func PLCRFetchBalance(address string) (*big.Int, error) {
	client, err := GetClientSession()
	if err != nil {
		return nil, err
	}

	// Fetch the PLCR contract's address
	plcrAddress, err := GetPLCRContractAddress()
	if err != nil {
		return nil, err
	}

	plcr, err := NewPLCRVoting(plcrAddress, client)
	if err != nil {
		return nil, err
	}

	balance, err := plcr.VoteTokenBalance(nil, common.HexToAddress(address))
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// PLCRWithdraw withdraws tokens from the contract and returns it to the user's
// multisig wallet
func PLCRWithdraw(multisigAddress string, amount *big.Int) (*types.Transaction, error) {
	// Fetch the PLCR contract's address
	plcrAddress, err := GetPLCRContractAddress()
	if err != nil {
		return nil, err
	}

	log.Printf("Withdrawing %d from PLCR contract for %s", GetHumanTokenAmount(amount).Int64(), multisigAddress)
	// Send off a request for voting rights for the given amount
	proxiedTX, err := newProxiedTransaction(
		plcrAddress,
		PLCRVotingABI,
		"withdrawVotingRights",
		amount,
	)
	if err != nil {
		return nil, err
	}

	return submitTransaction(multisigAddress, proxiedTX)
}
