package client

import (
	"context"
	"fmt"

	"github.com/NguyenHiu/lightning-exchange/app"
	App "github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/watcher/local"
	"perun.network/go-perun/wire"

	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	swallet "perun.network/go-perun/backend/ethereum/wallet/simple"
)

type AppClient struct {
	perunClient    *client.Client
	account        wallet.Address
	currencies     []channel.Asset
	stakes         []channel.Bal
	app            *app.VerifyApp
	channels       chan *VerifyChannel
	UseTrigger     bool
	TriggerChannel chan *App.Order
}

// SetupAppClient creates a new app client.
func SetupAppClient(
	bus wire.Bus, // bus is used of off-chain communication.
	w *swallet.Wallet, // w is the wallet used for signing transactions.
	acc common.Address, // acc is the address of the account to be used for signing transactions.
	nodeURL string, // nodeURL is the URL of the blockchain node.
	chainID uint64, // chainID is the identifier of the blockchain.
	adjudicator common.Address, // adjudicator is the address of the adjudicator.
	assets []ethwallet.Address, // assets are the address of the asset holder for our app channels.
	app *app.VerifyApp, // app is the channel app we want to set up the client with.
	stakes []channel.Bal, // stake is the balance the client is willing to fund the channel with.
	useTrigger bool,
) (*AppClient, error) {
	// Create Ethereum client and contract backend.
	cb, err := CreateContractBackend(nodeURL, chainID, w)
	if err != nil {
		return nil, fmt.Errorf("creating contract backend: %w", err)
	}

	// Validate contracts.
	err = ethchannel.ValidateAdjudicator(context.TODO(), cb, adjudicator)
	if err != nil {
		return nil, fmt.Errorf("validating adjudicator: %w", err)
	}
	err = ethchannel.ValidateAssetHolderETH(context.TODO(), cb, common.Address(assets[constants.ETH]), adjudicator)
	if err != nil {
		return nil, fmt.Errorf("validating asset holder: %w", err)
	}
	err = ethchannel.ValidateAssetHolderERC20(context.TODO(), cb, common.Address(assets[constants.GVN]), adjudicator, common.HexToAddress(constants.GAVIN_TOKEN_ADDRESS))
	if err != nil {
		return nil, fmt.Errorf("validating asset holder erc20: %w", err)
	}

	// Setup funder.
	funder := ethchannel.NewFunder(cb)
	dep := ethchannel.NewETHDepositor()
	ethAcc := accounts.Account{Address: acc}
	funder.RegisterAsset(assets[constants.ETH], dep, ethAcc)
	depERC20 := ethchannel.NewERC20Depositor(common.HexToAddress(constants.GAVIN_TOKEN_ADDRESS))
	funder.RegisterAsset(assets[constants.GVN], depERC20, ethAcc)

	// Setup adjudicator.
	adj := ethchannel.NewAdjudicator(cb, adjudicator, acc, ethAcc)

	// Setup dispute watcher.
	watcher, err := local.NewWatcher(adj)
	if err != nil {
		return nil, fmt.Errorf("intializing watcher: %w", err)
	}

	// Setup Perun client.
	waddr := ethwallet.AsWalletAddr(acc)
	perunClient, err := client.New(waddr, bus, funder, adj, w, watcher)
	if err != nil {
		return nil, errors.WithMessage(err, "creating client")
	}

	assetsCoverted := []channel.Asset{}
	for _, asset := range assets {
		assetsCoverted = append(assetsCoverted, &asset)
	}

	// Create client and start request handler.
	c := &AppClient{
		perunClient:    perunClient,
		account:        waddr,
		currencies:     assetsCoverted,
		stakes:         stakes,
		app:            app,
		channels:       make(chan *VerifyChannel, 1),
		UseTrigger:     useTrigger,
		TriggerChannel: make(chan *App.Order),
	}

	channel.RegisterApp(app)
	go perunClient.Handle(c, c)

	return c, nil
}

func (c *AppClient) EthWalletAddress() *ethwallet.Address {
	return ethwallet.AsWalletAddr(common.Address(c.account.Bytes()))
}

// startWatching starts the dispute watcher for the specified channel.
func (c *AppClient) startWatching(ch *client.Channel) {
	go func() {
		err := ch.Watch(c)
		if err != nil {
			fmt.Printf("Watcher returned with error: %v", err)
		}
	}()
}

// OpenAppChannel opens a new app channel with the specified peer.
func (c *AppClient) OpenAppChannel(peer wire.Address) *VerifyChannel {
	participants := []wire.Address{c.account, peer}

	// We create an initial allocation which defines the starting balances.
	initAlloc := channel.NewAllocation(2, c.currencies...)
	initAlloc.SetAssetBalances(c.currencies[constants.ETH], []channel.Bal{
		c.stakes[constants.ETH],
		c.stakes[constants.ETH],
	})
	initAlloc.SetAssetBalances(c.currencies[constants.GVN], []channel.Bal{
		c.stakes[constants.GVN],
		c.stakes[constants.GVN],
	})

	// Prepare the channel proposal by defining the channel parameters.
	challengeDuration := uint64(10) // On-chain challenge duration in seconds.

	withApp := client.WithApp(c.app, c.app.InitData())

	proposal, err := client.NewLedgerChannelProposal(
		challengeDuration,
		c.account,
		initAlloc,
		participants,
		withApp,
	)
	if err != nil {
		panic(err)
	}

	// Send the app channel proposal.
	ch, err := c.perunClient.ProposeChannel(context.TODO(), proposal)
	if err != nil {
		panic(err)
	}

	// Start the on-chain event watcher. It automatically handles disputes.
	c.startWatching(ch)

	return newVerifyChannel(ch)
}

// AcceptedChannel returns the next accepted app channel.
func (c *AppClient) AcceptedChannel() *VerifyChannel {
	return <-c.channels
}

// Shutdown gracefully shuts down the client.
func (c *AppClient) Shutdown() {
	c.perunClient.Close()
}
