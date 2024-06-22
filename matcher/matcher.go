package matcher

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/client"
	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/NguyenHiu/lightning-exchange/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"perun.network/go-perun/backend/ethereum/wallet"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/wire"
)

var _logger = logger.NewLogger("Matcher", logger.None, logger.None)

type ClientConfig struct {
	AppClient     *client.AppClient
	VerifyChannel *client.VerifyChannel
}

type MatcherOrder struct {
	Data  *app.Order
	Owner uuid.UUID
}

type Matcher struct {
	// Perun's Data
	ID            uuid.UUID
	ClientConfigs map[uuid.UUID]*ClientConfig // store traders' channel
	Adjudicator   common.Address
	AssetHolders  []wallet.Address
	App           *app.VerifyApp
	Stakes        []*big.Int

	// Gavin Address
	GavinAddress common.Address

	// Order Book
	BidOrders []*MatcherOrder
	AskOrders []*MatcherOrder

	// Super Matcher & Onchain Contract
	SuperMatcherURI string             // Super Matcher API Server
	OnchainInstance *onchain.Onchain   // Onchain Contract
	Auth            *bind.TransactOpts // authentication for writting to smart contract
	Client          *ethclient.Client  // for getting nonce & gas price

	Address    common.Address // address of an account used for interacting with onchain
	PrivateKey *ecdsa.PrivateKey

	Batches map[uuid.UUID]*Batch // be ready for providing a valid proof of any batch
	Mux     sync.Mutex
}

func NewMatcher(
	assetHolders []common.Address,
	adj, appAddr, onchainAddr common.Address,
	privateKey, superMatcherURI string,
	clientNode *ethclient.Client,
	chainID int64,
	gavinAddress common.Address,
) *Matcher {
	id, _ := uuid.NewRandom()
	verifierApp := app.NewVerifyApp(ethwallet.AsWalletAddr(appAddr))
	stakeETH := client.EthToWei(big.NewFloat(5))
	stakeGVN := big.NewInt(5)
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
		App:           verifierApp,
		Stakes:        []*big.Int{stakeETH, stakeGVN},

		BidOrders: []*MatcherOrder{},
		AskOrders: []*MatcherOrder{},

		SuperMatcherURI: superMatcherURI,
		OnchainInstance: instance,
		Auth:            auth,
		Client:          clientNode,

		Address:    crypto.PubkeyToAddress(_privateKey.PublicKey),
		PrivateKey: _privateKey,

		Batches: make(map[uuid.UUID]*Batch),

		GavinAddress: gavinAddress,
	}
}

func (m *Matcher) SetupClient(userID uuid.UUID) (wire.Bus, common.Address, []wallet.Address, *app.VerifyApp, []*big.Int) {
	bus := wire.NewLocalBus()
	appClient := util.SetupClient(bus, constants.CHAIN_URL, m.Adjudicator, m.AssetHolders, m.PrivateKey, m.App, m.Stakes, true, m.GavinAddress)
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
	go m.receiveOrder(userID)
	go m.goBatching()
	return true
}

func (m *Matcher) goBatching() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		batches := m.batching()
		for _, batch := range batches {
			batch.Sign(m.PrivateKey)
			m.SendBatch(batch)
		}
	}
}

func (m *Matcher) receiveOrder(userID uuid.UUID) {
	for order := range m.ClientConfigs[userID].AppClient.TriggerChannel {
		_logger.Info("[%v] Receive an order: {%v, %v, %v, %v}\n", m.ID.String()[:6], order.Price, order.Amount, order.Side, order.Owner)

		m.addOrder(&MatcherOrder{
			Data:  order,
			Owner: userID,
		})
	}
}

func (m *Matcher) Settle(userID uuid.UUID) {
	m.ClientConfigs[userID].VerifyChannel.Settle()
}

func (m *Matcher) Shutdown(userID uuid.UUID) {
	m.ClientConfigs[userID].AppClient.Shutdown()
}
