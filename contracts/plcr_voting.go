// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// PLCRVotingABI is the input ABI used to generate the binding from.
const PLCRVotingABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"INITIAL_POLL_NONCE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"voteTokenBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"pollMap\",\"outputs\":[{\"name\":\"commitEndDate\",\"type\":\"uint256\"},{\"name\":\"revealEndDate\",\"type\":\"uint256\"},{\"name\":\"voteQuorum\",\"type\":\"uint256\"},{\"name\":\"votesFor\",\"type\":\"uint256\"},{\"name\":\"votesAgainst\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pollNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"pollID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"numTokens\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"voter\",\"type\":\"address\"}],\"name\":\"_VoteCommitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"pollID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"numTokens\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"votesFor\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"votesAgainst\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"choice\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"salt\",\"type\":\"uint256\"}],\"name\":\"_VoteRevealed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"voteQuorum\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"commitEndDate\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"revealEndDate\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"pollID\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"creator\",\"type\":\"address\"}],\"name\":\"_PollCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"numTokens\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"voter\",\"type\":\"address\"}],\"name\":\"_VotingRightsGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"numTokens\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"voter\",\"type\":\"address\"}],\"name\":\"_VotingRightsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"pollID\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"voter\",\"type\":\"address\"}],\"name\":\"_TokensRescued\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"init\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_numTokens\",\"type\":\"uint256\"}],\"name\":\"requestVotingRights\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_numTokens\",\"type\":\"uint256\"}],\"name\":\"withdrawVotingRights\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"rescueTokens\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_pollIDs\",\"type\":\"uint256[]\"}],\"name\":\"rescueTokensInMultiplePolls\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_pollID\",\"type\":\"uint256\"},{\"name\":\"_secretHash\",\"type\":\"bytes32\"},{\"name\":\"_numTokens\",\"type\":\"uint256\"},{\"name\":\"_prevPollID\",\"type\":\"uint256\"}],\"name\":\"commitVote\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_pollIDs\",\"type\":\"uint256[]\"},{\"name\":\"_secretHashes\",\"type\":\"bytes32[]\"},{\"name\":\"_numsTokens\",\"type\":\"uint256[]\"},{\"name\":\"_prevPollIDs\",\"type\":\"uint256[]\"}],\"name\":\"commitVotes\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_prevID\",\"type\":\"uint256\"},{\"name\":\"_nextID\",\"type\":\"uint256\"},{\"name\":\"_voter\",\"type\":\"address\"},{\"name\":\"_numTokens\",\"type\":\"uint256\"}],\"name\":\"validPosition\",\"outputs\":[{\"name\":\"valid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_pollID\",\"type\":\"uint256\"},{\"name\":\"_voteOption\",\"type\":\"uint256\"},{\"name\":\"_salt\",\"type\":\"uint256\"}],\"name\":\"revealVote\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_pollIDs\",\"type\":\"uint256[]\"},{\"name\":\"_voteOptions\",\"type\":\"uint256[]\"},{\"name\":\"_salts\",\"type\":\"uint256[]\"}],\"name\":\"revealVotes\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_voter\",\"type\":\"address\"},{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"getNumPassingTokens\",\"outputs\":[{\"name\":\"correctVotes\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_voteQuorum\",\"type\":\"uint256\"},{\"name\":\"_commitDuration\",\"type\":\"uint256\"},{\"name\":\"_revealDuration\",\"type\":\"uint256\"}],\"name\":\"startPoll\",\"outputs\":[{\"name\":\"pollID\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"isPassed\",\"outputs\":[{\"name\":\"passed\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"getTotalNumberOfTokensForWinningOption\",\"outputs\":[{\"name\":\"numTokens\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"pollEnded\",\"outputs\":[{\"name\":\"ended\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"commitPeriodActive\",\"outputs\":[{\"name\":\"active\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"revealPeriodActive\",\"outputs\":[{\"name\":\"active\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_voter\",\"type\":\"address\"},{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"didCommit\",\"outputs\":[{\"name\":\"committed\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_voter\",\"type\":\"address\"},{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"didReveal\",\"outputs\":[{\"name\":\"revealed\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"pollExists\",\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_voter\",\"type\":\"address\"},{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"getCommitHash\",\"outputs\":[{\"name\":\"commitHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_voter\",\"type\":\"address\"},{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"getNumTokens\",\"outputs\":[{\"name\":\"numTokens\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_voter\",\"type\":\"address\"}],\"name\":\"getLastNode\",\"outputs\":[{\"name\":\"pollID\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_voter\",\"type\":\"address\"}],\"name\":\"getLockedTokens\",\"outputs\":[{\"name\":\"numTokens\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_voter\",\"type\":\"address\"},{\"name\":\"_numTokens\",\"type\":\"uint256\"},{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"getInsertPointForNumTokens\",\"outputs\":[{\"name\":\"prevNode\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_terminationDate\",\"type\":\"uint256\"}],\"name\":\"isExpired\",\"outputs\":[{\"name\":\"expired\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_user\",\"type\":\"address\"},{\"name\":\"_pollID\",\"type\":\"uint256\"}],\"name\":\"attrUUID\",\"outputs\":[{\"name\":\"UUID\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// PLCRVoting is an auto generated Go binding around an Ethereum contract.
type PLCRVoting struct {
	PLCRVotingCaller     // Read-only binding to the contract
	PLCRVotingTransactor // Write-only binding to the contract
	PLCRVotingFilterer   // Log filterer for contract events
}

