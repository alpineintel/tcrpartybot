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

// RegistryABI is the input ABI used to generate the binding from.
const RegistryABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"challenges\",\"outputs\":[{\"name\":\"rewardPool\",\"type\":\"uint256\"},{\"name\":\"challenger\",\"type\":\"address\"},{\"name\":\"resolved\",\"type\":\"bool\"},{\"name\":\"stake\",\"type\":\"uint256\"},{\"name\":\"totalTokens\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"listings\",\"outputs\":[{\"name\":\"applicationExpiry\",\"type\":\"uint256\"},{\"name\":\"whitelisted\",\"type\":\"bool\"},{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"unstakedDeposit\",\"type\":\"uint256\"},{\"name\":\"challengeID\",\"type\":\"uint256\"},{\"name\":\"exitTime\",\"type\":\"uint256\"},{\"name\":\"exitTimeExpiry\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"parameterizer\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"voting\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"listingHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"deposit\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"appEndDate\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"string\"},{\"indexed\":true,\"name\":\"applicant\",\"type\":\"address\"}],\"name\":\"_Application\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"listingHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"challengeID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"commitEndDate\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"revealEndDate\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"challenger\",\"type\":\"address\"}],\"name\":\"_Challenge\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"listingHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"added\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"newTotal\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"_Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"listingHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"withdrew\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"newTotal\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"_Withdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"listingHash\",\"type\":\"bytes32\"}],\"name\":\"_ApplicationWhitelisted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"listingHash\",\"type\":\"bytes32\"}],\"name\":\"_ApplicationRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"listingHash\",\"type\":\"bytes32\"}],\"name\":\"_ListingRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"listingHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"_ListingWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"listingHash\",\"type\":\"bytes32\"}],\"name\":\"_TouchAndRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"listingHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"challengeID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"rewardPool\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"totalTokens\",\"type\":\"uint256\"}],\"name\":\"_ChallengeFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"listingHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"challengeID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"rewardPool\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"totalTokens\",\"type\":\"uint256\"}],\"name\":\"_ChallengeSucceeded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"challengeID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"reward\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"voter\",\"type\":\"address\"}],\"name\":\"_RewardClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"listingHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"exitTime\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"exitDelayEndDate\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"_ExitInitialized\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_voting\",\"type\":\"address\"},{\"name\":\"_parameterizer\",\"type\":\"address\"},{\"name\":\"_name\",\"type\":\"string\"}],\"name\":\"init\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_listingHash\",\"type\":\"bytes32\"},{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_data\",\"type\":\"string\"}],\"name\":\"apply\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_listingHash\",\"type\":\"bytes32\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_listingHash\",\"type\":\"bytes32\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_listingHash\",\"type\":\"bytes32\"}],\"name\":\"initExit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_listingHash\",\"type\":\"bytes32\"}],\"name\":\"finalizeExit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_listingHash\",\"type\":\"bytes32\"},{\"name\":\"_data\",\"type\":\"string\"}],\"name\":\"challenge\",\"outputs\":[{\"name\":\"challengeID\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_listingHash\",\"type\":\"bytes32\"}],\"name\":\"updateStatus\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_listingHashes\",\"type\":\"bytes32[]\"}],\"name\":\"updateStatuses\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_challengeID\",\"type\":\"uint256\"}],\"name\":\"claimReward\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_challengeIDs\",\"type\":\"uint256[]\"}],\"name\":\"claimRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_voter\",\"type\":\"address\"},{\"name\":\"_challengeID\",\"type\":\"uint256\"}],\"name\":\"voterReward\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_listingHash\",\"type\":\"bytes32\"}],\"name\":\"canBeWhitelisted\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_listingHash\",\"type\":\"bytes32\"}],\"name\":\"isWhitelisted\",\"outputs\":[{\"name\":\"whitelisted\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_listingHash\",\"type\":\"bytes32\"}],\"name\":\"appWasMade\",\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_listingHash\",\"type\":\"bytes32\"}],\"name\":\"challengeExists\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_listingHash\",\"type\":\"bytes32\"}],\"name\":\"challengeCanBeResolved\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_challengeID\",\"type\":\"uint256\"}],\"name\":\"determineReward\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_challengeID\",\"type\":\"uint256\"},{\"name\":\"_voter\",\"type\":\"address\"}],\"name\":\"tokenClaims\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Registry is an auto generated Go binding around an Ethereum contract.
type Registry struct {
	RegistryCaller     // Read-only binding to the contract
	RegistryTransactor // Write-only binding to the contract
	RegistryFilterer   // Log filterer for contract events
}

// RegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RegistrySession struct {
	Contract     *Registry         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RegistryCallerSession struct {
	Contract *RegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// RegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RegistryTransactorSession struct {
	Contract     *RegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RegistryRaw struct {
	Contract *Registry // Generic contract binding to access the raw methods on
}

// RegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RegistryCallerRaw struct {
	Contract *RegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RegistryTransactorRaw struct {
	Contract *RegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegistry creates a new instance of Registry, bound to a specific deployed contract.
func NewRegistry(address common.Address, backend bind.ContractBackend) (*Registry, error) {
	contract, err := bindRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// NewRegistryCaller creates a new read-only instance of Registry, bound to a specific deployed contract.
func NewRegistryCaller(address common.Address, caller bind.ContractCaller) (*RegistryCaller, error) {
	contract, err := bindRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryCaller{contract: contract}, nil
}

// NewRegistryTransactor creates a new write-only instance of Registry, bound to a specific deployed contract.
func NewRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryTransactor, error) {
	contract, err := bindRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryTransactor{contract: contract}, nil
}

// NewRegistryFilterer creates a new log filterer instance of Registry, bound to a specific deployed contract.
func NewRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistryFilterer, error) {
	contract, err := bindRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistryFilterer{contract: contract}, nil
}

// bindRegistry binds a generic wrapper to an already deployed contract.
func bindRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.RegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transact(opts, method, params...)
}

// AppWasMade is a free data retrieval call binding the contract method 0x8cf8151f.
//
// Solidity: function appWasMade(_listingHash bytes32) constant returns(exists bool)
func (_Registry *RegistryCaller) AppWasMade(opts *bind.CallOpts, _listingHash [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "appWasMade", _listingHash)
	return *ret0, err
}

// AppWasMade is a free data retrieval call binding the contract method 0x8cf8151f.
//
// Solidity: function appWasMade(_listingHash bytes32) constant returns(exists bool)
func (_Registry *RegistrySession) AppWasMade(_listingHash [32]byte) (bool, error) {
	return _Registry.Contract.AppWasMade(&_Registry.CallOpts, _listingHash)
}

// AppWasMade is a free data retrieval call binding the contract method 0x8cf8151f.
//
// Solidity: function appWasMade(_listingHash bytes32) constant returns(exists bool)
func (_Registry *RegistryCallerSession) AppWasMade(_listingHash [32]byte) (bool, error) {
	return _Registry.Contract.AppWasMade(&_Registry.CallOpts, _listingHash)
}

// CanBeWhitelisted is a free data retrieval call binding the contract method 0x691a38ab.
//
// Solidity: function canBeWhitelisted(_listingHash bytes32) constant returns(bool)
func (_Registry *RegistryCaller) CanBeWhitelisted(opts *bind.CallOpts, _listingHash [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "canBeWhitelisted", _listingHash)
	return *ret0, err
}

