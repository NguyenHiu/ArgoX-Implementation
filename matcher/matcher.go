package matcher

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/data"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/NguyenHiu/lightning-exchange/orderApp"
	"github.com/NguyenHiu/lightning-exchange/orderClient"
	"github.com/NguyenHiu/lightning-exchange/supermatcher"
	"github.com/NguyenHiu/lightning-exchange/tradeApp"
	"github.com/NguyenHiu/lightning-exchange/tradeClient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"perun.network/go-perun/wire"

	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
)

var _logger = logger.NewLogger("Matcher", logger.Green, logger.None)

type ClientConfig struct {
	TradeAppClient *tradeClient.TradeAppClient
	TradeChannel   *tradeClient.TradeChannel
	OrderAppClient *orderClient.OrderAppClient
	OrderChannel   *orderClient.OrderChannel
}
type Matcher struct {
	// Perun's Data
	ID            uuid.UUID
	ClientConfigs map[uuid.UUID]*ClientConfig // store traders' channel
	Adjudicator   common.Address
	AssetHolders  []ethwallet.Address
	OrderApp      *orderApp.OrderApp
	TradeApp      *tradeApp.TradeApp
	Stakes        []*big.Int
	EmptyStakes   []*big.Int

	// Gavin Address
	GavinAddress common.Address

	// Order Book
	BidOrders         []*MatcherOrder
	AskOrders         []*MatcherOrder
	Orders            map[uuid.UUID]*tradeApp.Order
	ExecutedTrade     []*tradeApp.Trade
	mappingBidtoTrade map[uuid.UUID][]*tradeApp.Trade
	mappingAskToTrade map[uuid.UUID][]*tradeApp.Trade

	// Super Matcher & Onchain Contract
	OnchainInstance *onchain.Onchain   // Onchain Contract
	Auth            *bind.TransactOpts // authentication for writting to smart contract
	Client          *ethclient.Client  // for getting nonce & gas price

	Address    common.Address // address of an account used for interacting with onchain
	PrivateKey *ecdsa.PrivateKey

	Batches map[uuid.UUID]*Batch // be ready for providing a valid proof of any batch
	Mux     sync.Mutex

	SuperMatcherInstance *supermatcher.SuperMatcher

	// Store orders' data
	OrderStorage []*data.OrderData

	/* MATCHING ANALYSIS */
	CreateTime              map[uuid.UUID]int64
	TotalTimeLocal          int64
	TotalMatchedAmountLocal *big.Int
	NumberOfMatchedOrder    int64
	PriceCurveLocal         []*big.Int
	CurrentPrice            *big.Int
	IsGetPriceCurve         bool
	TotalProfitLocal        *big.Int
	/* MATCHING ANALYSIS */

}

func NewMatcher(
	assetHolders []common.Address,
	adj, appAddr, onchainAddr common.Address,
	privateKey string,
	clientNode *ethclient.Client,
	chainID int64,
	gavinAddress common.Address,
	supermatcherInstance *supermatcher.SuperMatcher,
) *Matcher {
	id, _ := uuid.NewRandom()
	_orderApp := orderApp.NewOrderApp(ethwallet.AsWalletAddr(appAddr))
	_tradeApp := tradeApp.NewTradeApp(ethwallet.AsWalletAddr(appAddr))
	stakeETH := big.NewInt(constants.NO_ETH_IN_CHANNEL)
	stakeGVN := big.NewInt(constants.NO_GVN_IN_CHANNEL)
	ethwalletAssetHolders := []ethwallet.Address{}
	for _, asset := range assetHolders {
		ethwalletAssetHolders = append(ethwalletAssetHolders, *ethwallet.AsWalletAddr(asset))
	}

	instance, err := onchain.NewOnchain(onchainAddr, clientNode)
	if err != nil {
		log.Fatal(err)
	}

	_privateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(_privateKey, big.NewInt(chainID))
	if err != nil {
		log.Fatal(err)
	}

	return &Matcher{
		ID:            id,
		ClientConfigs: make(map[uuid.UUID]*ClientConfig),
		Adjudicator:   adj,
		AssetHolders:  ethwalletAssetHolders,
		OrderApp:      _orderApp,
		TradeApp:      _tradeApp,
		Stakes:        []*big.Int{stakeETH, stakeGVN},
		EmptyStakes:   []*big.Int{new(big.Int), new(big.Int)},

		BidOrders:         []*MatcherOrder{},
		AskOrders:         []*MatcherOrder{},
		Orders:            make(map[uuid.UUID]*tradeApp.Order),
		ExecutedTrade:     []*tradeApp.Trade{},
		mappingBidtoTrade: make(map[uuid.UUID][]*tradeApp.Trade),
		mappingAskToTrade: make(map[uuid.UUID][]*tradeApp.Trade),

		OnchainInstance: instance,
		Auth:            auth,
		Client:          clientNode,

		Address:    crypto.PubkeyToAddress(_privateKey.PublicKey),
		PrivateKey: _privateKey,

		Batches: make(map[uuid.UUID]*Batch),

		GavinAddress:         gavinAddress,
		SuperMatcherInstance: supermatcherInstance,

		OrderStorage: make([]*data.OrderData, 0),

		/* MATCHING ANALYSIS */
		CreateTime:              make(map[uuid.UUID]int64),
		TotalTimeLocal:          0,
		TotalMatchedAmountLocal: new(big.Int),
		NumberOfMatchedOrder:    0,
		PriceCurveLocal:         []*big.Int{},
		CurrentPrice:            new(big.Int),
		IsGetPriceCurve:         false,
		TotalProfitLocal:        big.NewInt(0),
		/* MATCHING ANALYSIS */
	}
}

