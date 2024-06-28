package user

import (
	"github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/client"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/NguyenHiu/lightning-exchange/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wire"
)

var _logger = logger.NewLogger("User", logger.None, logger.None)

type User struct {
	ID            uuid.UUID
	PrivateKey    string
	AppClient     *client.AppClient
	VerifyChannel *client.VerifyChannel
}

func NewUser(privateKey string) *User {
	uuid, _ := uuid.NewRandom()
	return &User{
		ID:            uuid,
		PrivateKey:    privateKey,
		AppClient:     &client.AppClient{},
		VerifyChannel: &client.VerifyChannel{},
	}
}

func (u *User) SetupClient(
	bus wire.Bus,
	nodeURL string,
	adjudicator common.Address,
	assets []ethwallet.Address,
	app *app.VerifyApp,
	stakes []channel.Bal,
	gavinAddr common.Address,
) {
	_prvKey, _ := crypto.HexToECDSA(u.PrivateKey)
	u.AppClient = util.SetupClient(bus, nodeURL, adjudicator, assets, _prvKey, app, stakes, false, gavinAddr)
}

func (u *User) AcceptedChannel() {
	u.VerifyChannel = u.AppClient.AcceptedChannel()
}

func (u *User) SendNewOrder(newOrder *app.Order) {
	_logger.Info("Sending new ORDER...\n")
	u.VerifyChannel.SendNewOrder(newOrder)
}

func (u *User) Settle() {
	u.VerifyChannel.Settle()
}

func (u *User) Shutdown() {
	u.AppClient.Shutdown()
}