// CanBeWhitelisted is a free data retrieval call binding the contract method 0x691a38ab.
//
// Solidity: function canBeWhitelisted(_listingHash bytes32) constant returns(bool)
func (_Registry *RegistrySession) CanBeWhitelisted(_listingHash [32]byte) (bool, error) {
	return _Registry.Contract.CanBeWhitelisted(&_Registry.CallOpts, _listingHash)
}

// CanBeWhitelisted is a free data retrieval call binding the contract method 0x691a38ab.
//
// Solidity: function canBeWhitelisted(_listingHash bytes32) constant returns(bool)
func (_Registry *RegistryCallerSession) CanBeWhitelisted(_listingHash [32]byte) (bool, error) {
	return _Registry.Contract.CanBeWhitelisted(&_Registry.CallOpts, _listingHash)
}

// ChallengeCanBeResolved is a free data retrieval call binding the contract method 0x77609a41.
//
// Solidity: function challengeCanBeResolved(_listingHash bytes32) constant returns(bool)
func (_Registry *RegistryCaller) ChallengeCanBeResolved(opts *bind.CallOpts, _listingHash [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "challengeCanBeResolved", _listingHash)
	return *ret0, err
}

// ChallengeCanBeResolved is a free data retrieval call binding the contract method 0x77609a41.
//
// Solidity: function challengeCanBeResolved(_listingHash bytes32) constant returns(bool)
func (_Registry *RegistrySession) ChallengeCanBeResolved(_listingHash [32]byte) (bool, error) {
	return _Registry.Contract.ChallengeCanBeResolved(&_Registry.CallOpts, _listingHash)
}

// ChallengeCanBeResolved is a free data retrieval call binding the contract method 0x77609a41.
//
// Solidity: function challengeCanBeResolved(_listingHash bytes32) constant returns(bool)
func (_Registry *RegistryCallerSession) ChallengeCanBeResolved(_listingHash [32]byte) (bool, error) {
	return _Registry.Contract.ChallengeCanBeResolved(&_Registry.CallOpts, _listingHash)
}

// ChallengeExists is a free data retrieval call binding the contract method 0x1b7bbecb.
//
// Solidity: function challengeExists(_listingHash bytes32) constant returns(bool)
func (_Registry *RegistryCaller) ChallengeExists(opts *bind.CallOpts, _listingHash [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "challengeExists", _listingHash)
	return *ret0, err
}

// ChallengeExists is a free data retrieval call binding the contract method 0x1b7bbecb.
//
// Solidity: function challengeExists(_listingHash bytes32) constant returns(bool)
func (_Registry *RegistrySession) ChallengeExists(_listingHash [32]byte) (bool, error) {
	return _Registry.Contract.ChallengeExists(&_Registry.CallOpts, _listingHash)
}

// ChallengeExists is a free data retrieval call binding the contract method 0x1b7bbecb.
//
// Solidity: function challengeExists(_listingHash bytes32) constant returns(bool)
func (_Registry *RegistryCallerSession) ChallengeExists(_listingHash [32]byte) (bool, error) {
	return _Registry.Contract.ChallengeExists(&_Registry.CallOpts, _listingHash)
}

// Challenges is a free data retrieval call binding the contract method 0x8f1d3776.
//
// Solidity: function challenges( uint256) constant returns(rewardPool uint256, challenger address, resolved bool, stake uint256, totalTokens uint256)
func (_Registry *RegistryCaller) Challenges(opts *bind.CallOpts, arg0 *big.Int) (struct {
	RewardPool  *big.Int
	Challenger  common.Address
	Resolved    bool
	Stake       *big.Int
	TotalTokens *big.Int
}, error) {
	ret := new(struct {
		RewardPool  *big.Int
		Challenger  common.Address
		Resolved    bool
		Stake       *big.Int
		TotalTokens *big.Int
	})
	out := ret
	err := _Registry.contract.Call(opts, out, "challenges", arg0)
	return *ret, err
}

// Challenges is a free data retrieval call binding the contract method 0x8f1d3776.
//
// Solidity: function challenges( uint256) constant returns(rewardPool uint256, challenger address, resolved bool, stake uint256, totalTokens uint256)
func (_Registry *RegistrySession) Challenges(arg0 *big.Int) (struct {
	RewardPool  *big.Int
	Challenger  common.Address
	Resolved    bool
	Stake       *big.Int
	TotalTokens *big.Int
}, error) {
	return _Registry.Contract.Challenges(&_Registry.CallOpts, arg0)
}

// Challenges is a free data retrieval call binding the contract method 0x8f1d3776.
//
// Solidity: function challenges( uint256) constant returns(rewardPool uint256, challenger address, resolved bool, stake uint256, totalTokens uint256)
func (_Registry *RegistryCallerSession) Challenges(arg0 *big.Int) (struct {
	RewardPool  *big.Int
	Challenger  common.Address
	Resolved    bool
	Stake       *big.Int
	TotalTokens *big.Int
}, error) {
	return _Registry.Contract.Challenges(&_Registry.CallOpts, arg0)
}

// DetermineReward is a free data retrieval call binding the contract method 0xc8187cf1.
//
// Solidity: function determineReward(_challengeID uint256) constant returns(uint256)
func (_Registry *RegistryCaller) DetermineReward(opts *bind.CallOpts, _challengeID *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "determineReward", _challengeID)
	return *ret0, err
}

// DetermineReward is a free data retrieval call binding the contract method 0xc8187cf1.
//
// Solidity: function determineReward(_challengeID uint256) constant returns(uint256)
func (_Registry *RegistrySession) DetermineReward(_challengeID *big.Int) (*big.Int, error) {
	return _Registry.Contract.DetermineReward(&_Registry.CallOpts, _challengeID)
}

// DetermineReward is a free data retrieval call binding the contract method 0xc8187cf1.
//
// Solidity: function determineReward(_challengeID uint256) constant returns(uint256)
func (_Registry *RegistryCallerSession) DetermineReward(_challengeID *big.Int) (*big.Int, error) {
	return _Registry.Contract.DetermineReward(&_Registry.CallOpts, _challengeID)
}

// IsWhitelisted is a free data retrieval call binding the contract method 0x01a5e3fe.
//
// Solidity: function isWhitelisted(_listingHash bytes32) constant returns(whitelisted bool)
func (_Registry *RegistryCaller) IsWhitelisted(opts *bind.CallOpts, _listingHash [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "isWhitelisted", _listingHash)
	return *ret0, err
}

// IsWhitelisted is a free data retrieval call binding the contract method 0x01a5e3fe.
//
// Solidity: function isWhitelisted(_listingHash bytes32) constant returns(whitelisted bool)
func (_Registry *RegistrySession) IsWhitelisted(_listingHash [32]byte) (bool, error) {
	return _Registry.Contract.IsWhitelisted(&_Registry.CallOpts, _listingHash)
}

// IsWhitelisted is a free data retrieval call binding the contract method 0x01a5e3fe.
//
// Solidity: function isWhitelisted(_listingHash bytes32) constant returns(whitelisted bool)
func (_Registry *RegistryCallerSession) IsWhitelisted(_listingHash [32]byte) (bool, error) {
	return _Registry.Contract.IsWhitelisted(&_Registry.CallOpts, _listingHash)
}

