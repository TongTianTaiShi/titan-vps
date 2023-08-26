// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package filecoinbridge

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

// DescribePriceRequestDataDisk is an auto generated low-level Go binding around an user-defined struct.
type DescribePriceRequestDataDisk struct {
	Category         string
	PerformanceLevel string
	Size             int64
}

// IpcOrderInfo is an auto generated low-level Go binding around an user-defined struct.
type IpcOrderInfo struct {
	OrderID                 string
	UserID                  string
	Value                   string
	CreatedTime             string
	Type                    string
	PeriodUnit              string
	Period                  int32
	RegionId                string
	InstanceId              string
	InstanceType            string
	ImageId                 string
	Memory                  string
	MemoryUsed              string
	Cores                   int32
	CoresUsed               string
	SecurityGroupId         string
	InstanceChargeType      string
	InternetMaxBandwidthOut int32
	InternetMaxBandwidthIn  int32
	IpAddress               string
	TradePrice              string
	SystemDiskCategory      string
	OSType                  string
	InternetChargeType      string
	SystemDiskSize          int32
	DataDiskString          string
	DataDisk                []DescribePriceRequestDataDisk
}

// FvmMetaData contains all meta data concerning the Fvm contract.
var FvmMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"OrderID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"UserID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"Value\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"CreatedTime\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"Type\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"PeriodUnit\",\"type\":\"string\"},{\"internalType\":\"int32\",\"name\":\"Period\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"RegionId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"InstanceId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"InstanceType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"ImageId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"Memory\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"MemoryUsed\",\"type\":\"string\"},{\"internalType\":\"int32\",\"name\":\"Cores\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"CoresUsed\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"SecurityGroupId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"InstanceChargeType\",\"type\":\"string\"},{\"internalType\":\"int32\",\"name\":\"InternetMaxBandwidthOut\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"InternetMaxBandwidthIn\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"IpAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"TradePrice\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"SystemDiskCategory\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"OSType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"InternetChargeType\",\"type\":\"string\"},{\"internalType\":\"int32\",\"name\":\"SystemDiskSize\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DataDiskString\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"Category\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"PerformanceLevel\",\"type\":\"string\"},{\"internalType\":\"int64\",\"name\":\"Size\",\"type\":\"int64\"}],\"internalType\":\"structDescribePriceRequestDataDisk[]\",\"name\":\"DataDisk\",\"type\":\"tuple[]\"}],\"internalType\":\"structIpcOrderInfo\",\"name\":\"x\",\"type\":\"tuple\"}],\"name\":\"setOrderInfo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"orderID\",\"type\":\"string\"}],\"name\":\"getOrderInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"OrderID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"UserID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"Value\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"CreatedTime\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"Type\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"PeriodUnit\",\"type\":\"string\"},{\"internalType\":\"int32\",\"name\":\"Period\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"RegionId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"InstanceId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"InstanceType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"ImageId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"Memory\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"MemoryUsed\",\"type\":\"string\"},{\"internalType\":\"int32\",\"name\":\"Cores\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"CoresUsed\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"SecurityGroupId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"InstanceChargeType\",\"type\":\"string\"},{\"internalType\":\"int32\",\"name\":\"InternetMaxBandwidthOut\",\"type\":\"int32\"},{\"internalType\":\"int32\",\"name\":\"InternetMaxBandwidthIn\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"IpAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"TradePrice\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"SystemDiskCategory\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"OSType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"InternetChargeType\",\"type\":\"string\"},{\"internalType\":\"int32\",\"name\":\"SystemDiskSize\",\"type\":\"int32\"},{\"internalType\":\"string\",\"name\":\"DataDiskString\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"Category\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"PerformanceLevel\",\"type\":\"string\"},{\"internalType\":\"int64\",\"name\":\"Size\",\"type\":\"int64\"}],\"internalType\":\"structDescribePriceRequestDataDisk[]\",\"name\":\"DataDisk\",\"type\":\"tuple[]\"}],\"internalType\":\"structIpcOrderInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// FvmABI is the input ABI used to generate the binding from.
// Deprecated: Use FvmMetaData.ABI instead.
var FvmABI = FvmMetaData.ABI

// Fvm is an auto generated Go binding around an Ethereum contract.
type Fvm struct {
	FvmCaller     // Read-only binding to the contract
	FvmTransactor // Write-only binding to the contract
	FvmFilterer   // Log filterer for contract events
}