// PLCRVotingCaller is an auto generated read-only Go binding around an Ethereum contract.
type PLCRVotingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PLCRVotingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PLCRVotingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PLCRVotingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PLCRVotingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PLCRVotingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PLCRVotingSession struct {
	Contract     *PLCRVoting       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PLCRVotingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PLCRVotingCallerSession struct {
	Contract *PLCRVotingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// PLCRVotingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PLCRVotingTransactorSession struct {
	Contract     *PLCRVotingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// PLCRVotingRaw is an auto generated low-level Go binding around an Ethereum contract.
type PLCRVotingRaw struct {
	Contract *PLCRVoting // Generic contract binding to access the raw methods on
}

// PLCRVotingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PLCRVotingCallerRaw struct {
	Contract *PLCRVotingCaller // Generic read-only contract binding to access the raw methods on
}

// PLCRVotingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PLCRVotingTransactorRaw struct {
	Contract *PLCRVotingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPLCRVoting creates a new instance of PLCRVoting, bound to a specific deployed contract.
func NewPLCRVoting(address common.Address, backend bind.ContractBackend) (*PLCRVoting, error) {
	contract, err := bindPLCRVoting(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PLCRVoting{PLCRVotingCaller: PLCRVotingCaller{contract: contract}, PLCRVotingTransactor: PLCRVotingTransactor{contract: contract}, PLCRVotingFilterer: PLCRVotingFilterer{contract: contract}}, nil
}

// NewPLCRVotingCaller creates a new read-only instance of PLCRVoting, bound to a specific deployed contract.
func NewPLCRVotingCaller(address common.Address, caller bind.ContractCaller) (*PLCRVotingCaller, error) {
	contract, err := bindPLCRVoting(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PLCRVotingCaller{contract: contract}, nil
}

// NewPLCRVotingTransactor creates a new write-only instance of PLCRVoting, bound to a specific deployed contract.
func NewPLCRVotingTransactor(address common.Address, transactor bind.ContractTransactor) (*PLCRVotingTransactor, error) {
	contract, err := bindPLCRVoting(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PLCRVotingTransactor{contract: contract}, nil
}

// NewPLCRVotingFilterer creates a new log filterer instance of PLCRVoting, bound to a specific deployed contract.
func NewPLCRVotingFilterer(address common.Address, filterer bind.ContractFilterer) (*PLCRVotingFilterer, error) {
	contract, err := bindPLCRVoting(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PLCRVotingFilterer{contract: contract}, nil
}

// bindPLCRVoting binds a generic wrapper to an already deployed contract.
func bindPLCRVoting(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PLCRVotingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PLCRVoting *PLCRVotingRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PLCRVoting.Contract.PLCRVotingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PLCRVoting *PLCRVotingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PLCRVoting.Contract.PLCRVotingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PLCRVoting *PLCRVotingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PLCRVoting.Contract.PLCRVotingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PLCRVoting *PLCRVotingCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PLCRVoting.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PLCRVoting *PLCRVotingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PLCRVoting.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PLCRVoting *PLCRVotingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PLCRVoting.Contract.contract.Transact(opts, method, params...)
}

// INITIALPOLLNONCE is a free data retrieval call binding the contract method 0x2173a10f.
//
// Solidity: function INITIAL_POLL_NONCE() constant returns(uint256)
func (_PLCRVoting *PLCRVotingCaller) INITIALPOLLNONCE(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "INITIAL_POLL_NONCE")
	return *ret0, err
}

// INITIALPOLLNONCE is a free data retrieval call binding the contract method 0x2173a10f.
//
// Solidity: function INITIAL_POLL_NONCE() constant returns(uint256)
func (_PLCRVoting *PLCRVotingSession) INITIALPOLLNONCE() (*big.Int, error) {
	return _PLCRVoting.Contract.INITIALPOLLNONCE(&_PLCRVoting.CallOpts)
}

// INITIALPOLLNONCE is a free data retrieval call binding the contract method 0x2173a10f.
//
// Solidity: function INITIAL_POLL_NONCE() constant returns(uint256)
func (_PLCRVoting *PLCRVotingCallerSession) INITIALPOLLNONCE() (*big.Int, error) {
	return _PLCRVoting.Contract.INITIALPOLLNONCE(&_PLCRVoting.CallOpts)
}

// AttrUUID is a free data retrieval call binding the contract method 0xa1103f37.
//
// Solidity: function attrUUID(_user address, _pollID uint256) constant returns(UUID bytes32)
func (_PLCRVoting *PLCRVotingCaller) AttrUUID(opts *bind.CallOpts, _user common.Address, _pollID *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "attrUUID", _user, _pollID)
	return *ret0, err
}

// AttrUUID is a free data retrieval call binding the contract method 0xa1103f37.
//
// Solidity: function attrUUID(_user address, _pollID uint256) constant returns(UUID bytes32)
func (_PLCRVoting *PLCRVotingSession) AttrUUID(_user common.Address, _pollID *big.Int) ([32]byte, error) {
	return _PLCRVoting.Contract.AttrUUID(&_PLCRVoting.CallOpts, _user, _pollID)
}

// AttrUUID is a free data retrieval call binding the contract method 0xa1103f37.
//
// Solidity: function attrUUID(_user address, _pollID uint256) constant returns(UUID bytes32)
func (_PLCRVoting *PLCRVotingCallerSession) AttrUUID(_user common.Address, _pollID *big.Int) ([32]byte, error) {
	return _PLCRVoting.Contract.AttrUUID(&_PLCRVoting.CallOpts, _user, _pollID)
}

// CommitPeriodActive is a free data retrieval call binding the contract method 0xa4439dc5.
//
// Solidity: function commitPeriodActive(_pollID uint256) constant returns(active bool)
func (_PLCRVoting *PLCRVotingCaller) CommitPeriodActive(opts *bind.CallOpts, _pollID *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "commitPeriodActive", _pollID)
	return *ret0, err
}

// CommitPeriodActive is a free data retrieval call binding the contract method 0xa4439dc5.
//
// Solidity: function commitPeriodActive(_pollID uint256) constant returns(active bool)
func (_PLCRVoting *PLCRVotingSession) CommitPeriodActive(_pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.CommitPeriodActive(&_PLCRVoting.CallOpts, _pollID)
}

// CommitPeriodActive is a free data retrieval call binding the contract method 0xa4439dc5.
//
// Solidity: function commitPeriodActive(_pollID uint256) constant returns(active bool)
func (_PLCRVoting *PLCRVotingCallerSession) CommitPeriodActive(_pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.CommitPeriodActive(&_PLCRVoting.CallOpts, _pollID)
}

// DidCommit is a free data retrieval call binding the contract method 0x7f97e836.
//
// Solidity: function didCommit(_voter address, _pollID uint256) constant returns(committed bool)
func (_PLCRVoting *PLCRVotingCaller) DidCommit(opts *bind.CallOpts, _voter common.Address, _pollID *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "didCommit", _voter, _pollID)
	return *ret0, err
}

// DidCommit is a free data retrieval call binding the contract method 0x7f97e836.
//
// Solidity: function didCommit(_voter address, _pollID uint256) constant returns(committed bool)
func (_PLCRVoting *PLCRVotingSession) DidCommit(_voter common.Address, _pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.DidCommit(&_PLCRVoting.CallOpts, _voter, _pollID)
}

// DidCommit is a free data retrieval call binding the contract method 0x7f97e836.
//
// Solidity: function didCommit(_voter address, _pollID uint256) constant returns(committed bool)
func (_PLCRVoting *PLCRVotingCallerSession) DidCommit(_voter common.Address, _pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.DidCommit(&_PLCRVoting.CallOpts, _voter, _pollID)
}

// DidReveal is a free data retrieval call binding the contract method 0xaa7ca464.
//
// Solidity: function didReveal(_voter address, _pollID uint256) constant returns(revealed bool)
func (_PLCRVoting *PLCRVotingCaller) DidReveal(opts *bind.CallOpts, _voter common.Address, _pollID *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "didReveal", _voter, _pollID)
	return *ret0, err
}

// DidReveal is a free data retrieval call binding the contract method 0xaa7ca464.
//
// Solidity: function didReveal(_voter address, _pollID uint256) constant returns(revealed bool)
func (_PLCRVoting *PLCRVotingSession) DidReveal(_voter common.Address, _pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.DidReveal(&_PLCRVoting.CallOpts, _voter, _pollID)
}

// DidReveal is a free data retrieval call binding the contract method 0xaa7ca464.
//
// Solidity: function didReveal(_voter address, _pollID uint256) constant returns(revealed bool)
func (_PLCRVoting *PLCRVotingCallerSession) DidReveal(_voter common.Address, _pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.DidReveal(&_PLCRVoting.CallOpts, _voter, _pollID)
}

// GetCommitHash is a free data retrieval call binding the contract method 0xd901402b.
//
// Solidity: function getCommitHash(_voter address, _pollID uint256) constant returns(commitHash bytes32)
func (_PLCRVoting *PLCRVotingCaller) GetCommitHash(opts *bind.CallOpts, _voter common.Address, _pollID *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "getCommitHash", _voter, _pollID)
	return *ret0, err
}

// GetCommitHash is a free data retrieval call binding the contract method 0xd901402b.
//
// Solidity: function getCommitHash(_voter address, _pollID uint256) constant returns(commitHash bytes32)
func (_PLCRVoting *PLCRVotingSession) GetCommitHash(_voter common.Address, _pollID *big.Int) ([32]byte, error) {
	return _PLCRVoting.Contract.GetCommitHash(&_PLCRVoting.CallOpts, _voter, _pollID)
}

// GetCommitHash is a free data retrieval call binding the contract method 0xd901402b.
//
// Solidity: function getCommitHash(_voter address, _pollID uint256) constant returns(commitHash bytes32)
func (_PLCRVoting *PLCRVotingCallerSession) GetCommitHash(_voter common.Address, _pollID *big.Int) ([32]byte, error) {
	return _PLCRVoting.Contract.GetCommitHash(&_PLCRVoting.CallOpts, _voter, _pollID)
}

// GetInsertPointForNumTokens is a free data retrieval call binding the contract method 0x2c052031.
//
// Solidity: function getInsertPointForNumTokens(_voter address, _numTokens uint256, _pollID uint256) constant returns(prevNode uint256)
func (_PLCRVoting *PLCRVotingCaller) GetInsertPointForNumTokens(opts *bind.CallOpts, _voter common.Address, _numTokens *big.Int, _pollID *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "getInsertPointForNumTokens", _voter, _numTokens, _pollID)
	return *ret0, err
}

