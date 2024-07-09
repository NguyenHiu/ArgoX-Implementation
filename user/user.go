package user

import (
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

type User struct {
	ID             uuid.UUID
	PrivateKey     string
	OrderAppClient *orderClient.OrderAppClient
	OrderChannel   *orderClient.OrderChannel
	TradeAppClient *tradeClient.TradeAppClient
	TradeChannel   *tradeClient.TradeChannel
	Address        common.Address
}

func NewUser(privateKey string) *User {
	uuid, _ := uuid.NewRandom()
	privKey, _ := crypto.HexToECDSA(privateKey)

	return &User{
		ID:             uuid,
		PrivateKey:     privateKey,
		OrderAppClient: &orderClient.OrderAppClient{},
		OrderChannel:   &orderClient.OrderChannel{},
		TradeAppClient: &tradeClient.TradeAppClient{},
		TradeChannel:   &tradeClient.TradeChannel{},
		Address:        crypto.PubkeyToAddress(privKey.PublicKey),
	}
}

func (u *User) SetupClient(
	busOrder wire.Bus,
	busTrade wire.Bus,
	nodeURL string,
	adjudicator common.Address,
	assets []ethwallet.Address,
	_orderApp *orderApp.OrderApp,
	_tradeApp *tradeApp.TradeApp,
	stakes []channel.Bal,
	gavinAddr common.Address,
) {
	_prvKey, _ := crypto.HexToECDSA(u.PrivateKey)
	u.OrderAppClient = orderClient.SetupClient(busOrder, nodeURL, adjudicator, assets, _prvKey, _orderApp, stakes, gavinAddr)
	u.TradeAppClient = tradeClient.SetupClient(busTrade, nodeURL, adjudicator, assets, _prvKey, _tradeApp, stakes, gavinAddr)
}

func (u *User) AcceptedChannel() {
	u.OrderChannel = u.OrderAppClient.AcceptedChannel()
	u.TradeChannel = u.TradeAppClient.AcceptedChannel()
}

func (u *User) SendNewOrders(newOrders []*orderApp.Order) {
	_logger.Info("[%v] Sending new ORDER...\n", u.Address.String()[:5])
	u.OrderChannel.SendNewOrders(newOrders)
}

func (u *User) Settle() {
	// u.OrderChannel.Settle()
	u.TradeChannel.Settle()
}

func (u *User) Shutdown() {
	u.OrderAppClient.Shutdown()
	u.TradeAppClient.Shutdown()
}