// Listings is a free data retrieval call binding the contract method 0xc18b8db4.
//
// Solidity: function listings( bytes32) constant returns(applicationExpiry uint256, whitelisted bool, owner address, unstakedDeposit uint256, challengeID uint256, exitTime uint256, exitTimeExpiry uint256)
func (_Registry *RegistryCaller) Listings(opts *bind.CallOpts, arg0 [32]byte) (struct {
	ApplicationExpiry *big.Int
	Whitelisted       bool
	Owner             common.Address
	UnstakedDeposit   *big.Int
	ChallengeID       *big.Int
	ExitTime          *big.Int
	ExitTimeExpiry    *big.Int
}, error) {
	ret := new(struct {
		ApplicationExpiry *big.Int
		Whitelisted       bool
		Owner             common.Address
		UnstakedDeposit   *big.Int
		ChallengeID       *big.Int
		ExitTime          *big.Int
		ExitTimeExpiry    *big.Int
	})
	out := ret
	err := _Registry.contract.Call(opts, out, "listings", arg0)
	return *ret, err
}

// Listings is a free data retrieval call binding the contract method 0xc18b8db4.
//
// Solidity: function listings( bytes32) constant returns(applicationExpiry uint256, whitelisted bool, owner address, unstakedDeposit uint256, challengeID uint256, exitTime uint256, exitTimeExpiry uint256)
func (_Registry *RegistrySession) Listings(arg0 [32]byte) (struct {
	ApplicationExpiry *big.Int
	Whitelisted       bool
	Owner             common.Address
	UnstakedDeposit   *big.Int
	ChallengeID       *big.Int
	ExitTime          *big.Int
	ExitTimeExpiry    *big.Int
}, error) {
	return _Registry.Contract.Listings(&_Registry.CallOpts, arg0)
}

// Listings is a free data retrieval call binding the contract method 0xc18b8db4.
//
// Solidity: function listings( bytes32) constant returns(applicationExpiry uint256, whitelisted bool, owner address, unstakedDeposit uint256, challengeID uint256, exitTime uint256, exitTimeExpiry uint256)
func (_Registry *RegistryCallerSession) Listings(arg0 [32]byte) (struct {
	ApplicationExpiry *big.Int
	Whitelisted       bool
	Owner             common.Address
	UnstakedDeposit   *big.Int
	ChallengeID       *big.Int
	ExitTime          *big.Int
	ExitTimeExpiry    *big.Int
}, error) {
	return _Registry.Contract.Listings(&_Registry.CallOpts, arg0)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_Registry *RegistryCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_Registry *RegistrySession) Name() (string, error) {
	return _Registry.Contract.Name(&_Registry.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_Registry *RegistryCallerSession) Name() (string, error) {
	return _Registry.Contract.Name(&_Registry.CallOpts)
}

// Parameterizer is a free data retrieval call binding the contract method 0xe1e3f915.
//
// Solidity: function parameterizer() constant returns(address)
func (_Registry *RegistryCaller) Parameterizer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "parameterizer")
	return *ret0, err
}

// Parameterizer is a free data retrieval call binding the contract method 0xe1e3f915.
//
// Solidity: function parameterizer() constant returns(address)
func (_Registry *RegistrySession) Parameterizer() (common.Address, error) {
	return _Registry.Contract.Parameterizer(&_Registry.CallOpts)
}

// Parameterizer is a free data retrieval call binding the contract method 0xe1e3f915.
//
// Solidity: function parameterizer() constant returns(address)
func (_Registry *RegistryCallerSession) Parameterizer() (common.Address, error) {
	return _Registry.Contract.Parameterizer(&_Registry.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Registry *RegistryCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "token")
	return *ret0, err
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Registry *RegistrySession) Token() (common.Address, error) {
	return _Registry.Contract.Token(&_Registry.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Registry *RegistryCallerSession) Token() (common.Address, error) {
	return _Registry.Contract.Token(&_Registry.CallOpts)
}

// TokenClaims is a free data retrieval call binding the contract method 0xa5ba3b1e.
//
// Solidity: function tokenClaims(_challengeID uint256, _voter address) constant returns(bool)
func (_Registry *RegistryCaller) TokenClaims(opts *bind.CallOpts, _challengeID *big.Int, _voter common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "tokenClaims", _challengeID, _voter)
	return *ret0, err
}

// TokenClaims is a free data retrieval call binding the contract method 0xa5ba3b1e.
//
// Solidity: function tokenClaims(_challengeID uint256, _voter address) constant returns(bool)
func (_Registry *RegistrySession) TokenClaims(_challengeID *big.Int, _voter common.Address) (bool, error) {
	return _Registry.Contract.TokenClaims(&_Registry.CallOpts, _challengeID, _voter)
}

// TokenClaims is a free data retrieval call binding the contract method 0xa5ba3b1e.
//
// Solidity: function tokenClaims(_challengeID uint256, _voter address) constant returns(bool)
func (_Registry *RegistryCallerSession) TokenClaims(_challengeID *big.Int, _voter common.Address) (bool, error) {
	return _Registry.Contract.TokenClaims(&_Registry.CallOpts, _challengeID, _voter)
}

// VoterReward is a free data retrieval call binding the contract method 0xbd0ae405.
//
// Solidity: function voterReward(_voter address, _challengeID uint256) constant returns(uint256)
func (_Registry *RegistryCaller) VoterReward(opts *bind.CallOpts, _voter common.Address, _challengeID *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "voterReward", _voter, _challengeID)
	return *ret0, err
}

// VoterReward is a free data retrieval call binding the contract method 0xbd0ae405.
//
// Solidity: function voterReward(_voter address, _challengeID uint256) constant returns(uint256)
func (_Registry *RegistrySession) VoterReward(_voter common.Address, _challengeID *big.Int) (*big.Int, error) {
	return _Registry.Contract.VoterReward(&_Registry.CallOpts, _voter, _challengeID)
}

// VoterReward is a free data retrieval call binding the contract method 0xbd0ae405.
//
// Solidity: function voterReward(_voter address, _challengeID uint256) constant returns(uint256)
func (_Registry *RegistryCallerSession) VoterReward(_voter common.Address, _challengeID *big.Int) (*big.Int, error) {
	return _Registry.Contract.VoterReward(&_Registry.CallOpts, _voter, _challengeID)
}

// Voting is a free data retrieval call binding the contract method 0xfce1ccca.
//
// Solidity: function voting() constant returns(address)
func (_Registry *RegistryCaller) Voting(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "voting")
	return *ret0, err
}

// Voting is a free data retrieval call binding the contract method 0xfce1ccca.
//
// Solidity: function voting() constant returns(address)
func (_Registry *RegistrySession) Voting() (common.Address, error) {
	return _Registry.Contract.Voting(&_Registry.CallOpts)
}

// Voting is a free data retrieval call binding the contract method 0xfce1ccca.
//
// Solidity: function voting() constant returns(address)
func (_Registry *RegistryCallerSession) Voting() (common.Address, error) {
	return _Registry.Contract.Voting(&_Registry.CallOpts)
}

// Apply is a paid mutator transaction binding the contract method 0x89bb55c7.
//
// Solidity: function apply(_listingHash bytes32, _amount uint256, _data string) returns()
func (_Registry *RegistryTransactor) Apply(opts *bind.TransactOpts, _listingHash [32]byte, _amount *big.Int, _data string) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "apply", _listingHash, _amount, _data)
}

// Apply is a paid mutator transaction binding the contract method 0x89bb55c7.
//
// Solidity: function apply(_listingHash bytes32, _amount uint256, _data string) returns()
func (_Registry *RegistrySession) Apply(_listingHash [32]byte, _amount *big.Int, _data string) (*types.Transaction, error) {
	return _Registry.Contract.Apply(&_Registry.TransactOpts, _listingHash, _amount, _data)
}