// GetInsertPointForNumTokens is a free data retrieval call binding the contract method 0x2c052031.
//
// Solidity: function getInsertPointForNumTokens(_voter address, _numTokens uint256, _pollID uint256) constant returns(prevNode uint256)
func (_PLCRVoting *PLCRVotingSession) GetInsertPointForNumTokens(_voter common.Address, _numTokens *big.Int, _pollID *big.Int) (*big.Int, error) {
	return _PLCRVoting.Contract.GetInsertPointForNumTokens(&_PLCRVoting.CallOpts, _voter, _numTokens, _pollID)
}

// GetInsertPointForNumTokens is a free data retrieval call binding the contract method 0x2c052031.
//
// Solidity: function getInsertPointForNumTokens(_voter address, _numTokens uint256, _pollID uint256) constant returns(prevNode uint256)
func (_PLCRVoting *PLCRVotingCallerSession) GetInsertPointForNumTokens(_voter common.Address, _numTokens *big.Int, _pollID *big.Int) (*big.Int, error) {
	return _PLCRVoting.Contract.GetInsertPointForNumTokens(&_PLCRVoting.CallOpts, _voter, _numTokens, _pollID)
}

// GetLastNode is a free data retrieval call binding the contract method 0x427fa1d2.
//
// Solidity: function getLastNode(_voter address) constant returns(pollID uint256)
func (_PLCRVoting *PLCRVotingCaller) GetLastNode(opts *bind.CallOpts, _voter common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "getLastNode", _voter)
	return *ret0, err
}

// GetLastNode is a free data retrieval call binding the contract method 0x427fa1d2.
//
// Solidity: function getLastNode(_voter address) constant returns(pollID uint256)
func (_PLCRVoting *PLCRVotingSession) GetLastNode(_voter common.Address) (*big.Int, error) {
	return _PLCRVoting.Contract.GetLastNode(&_PLCRVoting.CallOpts, _voter)
}

// GetLastNode is a free data retrieval call binding the contract method 0x427fa1d2.
//
// Solidity: function getLastNode(_voter address) constant returns(pollID uint256)
func (_PLCRVoting *PLCRVotingCallerSession) GetLastNode(_voter common.Address) (*big.Int, error) {
	return _PLCRVoting.Contract.GetLastNode(&_PLCRVoting.CallOpts, _voter)
}

// GetLockedTokens is a free data retrieval call binding the contract method 0x6b2d95d4.
//
// Solidity: function getLockedTokens(_voter address) constant returns(numTokens uint256)
func (_PLCRVoting *PLCRVotingCaller) GetLockedTokens(opts *bind.CallOpts, _voter common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "getLockedTokens", _voter)
	return *ret0, err
}

// GetLockedTokens is a free data retrieval call binding the contract method 0x6b2d95d4.
//
// Solidity: function getLockedTokens(_voter address) constant returns(numTokens uint256)
func (_PLCRVoting *PLCRVotingSession) GetLockedTokens(_voter common.Address) (*big.Int, error) {
	return _PLCRVoting.Contract.GetLockedTokens(&_PLCRVoting.CallOpts, _voter)
}

// GetLockedTokens is a free data retrieval call binding the contract method 0x6b2d95d4.
//
// Solidity: function getLockedTokens(_voter address) constant returns(numTokens uint256)
func (_PLCRVoting *PLCRVotingCallerSession) GetLockedTokens(_voter common.Address) (*big.Int, error) {
	return _PLCRVoting.Contract.GetLockedTokens(&_PLCRVoting.CallOpts, _voter)
}

// GetNumPassingTokens is a free data retrieval call binding the contract method 0x0c03fbd7.
//
// Solidity: function getNumPassingTokens(_voter address, _pollID uint256) constant returns(correctVotes uint256)
func (_PLCRVoting *PLCRVotingCaller) GetNumPassingTokens(opts *bind.CallOpts, _voter common.Address, _pollID *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "getNumPassingTokens", _voter, _pollID)
	return *ret0, err
}

// GetNumPassingTokens is a free data retrieval call binding the contract method 0x0c03fbd7.
//
// Solidity: function getNumPassingTokens(_voter address, _pollID uint256) constant returns(correctVotes uint256)
func (_PLCRVoting *PLCRVotingSession) GetNumPassingTokens(_voter common.Address, _pollID *big.Int) (*big.Int, error) {
	return _PLCRVoting.Contract.GetNumPassingTokens(&_PLCRVoting.CallOpts, _voter, _pollID)
}

// GetNumPassingTokens is a free data retrieval call binding the contract method 0x0c03fbd7.
//
// Solidity: function getNumPassingTokens(_voter address, _pollID uint256) constant returns(correctVotes uint256)
func (_PLCRVoting *PLCRVotingCallerSession) GetNumPassingTokens(_voter common.Address, _pollID *big.Int) (*big.Int, error) {
	return _PLCRVoting.Contract.GetNumPassingTokens(&_PLCRVoting.CallOpts, _voter, _pollID)
}

// GetNumTokens is a free data retrieval call binding the contract method 0xd1382092.
//
// Solidity: function getNumTokens(_voter address, _pollID uint256) constant returns(numTokens uint256)
func (_PLCRVoting *PLCRVotingCaller) GetNumTokens(opts *bind.CallOpts, _voter common.Address, _pollID *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "getNumTokens", _voter, _pollID)
	return *ret0, err
}

// GetNumTokens is a free data retrieval call binding the contract method 0xd1382092.
//
// Solidity: function getNumTokens(_voter address, _pollID uint256) constant returns(numTokens uint256)
func (_PLCRVoting *PLCRVotingSession) GetNumTokens(_voter common.Address, _pollID *big.Int) (*big.Int, error) {
	return _PLCRVoting.Contract.GetNumTokens(&_PLCRVoting.CallOpts, _voter, _pollID)
}

// GetNumTokens is a free data retrieval call binding the contract method 0xd1382092.
//
// Solidity: function getNumTokens(_voter address, _pollID uint256) constant returns(numTokens uint256)
func (_PLCRVoting *PLCRVotingCallerSession) GetNumTokens(_voter common.Address, _pollID *big.Int) (*big.Int, error) {
	return _PLCRVoting.Contract.GetNumTokens(&_PLCRVoting.CallOpts, _voter, _pollID)
}

// GetTotalNumberOfTokensForWinningOption is a free data retrieval call binding the contract method 0x053e71a6.
//
// Solidity: function getTotalNumberOfTokensForWinningOption(_pollID uint256) constant returns(numTokens uint256)
func (_PLCRVoting *PLCRVotingCaller) GetTotalNumberOfTokensForWinningOption(opts *bind.CallOpts, _pollID *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "getTotalNumberOfTokensForWinningOption", _pollID)
	return *ret0, err
}

// GetTotalNumberOfTokensForWinningOption is a free data retrieval call binding the contract method 0x053e71a6.
//
// Solidity: function getTotalNumberOfTokensForWinningOption(_pollID uint256) constant returns(numTokens uint256)
func (_PLCRVoting *PLCRVotingSession) GetTotalNumberOfTokensForWinningOption(_pollID *big.Int) (*big.Int, error) {
	return _PLCRVoting.Contract.GetTotalNumberOfTokensForWinningOption(&_PLCRVoting.CallOpts, _pollID)
}

// GetTotalNumberOfTokensForWinningOption is a free data retrieval call binding the contract method 0x053e71a6.
//
// Solidity: function getTotalNumberOfTokensForWinningOption(_pollID uint256) constant returns(numTokens uint256)
func (_PLCRVoting *PLCRVotingCallerSession) GetTotalNumberOfTokensForWinningOption(_pollID *big.Int) (*big.Int, error) {
	return _PLCRVoting.Contract.GetTotalNumberOfTokensForWinningOption(&_PLCRVoting.CallOpts, _pollID)
}