// FvmCaller is an auto generated read-only Go binding around an Ethereum contract.
type FvmCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FvmTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FvmTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FvmFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FvmFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FvmSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FvmSession struct {
	Contract     *Fvm              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FvmCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FvmCallerSession struct {
	Contract *FvmCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FvmTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FvmTransactorSession struct {
	Contract     *FvmTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FvmRaw is an auto generated low-level Go binding around an Ethereum contract.
type FvmRaw struct {
	Contract *Fvm // Generic contract binding to access the raw methods on
}

// FvmCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FvmCallerRaw struct {
	Contract *FvmCaller // Generic read-only contract binding to access the raw methods on
}

// FvmTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FvmTransactorRaw struct {
	Contract *FvmTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFvm creates a new instance of Fvm, bound to a specific deployed contract.
func NewFvm(address common.Address, backend bind.ContractBackend) (*Fvm, error) {
	contract, err := bindFvm(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Fvm{FvmCaller: FvmCaller{contract: contract}, FvmTransactor: FvmTransactor{contract: contract}, FvmFilterer: FvmFilterer{contract: contract}}, nil
}

// NewFvmCaller creates a new read-only instance of Fvm, bound to a specific deployed contract.
func NewFvmCaller(address common.Address, caller bind.ContractCaller) (*FvmCaller, error) {
	contract, err := bindFvm(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FvmCaller{contract: contract}, nil
}

// NewFvmTransactor creates a new write-only instance of Fvm, bound to a specific deployed contract.
func NewFvmTransactor(address common.Address, transactor bind.ContractTransactor) (*FvmTransactor, error) {
	contract, err := bindFvm(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FvmTransactor{contract: contract}, nil
}

// NewFvmFilterer creates a new log filterer instance of Fvm, bound to a specific deployed contract.
func NewFvmFilterer(address common.Address, filterer bind.ContractFilterer) (*FvmFilterer, error) {
	contract, err := bindFvm(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FvmFilterer{contract: contract}, nil
}

// bindFvm binds a generic wrapper to an already deployed contract.
func bindFvm(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FvmMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Fvm *FvmRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Fvm.Contract.FvmCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Fvm *FvmRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fvm.Contract.FvmTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Fvm *FvmRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Fvm.Contract.FvmTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Fvm *FvmCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Fvm.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Fvm *FvmTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fvm.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Fvm *FvmTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Fvm.Contract.contract.Transact(opts, method, params...)
}

// GetOrderInfo is a free data retrieval call binding the contract method 0x15eb3f30.
//
// Solidity: function getOrderInfo(string orderID) view returns((string,string,string,string,string,string,int32,string,string,string,string,string,string,int32,string,string,string,int32,int32,string,string,string,string,string,int32,string,(string,string,int64)[]))
func (_Fvm *FvmCaller) GetOrderInfo(opts *bind.CallOpts, orderID string) (IpcOrderInfo, error) {
	var out []interface{}
	err := _Fvm.contract.Call(opts, &out, "getOrderInfo", orderID)

	if err != nil {
		return *new(IpcOrderInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IpcOrderInfo)).(*IpcOrderInfo)

	return out0, err

}

// GetOrderInfo is a free data retrieval call binding the contract method 0x15eb3f30.
//
// Solidity: function getOrderInfo(string orderID) view returns((string,string,string,string,string,string,int32,string,string,string,string,string,string,int32,string,string,string,int32,int32,string,string,string,string,string,int32,string,(string,string,int64)[]))
func (_Fvm *FvmSession) GetOrderInfo(orderID string) (IpcOrderInfo, error) {
	return _Fvm.Contract.GetOrderInfo(&_Fvm.CallOpts, orderID)
}

// GetOrderInfo is a free data retrieval call binding the contract method 0x15eb3f30.
//
// Solidity: function getOrderInfo(string orderID) view returns((string,string,string,string,string,string,int32,string,string,string,string,string,string,int32,string,string,string,int32,int32,string,string,string,string,string,int32,string,(string,string,int64)[]))
func (_Fvm *FvmCallerSession) GetOrderInfo(orderID string) (IpcOrderInfo, error) {
	return _Fvm.Contract.GetOrderInfo(&_Fvm.CallOpts, orderID)
}

// SetOrderInfo is a paid mutator transaction binding the contract method 0x4aa7d79b.
//
// Solidity: function setOrderInfo((string,string,string,string,string,string,int32,string,string,string,string,string,string,int32,string,string,string,int32,int32,string,string,string,string,string,int32,string,(string,string,int64)[]) x) returns()
func (_Fvm *FvmTransactor) SetOrderInfo(opts *bind.TransactOpts, x IpcOrderInfo) (*types.Transaction, error) {
	return _Fvm.contract.Transact(opts, "setOrderInfo", x)
}

// SetOrderInfo is a paid mutator transaction binding the contract method 0x4aa7d79b.
//
// Solidity: function setOrderInfo((string,string,string,string,string,string,int32,string,string,string,string,string,string,int32,string,string,string,int32,int32,string,string,string,string,string,int32,string,(string,string,int64)[]) x) returns()
func (_Fvm *FvmSession) SetOrderInfo(x IpcOrderInfo) (*types.Transaction, error) {
	return _Fvm.Contract.SetOrderInfo(&_Fvm.TransactOpts, x)
}

// SetOrderInfo is a paid mutator transaction binding the contract method 0x4aa7d79b.
//
// Solidity: function setOrderInfo((string,string,string,string,string,string,int32,string,string,string,string,string,string,int32,string,string,string,int32,int32,string,string,string,string,string,int32,string,(string,string,int64)[]) x) returns()
func (_Fvm *FvmTransactorSession) SetOrderInfo(x IpcOrderInfo) (*types.Transaction, error) {
	return _Fvm.Contract.SetOrderInfo(&_Fvm.TransactOpts, x)
}