// Apply is a paid mutator transaction binding the contract method 0x89bb55c7.
//
// Solidity: function apply(_listingHash bytes32, _amount uint256, _data string) returns()
func (_Registry *RegistryTransactorSession) Apply(_listingHash [32]byte, _amount *big.Int, _data string) (*types.Transaction, error) {
	return _Registry.Contract.Apply(&_Registry.TransactOpts, _listingHash, _amount, _data)
}

// Challenge is a paid mutator transaction binding the contract method 0x43cffefe.
//
// Solidity: function challenge(_listingHash bytes32, _data string) returns(challengeID uint256)
func (_Registry *RegistryTransactor) Challenge(opts *bind.TransactOpts, _listingHash [32]byte, _data string) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "challenge", _listingHash, _data)
}

// Challenge is a paid mutator transaction binding the contract method 0x43cffefe.
//
// Solidity: function challenge(_listingHash bytes32, _data string) returns(challengeID uint256)
func (_Registry *RegistrySession) Challenge(_listingHash [32]byte, _data string) (*types.Transaction, error) {
	return _Registry.Contract.Challenge(&_Registry.TransactOpts, _listingHash, _data)
}

// Challenge is a paid mutator transaction binding the contract method 0x43cffefe.
//
// Solidity: function challenge(_listingHash bytes32, _data string) returns(challengeID uint256)
func (_Registry *RegistryTransactorSession) Challenge(_listingHash [32]byte, _data string) (*types.Transaction, error) {
	return _Registry.Contract.Challenge(&_Registry.TransactOpts, _listingHash, _data)
}

// ClaimReward is a paid mutator transaction binding the contract method 0xae169a50.
//
// Solidity: function claimReward(_challengeID uint256) returns()
func (_Registry *RegistryTransactor) ClaimReward(opts *bind.TransactOpts, _challengeID *big.Int) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "claimReward", _challengeID)
}

// ClaimReward is a paid mutator transaction binding the contract method 0xae169a50.
//
// Solidity: function claimReward(_challengeID uint256) returns()
func (_Registry *RegistrySession) ClaimReward(_challengeID *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.ClaimReward(&_Registry.TransactOpts, _challengeID)
}

// ClaimReward is a paid mutator transaction binding the contract method 0xae169a50.
//
// Solidity: function claimReward(_challengeID uint256) returns()
func (_Registry *RegistryTransactorSession) ClaimReward(_challengeID *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.ClaimReward(&_Registry.TransactOpts, _challengeID)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x5eac6239.
//
// Solidity: function claimRewards(_challengeIDs uint256[]) returns()
func (_Registry *RegistryTransactor) ClaimRewards(opts *bind.TransactOpts, _challengeIDs []*big.Int) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "claimRewards", _challengeIDs)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x5eac6239.
//
// Solidity: function claimRewards(_challengeIDs uint256[]) returns()
func (_Registry *RegistrySession) ClaimRewards(_challengeIDs []*big.Int) (*types.Transaction, error) {
	return _Registry.Contract.ClaimRewards(&_Registry.TransactOpts, _challengeIDs)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x5eac6239.
//
// Solidity: function claimRewards(_challengeIDs uint256[]) returns()
func (_Registry *RegistryTransactorSession) ClaimRewards(_challengeIDs []*big.Int) (*types.Transaction, error) {
	return _Registry.Contract.ClaimRewards(&_Registry.TransactOpts, _challengeIDs)
}

// Deposit is a paid mutator transaction binding the contract method 0x1de26e16.
//
// Solidity: function deposit(_listingHash bytes32, _amount uint256) returns()
func (_Registry *RegistryTransactor) Deposit(opts *bind.TransactOpts, _listingHash [32]byte, _amount *big.Int) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "deposit", _listingHash, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x1de26e16.
//
// Solidity: function deposit(_listingHash bytes32, _amount uint256) returns()
func (_Registry *RegistrySession) Deposit(_listingHash [32]byte, _amount *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Deposit(&_Registry.TransactOpts, _listingHash, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x1de26e16.
//
// Solidity: function deposit(_listingHash bytes32, _amount uint256) returns()
func (_Registry *RegistryTransactorSession) Deposit(_listingHash [32]byte, _amount *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Deposit(&_Registry.TransactOpts, _listingHash, _amount)
}

// FinalizeExit is a paid mutator transaction binding the contract method 0x0960db7c.
//
// Solidity: function finalizeExit(_listingHash bytes32) returns()
func (_Registry *RegistryTransactor) FinalizeExit(opts *bind.TransactOpts, _listingHash [32]byte) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "finalizeExit", _listingHash)
}

// FinalizeExit is a paid mutator transaction binding the contract method 0x0960db7c.
//
// Solidity: function finalizeExit(_listingHash bytes32) returns()
func (_Registry *RegistrySession) FinalizeExit(_listingHash [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.FinalizeExit(&_Registry.TransactOpts, _listingHash)
}

// FinalizeExit is a paid mutator transaction binding the contract method 0x0960db7c.
//
// Solidity: function finalizeExit(_listingHash bytes32) returns()
func (_Registry *RegistryTransactorSession) FinalizeExit(_listingHash [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.FinalizeExit(&_Registry.TransactOpts, _listingHash)
}

// Init is a paid mutator transaction binding the contract method 0x24f91d83.
//
// Solidity: function init(_token address, _voting address, _parameterizer address, _name string) returns()
func (_Registry *RegistryTransactor) Init(opts *bind.TransactOpts, _token common.Address, _voting common.Address, _parameterizer common.Address, _name string) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "init", _token, _voting, _parameterizer, _name)
}

// Init is a paid mutator transaction binding the contract method 0x24f91d83.
//
// Solidity: function init(_token address, _voting address, _parameterizer address, _name string) returns()
func (_Registry *RegistrySession) Init(_token common.Address, _voting common.Address, _parameterizer common.Address, _name string) (*types.Transaction, error) {
	return _Registry.Contract.Init(&_Registry.TransactOpts, _token, _voting, _parameterizer, _name)
}

// Init is a paid mutator transaction binding the contract method 0x24f91d83.
//
// Solidity: function init(_token address, _voting address, _parameterizer address, _name string) returns()
func (_Registry *RegistryTransactorSession) Init(_token common.Address, _voting common.Address, _parameterizer common.Address, _name string) (*types.Transaction, error) {
	return _Registry.Contract.Init(&_Registry.TransactOpts, _token, _voting, _parameterizer, _name)
}

// InitExit is a paid mutator transaction binding the contract method 0x07b99366.
//
// Solidity: function initExit(_listingHash bytes32) returns()
func (_Registry *RegistryTransactor) InitExit(opts *bind.TransactOpts, _listingHash [32]byte) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "initExit", _listingHash)
}

// InitExit is a paid mutator transaction binding the contract method 0x07b99366.
//
// Solidity: function initExit(_listingHash bytes32) returns()
func (_Registry *RegistrySession) InitExit(_listingHash [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.InitExit(&_Registry.TransactOpts, _listingHash)
}

// InitExit is a paid mutator transaction binding the contract method 0x07b99366.
//
// Solidity: function initExit(_listingHash bytes32) returns()
func (_Registry *RegistryTransactorSession) InitExit(_listingHash [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.InitExit(&_Registry.TransactOpts, _listingHash)
}

// UpdateStatus is a paid mutator transaction binding the contract method 0x8a59eb56.
//
// Solidity: function updateStatus(_listingHash bytes32) returns()
func (_Registry *RegistryTransactor) UpdateStatus(opts *bind.TransactOpts, _listingHash [32]byte) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "updateStatus", _listingHash)
}

