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

// MultiSigWalletFactoryABI is the input ABI used to generate the binding from.
const MultiSigWalletFactoryABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"isInstantiation\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"instantiations\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"creator\",\"type\":\"address\"}],\"name\":\"getInstantiationCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"instantiation\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"identifier\",\"type\":\"uint256\"}],\"name\":\"ContractInstantiation\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_owners\",\"type\":\"address[]\"},{\"name\":\"_required\",\"type\":\"uint256\"},{\"name\":\"identifier\",\"type\":\"uint256\"}],\"name\":\"create\",\"outputs\":[{\"name\":\"wallet\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// MultiSigWalletFactory is an auto generated Go binding around an Ethereum contract.
type MultiSigWalletFactory struct {
	MultiSigWalletFactoryCaller     // Read-only binding to the contract
	MultiSigWalletFactoryTransactor // Write-only binding to the contract
	MultiSigWalletFactoryFilterer   // Log filterer for contract events
}

// MultiSigWalletFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultiSigWalletFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiSigWalletFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultiSigWalletFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiSigWalletFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultiSigWalletFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiSigWalletFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultiSigWalletFactorySession struct {
	Contract     *MultiSigWalletFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// MultiSigWalletFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultiSigWalletFactoryCallerSession struct {
	Contract *MultiSigWalletFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// MultiSigWalletFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultiSigWalletFactoryTransactorSession struct {
	Contract     *MultiSigWalletFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// MultiSigWalletFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultiSigWalletFactoryRaw struct {
	Contract *MultiSigWalletFactory // Generic contract binding to access the raw methods on
}

// MultiSigWalletFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultiSigWalletFactoryCallerRaw struct {
	Contract *MultiSigWalletFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// MultiSigWalletFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultiSigWalletFactoryTransactorRaw struct {
	Contract *MultiSigWalletFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultiSigWalletFactory creates a new instance of MultiSigWalletFactory, bound to a specific deployed contract.
func NewMultiSigWalletFactory(address common.Address, backend bind.ContractBackend) (*MultiSigWalletFactory, error) {
	contract, err := bindMultiSigWalletFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultiSigWalletFactory{MultiSigWalletFactoryCaller: MultiSigWalletFactoryCaller{contract: contract}, MultiSigWalletFactoryTransactor: MultiSigWalletFactoryTransactor{contract: contract}, MultiSigWalletFactoryFilterer: MultiSigWalletFactoryFilterer{contract: contract}}, nil
}

// NewMultiSigWalletFactoryCaller creates a new read-only instance of MultiSigWalletFactory, bound to a specific deployed contract.
func NewMultiSigWalletFactoryCaller(address common.Address, caller bind.ContractCaller) (*MultiSigWalletFactoryCaller, error) {
	contract, err := bindMultiSigWalletFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiSigWalletFactoryCaller{contract: contract}, nil
}

// NewMultiSigWalletFactoryTransactor creates a new write-only instance of MultiSigWalletFactory, bound to a specific deployed contract.
func NewMultiSigWalletFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*MultiSigWalletFactoryTransactor, error) {
	contract, err := bindMultiSigWalletFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultiSigWalletFactoryTransactor{contract: contract}, nil
}

// NewMultiSigWalletFactoryFilterer creates a new log filterer instance of MultiSigWalletFactory, bound to a specific deployed contract.
func NewMultiSigWalletFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*MultiSigWalletFactoryFilterer, error) {
	contract, err := bindMultiSigWalletFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultiSigWalletFactoryFilterer{contract: contract}, nil
}

// bindMultiSigWalletFactory binds a generic wrapper to an already deployed contract.
func bindMultiSigWalletFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MultiSigWalletFactoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiSigWalletFactory *MultiSigWalletFactoryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MultiSigWalletFactory.Contract.MultiSigWalletFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiSigWalletFactory *MultiSigWalletFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiSigWalletFactory.Contract.MultiSigWalletFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiSigWalletFactory *MultiSigWalletFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiSigWalletFactory.Contract.MultiSigWalletFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiSigWalletFactory *MultiSigWalletFactoryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MultiSigWalletFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiSigWalletFactory *MultiSigWalletFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiSigWalletFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiSigWalletFactory *MultiSigWalletFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiSigWalletFactory.Contract.contract.Transact(opts, method, params...)
}

