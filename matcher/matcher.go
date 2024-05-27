package matcher

import (
	"fmt"
	"math/big"

	"github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/client"
	"github.com/NguyenHiu/lightning-exchange/constants"
	utils "github.com/NguyenHiu/lightning-exchange/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"perun.network/go-perun/backend/ethereum/wallet"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/wire"
)

type ClientConfig struct {
	AppClient     *client.AppClient
	VerifyChannel *client.VerifyChannel
}

type Matcher struct {
	ID            uuid.UUID
	ClientConfigs map[uuid.UUID]*ClientConfig
	Adjudicator   common.Address
	AssetHolders  []wallet.Address
	App           *app.VerifyApp
	Stakes        []*big.Int
}

func NewMatcher(chainURL string, chain int64, privateKey string) *Matcher {
	id, _ := uuid.NewRandom()
	fmt.Println("Deploying Smart Contracts...")
	adj, assetHolders, appAddr := utils.DeployContracts(constants.CHAIN_URL, constants.CHAIN_ID, constants.KEY_DEPLOYER)
	verifierApp := app.NewVerifyApp(ethwallet.AsWalletAddr(appAddr))
	stakeETH := client.EthToWei(big.NewFloat(5))
	stakeGVN := big.NewInt(5)
	ethwalletAssetHolders := []ethwallet.Address{}
	for _, asset := range assetHolders {
		ethwalletAssetHolders = append(ethwalletAssetHolders, *ethwallet.AsWalletAddr(asset))
	}
	return &Matcher{
		ID:            id,
		ClientConfigs: make(map[uuid.UUID]*ClientConfig),
		Adjudicator:   adj,
		AssetHolders:  ethwalletAssetHolders,
		App:           verifierApp,
		Stakes:        []*big.Int{stakeETH, stakeGVN},
	}
}

func (m *Matcher) SetupClient(userID uuid.UUID) (wire.Bus, common.Address, []wallet.Address, *app.VerifyApp, []*big.Int) {
	bus := wire.NewLocalBus()
	appClient := utils.SetupClient(bus, constants.CHAIN_URL, m.Adjudicator, m.AssetHolders, constants.KEY_MATCHER, m.App, m.Stakes)
	m.ClientConfigs[userID] = &ClientConfig{
		AppClient:     appClient,
		VerifyChannel: &client.VerifyChannel{},
	}
	return bus, m.Adjudicator, m.AssetHolders, m.App, m.Stakes
}

func (m *Matcher) OpenAppChannel(userID uuid.UUID, userPeer wire.Address) bool {
	user, ok := m.ClientConfigs[userID]
	if !ok {
		return false
	}
	m.ClientConfigs[userID].VerifyChannel = user.AppClient.OpenAppChannel(userPeer)
	return true
}

func (m *Matcher) Settle(userID uuid.UUID) {
	m.ClientConfigs[userID].VerifyChannel.Settle()
}

func (m *Matcher) Shutdown(userID uuid.UUID) {
	m.ClientConfigs[userID].AppClient.Shutdown()
}
