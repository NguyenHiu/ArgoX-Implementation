package user

import (
	"github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/client"
	utils "github.com/NguyenHiu/lightning-exchange/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/wire"
)

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
) {
	u.AppClient = utils.SetupClient(bus, nodeURL, adjudicator, assets, u.PrivateKey, app, stakes)
}

func (u *User) AcceptedChannel() {
	u.VerifyChannel = u.AppClient.AcceptedChannel()
}

func (u *User) SendNewOrder(newOrder *app.Order) {
	u.VerifyChannel.SendNewOrder(newOrder)
}

func (u *User) Settle() {
	u.VerifyChannel.Settle()
}

func (u *User) Shutdown() {
	u.AppClient.Shutdown()
}