// GetInstantiationCount is a free data retrieval call binding the contract method 0x8f838478.
//
// Solidity: function getInstantiationCount(creator address) constant returns(uint256)
func (_MultiSigWalletFactory *MultiSigWalletFactoryCaller) GetInstantiationCount(opts *bind.CallOpts, creator common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MultiSigWalletFactory.contract.Call(opts, out, "getInstantiationCount", creator)
	return *ret0, err
}

// GetInstantiationCount is a free data retrieval call binding the contract method 0x8f838478.
//
// Solidity: function getInstantiationCount(creator address) constant returns(uint256)
func (_MultiSigWalletFactory *MultiSigWalletFactorySession) GetInstantiationCount(creator common.Address) (*big.Int, error) {
	return _MultiSigWalletFactory.Contract.GetInstantiationCount(&_MultiSigWalletFactory.CallOpts, creator)
}

// GetInstantiationCount is a free data retrieval call binding the contract method 0x8f838478.
//
// Solidity: function getInstantiationCount(creator address) constant returns(uint256)
func (_MultiSigWalletFactory *MultiSigWalletFactoryCallerSession) GetInstantiationCount(creator common.Address) (*big.Int, error) {
	return _MultiSigWalletFactory.Contract.GetInstantiationCount(&_MultiSigWalletFactory.CallOpts, creator)
}

// Instantiations is a free data retrieval call binding the contract method 0x57183c82.
//
// Solidity: function instantiations( address,  uint256) constant returns(address)
func (_MultiSigWalletFactory *MultiSigWalletFactoryCaller) Instantiations(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _MultiSigWalletFactory.contract.Call(opts, out, "instantiations", arg0, arg1)
	return *ret0, err
}

// Instantiations is a free data retrieval call binding the contract method 0x57183c82.
//
// Solidity: function instantiations( address,  uint256) constant returns(address)
func (_MultiSigWalletFactory *MultiSigWalletFactorySession) Instantiations(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _MultiSigWalletFactory.Contract.Instantiations(&_MultiSigWalletFactory.CallOpts, arg0, arg1)
}

// Instantiations is a free data retrieval call binding the contract method 0x57183c82.
//
// Solidity: function instantiations( address,  uint256) constant returns(address)
func (_MultiSigWalletFactory *MultiSigWalletFactoryCallerSession) Instantiations(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _MultiSigWalletFactory.Contract.Instantiations(&_MultiSigWalletFactory.CallOpts, arg0, arg1)
}

// IsInstantiation is a free data retrieval call binding the contract method 0x2f4f3316.
//
// Solidity: function isInstantiation( address) constant returns(bool)
func (_MultiSigWalletFactory *MultiSigWalletFactoryCaller) IsInstantiation(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _MultiSigWalletFactory.contract.Call(opts, out, "isInstantiation", arg0)
	return *ret0, err
}

// IsInstantiation is a free data retrieval call binding the contract method 0x2f4f3316.
//
// Solidity: function isInstantiation( address) constant returns(bool)
func (_MultiSigWalletFactory *MultiSigWalletFactorySession) IsInstantiation(arg0 common.Address) (bool, error) {
	return _MultiSigWalletFactory.Contract.IsInstantiation(&_MultiSigWalletFactory.CallOpts, arg0)
}

// IsInstantiation is a free data retrieval call binding the contract method 0x2f4f3316.
//
// Solidity: function isInstantiation( address) constant returns(bool)
func (_MultiSigWalletFactory *MultiSigWalletFactoryCallerSession) IsInstantiation(arg0 common.Address) (bool, error) {
	return _MultiSigWalletFactory.Contract.IsInstantiation(&_MultiSigWalletFactory.CallOpts, arg0)
}

