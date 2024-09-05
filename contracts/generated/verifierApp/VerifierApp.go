// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifierApp

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ChannelAllocation is an auto generated low-level Go binding around an user-defined struct.
type ChannelAllocation struct {
	Assets   []ChannelAsset
	Balances [][]*big.Int
	Locked   []ChannelSubAlloc
}

// ChannelAsset is an auto generated low-level Go binding around an user-defined struct.
type ChannelAsset struct {
	ChainID *big.Int
	Holder  common.Address
}

// ChannelParams is an auto generated low-level Go binding around an user-defined struct.
type ChannelParams struct {
	ChallengeDuration *big.Int
	Nonce             *big.Int
	Participants      []common.Address
	App               common.Address
	LedgerChannel     bool
	VirtualChannel    bool
}

// ChannelState is an auto generated low-level Go binding around an user-defined struct.
type ChannelState struct {
	ChannelID [32]byte
	Version   uint64
	Outcome   ChannelAllocation
	AppData   []byte
	IsFinal   bool
}

// ChannelSubAlloc is an auto generated low-level Go binding around an user-defined struct.
type ChannelSubAlloc struct {
	ID       [32]byte
	Balances []*big.Int
	IndexMap []uint16
}

// VerifierAppMetaData contains all meta data concerning the VerifierApp contract.
var VerifierAppMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"challengeDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"app\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"ledgerChannel\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"virtualChannel\",\"type\":\"bool\"}],\"internalType\":\"structChannel.Params\",\"name\":\"params\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"holder\",\"type\":\"address\"}],\"internalType\":\"structChannel.Asset[]\",\"name\":\"assets\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"from\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"channelID\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"holder\",\"type\":\"address\"}],\"internalType\":\"structChannel.Asset[]\",\"name\":\"assets\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[][]\",\"name\":\"balances\",\"type\":\"uint256[][]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"indexMap\",\"type\":\"uint16[]\"}],\"internalType\":\"structChannel.SubAlloc[]\",\"name\":\"locked\",\"type\":\"tuple[]\"}],\"internalType\":\"structChannel.Allocation\",\"name\":\"outcome\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structChannel.State\",\"name\":\"to\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"actorIdx\",\"type\":\"uint256\"}],\"name\":\"validTransition\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b506101b88061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610029575f3560e01c80636d7eba0d1461002d575b5f80fd5b610047600480360381019061004291906100ca565b610049565b005b50505050565b5f80fd5b5f80fd5b5f80fd5b5f60c082840312156100705761006f610057565b5b81905092915050565b5f60a0828403121561008e5761008d610057565b5b81905092915050565b5f819050919050565b6100a981610097565b81146100b3575f80fd5b50565b5f813590506100c4816100a0565b92915050565b5f805f80608085870312156100e2576100e161004f565b5b5f85013567ffffffffffffffff8111156100ff576100fe610053565b5b61010b8782880161005b565b945050602085013567ffffffffffffffff81111561012c5761012b610053565b5b61013887828801610079565b935050604085013567ffffffffffffffff81111561015957610158610053565b5b61016587828801610079565b9250506060610176878288016100b6565b9150509295919450925056fea264697066735822122078e36d860aa0c55dc7128cf72e7cf365208a220908392635806a7043ac80be6764736f6c63430008180033",
}

// VerifierAppABI is the input ABI used to generate the binding from.
// Deprecated: Use VerifierAppMetaData.ABI instead.
var VerifierAppABI = VerifierAppMetaData.ABI

// VerifierAppBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VerifierAppMetaData.Bin instead.
var VerifierAppBin = VerifierAppMetaData.Bin

// DeployVerifierApp deploys a new Ethereum contract, binding an instance of VerifierApp to it.
func DeployVerifierApp(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VerifierApp, error) {
	parsed, err := VerifierAppMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VerifierAppBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VerifierApp{VerifierAppCaller: VerifierAppCaller{contract: contract}, VerifierAppTransactor: VerifierAppTransactor{contract: contract}, VerifierAppFilterer: VerifierAppFilterer{contract: contract}}, nil
}

// VerifierApp is an auto generated Go binding around an Ethereum contract.
type VerifierApp struct {
	VerifierAppCaller     // Read-only binding to the contract
	VerifierAppTransactor // Write-only binding to the contract
	VerifierAppFilterer   // Log filterer for contract events
}