func (m *Matcher) NewTrade(bid, ask uuid.UUID, price, amount *big.Int) *tradeApp.Trade {

	id, _ := uuid.NewRandom()
	executedtrade := &tradeApp.Trade{
		TradeID:   id,
		BidOrder:  bid,
		AskOrder:  ask,
		Price:     price,
		Amount:    amount,
		Owner:     m.Address,
		Signature: []byte{},
	}

	executedtrade.Sign(m.PrivateKey)

	m.ExecutedTrade = append(m.ExecutedTrade, executedtrade)
	m.mappingBidtoTrade[bid] = append(m.mappingBidtoTrade[bid], executedtrade)
	m.mappingAskToTrade[ask] = append(m.mappingAskToTrade[ask], executedtrade)

	return executedtrade
}

func (m *Matcher) GetPriceCurve() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		if m.IsGetPriceCurve {
			m.PriceCurveLocal = append(m.PriceCurveLocal, m.CurrentPrice)
		}
	}
}

// Create 2 channels: one for receiving orders, one for sending message
func (m *Matcher) SetupClient(userID uuid.UUID) (wire.Bus, wire.Bus) {
	orderBus := wire.NewLocalBus()
	tradeBus := wire.NewLocalBus()
	orderAppClient := orderClient.SetupClient(orderBus, constants.CHAIN_URL, m.Adjudicator, m.AssetHolders, m.PrivateKey, m.OrderApp, m.EmptyStakes, m.GavinAddress)
	tradeAppClient := tradeClient.SetupClient(tradeBus, constants.CHAIN_URL, m.Adjudicator, m.AssetHolders, m.PrivateKey, m.TradeApp, m.Stakes, m.GavinAddress)
	m.ClientConfigs[userID] = &ClientConfig{
		TradeAppClient: tradeAppClient,
		TradeChannel:   &tradeClient.TradeChannel{},
		OrderAppClient: orderAppClient,
		OrderChannel:   &orderClient.OrderChannel{},
	}
	return orderBus, tradeBus
}

func (m *Matcher) OpenAppChannel(userID uuid.UUID, userPeer wire.Address) bool {
	user, ok := m.ClientConfigs[userID]
	if !ok {
		return false
	}
	m.ClientConfigs[userID].OrderChannel = user.OrderAppClient.OpenAppChannel(userPeer)
	m.ClientConfigs[userID].TradeChannel = user.TradeAppClient.OpenAppChannel(userPeer)
	go m.receiveOrder(userID)
	// go m.goBatching()
	return true
}

// func (m *Matcher) goBatching() {
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		batches := m.batching()
// 		for _, batch := range batches {
// 			batch.Sign(m.PrivateKey)
// 			m.SendBatch(batch)
// 		}
// 	}
// }

func (m *Matcher) receiveOrder(userID uuid.UUID) {
	for orders := range m.ClientConfigs[userID].OrderAppClient.TriggerChannel {
		for _, order := range orders {
			if _, ok := m.Orders[order.OrderID]; ok {
				log.Fatal("receiving the same order\n")
			}

			if order.OrderID == tradeApp.EndID {
				endTrade, _ := tradeApp.EndTrade(m.PrivateKey)
				m.ClientConfigs[userID].TradeChannel.SendNewTrades([]*tradeApp.Trade{endTrade}, nil, nil, false)
				continue
			}

			_side := "bid"
			if order.Side == constants.ASK {
				_side = "ask"
			}
			_logger.Info("[%v::%v] Receive an order::%v::%v, price: %v, amount: %v, %v\n", m.ID.String()[:5], m.Address.String()[:5], order.OrderID.String()[:6], order.Owner.String()[:5], order.Price, order.Amount, _side)

			m.CreateTime[order.OrderID] = time.Now().Unix()

			_order := order.Clone()
			m.Orders[_order.OrderID] = &tradeApp.Order{
				OrderID:   _order.OrderID,
				Price:     _order.Price,
				Amount:    _order.Amount,
				Side:      _order.Side,
				Owner:     _order.Owner,
				Signature: _order.Signature,
			}

			__order := order.Clone()
			_newOrder := &MatcherOrder{
				Data: &ShadowOrder{
					Price:  __order.Price,
					Amount: __order.Amount,
					Side:   __order.Side,
					From:   __order.OrderID,
				},
				Owner: userID,
			}

			if _newOrder.Data.Amount.Cmp(order.Amount) != 0 {
				_logger.Debug("Damn, it;s fking wrong\n")
				log.Fatal("SUPER ERROR")
			}

			if _newOrder.Data.From.String() != order.OrderID.String() {
				_logger.Debug("Damn, it;s fking wrong (ID) \n")
				log.Fatal("SUPER ERROR")
			}

			m.addOrder(_newOrder)

			// Used to export orders to file
			m.OrderStorage = append(m.OrderStorage, &data.OrderData{
				Price:  int(__order.Price.Int64()),
				Amount: int(__order.Amount.Int64()),
				Side:   __order.Side,
			})
		}
	}
}

func (m *Matcher) Settle(userID uuid.UUID) {
	if _config, _ok := m.ClientConfigs[userID]; _ok {
		_config.TradeChannel.Settle()
		_config.OrderChannel.Settle()
	}
}

func (m *Matcher) Shutdown(userID uuid.UUID) {
	if _config, _ok := m.ClientConfigs[userID]; _ok {
		_config.TradeAppClient.Shutdown()
		_config.OrderAppClient.Shutdown()
	}
}

func (m *Matcher) ExportOrdersData(filename string) error {
	return data.SaveOrders(m.OrderStorage, filename)
}