// UpdateStatus is a paid mutator transaction binding the contract method 0x8a59eb56.
//
// Solidity: function updateStatus(_listingHash bytes32) returns()
func (_Registry *RegistrySession) UpdateStatus(_listingHash [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.UpdateStatus(&_Registry.TransactOpts, _listingHash)
}

// UpdateStatus is a paid mutator transaction binding the contract method 0x8a59eb56.
//
// Solidity: function updateStatus(_listingHash bytes32) returns()
func (_Registry *RegistryTransactorSession) UpdateStatus(_listingHash [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.UpdateStatus(&_Registry.TransactOpts, _listingHash)
}

// UpdateStatuses is a paid mutator transaction binding the contract method 0x8c82dccb.
//
// Solidity: function updateStatuses(_listingHashes bytes32[]) returns()
func (_Registry *RegistryTransactor) UpdateStatuses(opts *bind.TransactOpts, _listingHashes [][32]byte) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "updateStatuses", _listingHashes)
}

// UpdateStatuses is a paid mutator transaction binding the contract method 0x8c82dccb.
//
// Solidity: function updateStatuses(_listingHashes bytes32[]) returns()
func (_Registry *RegistrySession) UpdateStatuses(_listingHashes [][32]byte) (*types.Transaction, error) {
	return _Registry.Contract.UpdateStatuses(&_Registry.TransactOpts, _listingHashes)
}

// UpdateStatuses is a paid mutator transaction binding the contract method 0x8c82dccb.
//
// Solidity: function updateStatuses(_listingHashes bytes32[]) returns()
func (_Registry *RegistryTransactorSession) UpdateStatuses(_listingHashes [][32]byte) (*types.Transaction, error) {
	return _Registry.Contract.UpdateStatuses(&_Registry.TransactOpts, _listingHashes)
}