// IsExpired is a free data retrieval call binding the contract method 0xd9548e53.
//
// Solidity: function isExpired(_terminationDate uint256) constant returns(expired bool)
func (_PLCRVoting *PLCRVotingCaller) IsExpired(opts *bind.CallOpts, _terminationDate *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "isExpired", _terminationDate)
	return *ret0, err
}

// IsExpired is a free data retrieval call binding the contract method 0xd9548e53.
//
// Solidity: function isExpired(_terminationDate uint256) constant returns(expired bool)
func (_PLCRVoting *PLCRVotingSession) IsExpired(_terminationDate *big.Int) (bool, error) {
	return _PLCRVoting.Contract.IsExpired(&_PLCRVoting.CallOpts, _terminationDate)
}

// IsExpired is a free data retrieval call binding the contract method 0xd9548e53.
//
// Solidity: function isExpired(_terminationDate uint256) constant returns(expired bool)
func (_PLCRVoting *PLCRVotingCallerSession) IsExpired(_terminationDate *big.Int) (bool, error) {
	return _PLCRVoting.Contract.IsExpired(&_PLCRVoting.CallOpts, _terminationDate)
}

// IsPassed is a free data retrieval call binding the contract method 0x49403183.
//
// Solidity: function isPassed(_pollID uint256) constant returns(passed bool)
func (_PLCRVoting *PLCRVotingCaller) IsPassed(opts *bind.CallOpts, _pollID *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "isPassed", _pollID)
	return *ret0, err
}

// IsPassed is a free data retrieval call binding the contract method 0x49403183.
//
// Solidity: function isPassed(_pollID uint256) constant returns(passed bool)
func (_PLCRVoting *PLCRVotingSession) IsPassed(_pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.IsPassed(&_PLCRVoting.CallOpts, _pollID)
}

// IsPassed is a free data retrieval call binding the contract method 0x49403183.
//
// Solidity: function isPassed(_pollID uint256) constant returns(passed bool)
func (_PLCRVoting *PLCRVotingCallerSession) IsPassed(_pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.IsPassed(&_PLCRVoting.CallOpts, _pollID)
}

// PollEnded is a free data retrieval call binding the contract method 0xee684830.
//
// Solidity: function pollEnded(_pollID uint256) constant returns(ended bool)
func (_PLCRVoting *PLCRVotingCaller) PollEnded(opts *bind.CallOpts, _pollID *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "pollEnded", _pollID)
	return *ret0, err
}

// PollEnded is a free data retrieval call binding the contract method 0xee684830.
//
// Solidity: function pollEnded(_pollID uint256) constant returns(ended bool)
func (_PLCRVoting *PLCRVotingSession) PollEnded(_pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.PollEnded(&_PLCRVoting.CallOpts, _pollID)
}

// PollEnded is a free data retrieval call binding the contract method 0xee684830.
//
// Solidity: function pollEnded(_pollID uint256) constant returns(ended bool)
func (_PLCRVoting *PLCRVotingCallerSession) PollEnded(_pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.PollEnded(&_PLCRVoting.CallOpts, _pollID)
}

// PollExists is a free data retrieval call binding the contract method 0x88d21ff3.
//
// Solidity: function pollExists(_pollID uint256) constant returns(exists bool)
func (_PLCRVoting *PLCRVotingCaller) PollExists(opts *bind.CallOpts, _pollID *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "pollExists", _pollID)
	return *ret0, err
}

// PollExists is a free data retrieval call binding the contract method 0x88d21ff3.
//
// Solidity: function pollExists(_pollID uint256) constant returns(exists bool)
func (_PLCRVoting *PLCRVotingSession) PollExists(_pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.PollExists(&_PLCRVoting.CallOpts, _pollID)
}

// PollExists is a free data retrieval call binding the contract method 0x88d21ff3.
//
// Solidity: function pollExists(_pollID uint256) constant returns(exists bool)
func (_PLCRVoting *PLCRVotingCallerSession) PollExists(_pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.PollExists(&_PLCRVoting.CallOpts, _pollID)
}

// PollMap is a free data retrieval call binding the contract method 0x6148fed5.
//
// Solidity: function pollMap( uint256) constant returns(commitEndDate uint256, revealEndDate uint256, voteQuorum uint256, votesFor uint256, votesAgainst uint256)
func (_PLCRVoting *PLCRVotingCaller) PollMap(opts *bind.CallOpts, arg0 *big.Int) (struct {
	CommitEndDate *big.Int
	RevealEndDate *big.Int
	VoteQuorum    *big.Int
	VotesFor      *big.Int
	VotesAgainst  *big.Int
}, error) {
	ret := new(struct {
		CommitEndDate *big.Int
		RevealEndDate *big.Int
		VoteQuorum    *big.Int
		VotesFor      *big.Int
		VotesAgainst  *big.Int
	})
	out := ret
	err := _PLCRVoting.contract.Call(opts, out, "pollMap", arg0)
	return *ret, err
}

// PollMap is a free data retrieval call binding the contract method 0x6148fed5.
//
// Solidity: function pollMap( uint256) constant returns(commitEndDate uint256, revealEndDate uint256, voteQuorum uint256, votesFor uint256, votesAgainst uint256)
func (_PLCRVoting *PLCRVotingSession) PollMap(arg0 *big.Int) (struct {
	CommitEndDate *big.Int
	RevealEndDate *big.Int
	VoteQuorum    *big.Int
	VotesFor      *big.Int
	VotesAgainst  *big.Int
}, error) {
	return _PLCRVoting.Contract.PollMap(&_PLCRVoting.CallOpts, arg0)
}

// PollMap is a free data retrieval call binding the contract method 0x6148fed5.
//
// Solidity: function pollMap( uint256) constant returns(commitEndDate uint256, revealEndDate uint256, voteQuorum uint256, votesFor uint256, votesAgainst uint256)
func (_PLCRVoting *PLCRVotingCallerSession) PollMap(arg0 *big.Int) (struct {
	CommitEndDate *big.Int
	RevealEndDate *big.Int
	VoteQuorum    *big.Int
	VotesFor      *big.Int
	VotesAgainst  *big.Int
}, error) {
	return _PLCRVoting.Contract.PollMap(&_PLCRVoting.CallOpts, arg0)
}

// PollNonce is a free data retrieval call binding the contract method 0x97508f36.
//
// Solidity: function pollNonce() constant returns(uint256)
func (_PLCRVoting *PLCRVotingCaller) PollNonce(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "pollNonce")
	return *ret0, err
}

// PollNonce is a free data retrieval call binding the contract method 0x97508f36.
//
// Solidity: function pollNonce() constant returns(uint256)
func (_PLCRVoting *PLCRVotingSession) PollNonce() (*big.Int, error) {
	return _PLCRVoting.Contract.PollNonce(&_PLCRVoting.CallOpts)
}

// PollNonce is a free data retrieval call binding the contract method 0x97508f36.
//
// Solidity: function pollNonce() constant returns(uint256)
func (_PLCRVoting *PLCRVotingCallerSession) PollNonce() (*big.Int, error) {
	return _PLCRVoting.Contract.PollNonce(&_PLCRVoting.CallOpts)
}

// RevealPeriodActive is a free data retrieval call binding the contract method 0x441c77c0.
//
// Solidity: function revealPeriodActive(_pollID uint256) constant returns(active bool)
func (_PLCRVoting *PLCRVotingCaller) RevealPeriodActive(opts *bind.CallOpts, _pollID *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "revealPeriodActive", _pollID)
	return *ret0, err
}

// RevealPeriodActive is a free data retrieval call binding the contract method 0x441c77c0.
//
// Solidity: function revealPeriodActive(_pollID uint256) constant returns(active bool)
func (_PLCRVoting *PLCRVotingSession) RevealPeriodActive(_pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.RevealPeriodActive(&_PLCRVoting.CallOpts, _pollID)
}