// Create is a paid mutator transaction binding the contract method 0x53d9d910.
//
// Solidity: function create(_owners address[], _required uint256, identifier uint256) returns(wallet address)
func (_MultiSigWalletFactory *MultiSigWalletFactoryTransactor) Create(opts *bind.TransactOpts, _owners []common.Address, _required *big.Int, identifier *big.Int) (*types.Transaction, error) {
	return _MultiSigWalletFactory.contract.Transact(opts, "create", _owners, _required, identifier)
}

// Create is a paid mutator transaction binding the contract method 0x53d9d910.
//
// Solidity: function create(_owners address[], _required uint256, identifier uint256) returns(wallet address)
func (_MultiSigWalletFactory *MultiSigWalletFactorySession) Create(_owners []common.Address, _required *big.Int, identifier *big.Int) (*types.Transaction, error) {
	return _MultiSigWalletFactory.Contract.Create(&_MultiSigWalletFactory.TransactOpts, _owners, _required, identifier)
}

// Create is a paid mutator transaction binding the contract method 0x53d9d910.
//
// Solidity: function create(_owners address[], _required uint256, identifier uint256) returns(wallet address)
func (_MultiSigWalletFactory *MultiSigWalletFactoryTransactorSession) Create(_owners []common.Address, _required *big.Int, identifier *big.Int) (*types.Transaction, error) {
	return _MultiSigWalletFactory.Contract.Create(&_MultiSigWalletFactory.TransactOpts, _owners, _required, identifier)
}

// MultiSigWalletFactoryContractInstantiationIterator is returned from FilterContractInstantiation and is used to iterate over the raw logs and unpacked data for ContractInstantiation events raised by the MultiSigWalletFactory contract.
type MultiSigWalletFactoryContractInstantiationIterator struct {
	Event *MultiSigWalletFactoryContractInstantiation // Event containing the contract specifics and raw log

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
func (it *MultiSigWalletFactoryContractInstantiationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiSigWalletFactoryContractInstantiation)
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
		it.Event = new(MultiSigWalletFactoryContractInstantiation)
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
func (it *MultiSigWalletFactoryContractInstantiationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiSigWalletFactoryContractInstantiationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiSigWalletFactoryContractInstantiation represents a ContractInstantiation event raised by the MultiSigWalletFactory contract.
type MultiSigWalletFactoryContractInstantiation struct {
	Sender        common.Address
	Instantiation common.Address
	Identifier    *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterContractInstantiation is a free log retrieval operation binding the contract event 0xb004d2590bf6eda75630c427abba336b8b58444d1cd8634dd89c811f5fcb5a11.
//
// Solidity: e ContractInstantiation(sender address, instantiation address, identifier uint256)
func (_MultiSigWalletFactory *MultiSigWalletFactoryFilterer) FilterContractInstantiation(opts *bind.FilterOpts) (*MultiSigWalletFactoryContractInstantiationIterator, error) {

	logs, sub, err := _MultiSigWalletFactory.contract.FilterLogs(opts, "ContractInstantiation")
	if err != nil {
		return nil, err
	}
	return &MultiSigWalletFactoryContractInstantiationIterator{contract: _MultiSigWalletFactory.contract, event: "ContractInstantiation", logs: logs, sub: sub}, nil
}

// WatchContractInstantiation is a free log subscription operation binding the contract event 0xb004d2590bf6eda75630c427abba336b8b58444d1cd8634dd89c811f5fcb5a11.
//
// Solidity: e ContractInstantiation(sender address, instantiation address, identifier uint256)
func (_MultiSigWalletFactory *MultiSigWalletFactoryFilterer) WatchContractInstantiation(opts *bind.WatchOpts, sink chan<- *MultiSigWalletFactoryContractInstantiation) (event.Subscription, error) {

	logs, sub, err := _MultiSigWalletFactory.contract.WatchLogs(opts, "ContractInstantiation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiSigWalletFactoryContractInstantiation)
				if err := _MultiSigWalletFactory.contract.UnpackLog(event, "ContractInstantiation", log); err != nil {
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
