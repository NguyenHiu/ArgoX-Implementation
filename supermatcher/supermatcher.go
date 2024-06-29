package supermatcher

import (
	"math/big"
	"sync"
	"time"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

// var _logger = logger.NewLogger("\033[1mSuper Matcher\033[0m", logger.None, logger.None)
var _logger = logger.NewLogger("Super Matcher", logger.White, logger.Bold)

type SuperMatcher struct {
	OnchainInstance *onchain.Onchain
	Auth            *bind.TransactOpts
	Client          *ethclient.Client
	Address         common.Address
	Batches         []*Batch
	Orders          map[uuid.UUID]bool

	Port  int
	Mutex sync.Mutex
}

func NewSuperMatcher(onchain *onchain.Onchain, privateKeyHex string, port int, chainID int) (*SuperMatcher, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}
	addr := crypto.PubkeyToAddress(privateKey.PublicKey)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(chainID)))
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(constants.CHAIN_URL)
	if err != nil {
		return nil, err
	}

	sm := &SuperMatcher{
		OnchainInstance: onchain,
		Auth:            auth,
		Client:          client,
		Address:         addr,
		Batches:         []*Batch{},
		Orders:          make(map[uuid.UUID]bool),
		Port:            port,
	}

	go sm.processing()

	return sm, nil
}

func (sm *SuperMatcher) isExists(order *Order) bool {
	_, ok := sm.Orders[order.OrderID]
	return ok
}

func (sm *SuperMatcher) addOrder(order *Order) {
	sm.Orders[order.OrderID] = true
}

func (sm *SuperMatcher) Process() {
	sm.Mutex.Lock()

	// Check if having any batch to be processed
	if len(sm.Batches) == 0 {
		sm.Mutex.Unlock()
		return
	}

	// Get the first batch in queue
	batch := sm.Batches[0]
	sm.Batches = sm.Batches[1:]
	sm.Mutex.Unlock()

	// The batch is already validated when being appended to the sm.Batches

	// Filter orders in the batch
	validOrders := []*Order{}
	for idx, order := range batch.Orders {
		if sm.isExists(order) {
			_logger.Debug("order at %v has already existed\n", idx)
			// fmt.Printf("order at %v has already existed\n", idx)
		} else {
			sm.addOrder(order)
			validOrders = append(validOrders, order)
		}
	}

	// If the batch is empty, stop
	if len(validOrders) == 0 {
		_logger.Debug("batch (%v) is empty\n", batch.BatchID)
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

	// _logger.Debug("Batch::%v is valid\n", batch.BatchID.String()[:6])

	// Send batch to smart contract
	sm.SendBatch(batch)

}

func (sm *SuperMatcher) CheckValidBatch(batch *Batch) bool {
	// 1. valid owner: owner is a matcher & the signature is valid
	if !sm.isMatcher(batch.Owner) {
		_logger.Debug("Invalid Matcher\n")
		return false
	}

	if !batch.IsValidSignature() {
		_logger.Debug("Invalid Batch's Signature\n")
		return false
	}

	// 2. check signatures of orders in the batch
	for idx, order := range batch.Orders {
		if !order.IsValidSignature() {
			_logger.Debug("Invalid Order's Signature at %v \n", idx)
			return false
		}
	}

	return true
}

func (sm *SuperMatcher) processing() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		sm.Process()
	}
}

func (sm *SuperMatcher) AddBatch(batch *Batch) {
	if sm.CheckValidBatch(batch) {
		sm.Mutex.Lock()
		defer sm.Mutex.Unlock()
		sm.Batches = append(sm.Batches, batch)
	}
}