// RevealPeriodActive is a free data retrieval call binding the contract method 0x441c77c0.
//
// Solidity: function revealPeriodActive(_pollID uint256) constant returns(active bool)
func (_PLCRVoting *PLCRVotingCallerSession) RevealPeriodActive(_pollID *big.Int) (bool, error) {
	return _PLCRVoting.Contract.RevealPeriodActive(&_PLCRVoting.CallOpts, _pollID)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_PLCRVoting *PLCRVotingCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "token")
	return *ret0, err
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_PLCRVoting *PLCRVotingSession) Token() (common.Address, error) {
	return _PLCRVoting.Contract.Token(&_PLCRVoting.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_PLCRVoting *PLCRVotingCallerSession) Token() (common.Address, error) {
	return _PLCRVoting.Contract.Token(&_PLCRVoting.CallOpts)
}

// ValidPosition is a free data retrieval call binding the contract method 0x819b0293.
//
// Solidity: function validPosition(_prevID uint256, _nextID uint256, _voter address, _numTokens uint256) constant returns(valid bool)
func (_PLCRVoting *PLCRVotingCaller) ValidPosition(opts *bind.CallOpts, _prevID *big.Int, _nextID *big.Int, _voter common.Address, _numTokens *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "validPosition", _prevID, _nextID, _voter, _numTokens)
	return *ret0, err
}

// ValidPosition is a free data retrieval call binding the contract method 0x819b0293.
//
// Solidity: function validPosition(_prevID uint256, _nextID uint256, _voter address, _numTokens uint256) constant returns(valid bool)
func (_PLCRVoting *PLCRVotingSession) ValidPosition(_prevID *big.Int, _nextID *big.Int, _voter common.Address, _numTokens *big.Int) (bool, error) {
	return _PLCRVoting.Contract.ValidPosition(&_PLCRVoting.CallOpts, _prevID, _nextID, _voter, _numTokens)
}

// ValidPosition is a free data retrieval call binding the contract method 0x819b0293.
//
// Solidity: function validPosition(_prevID uint256, _nextID uint256, _voter address, _numTokens uint256) constant returns(valid bool)
func (_PLCRVoting *PLCRVotingCallerSession) ValidPosition(_prevID *big.Int, _nextID *big.Int, _voter common.Address, _numTokens *big.Int) (bool, error) {
	return _PLCRVoting.Contract.ValidPosition(&_PLCRVoting.CallOpts, _prevID, _nextID, _voter, _numTokens)
}

// VoteTokenBalance is a free data retrieval call binding the contract method 0x3b930294.
//
// Solidity: function voteTokenBalance( address) constant returns(uint256)
func (_PLCRVoting *PLCRVotingCaller) VoteTokenBalance(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PLCRVoting.contract.Call(opts, out, "voteTokenBalance", arg0)
	return *ret0, err
}

// VoteTokenBalance is a free data retrieval call binding the contract method 0x3b930294.
//
// Solidity: function voteTokenBalance( address) constant returns(uint256)
func (_PLCRVoting *PLCRVotingSession) VoteTokenBalance(arg0 common.Address) (*big.Int, error) {
	return _PLCRVoting.Contract.VoteTokenBalance(&_PLCRVoting.CallOpts, arg0)
}

// VoteTokenBalance is a free data retrieval call binding the contract method 0x3b930294.
//
// Solidity: function voteTokenBalance( address) constant returns(uint256)
func (_PLCRVoting *PLCRVotingCallerSession) VoteTokenBalance(arg0 common.Address) (*big.Int, error) {
	return _PLCRVoting.Contract.VoteTokenBalance(&_PLCRVoting.CallOpts, arg0)
}

// CommitVote is a paid mutator transaction binding the contract method 0x6cbf9c5e.
//
// Solidity: function commitVote(_pollID uint256, _secretHash bytes32, _numTokens uint256, _prevPollID uint256) returns()
func (_PLCRVoting *PLCRVotingTransactor) CommitVote(opts *bind.TransactOpts, _pollID *big.Int, _secretHash [32]byte, _numTokens *big.Int, _prevPollID *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.contract.Transact(opts, "commitVote", _pollID, _secretHash, _numTokens, _prevPollID)
}

// CommitVote is a paid mutator transaction binding the contract method 0x6cbf9c5e.
//
// Solidity: function commitVote(_pollID uint256, _secretHash bytes32, _numTokens uint256, _prevPollID uint256) returns()
func (_PLCRVoting *PLCRVotingSession) CommitVote(_pollID *big.Int, _secretHash [32]byte, _numTokens *big.Int, _prevPollID *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.CommitVote(&_PLCRVoting.TransactOpts, _pollID, _secretHash, _numTokens, _prevPollID)
}

// CommitVote is a paid mutator transaction binding the contract method 0x6cbf9c5e.
//
// Solidity: function commitVote(_pollID uint256, _secretHash bytes32, _numTokens uint256, _prevPollID uint256) returns()
func (_PLCRVoting *PLCRVotingTransactorSession) CommitVote(_pollID *big.Int, _secretHash [32]byte, _numTokens *big.Int, _prevPollID *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.CommitVote(&_PLCRVoting.TransactOpts, _pollID, _secretHash, _numTokens, _prevPollID)
}

// CommitVotes is a paid mutator transaction binding the contract method 0x3ec36b99.
//
// Solidity: function commitVotes(_pollIDs uint256[], _secretHashes bytes32[], _numsTokens uint256[], _prevPollIDs uint256[]) returns()
func (_PLCRVoting *PLCRVotingTransactor) CommitVotes(opts *bind.TransactOpts, _pollIDs []*big.Int, _secretHashes [][32]byte, _numsTokens []*big.Int, _prevPollIDs []*big.Int) (*types.Transaction, error) {
	return _PLCRVoting.contract.Transact(opts, "commitVotes", _pollIDs, _secretHashes, _numsTokens, _prevPollIDs)
}

// CommitVotes is a paid mutator transaction binding the contract method 0x3ec36b99.
//
// Solidity: function commitVotes(_pollIDs uint256[], _secretHashes bytes32[], _numsTokens uint256[], _prevPollIDs uint256[]) returns()
func (_PLCRVoting *PLCRVotingSession) CommitVotes(_pollIDs []*big.Int, _secretHashes [][32]byte, _numsTokens []*big.Int, _prevPollIDs []*big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.CommitVotes(&_PLCRVoting.TransactOpts, _pollIDs, _secretHashes, _numsTokens, _prevPollIDs)
}

