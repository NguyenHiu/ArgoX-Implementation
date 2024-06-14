package supermatcher

import (
	"fmt"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"
)

type SuperMatcher struct {
	OnchainInstance *onchain.Onchain
	Auth            *bind.TransactOpts
	Client          *ethclient.Client
	Address         common.Address
	Batches         []*Batch
	Orders          map[uuid.UUID]bool

	Port int
}

func (sm *SuperMatcher) isExists(order *Order) bool {
	_, ok := sm.Orders[order.OrderID]
	return ok
}

func (sm *SuperMatcher) addOrder(order *Order) {
	sm.Orders[order.OrderID] = true
}

func (sm *SuperMatcher) Process() {
	// Check if having any batch to be processed
	if len(sm.Batches) == 0 {
		return
	}

	// Get the first batch in queue
	batch := sm.Batches[0]
	sm.Batches = sm.Batches[1:]

	// The batch is already validated when being appended to the sm.Batches

	// Filter orders in the batch
	validOrders := []*Order{}
	for idx, order := range batch.Orders {
		if sm.isExists(order) {
			fmt.Printf("order at %v has already existed\n", idx)
		} else {
			sm.addOrder(order)
			validOrders = append(validOrders, order)
		}
	}

	// If the batch is empty, stop
	if len(validOrders) == 0 {
		fmt.Printf("batch (%v) is empty", batch.BatchID)
		return
	}

	// Update `amount` & `orders` of the batch
	if len(validOrders) != len(batch.Orders) {
		_amount := &big.Int{}
		for _, order := range validOrders {
			_amount = new(big.Int).Add(_amount, order.Amount)
		}
		batch.Amount = _amount
		batch.Orders = validOrders
	}

	// Send batch to smart contract

}

func (sm *SuperMatcher) CheckValidBatch(batch *Batch) bool {
	// 1. valid owner: owner is a matcher & the signature is valid
	if !sm.isMatcher(batch.Owner) {
		log.Error("Invalid Matcher")
		return false
	}

	if !batch.IsValidSignature() {
		log.Error("Invalid Batch's Signature")
		return false
	}

	// 2. check signatures of orders in the batch
	for idx, order := range batch.Orders {
		if !order.IsValidSignature() {
			log.Error("Invalid Order's Signature at", idx)
			return false
		}
	}

	return true
}

// func (sm *SuperMatcher)

func (sm *SuperMatcher) AddBatch(batch *Batch) {
	if sm.CheckValidBatch(batch) {
		sm.Batches = append(sm.Batches, batch)
	}
}
