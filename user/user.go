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

var _logger = logger.NewLogger("User", logger.Magenta, logger.None)

type User struct {
	ID            uuid.UUID
	PrivateKey    string
	AppClient     *client.AppClient
	VerifyChannel *client.VerifyChannel
	Address       common.Address
}

func NewUser(privateKey string) *User {
	uuid, _ := uuid.NewRandom()
	privKey, _ := crypto.HexToECDSA(privateKey)

	return &User{
		ID:            uuid,
		PrivateKey:    privateKey,
		AppClient:     &client.AppClient{},
		VerifyChannel: &client.VerifyChannel{},
		Address:       crypto.PubkeyToAddress(privKey.PublicKey),
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

func (u *User) SendNewOrders(newOrders []*app.Order) {
	_logger.Info("[%v] Sending new ORDER...\n", u.Address.String()[:5])
	u.VerifyChannel.SendNewOrders(newOrders)
}

func (u *User) Settle() {
	u.VerifyChannel.Settle()
}

func (u *User) Shutdown() {
	u.AppClient.Shutdown()
}