// CommitVotes is a paid mutator transaction binding the contract method 0x3ec36b99.
//
// Solidity: function commitVotes(_pollIDs uint256[], _secretHashes bytes32[], _numsTokens uint256[], _prevPollIDs uint256[]) returns()
func (_PLCRVoting *PLCRVotingTransactorSession) CommitVotes(_pollIDs []*big.Int, _secretHashes [][32]byte, _numsTokens []*big.Int, _prevPollIDs []*big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.CommitVotes(&_PLCRVoting.TransactOpts, _pollIDs, _secretHashes, _numsTokens, _prevPollIDs)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(_token address) returns()
func (_PLCRVoting *PLCRVotingTransactor) Init(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _PLCRVoting.contract.Transact(opts, "init", _token)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(_token address) returns()
func (_PLCRVoting *PLCRVotingSession) Init(_token common.Address) (*types.Transaction, error) {
	return _PLCRVoting.Contract.Init(&_PLCRVoting.TransactOpts, _token)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(_token address) returns()
func (_PLCRVoting *PLCRVotingTransactorSession) Init(_token common.Address) (*types.Transaction, error) {
	return _PLCRVoting.Contract.Init(&_PLCRVoting.TransactOpts, _token)
}

// RequestVotingRights is a paid mutator transaction binding the contract method 0xa25236fe.
//
// Solidity: function requestVotingRights(_numTokens uint256) returns()
func (_PLCRVoting *PLCRVotingTransactor) RequestVotingRights(opts *bind.TransactOpts, _numTokens *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.contract.Transact(opts, "requestVotingRights", _numTokens)
}

// RequestVotingRights is a paid mutator transaction binding the contract method 0xa25236fe.
//
// Solidity: function requestVotingRights(_numTokens uint256) returns()
func (_PLCRVoting *PLCRVotingSession) RequestVotingRights(_numTokens *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.RequestVotingRights(&_PLCRVoting.TransactOpts, _numTokens)
}

// RequestVotingRights is a paid mutator transaction binding the contract method 0xa25236fe.
//
// Solidity: function requestVotingRights(_numTokens uint256) returns()
func (_PLCRVoting *PLCRVotingTransactorSession) RequestVotingRights(_numTokens *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.RequestVotingRights(&_PLCRVoting.TransactOpts, _numTokens)
}

// RescueTokens is a paid mutator transaction binding the contract method 0x97603560.
//
// Solidity: function rescueTokens(_pollID uint256) returns()
func (_PLCRVoting *PLCRVotingTransactor) RescueTokens(opts *bind.TransactOpts, _pollID *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.contract.Transact(opts, "rescueTokens", _pollID)
}

// RescueTokens is a paid mutator transaction binding the contract method 0x97603560.
//
// Solidity: function rescueTokens(_pollID uint256) returns()
func (_PLCRVoting *PLCRVotingSession) RescueTokens(_pollID *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.RescueTokens(&_PLCRVoting.TransactOpts, _pollID)
}

// RescueTokens is a paid mutator transaction binding the contract method 0x97603560.
//
// Solidity: function rescueTokens(_pollID uint256) returns()
func (_PLCRVoting *PLCRVotingTransactorSession) RescueTokens(_pollID *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.RescueTokens(&_PLCRVoting.TransactOpts, _pollID)
}

// RescueTokensInMultiplePolls is a paid mutator transaction binding the contract method 0xbb11ed7e.
//
// Solidity: function rescueTokensInMultiplePolls(_pollIDs uint256[]) returns()
func (_PLCRVoting *PLCRVotingTransactor) RescueTokensInMultiplePolls(opts *bind.TransactOpts, _pollIDs []*big.Int) (*types.Transaction, error) {
	return _PLCRVoting.contract.Transact(opts, "rescueTokensInMultiplePolls", _pollIDs)
}

// RescueTokensInMultiplePolls is a paid mutator transaction binding the contract method 0xbb11ed7e.
//
// Solidity: function rescueTokensInMultiplePolls(_pollIDs uint256[]) returns()
func (_PLCRVoting *PLCRVotingSession) RescueTokensInMultiplePolls(_pollIDs []*big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.RescueTokensInMultiplePolls(&_PLCRVoting.TransactOpts, _pollIDs)
}

// RescueTokensInMultiplePolls is a paid mutator transaction binding the contract method 0xbb11ed7e.
//
// Solidity: function rescueTokensInMultiplePolls(_pollIDs uint256[]) returns()
func (_PLCRVoting *PLCRVotingTransactorSession) RescueTokensInMultiplePolls(_pollIDs []*big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.RescueTokensInMultiplePolls(&_PLCRVoting.TransactOpts, _pollIDs)
}

// RevealVote is a paid mutator transaction binding the contract method 0xb11d8bb8.
//
// Solidity: function revealVote(_pollID uint256, _voteOption uint256, _salt uint256) returns()
func (_PLCRVoting *PLCRVotingTransactor) RevealVote(opts *bind.TransactOpts, _pollID *big.Int, _voteOption *big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.contract.Transact(opts, "revealVote", _pollID, _voteOption, _salt)
}

// RevealVote is a paid mutator transaction binding the contract method 0xb11d8bb8.
//
// Solidity: function revealVote(_pollID uint256, _voteOption uint256, _salt uint256) returns()
func (_PLCRVoting *PLCRVotingSession) RevealVote(_pollID *big.Int, _voteOption *big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.RevealVote(&_PLCRVoting.TransactOpts, _pollID, _voteOption, _salt)
}

// RevealVote is a paid mutator transaction binding the contract method 0xb11d8bb8.
//
// Solidity: function revealVote(_pollID uint256, _voteOption uint256, _salt uint256) returns()
func (_PLCRVoting *PLCRVotingTransactorSession) RevealVote(_pollID *big.Int, _voteOption *big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.RevealVote(&_PLCRVoting.TransactOpts, _pollID, _voteOption, _salt)
}

// RevealVotes is a paid mutator transaction binding the contract method 0x8090f92e.
//
// Solidity: function revealVotes(_pollIDs uint256[], _voteOptions uint256[], _salts uint256[]) returns()
func (_PLCRVoting *PLCRVotingTransactor) RevealVotes(opts *bind.TransactOpts, _pollIDs []*big.Int, _voteOptions []*big.Int, _salts []*big.Int) (*types.Transaction, error) {
	return _PLCRVoting.contract.Transact(opts, "revealVotes", _pollIDs, _voteOptions, _salts)
}

// RevealVotes is a paid mutator transaction binding the contract method 0x8090f92e.
//
// Solidity: function revealVotes(_pollIDs uint256[], _voteOptions uint256[], _salts uint256[]) returns()
func (_PLCRVoting *PLCRVotingSession) RevealVotes(_pollIDs []*big.Int, _voteOptions []*big.Int, _salts []*big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.RevealVotes(&_PLCRVoting.TransactOpts, _pollIDs, _voteOptions, _salts)
}

// RevealVotes is a paid mutator transaction binding the contract method 0x8090f92e.
//
// Solidity: function revealVotes(_pollIDs uint256[], _voteOptions uint256[], _salts uint256[]) returns()
func (_PLCRVoting *PLCRVotingTransactorSession) RevealVotes(_pollIDs []*big.Int, _voteOptions []*big.Int, _salts []*big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.RevealVotes(&_PLCRVoting.TransactOpts, _pollIDs, _voteOptions, _salts)
}

// StartPoll is a paid mutator transaction binding the contract method 0x32ed3d60.
//
// Solidity: function startPoll(_voteQuorum uint256, _commitDuration uint256, _revealDuration uint256) returns(pollID uint256)
func (_PLCRVoting *PLCRVotingTransactor) StartPoll(opts *bind.TransactOpts, _voteQuorum *big.Int, _commitDuration *big.Int, _revealDuration *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.contract.Transact(opts, "startPoll", _voteQuorum, _commitDuration, _revealDuration)
}

// StartPoll is a paid mutator transaction binding the contract method 0x32ed3d60.
//
// Solidity: function startPoll(_voteQuorum uint256, _commitDuration uint256, _revealDuration uint256) returns(pollID uint256)
func (_PLCRVoting *PLCRVotingSession) StartPoll(_voteQuorum *big.Int, _commitDuration *big.Int, _revealDuration *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.StartPoll(&_PLCRVoting.TransactOpts, _voteQuorum, _commitDuration, _revealDuration)
}

// StartPoll is a paid mutator transaction binding the contract method 0x32ed3d60.
//
// Solidity: function startPoll(_voteQuorum uint256, _commitDuration uint256, _revealDuration uint256) returns(pollID uint256)
func (_PLCRVoting *PLCRVotingTransactorSession) StartPoll(_voteQuorum *big.Int, _commitDuration *big.Int, _revealDuration *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.StartPoll(&_PLCRVoting.TransactOpts, _voteQuorum, _commitDuration, _revealDuration)
}

// WithdrawVotingRights is a paid mutator transaction binding the contract method 0xe7b1d43c.
//
// Solidity: function withdrawVotingRights(_numTokens uint256) returns()
func (_PLCRVoting *PLCRVotingTransactor) WithdrawVotingRights(opts *bind.TransactOpts, _numTokens *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.contract.Transact(opts, "withdrawVotingRights", _numTokens)
}

// WithdrawVotingRights is a paid mutator transaction binding the contract method 0xe7b1d43c.
//
// Solidity: function withdrawVotingRights(_numTokens uint256) returns()
func (_PLCRVoting *PLCRVotingSession) WithdrawVotingRights(_numTokens *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.WithdrawVotingRights(&_PLCRVoting.TransactOpts, _numTokens)
}

// WithdrawVotingRights is a paid mutator transaction binding the contract method 0xe7b1d43c.
//
// Solidity: function withdrawVotingRights(_numTokens uint256) returns()
func (_PLCRVoting *PLCRVotingTransactorSession) WithdrawVotingRights(_numTokens *big.Int) (*types.Transaction, error) {
	return _PLCRVoting.Contract.WithdrawVotingRights(&_PLCRVoting.TransactOpts, _numTokens)
}

// PLCRVotingPollCreatedIterator is returned from FilterPollCreated and is used to iterate over the raw logs and unpacked data for PollCreated events raised by the PLCRVoting contract.
type PLCRVotingPollCreatedIterator struct {
	Event *PLCRVotingPollCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PLCRVotingPollCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PLCRVotingPollCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PLCRVotingPollCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PLCRVotingPollCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PLCRVotingPollCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PLCRVotingPollCreated represents a PollCreated event raised by the PLCRVoting contract.
type PLCRVotingPollCreated struct {
	VoteQuorum    *big.Int
	CommitEndDate *big.Int
	RevealEndDate *big.Int
	PollID        *big.Int
	Creator       common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterPollCreated is a free log retrieval operation binding the contract event 0x404f1f1c229d9eb2a949e7584da6ffde9d059ef2169f487ca815434cce0640d0.
//
// Solidity: e _PollCreated(voteQuorum uint256, commitEndDate uint256, revealEndDate uint256, pollID indexed uint256, creator indexed address)
func (_PLCRVoting *PLCRVotingFilterer) FilterPollCreated(opts *bind.FilterOpts, pollID []*big.Int, creator []common.Address) (*PLCRVotingPollCreatedIterator, error) {

	var pollIDRule []interface{}
	for _, pollIDItem := range pollID {
		pollIDRule = append(pollIDRule, pollIDItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _PLCRVoting.contract.FilterLogs(opts, "_PollCreated", pollIDRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &PLCRVotingPollCreatedIterator{contract: _PLCRVoting.contract, event: "_PollCreated", logs: logs, sub: sub}, nil
}

// WatchPollCreated is a free log subscription operation binding the contract event 0x404f1f1c229d9eb2a949e7584da6ffde9d059ef2169f487ca815434cce0640d0.
//
// Solidity: e _PollCreated(voteQuorum uint256, commitEndDate uint256, revealEndDate uint256, pollID indexed uint256, creator indexed address)
func (_PLCRVoting *PLCRVotingFilterer) WatchPollCreated(opts *bind.WatchOpts, sink chan<- *PLCRVotingPollCreated, pollID []*big.Int, creator []common.Address) (event.Subscription, error) {

	var pollIDRule []interface{}
	for _, pollIDItem := range pollID {
		pollIDRule = append(pollIDRule, pollIDItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _PLCRVoting.contract.WatchLogs(opts, "_PollCreated", pollIDRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PLCRVotingPollCreated)
				if err := _PLCRVoting.contract.UnpackLog(event, "_PollCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PLCRVotingTokensRescuedIterator is returned from FilterTokensRescued and is used to iterate over the raw logs and unpacked data for TokensRescued events raised by the PLCRVoting contract.
type PLCRVotingTokensRescuedIterator struct {
	Event *PLCRVotingTokensRescued // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PLCRVotingTokensRescuedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PLCRVotingTokensRescued)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PLCRVotingTokensRescued)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PLCRVotingTokensRescuedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PLCRVotingTokensRescuedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PLCRVotingTokensRescued represents a TokensRescued event raised by the PLCRVoting contract.
type PLCRVotingTokensRescued struct {
	PollID *big.Int
	Voter  common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokensRescued is a free log retrieval operation binding the contract event 0x402507661c8c8cb90e0a796450b8bdd28b6c516f05279c0cd29e84c344e1699a.
//
// Solidity: e _TokensRescued(pollID indexed uint256, voter indexed address)
func (_PLCRVoting *PLCRVotingFilterer) FilterTokensRescued(opts *bind.FilterOpts, pollID []*big.Int, voter []common.Address) (*PLCRVotingTokensRescuedIterator, error) {

	var pollIDRule []interface{}
	for _, pollIDItem := range pollID {
		pollIDRule = append(pollIDRule, pollIDItem)
	}
	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _PLCRVoting.contract.FilterLogs(opts, "_TokensRescued", pollIDRule, voterRule)
	if err != nil {
		return nil, err
	}
	return &PLCRVotingTokensRescuedIterator{contract: _PLCRVoting.contract, event: "_TokensRescued", logs: logs, sub: sub}, nil
}

// WatchTokensRescued is a free log subscription operation binding the contract event 0x402507661c8c8cb90e0a796450b8bdd28b6c516f05279c0cd29e84c344e1699a.
//
// Solidity: e _TokensRescued(pollID indexed uint256, voter indexed address)
func (_PLCRVoting *PLCRVotingFilterer) WatchTokensRescued(opts *bind.WatchOpts, sink chan<- *PLCRVotingTokensRescued, pollID []*big.Int, voter []common.Address) (event.Subscription, error) {

	var pollIDRule []interface{}
	for _, pollIDItem := range pollID {
		pollIDRule = append(pollIDRule, pollIDItem)
	}
	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _PLCRVoting.contract.WatchLogs(opts, "_TokensRescued", pollIDRule, voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PLCRVotingTokensRescued)
				if err := _PLCRVoting.contract.UnpackLog(event, "_TokensRescued", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PLCRVotingVoteCommittedIterator is returned from FilterVoteCommitted and is used to iterate over the raw logs and unpacked data for VoteCommitted events raised by the PLCRVoting contract.
type PLCRVotingVoteCommittedIterator struct {
	Event *PLCRVotingVoteCommitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PLCRVotingVoteCommittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PLCRVotingVoteCommitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PLCRVotingVoteCommitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PLCRVotingVoteCommittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PLCRVotingVoteCommittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PLCRVotingVoteCommitted represents a VoteCommitted event raised by the PLCRVoting contract.
type PLCRVotingVoteCommitted struct {
	PollID    *big.Int
	NumTokens *big.Int
	Voter     common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterVoteCommitted is a free log retrieval operation binding the contract event 0xea7979e4280d7e6bffc1c7d83a1ac99f16d02ecc14465ce41016226783b663d7.
//
// Solidity: e _VoteCommitted(pollID indexed uint256, numTokens uint256, voter indexed address)
func (_PLCRVoting *PLCRVotingFilterer) FilterVoteCommitted(opts *bind.FilterOpts, pollID []*big.Int, voter []common.Address) (*PLCRVotingVoteCommittedIterator, error) {

	var pollIDRule []interface{}
	for _, pollIDItem := range pollID {
		pollIDRule = append(pollIDRule, pollIDItem)
	}

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _PLCRVoting.contract.FilterLogs(opts, "_VoteCommitted", pollIDRule, voterRule)
	if err != nil {
		return nil, err
	}
	return &PLCRVotingVoteCommittedIterator{contract: _PLCRVoting.contract, event: "_VoteCommitted", logs: logs, sub: sub}, nil
}

// WatchVoteCommitted is a free log subscription operation binding the contract event 0xea7979e4280d7e6bffc1c7d83a1ac99f16d02ecc14465ce41016226783b663d7.
//
// Solidity: e _VoteCommitted(pollID indexed uint256, numTokens uint256, voter indexed address)
func (_PLCRVoting *PLCRVotingFilterer) WatchVoteCommitted(opts *bind.WatchOpts, sink chan<- *PLCRVotingVoteCommitted, pollID []*big.Int, voter []common.Address) (event.Subscription, error) {

	var pollIDRule []interface{}
	for _, pollIDItem := range pollID {
		pollIDRule = append(pollIDRule, pollIDItem)
	}

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _PLCRVoting.contract.WatchLogs(opts, "_VoteCommitted", pollIDRule, voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PLCRVotingVoteCommitted)
				if err := _PLCRVoting.contract.UnpackLog(event, "_VoteCommitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PLCRVotingVoteRevealedIterator is returned from FilterVoteRevealed and is used to iterate over the raw logs and unpacked data for VoteRevealed events raised by the PLCRVoting contract.
type PLCRVotingVoteRevealedIterator struct {
	Event *PLCRVotingVoteRevealed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PLCRVotingVoteRevealedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PLCRVotingVoteRevealed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PLCRVotingVoteRevealed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PLCRVotingVoteRevealedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PLCRVotingVoteRevealedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PLCRVotingVoteRevealed represents a VoteRevealed event raised by the PLCRVoting contract.
type PLCRVotingVoteRevealed struct {
	PollID       *big.Int
	NumTokens    *big.Int
	VotesFor     *big.Int
	VotesAgainst *big.Int
	Choice       *big.Int
	Voter        common.Address
	Salt         *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterVoteRevealed is a free log retrieval operation binding the contract event 0x9b19aaec524fad29c0ced9b9973a15e3045d7c3be156d71394ab40f0d5f119ff.
//
// Solidity: e _VoteRevealed(pollID indexed uint256, numTokens uint256, votesFor uint256, votesAgainst uint256, choice indexed uint256, voter indexed address, salt uint256)
func (_PLCRVoting *PLCRVotingFilterer) FilterVoteRevealed(opts *bind.FilterOpts, pollID []*big.Int, choice []*big.Int, voter []common.Address) (*PLCRVotingVoteRevealedIterator, error) {

	var pollIDRule []interface{}
	for _, pollIDItem := range pollID {
		pollIDRule = append(pollIDRule, pollIDItem)
	}

	var choiceRule []interface{}
	for _, choiceItem := range choice {
		choiceRule = append(choiceRule, choiceItem)
	}
	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _PLCRVoting.contract.FilterLogs(opts, "_VoteRevealed", pollIDRule, choiceRule, voterRule)
	if err != nil {
		return nil, err
	}
	return &PLCRVotingVoteRevealedIterator{contract: _PLCRVoting.contract, event: "_VoteRevealed", logs: logs, sub: sub}, nil
}

// WatchVoteRevealed is a free log subscription operation binding the contract event 0x9b19aaec524fad29c0ced9b9973a15e3045d7c3be156d71394ab40f0d5f119ff.
//
// Solidity: e _VoteRevealed(pollID indexed uint256, numTokens uint256, votesFor uint256, votesAgainst uint256, choice indexed uint256, voter indexed address, salt uint256)
func (_PLCRVoting *PLCRVotingFilterer) WatchVoteRevealed(opts *bind.WatchOpts, sink chan<- *PLCRVotingVoteRevealed, pollID []*big.Int, choice []*big.Int, voter []common.Address) (event.Subscription, error) {

	var pollIDRule []interface{}
	for _, pollIDItem := range pollID {
		pollIDRule = append(pollIDRule, pollIDItem)
	}

	var choiceRule []interface{}
	for _, choiceItem := range choice {
		choiceRule = append(choiceRule, choiceItem)
	}
	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _PLCRVoting.contract.WatchLogs(opts, "_VoteRevealed", pollIDRule, choiceRule, voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PLCRVotingVoteRevealed)
				if err := _PLCRVoting.contract.UnpackLog(event, "_VoteRevealed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PLCRVotingVotingRightsGrantedIterator is returned from FilterVotingRightsGranted and is used to iterate over the raw logs and unpacked data for VotingRightsGranted events raised by the PLCRVoting contract.
type PLCRVotingVotingRightsGrantedIterator struct {
	Event *PLCRVotingVotingRightsGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PLCRVotingVotingRightsGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PLCRVotingVotingRightsGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PLCRVotingVotingRightsGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PLCRVotingVotingRightsGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PLCRVotingVotingRightsGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PLCRVotingVotingRightsGranted represents a VotingRightsGranted event raised by the PLCRVoting contract.
type PLCRVotingVotingRightsGranted struct {
	NumTokens *big.Int
	Voter     common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterVotingRightsGranted is a free log retrieval operation binding the contract event 0xf7aaf024511d9982df8cd0d437c71c30106e6848cd1ba3d288d7a9c0e276aeda.
//
// Solidity: e _VotingRightsGranted(numTokens uint256, voter indexed address)
func (_PLCRVoting *PLCRVotingFilterer) FilterVotingRightsGranted(opts *bind.FilterOpts, voter []common.Address) (*PLCRVotingVotingRightsGrantedIterator, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _PLCRVoting.contract.FilterLogs(opts, "_VotingRightsGranted", voterRule)
	if err != nil {
		return nil, err
	}
	return &PLCRVotingVotingRightsGrantedIterator{contract: _PLCRVoting.contract, event: "_VotingRightsGranted", logs: logs, sub: sub}, nil
}

// WatchVotingRightsGranted is a free log subscription operation binding the contract event 0xf7aaf024511d9982df8cd0d437c71c30106e6848cd1ba3d288d7a9c0e276aeda.
//
// Solidity: e _VotingRightsGranted(numTokens uint256, voter indexed address)
func (_PLCRVoting *PLCRVotingFilterer) WatchVotingRightsGranted(opts *bind.WatchOpts, sink chan<- *PLCRVotingVotingRightsGranted, voter []common.Address) (event.Subscription, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _PLCRVoting.contract.WatchLogs(opts, "_VotingRightsGranted", voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PLCRVotingVotingRightsGranted)
				if err := _PLCRVoting.contract.UnpackLog(event, "_VotingRightsGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PLCRVotingVotingRightsWithdrawnIterator is returned from FilterVotingRightsWithdrawn and is used to iterate over the raw logs and unpacked data for VotingRightsWithdrawn events raised by the PLCRVoting contract.
type PLCRVotingVotingRightsWithdrawnIterator struct {
	Event *PLCRVotingVotingRightsWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PLCRVotingVotingRightsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PLCRVotingVotingRightsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PLCRVotingVotingRightsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PLCRVotingVotingRightsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PLCRVotingVotingRightsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PLCRVotingVotingRightsWithdrawn represents a VotingRightsWithdrawn event raised by the PLCRVoting contract.
type PLCRVotingVotingRightsWithdrawn struct {
	NumTokens *big.Int
	Voter     common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterVotingRightsWithdrawn is a free log retrieval operation binding the contract event 0xfaeb7dbb9992397d26ea1944efd40c80b40f702faf69b46c67ad10aba68ccb79.
//
// Solidity: e _VotingRightsWithdrawn(numTokens uint256, voter indexed address)
func (_PLCRVoting *PLCRVotingFilterer) FilterVotingRightsWithdrawn(opts *bind.FilterOpts, voter []common.Address) (*PLCRVotingVotingRightsWithdrawnIterator, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _PLCRVoting.contract.FilterLogs(opts, "_VotingRightsWithdrawn", voterRule)
	if err != nil {
		return nil, err
	}
	return &PLCRVotingVotingRightsWithdrawnIterator{contract: _PLCRVoting.contract, event: "_VotingRightsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchVotingRightsWithdrawn is a free log subscription operation binding the contract event 0xfaeb7dbb9992397d26ea1944efd40c80b40f702faf69b46c67ad10aba68ccb79.
//
// Solidity: e _VotingRightsWithdrawn(numTokens uint256, voter indexed address)
func (_PLCRVoting *PLCRVotingFilterer) WatchVotingRightsWithdrawn(opts *bind.WatchOpts, sink chan<- *PLCRVotingVotingRightsWithdrawn, voter []common.Address) (event.Subscription, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _PLCRVoting.contract.WatchLogs(opts, "_VotingRightsWithdrawn", voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PLCRVotingVotingRightsWithdrawn)
				if err := _PLCRVoting.contract.UnpackLog(event, "_VotingRightsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
