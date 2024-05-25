package main

import (
	"fmt"
	"log"
	"math/big"

	App "github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/client"
	"github.com/ethereum/go-ethereum/crypto"
	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/wire"
)

const (
	chainURL = "ws://127.0.0.1:8545"
	chainID  = 1337

	// Private keys.
	keyDeployer = "79ea8f62d97bc0591a4224c1725fca6b00de5b2cea286fe2e0bb35c5e76be46e"
	keyAlice    = "1af2e950272dd403de7a5760d41c6e44d92b6d02797e51810795ff03cc2cda4f"
	keyBob      = "f63d7d8e930bccd74e93cf5662fde2c28fd8be95edb70c73f1bdd863d07f412e"
)

func main() {
	// Deploy contracts.
	log.Println("Deploying contracts.")
	adjudicator, assetHolder, appAddress := deployContracts(chainURL, chainID, keyDeployer)
	asset := *ethwallet.AsWalletAddr(assetHolder)
	app := App.NewVerifyApp(ethwallet.AsWalletAddr(appAddress))

	// Setup clients.
	log.Println("Setting up clients.")
	bus := wire.NewLocalBus() // Message bus used for off-chain communication.
	stake := client.EthToWei(big.NewFloat(5))
	alice := SetupClient(bus, chainURL, adjudicator, asset, keyAlice, app, stake)
	bob := SetupClient(bus, chainURL, adjudicator, asset, keyBob, app, stake)

	// Print balances before transactions.
	l := newBalanceLogger(chainURL)
	l.LogBalances(alice, bob)

	// Open app channel and play.
	log.Println("Opening channel.")
	appAlice := alice.OpenAppChannel(bob.WireAddress())
	appBob := bob.AcceptedChannel()

	newOrder := App.NewOrder(client.EthToWei(big.NewFloat(5)).Int64(), 5, App.BID, alice.WalletAddressAsEthwallet(), "P")
	alicePrvKey, err := crypto.HexToECDSA(keyAlice)
	if err != nil {
		panic(err)
	}
	newOrder.Sign(*alicePrvKey)
	appAlice.SendNewOrder(&newOrder)

	lastOrder := App.NewOrder(0, 0, App.BID, alice.WalletAddressAsEthwallet(), "F")
	lastOrder.Sign(*alicePrvKey)
	appAlice.SendNewOrder(&lastOrder)

	// Payout.
	fmt.Println("Settle")
	appAlice.Settle()
	appBob.Settle()

	// Print balances after transactions.
	fmt.Println("LogBalances")
	l.LogBalances(alice, bob)

	// Cleanup.
	fmt.Println("Shutdown")
	alice.Shutdown()
	bob.Shutdown()
}
