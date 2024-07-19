package supermatcher

import (
	"math/big"
	"sync"

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
	Orders          map[uuid.UUID][]*ExpandOrder
	MatchedOrders   map[uuid.UUID]*big.Int
	Mutex           sync.Mutex
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
		Orders:          make(map[uuid.UUID][]*ExpandOrder),
		MatchedOrders:   make(map[uuid.UUID]*big.Int),
	}

	return sm, nil
}

func (sm *SuperMatcher) isExists(order *ExpandOrder) bool {
	expandOrders, ok := sm.Orders[order.OriginalOrder.OrderID]
	if !ok {
		return false
	}

	for _, expandOrder := range expandOrders {
		if expandOrder.Equal(order) {
			return true
		}
	}
	return false
}

func (sm *SuperMatcher) addOrder(order *ExpandOrder) {
	sm.Orders[order.OriginalOrder.OrderID] = append(sm.Orders[order.OriginalOrder.OrderID], &ExpandOrder{
		ShadowOrder:   order.ShadowOrder.Clone(),
		Trades:        order.Trades,
		OriginalOrder: order.OriginalOrder,
	})
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

	// _logger.Debug("Batch::%v is valid\n", batch.BatchID.String()[:6])

	// Send batch to smart contract
	sm.SendBatch(batch)
}

func (sm *SuperMatcher) CheckValidBatch(batch *Batch) bool {
	// 1. valid owner: owner is a matcher & the signature is valid
	if !sm.isMatcher(batch.Owner) {
		_logger.Debug("Batch::%v (Invalid Matcher)\n", batch.BatchID)
		return false
	}

	if !batch.IsValidSignature() {
		_logger.Debug("Batch::%v (Invalid Batch's Signature)\n", batch.BatchID)
		return false
	}

	// 2. check signatures of orders in the batch
	for idx, order := range batch.Orders {
		if !order.IsValidOrder(batch.Owner) {
			_logger.Debug("Batch::%v (Invalid Order at %v) \n", batch.BatchID, idx)
			return false
		}
	}

	return true
}

// Return values:
//   - "OK"
//   - "REMOVE"
//   - "RESIGN"
//   - "INVALID_BATCH"
func (sm *SuperMatcher) AddBatch(batch *Batch) (string, []*ExpandOrder) {
	if sm.CheckValidBatch(batch) {
		sm.Mutex.Lock()
		defer sm.Mutex.Unlock()
		_logger.Info("Get valid batch::%v\n", batch.BatchID.String())

		// Filter orders in the batch
		validOrders := []*ExpandOrder{}
		for idx, order := range batch.Orders {
			if sm.isExists(order) {
				_logger.Debug("Order::%v at %v has already existed (total: %v)\n", order.OriginalOrder.OrderID.String(), idx, len(batch.Orders))
			} else {
				sm.addOrder(order)
				validOrders = append(validOrders, order)
			}
		}

		// If the batch is empty, stop
		if len(validOrders) == 0 {
			_logger.Debug("Batch (%v) is empty\n", batch.BatchID)
			return "REMOVE", nil
		}

		// Update `amount` & `orders` of the batch
		if len(validOrders) != len(batch.Orders) {
			return "RESIGN", validOrders
		}

		sm.Batches = append(sm.Batches, batch)
		return "OK", nil
	}

	return "INVALID_BATCH", nil
}

func (sm *SuperMatcher) GetLeftAmount(id uuid.UUID) *big.Int {
	_leftAmount, ok := sm.MatchedOrders[id]
	if ok {
		return new(big.Int).Set(_leftAmount)
	}
	return big.NewInt(-1)
}

func (sm *SuperMatcher) MatchAnOrder(id uuid.UUID, leftAmount *big.Int) bool {
	_leftAmount, ok := sm.MatchedOrders[id]
	if !ok {
		sm.MatchedOrders[id] = new(big.Int).Set(leftAmount)
		return true
	}

	if leftAmount.Cmp(_leftAmount) != -1 {
		return false
	}
	sm.MatchedOrders[id] = new(big.Int).Set(leftAmount)
	return true
}