// Withdraw is a paid mutator transaction binding the contract method 0x040cf020.
//
// Solidity: function withdraw(_listingHash bytes32, _amount uint256) returns()
func (_Registry *RegistryTransactor) Withdraw(opts *bind.TransactOpts, _listingHash [32]byte, _amount *big.Int) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "withdraw", _listingHash, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x040cf020.
//
// Solidity: function withdraw(_listingHash bytes32, _amount uint256) returns()
func (_Registry *RegistrySession) Withdraw(_listingHash [32]byte, _amount *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Withdraw(&_Registry.TransactOpts, _listingHash, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x040cf020.
//
// Solidity: function withdraw(_listingHash bytes32, _amount uint256) returns()
func (_Registry *RegistryTransactorSession) Withdraw(_listingHash [32]byte, _amount *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Withdraw(&_Registry.TransactOpts, _listingHash, _amount)
}

// RegistryApplicationIterator is returned from FilterApplication and is used to iterate over the raw logs and unpacked data for Application events raised by the Registry contract.
type RegistryApplicationIterator struct {
	Event *RegistryApplication // Event containing the contract specifics and raw log

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
func (it *RegistryApplicationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryApplication)
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
		it.Event = new(RegistryApplication)
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
func (it *RegistryApplicationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryApplicationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryApplication represents a Application event raised by the Registry contract.
type RegistryApplication struct {
	ListingHash [32]byte
	Deposit     *big.Int
	AppEndDate  *big.Int
	Data        string
	Applicant   common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterApplication is a free log retrieval operation binding the contract event 0xa27f550c3c7a7c6d8369e5383fdc7a3b4850d8ce9e20066f9d496f6989f00864.
//
// Solidity: e _Application(listingHash indexed bytes32, deposit uint256, appEndDate uint256, data string, applicant indexed address)
func (_Registry *RegistryFilterer) FilterApplication(opts *bind.FilterOpts, listingHash [][32]byte, applicant []common.Address) (*RegistryApplicationIterator, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	var applicantRule []interface{}
	for _, applicantItem := range applicant {
		applicantRule = append(applicantRule, applicantItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "_Application", listingHashRule, applicantRule)
	if err != nil {
		return nil, err
	}
	return &RegistryApplicationIterator{contract: _Registry.contract, event: "_Application", logs: logs, sub: sub}, nil
}

// WatchApplication is a free log subscription operation binding the contract event 0xa27f550c3c7a7c6d8369e5383fdc7a3b4850d8ce9e20066f9d496f6989f00864.
//
// Solidity: e _Application(listingHash indexed bytes32, deposit uint256, appEndDate uint256, data string, applicant indexed address)
func (_Registry *RegistryFilterer) WatchApplication(opts *bind.WatchOpts, sink chan<- *RegistryApplication, listingHash [][32]byte, applicant []common.Address) (event.Subscription, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	var applicantRule []interface{}
	for _, applicantItem := range applicant {
		applicantRule = append(applicantRule, applicantItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "_Application", listingHashRule, applicantRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryApplication)
				if err := _Registry.contract.UnpackLog(event, "_Application", log); err != nil {
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

// RegistryApplicationRemovedIterator is returned from FilterApplicationRemoved and is used to iterate over the raw logs and unpacked data for ApplicationRemoved events raised by the Registry contract.
type RegistryApplicationRemovedIterator struct {
	Event *RegistryApplicationRemoved // Event containing the contract specifics and raw log

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
func (it *RegistryApplicationRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryApplicationRemoved)
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
		it.Event = new(RegistryApplicationRemoved)
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
func (it *RegistryApplicationRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryApplicationRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryApplicationRemoved represents a ApplicationRemoved event raised by the Registry contract.
type RegistryApplicationRemoved struct {
	ListingHash [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterApplicationRemoved is a free log retrieval operation binding the contract event 0x2e5ec035f6eac8ff1cf7cdf36cfeca7c85413f9f67652dc2c13d20f337204a26.
//
// Solidity: e _ApplicationRemoved(listingHash indexed bytes32)
func (_Registry *RegistryFilterer) FilterApplicationRemoved(opts *bind.FilterOpts, listingHash [][32]byte) (*RegistryApplicationRemovedIterator, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "_ApplicationRemoved", listingHashRule)
	if err != nil {
		return nil, err
	}
	return &RegistryApplicationRemovedIterator{contract: _Registry.contract, event: "_ApplicationRemoved", logs: logs, sub: sub}, nil
}

// WatchApplicationRemoved is a free log subscription operation binding the contract event 0x2e5ec035f6eac8ff1cf7cdf36cfeca7c85413f9f67652dc2c13d20f337204a26.
//
// Solidity: e _ApplicationRemoved(listingHash indexed bytes32)
func (_Registry *RegistryFilterer) WatchApplicationRemoved(opts *bind.WatchOpts, sink chan<- *RegistryApplicationRemoved, listingHash [][32]byte) (event.Subscription, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "_ApplicationRemoved", listingHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryApplicationRemoved)
				if err := _Registry.contract.UnpackLog(event, "_ApplicationRemoved", log); err != nil {
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

// RegistryApplicationWhitelistedIterator is returned from FilterApplicationWhitelisted and is used to iterate over the raw logs and unpacked data for ApplicationWhitelisted events raised by the Registry contract.
type RegistryApplicationWhitelistedIterator struct {
	Event *RegistryApplicationWhitelisted // Event containing the contract specifics and raw log

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
func (it *RegistryApplicationWhitelistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryApplicationWhitelisted)
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
		it.Event = new(RegistryApplicationWhitelisted)
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
func (it *RegistryApplicationWhitelistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryApplicationWhitelistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryApplicationWhitelisted represents a ApplicationWhitelisted event raised by the Registry contract.
type RegistryApplicationWhitelisted struct {
	ListingHash [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterApplicationWhitelisted is a free log retrieval operation binding the contract event 0xa7bc1d57d9006d9d248707c7b6828c1bab8c51719cc06d78c82a3ee891ef967c.
//
// Solidity: e _ApplicationWhitelisted(listingHash indexed bytes32)
func (_Registry *RegistryFilterer) FilterApplicationWhitelisted(opts *bind.FilterOpts, listingHash [][32]byte) (*RegistryApplicationWhitelistedIterator, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "_ApplicationWhitelisted", listingHashRule)
	if err != nil {
		return nil, err
	}
	return &RegistryApplicationWhitelistedIterator{contract: _Registry.contract, event: "_ApplicationWhitelisted", logs: logs, sub: sub}, nil
}

// WatchApplicationWhitelisted is a free log subscription operation binding the contract event 0xa7bc1d57d9006d9d248707c7b6828c1bab8c51719cc06d78c82a3ee891ef967c.
//
// Solidity: e _ApplicationWhitelisted(listingHash indexed bytes32)
func (_Registry *RegistryFilterer) WatchApplicationWhitelisted(opts *bind.WatchOpts, sink chan<- *RegistryApplicationWhitelisted, listingHash [][32]byte) (event.Subscription, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "_ApplicationWhitelisted", listingHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryApplicationWhitelisted)
				if err := _Registry.contract.UnpackLog(event, "_ApplicationWhitelisted", log); err != nil {
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

// RegistryChallengeIterator is returned from FilterChallenge and is used to iterate over the raw logs and unpacked data for Challenge events raised by the Registry contract.
type RegistryChallengeIterator struct {
	Event *RegistryChallenge // Event containing the contract specifics and raw log

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
func (it *RegistryChallengeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryChallenge)
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
		it.Event = new(RegistryChallenge)
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
func (it *RegistryChallengeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryChallengeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryChallenge represents a Challenge event raised by the Registry contract.
type RegistryChallenge struct {
	ListingHash   [32]byte
	ChallengeID   *big.Int
	Data          string
	CommitEndDate *big.Int
	RevealEndDate *big.Int
	Challenger    common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterChallenge is a free log retrieval operation binding the contract event 0xf98a08756a3603420a080d66764f73deb1e30896c315cfed03e17f88f5eb30f7.
//
// Solidity: e _Challenge(listingHash indexed bytes32, challengeID uint256, data string, commitEndDate uint256, revealEndDate uint256, challenger indexed address)
func (_Registry *RegistryFilterer) FilterChallenge(opts *bind.FilterOpts, listingHash [][32]byte, challenger []common.Address) (*RegistryChallengeIterator, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "_Challenge", listingHashRule, challengerRule)
	if err != nil {
		return nil, err
	}
	return &RegistryChallengeIterator{contract: _Registry.contract, event: "_Challenge", logs: logs, sub: sub}, nil
}

// WatchChallenge is a free log subscription operation binding the contract event 0xf98a08756a3603420a080d66764f73deb1e30896c315cfed03e17f88f5eb30f7.
//
// Solidity: e _Challenge(listingHash indexed bytes32, challengeID uint256, data string, commitEndDate uint256, revealEndDate uint256, challenger indexed address)
func (_Registry *RegistryFilterer) WatchChallenge(opts *bind.WatchOpts, sink chan<- *RegistryChallenge, listingHash [][32]byte, challenger []common.Address) (event.Subscription, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	var challengerRule []interface{}
	for _, challengerItem := range challenger {
		challengerRule = append(challengerRule, challengerItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "_Challenge", listingHashRule, challengerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryChallenge)
				if err := _Registry.contract.UnpackLog(event, "_Challenge", log); err != nil {
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

// RegistryChallengeFailedIterator is returned from FilterChallengeFailed and is used to iterate over the raw logs and unpacked data for ChallengeFailed events raised by the Registry contract.
type RegistryChallengeFailedIterator struct {
	Event *RegistryChallengeFailed // Event containing the contract specifics and raw log

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
func (it *RegistryChallengeFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryChallengeFailed)
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
		it.Event = new(RegistryChallengeFailed)
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
func (it *RegistryChallengeFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryChallengeFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryChallengeFailed represents a ChallengeFailed event raised by the Registry contract.
type RegistryChallengeFailed struct {
	ListingHash [32]byte
	ChallengeID *big.Int
	RewardPool  *big.Int
	TotalTokens *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterChallengeFailed is a free log retrieval operation binding the contract event 0xc4497224aa78dd50c9b3e344aab02596201ca1e6dca4057a91a6c02f83f4f6c1.
//
// Solidity: e _ChallengeFailed(listingHash indexed bytes32, challengeID indexed uint256, rewardPool uint256, totalTokens uint256)
func (_Registry *RegistryFilterer) FilterChallengeFailed(opts *bind.FilterOpts, listingHash [][32]byte, challengeID []*big.Int) (*RegistryChallengeFailedIterator, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}
	var challengeIDRule []interface{}
	for _, challengeIDItem := range challengeID {
		challengeIDRule = append(challengeIDRule, challengeIDItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "_ChallengeFailed", listingHashRule, challengeIDRule)
	if err != nil {
		return nil, err
	}
	return &RegistryChallengeFailedIterator{contract: _Registry.contract, event: "_ChallengeFailed", logs: logs, sub: sub}, nil
}

// WatchChallengeFailed is a free log subscription operation binding the contract event 0xc4497224aa78dd50c9b3e344aab02596201ca1e6dca4057a91a6c02f83f4f6c1.
//
// Solidity: e _ChallengeFailed(listingHash indexed bytes32, challengeID indexed uint256, rewardPool uint256, totalTokens uint256)
func (_Registry *RegistryFilterer) WatchChallengeFailed(opts *bind.WatchOpts, sink chan<- *RegistryChallengeFailed, listingHash [][32]byte, challengeID []*big.Int) (event.Subscription, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}
	var challengeIDRule []interface{}
	for _, challengeIDItem := range challengeID {
		challengeIDRule = append(challengeIDRule, challengeIDItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "_ChallengeFailed", listingHashRule, challengeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryChallengeFailed)
				if err := _Registry.contract.UnpackLog(event, "_ChallengeFailed", log); err != nil {
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

// RegistryChallengeSucceededIterator is returned from FilterChallengeSucceeded and is used to iterate over the raw logs and unpacked data for ChallengeSucceeded events raised by the Registry contract.
type RegistryChallengeSucceededIterator struct {
	Event *RegistryChallengeSucceeded // Event containing the contract specifics and raw log

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
func (it *RegistryChallengeSucceededIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryChallengeSucceeded)
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
		it.Event = new(RegistryChallengeSucceeded)
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
func (it *RegistryChallengeSucceededIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryChallengeSucceededIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryChallengeSucceeded represents a ChallengeSucceeded event raised by the Registry contract.
type RegistryChallengeSucceeded struct {
	ListingHash [32]byte
	ChallengeID *big.Int
	RewardPool  *big.Int
	TotalTokens *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterChallengeSucceeded is a free log retrieval operation binding the contract event 0x362a12431f779a2baff4f77f75ba7960ae993a5c41b425df11f7fd0af2b9cbe6.
//
// Solidity: e _ChallengeSucceeded(listingHash indexed bytes32, challengeID indexed uint256, rewardPool uint256, totalTokens uint256)
func (_Registry *RegistryFilterer) FilterChallengeSucceeded(opts *bind.FilterOpts, listingHash [][32]byte, challengeID []*big.Int) (*RegistryChallengeSucceededIterator, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}
	var challengeIDRule []interface{}
	for _, challengeIDItem := range challengeID {
		challengeIDRule = append(challengeIDRule, challengeIDItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "_ChallengeSucceeded", listingHashRule, challengeIDRule)
	if err != nil {
		return nil, err
	}
	return &RegistryChallengeSucceededIterator{contract: _Registry.contract, event: "_ChallengeSucceeded", logs: logs, sub: sub}, nil
}

// WatchChallengeSucceeded is a free log subscription operation binding the contract event 0x362a12431f779a2baff4f77f75ba7960ae993a5c41b425df11f7fd0af2b9cbe6.
//
// Solidity: e _ChallengeSucceeded(listingHash indexed bytes32, challengeID indexed uint256, rewardPool uint256, totalTokens uint256)
func (_Registry *RegistryFilterer) WatchChallengeSucceeded(opts *bind.WatchOpts, sink chan<- *RegistryChallengeSucceeded, listingHash [][32]byte, challengeID []*big.Int) (event.Subscription, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}
	var challengeIDRule []interface{}
	for _, challengeIDItem := range challengeID {
		challengeIDRule = append(challengeIDRule, challengeIDItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "_ChallengeSucceeded", listingHashRule, challengeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryChallengeSucceeded)
				if err := _Registry.contract.UnpackLog(event, "_ChallengeSucceeded", log); err != nil {
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

// RegistryDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Registry contract.
type RegistryDepositIterator struct {
	Event *RegistryDeposit // Event containing the contract specifics and raw log

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
func (it *RegistryDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryDeposit)
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
		it.Event = new(RegistryDeposit)
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
func (it *RegistryDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryDeposit represents a Deposit event raised by the Registry contract.
type RegistryDeposit struct {
	ListingHash [32]byte
	Added       *big.Int
	NewTotal    *big.Int
	Owner       common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xf63fcfb210c709154f0260502b2586fcee5550d17dc828be3127ccdedec620ab.
//
// Solidity: e _Deposit(listingHash indexed bytes32, added uint256, newTotal uint256, owner indexed address)
func (_Registry *RegistryFilterer) FilterDeposit(opts *bind.FilterOpts, listingHash [][32]byte, owner []common.Address) (*RegistryDepositIterator, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "_Deposit", listingHashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &RegistryDepositIterator{contract: _Registry.contract, event: "_Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xf63fcfb210c709154f0260502b2586fcee5550d17dc828be3127ccdedec620ab.
//
// Solidity: e _Deposit(listingHash indexed bytes32, added uint256, newTotal uint256, owner indexed address)
func (_Registry *RegistryFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *RegistryDeposit, listingHash [][32]byte, owner []common.Address) (event.Subscription, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "_Deposit", listingHashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryDeposit)
				if err := _Registry.contract.UnpackLog(event, "_Deposit", log); err != nil {
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

// RegistryExitInitializedIterator is returned from FilterExitInitialized and is used to iterate over the raw logs and unpacked data for ExitInitialized events raised by the Registry contract.
type RegistryExitInitializedIterator struct {
	Event *RegistryExitInitialized // Event containing the contract specifics and raw log

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
func (it *RegistryExitInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryExitInitialized)
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
		it.Event = new(RegistryExitInitialized)
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
func (it *RegistryExitInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryExitInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryExitInitialized represents a ExitInitialized event raised by the Registry contract.
type RegistryExitInitialized struct {
	ListingHash      [32]byte
	ExitTime         *big.Int
	ExitDelayEndDate *big.Int
	Owner            common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterExitInitialized is a free log retrieval operation binding the contract event 0x4b137a01f77b8f1b4ccaca2abac799475d550c8adc298e399e75ee95d317b546.
//
// Solidity: e _ExitInitialized(listingHash indexed bytes32, exitTime uint256, exitDelayEndDate uint256, owner indexed address)
func (_Registry *RegistryFilterer) FilterExitInitialized(opts *bind.FilterOpts, listingHash [][32]byte, owner []common.Address) (*RegistryExitInitializedIterator, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "_ExitInitialized", listingHashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &RegistryExitInitializedIterator{contract: _Registry.contract, event: "_ExitInitialized", logs: logs, sub: sub}, nil
}

// WatchExitInitialized is a free log subscription operation binding the contract event 0x4b137a01f77b8f1b4ccaca2abac799475d550c8adc298e399e75ee95d317b546.
//
// Solidity: e _ExitInitialized(listingHash indexed bytes32, exitTime uint256, exitDelayEndDate uint256, owner indexed address)
func (_Registry *RegistryFilterer) WatchExitInitialized(opts *bind.WatchOpts, sink chan<- *RegistryExitInitialized, listingHash [][32]byte, owner []common.Address) (event.Subscription, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "_ExitInitialized", listingHashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryExitInitialized)
				if err := _Registry.contract.UnpackLog(event, "_ExitInitialized", log); err != nil {
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

// RegistryListingRemovedIterator is returned from FilterListingRemoved and is used to iterate over the raw logs and unpacked data for ListingRemoved events raised by the Registry contract.
type RegistryListingRemovedIterator struct {
	Event *RegistryListingRemoved // Event containing the contract specifics and raw log

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
func (it *RegistryListingRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryListingRemoved)
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
		it.Event = new(RegistryListingRemoved)
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
func (it *RegistryListingRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryListingRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryListingRemoved represents a ListingRemoved event raised by the Registry contract.
type RegistryListingRemoved struct {
	ListingHash [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterListingRemoved is a free log retrieval operation binding the contract event 0xd1ffb796b7108387b2f02adf47b4b81a1690cf2a190422c87a4f670780103e63.
//
// Solidity: e _ListingRemoved(listingHash indexed bytes32)
func (_Registry *RegistryFilterer) FilterListingRemoved(opts *bind.FilterOpts, listingHash [][32]byte) (*RegistryListingRemovedIterator, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "_ListingRemoved", listingHashRule)
	if err != nil {
		return nil, err
	}
	return &RegistryListingRemovedIterator{contract: _Registry.contract, event: "_ListingRemoved", logs: logs, sub: sub}, nil
}

// WatchListingRemoved is a free log subscription operation binding the contract event 0xd1ffb796b7108387b2f02adf47b4b81a1690cf2a190422c87a4f670780103e63.
//
// Solidity: e _ListingRemoved(listingHash indexed bytes32)
func (_Registry *RegistryFilterer) WatchListingRemoved(opts *bind.WatchOpts, sink chan<- *RegistryListingRemoved, listingHash [][32]byte) (event.Subscription, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "_ListingRemoved", listingHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryListingRemoved)
				if err := _Registry.contract.UnpackLog(event, "_ListingRemoved", log); err != nil {
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

// RegistryListingWithdrawnIterator is returned from FilterListingWithdrawn and is used to iterate over the raw logs and unpacked data for ListingWithdrawn events raised by the Registry contract.
type RegistryListingWithdrawnIterator struct {
	Event *RegistryListingWithdrawn // Event containing the contract specifics and raw log

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
func (it *RegistryListingWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryListingWithdrawn)
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
		it.Event = new(RegistryListingWithdrawn)
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
func (it *RegistryListingWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryListingWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryListingWithdrawn represents a ListingWithdrawn event raised by the Registry contract.
type RegistryListingWithdrawn struct {
	ListingHash [32]byte
	Owner       common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterListingWithdrawn is a free log retrieval operation binding the contract event 0x7d16ed55582bcf69d7bb762cda5b82042371fba7de3a8ebea4517079d892f522.
//
// Solidity: e _ListingWithdrawn(listingHash indexed bytes32, owner indexed address)
func (_Registry *RegistryFilterer) FilterListingWithdrawn(opts *bind.FilterOpts, listingHash [][32]byte, owner []common.Address) (*RegistryListingWithdrawnIterator, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "_ListingWithdrawn", listingHashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &RegistryListingWithdrawnIterator{contract: _Registry.contract, event: "_ListingWithdrawn", logs: logs, sub: sub}, nil
}

// WatchListingWithdrawn is a free log subscription operation binding the contract event 0x7d16ed55582bcf69d7bb762cda5b82042371fba7de3a8ebea4517079d892f522.
//
// Solidity: e _ListingWithdrawn(listingHash indexed bytes32, owner indexed address)
func (_Registry *RegistryFilterer) WatchListingWithdrawn(opts *bind.WatchOpts, sink chan<- *RegistryListingWithdrawn, listingHash [][32]byte, owner []common.Address) (event.Subscription, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "_ListingWithdrawn", listingHashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryListingWithdrawn)
				if err := _Registry.contract.UnpackLog(event, "_ListingWithdrawn", log); err != nil {
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

// RegistryRewardClaimedIterator is returned from FilterRewardClaimed and is used to iterate over the raw logs and unpacked data for RewardClaimed events raised by the Registry contract.
type RegistryRewardClaimedIterator struct {
	Event *RegistryRewardClaimed // Event containing the contract specifics and raw log

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
func (it *RegistryRewardClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryRewardClaimed)
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
		it.Event = new(RegistryRewardClaimed)
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
func (it *RegistryRewardClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryRewardClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryRewardClaimed represents a RewardClaimed event raised by the Registry contract.
type RegistryRewardClaimed struct {
	ChallengeID *big.Int
	Reward      *big.Int
	Voter       common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRewardClaimed is a free log retrieval operation binding the contract event 0x6f4c982acc31b0af2cf1dc1556f21c0325d893782d65e83c68a5534a33f59957.
//
// Solidity: e _RewardClaimed(challengeID indexed uint256, reward uint256, voter indexed address)
func (_Registry *RegistryFilterer) FilterRewardClaimed(opts *bind.FilterOpts, challengeID []*big.Int, voter []common.Address) (*RegistryRewardClaimedIterator, error) {

	var challengeIDRule []interface{}
	for _, challengeIDItem := range challengeID {
		challengeIDRule = append(challengeIDRule, challengeIDItem)
	}

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "_RewardClaimed", challengeIDRule, voterRule)
	if err != nil {
		return nil, err
	}
	return &RegistryRewardClaimedIterator{contract: _Registry.contract, event: "_RewardClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardClaimed is a free log subscription operation binding the contract event 0x6f4c982acc31b0af2cf1dc1556f21c0325d893782d65e83c68a5534a33f59957.
//
// Solidity: e _RewardClaimed(challengeID indexed uint256, reward uint256, voter indexed address)
func (_Registry *RegistryFilterer) WatchRewardClaimed(opts *bind.WatchOpts, sink chan<- *RegistryRewardClaimed, challengeID []*big.Int, voter []common.Address) (event.Subscription, error) {

	var challengeIDRule []interface{}
	for _, challengeIDItem := range challengeID {
		challengeIDRule = append(challengeIDRule, challengeIDItem)
	}

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "_RewardClaimed", challengeIDRule, voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryRewardClaimed)
				if err := _Registry.contract.UnpackLog(event, "_RewardClaimed", log); err != nil {
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

// RegistryTouchAndRemovedIterator is returned from FilterTouchAndRemoved and is used to iterate over the raw logs and unpacked data for TouchAndRemoved events raised by the Registry contract.
type RegistryTouchAndRemovedIterator struct {
	Event *RegistryTouchAndRemoved // Event containing the contract specifics and raw log

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
func (it *RegistryTouchAndRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryTouchAndRemoved)
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
		it.Event = new(RegistryTouchAndRemoved)
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
func (it *RegistryTouchAndRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryTouchAndRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryTouchAndRemoved represents a TouchAndRemoved event raised by the Registry contract.
type RegistryTouchAndRemoved struct {
	ListingHash [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTouchAndRemoved is a free log retrieval operation binding the contract event 0x4a9ee335af9e32f32f2229943dc7a0d3b5adf7e4c5c4062b372eae8c476d9286.
//
// Solidity: e _TouchAndRemoved(listingHash indexed bytes32)
func (_Registry *RegistryFilterer) FilterTouchAndRemoved(opts *bind.FilterOpts, listingHash [][32]byte) (*RegistryTouchAndRemovedIterator, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "_TouchAndRemoved", listingHashRule)
	if err != nil {
		return nil, err
	}
	return &RegistryTouchAndRemovedIterator{contract: _Registry.contract, event: "_TouchAndRemoved", logs: logs, sub: sub}, nil
}

// WatchTouchAndRemoved is a free log subscription operation binding the contract event 0x4a9ee335af9e32f32f2229943dc7a0d3b5adf7e4c5c4062b372eae8c476d9286.
//
// Solidity: e _TouchAndRemoved(listingHash indexed bytes32)
func (_Registry *RegistryFilterer) WatchTouchAndRemoved(opts *bind.WatchOpts, sink chan<- *RegistryTouchAndRemoved, listingHash [][32]byte) (event.Subscription, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "_TouchAndRemoved", listingHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryTouchAndRemoved)
				if err := _Registry.contract.UnpackLog(event, "_TouchAndRemoved", log); err != nil {
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

// RegistryWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the Registry contract.
type RegistryWithdrawalIterator struct {
	Event *RegistryWithdrawal // Event containing the contract specifics and raw log

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
func (it *RegistryWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryWithdrawal)
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
		it.Event = new(RegistryWithdrawal)
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
func (it *RegistryWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryWithdrawal represents a Withdrawal event raised by the Registry contract.
type RegistryWithdrawal struct {
	ListingHash [32]byte
	Withdrew    *big.Int
	NewTotal    *big.Int
	Owner       common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x9d9ed58779badf90c56d72f3b54def9f73dc875d8f86416c8334b55328c6c106.
//
// Solidity: e _Withdrawal(listingHash indexed bytes32, withdrew uint256, newTotal uint256, owner indexed address)
func (_Registry *RegistryFilterer) FilterWithdrawal(opts *bind.FilterOpts, listingHash [][32]byte, owner []common.Address) (*RegistryWithdrawalIterator, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "_Withdrawal", listingHashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &RegistryWithdrawalIterator{contract: _Registry.contract, event: "_Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x9d9ed58779badf90c56d72f3b54def9f73dc875d8f86416c8334b55328c6c106.
//
// Solidity: e _Withdrawal(listingHash indexed bytes32, withdrew uint256, newTotal uint256, owner indexed address)
func (_Registry *RegistryFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *RegistryWithdrawal, listingHash [][32]byte, owner []common.Address) (event.Subscription, error) {

	var listingHashRule []interface{}
	for _, listingHashItem := range listingHash {
		listingHashRule = append(listingHashRule, listingHashItem)
	}

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "_Withdrawal", listingHashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryWithdrawal)
				if err := _Registry.contract.UnpackLog(event, "_Withdrawal", log); err != nil {
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