// VerifierAppCaller is an auto generated read-only Go binding around an Ethereum contract.
type VerifierAppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierAppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VerifierAppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierAppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VerifierAppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VerifierAppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VerifierAppSession struct {
	Contract     *VerifierApp      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VerifierAppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VerifierAppCallerSession struct {
	Contract *VerifierAppCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// VerifierAppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VerifierAppTransactorSession struct {
	Contract     *VerifierAppTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// VerifierAppRaw is an auto generated low-level Go binding around an Ethereum contract.
type VerifierAppRaw struct {
	Contract *VerifierApp // Generic contract binding to access the raw methods on
}

// VerifierAppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VerifierAppCallerRaw struct {
	Contract *VerifierAppCaller // Generic read-only contract binding to access the raw methods on
}

// VerifierAppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VerifierAppTransactorRaw struct {
	Contract *VerifierAppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVerifierApp creates a new instance of VerifierApp, bound to a specific deployed contract.
func NewVerifierApp(address common.Address, backend bind.ContractBackend) (*VerifierApp, error) {
	contract, err := bindVerifierApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VerifierApp{VerifierAppCaller: VerifierAppCaller{contract: contract}, VerifierAppTransactor: VerifierAppTransactor{contract: contract}, VerifierAppFilterer: VerifierAppFilterer{contract: contract}}, nil
}

// NewVerifierAppCaller creates a new read-only instance of VerifierApp, bound to a specific deployed contract.
func NewVerifierAppCaller(address common.Address, caller bind.ContractCaller) (*VerifierAppCaller, error) {
	contract, err := bindVerifierApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierAppCaller{contract: contract}, nil
}

// NewVerifierAppTransactor creates a new write-only instance of VerifierApp, bound to a specific deployed contract.
func NewVerifierAppTransactor(address common.Address, transactor bind.ContractTransactor) (*VerifierAppTransactor, error) {
	contract, err := bindVerifierApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VerifierAppTransactor{contract: contract}, nil
}

// NewVerifierAppFilterer creates a new log filterer instance of VerifierApp, bound to a specific deployed contract.
func NewVerifierAppFilterer(address common.Address, filterer bind.ContractFilterer) (*VerifierAppFilterer, error) {
	contract, err := bindVerifierApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VerifierAppFilterer{contract: contract}, nil
}

// bindVerifierApp binds a generic wrapper to an already deployed contract.
func bindVerifierApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VerifierAppMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VerifierApp *VerifierAppRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifierApp.Contract.VerifierAppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VerifierApp *VerifierAppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifierApp.Contract.VerifierAppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VerifierApp *VerifierAppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifierApp.Contract.VerifierAppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VerifierApp *VerifierAppCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VerifierApp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VerifierApp *VerifierAppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VerifierApp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VerifierApp *VerifierAppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VerifierApp.Contract.contract.Transact(opts, method, params...)
}

// ValidTransition is a free data retrieval call binding the contract method 0x6d7eba0d.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,((uint256,address)[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,((uint256,address)[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 actorIdx) pure returns()
func (_VerifierApp *VerifierAppCaller) ValidTransition(opts *bind.CallOpts, params ChannelParams, from ChannelState, to ChannelState, actorIdx *big.Int) error {
	var out []interface{}
	err := _VerifierApp.contract.Call(opts, &out, "validTransition", params, from, to, actorIdx)

	if err != nil {
		return err
	}

	return err

}

// ValidTransition is a free data retrieval call binding the contract method 0x6d7eba0d.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,((uint256,address)[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,((uint256,address)[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 actorIdx) pure returns()
func (_VerifierApp *VerifierAppSession) ValidTransition(params ChannelParams, from ChannelState, to ChannelState, actorIdx *big.Int) error {
	return _VerifierApp.Contract.ValidTransition(&_VerifierApp.CallOpts, params, from, to, actorIdx)
}

// ValidTransition is a free data retrieval call binding the contract method 0x6d7eba0d.
//
// Solidity: function validTransition((uint256,uint256,address[],address,bool,bool) params, (bytes32,uint64,((uint256,address)[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) from, (bytes32,uint64,((uint256,address)[],uint256[][],(bytes32,uint256[],uint16[])[]),bytes,bool) to, uint256 actorIdx) pure returns()
func (_VerifierApp *VerifierAppCallerSession) ValidTransition(params ChannelParams, from ChannelState, to ChannelState, actorIdx *big.Int) error {
	return _VerifierApp.Contract.ValidTransition(&_VerifierApp.CallOpts, params, from, to, actorIdx)
}
