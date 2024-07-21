package user

import (
	"sync"

	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/NguyenHiu/lightning-exchange/orderApp"
	"github.com/NguyenHiu/lightning-exchange/orderClient"
	"github.com/NguyenHiu/lightning-exchange/tradeApp"
	"github.com/NguyenHiu/lightning-exchange/tradeClient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wire"

	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
)

var _logger = logger.NewLogger("User", logger.Magenta, logger.None)

type Connection struct {
	OrderAppClient *orderClient.OrderAppClient
	OrderChannel   *orderClient.OrderChannel
	TradeAppClient *tradeClient.TradeAppClient
	TradeChannel   *tradeClient.TradeChannel
	IsBlocked      bool
	Mux            sync.Mutex
}

func NewConnection() *Connection {
	return &Connection{
		OrderAppClient: &orderClient.OrderAppClient{},
		OrderChannel:   &orderClient.OrderChannel{},
		TradeAppClient: &tradeClient.TradeAppClient{},
		TradeChannel:   &tradeClient.TradeChannel{},
		IsBlocked:      false,
		Mux:            sync.Mutex{},
	}
}

type User struct {
	ID          uuid.UUID
	PrivateKey  string
	Address     common.Address
	Connections map[uuid.UUID]*Connection
}

func NewUser(privateKey string) *User {
	_uuid, _ := uuid.NewRandom()
	privKey, _ := crypto.HexToECDSA(privateKey)

	return &User{
		ID:          _uuid,
		PrivateKey:  privateKey,
		Address:     crypto.PubkeyToAddress(privKey.PublicKey),
		Connections: make(map[uuid.UUID]*Connection),
	}
}

func (u *User) SetupClient(
	matcherID uuid.UUID,
	busOrder wire.Bus,
	busTrade wire.Bus,
	nodeURL string,
	adjudicator common.Address,
	assets []ethwallet.Address,
	_orderApp *orderApp.OrderApp,
	_tradeApp *tradeApp.TradeApp,
	stakes []channel.Bal,
	emptyStake []channel.Bal,
	gavinAddr common.Address,
) {
	_prvKey, _ := crypto.HexToECDSA(u.PrivateKey)
	u.Connections[matcherID] = NewConnection()
	u.Connections[matcherID].OrderAppClient = orderClient.SetupClient(busOrder, nodeURL, adjudicator, assets, _prvKey, _orderApp, emptyStake, gavinAddr)
	u.Connections[matcherID].TradeAppClient = tradeClient.SetupClient(busTrade, nodeURL, adjudicator, assets, _prvKey, _tradeApp, stakes, gavinAddr)
}

func (u *User) AcceptedChannel(matcherID uuid.UUID) {
	u.Connections[matcherID].OrderChannel = u.Connections[matcherID].OrderAppClient.AcceptedChannel()
	u.Connections[matcherID].TradeChannel = u.Connections[matcherID].TradeAppClient.AcceptedChannel()
	u.Connections[matcherID].IsBlocked = false
}

func (u *User) AcceptedChannelAll() {
	for _, v := range u.Connections {
		v.OrderChannel = v.OrderAppClient.AcceptedChannel()
		v.TradeChannel = v.TradeAppClient.AcceptedChannel()
		v.IsBlocked = false
	}
}

func (u *User) SendNewOrders(matcherID uuid.UUID, newOrders []*orderApp.Order) {
	_logger.Info("[%v] Sending new ORDER...\n", u.Address.String()[:5])
	u.Connections[matcherID].OrderChannel.SendNewOrders(newOrders)
}

func (u *User) Settle(matcherID uuid.UUID) {
	u.Connections[matcherID].TradeChannel.Settle()
	u.Connections[matcherID].OrderChannel.Settle()
}

func (u *User) SettleAll() {
	for _, v := range u.Connections {
		v.TradeChannel.Settle()
		v.OrderChannel.Settle()
	}
}

func (u *User) Shutdown(matcherID uuid.UUID) {
	u.Connections[matcherID].OrderAppClient.Shutdown()
	u.Connections[matcherID].TradeAppClient.Shutdown()
}

func (u *User) ShutdownAll() {
	for _, v := range u.Connections {
		v.OrderAppClient.Shutdown()
		v.TradeAppClient.Shutdown()
	}
}

// func (u *User) prepareNonceAndGasPrice(value float64, gasLimit int) {
// 	nodeClient, _ := ethclient.Dial(constants.CHAIN_URL)

// 	nonce, err := nodeClient.PendingNonceAt(context.Background(), u.Address)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	gasPrice, err := nodeClient.SuggestGasPrice(context.Background())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	m.Auth.Nonce = big.NewInt(int64(nonce))
// 	m.Auth.GasPrice = gasPrice
// 	m.Auth.Value = orderClient.EthToWei(big.NewFloat(float64(value)))
// 	m.Auth.GasLimit = uint64(gasLimit)
// }
