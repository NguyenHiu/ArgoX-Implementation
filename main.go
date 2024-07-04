package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	TOKEN "github.com/NguyenHiu/lightning-exchange/contracts/generated/token"
	"github.com/NguyenHiu/lightning-exchange/data"
	"github.com/NguyenHiu/lightning-exchange/deploy"
	"github.com/NguyenHiu/lightning-exchange/listener"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/NguyenHiu/lightning-exchange/matcher"
	"github.com/NguyenHiu/lightning-exchange/reporter"
	"github.com/NguyenHiu/lightning-exchange/supermatcher"
	"github.com/NguyenHiu/lightning-exchange/user"
	"github.com/NguyenHiu/lightning-exchange/worker"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var _logger = logger.NewLogger("Main", logger.Red, logger.Bold)

func getContracts() (common.Address, common.Address, common.Address, []common.Address, common.Address) {
	token, err := data.Get("token")
	if err != nil {
		log.Fatal(err)
	}

	onchain, err := data.Get("onchain")
	if err != nil {
		log.Fatal(err)
	}

	adj, err := data.Get("adj")
	if err != nil {
		log.Fatal(err)
	}

	ethHolder, err := data.Get("ethholder")
	if err != nil {
		log.Fatal(err)
	}
	gavHolder, err := data.Get("gvnholder")
	if err != nil {
		log.Fatal(err)
	}
	assetHolders := []common.Address{ethHolder, gavHolder}

	appAddr, err := data.Get("appaddr")
	if err != nil {
		log.Fatal(err)
	}

	return token, onchain, adj, assetHolders, appAddr
}

func StartSuperMatcher(onchainAddr common.Address, client *ethclient.Client, privateKeyHex string, port int) {
	onchainInstance, err := onchain.NewOnchain(onchainAddr, client)
	if err != nil {
		log.Fatal(err)
	}

	sm, err := supermatcher.NewSuperMatcher(onchainInstance, privateKeyHex, port, constants.CHAIN_ID)
	if err != nil {
		log.Fatal(err)
	}

	go sm.SetupHTTPServer()

	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		sm.Process()
	}
}

// TODO: Simulation
func main() {
	// Deploy contracts
	deploy.DeployContracts()
	token, onchain, adj, assetHolders, appAddr := getContracts()

	// Listen events from onchain contract
	go listener.StartListener(onchain)

	// Create a Reporter
	rp, err := reporter.NewReporter(onchain, constants.KEY_REPORTER, constants.CHAIN_ID)
	if err != nil {
		_logger.Error("Create reporter error, err: %v\n", err)
	}
	// Start the reporter
	rp.Listening()
	rp.Reporting()

	clientNode, err := ethclient.Dial(constants.CHAIN_URL)
	if err != nil {
		log.Fatal(err)
	}

	worker := worker.NewWorker(onchain, constants.KEY_WORKER, clientNode)
	worker.Listening()

	// Start Super Matcher
	go StartSuperMatcher(onchain, clientNode, constants.KEY_SUPER_MATCHER, constants.SUPER_MATCHER_PORT)

	superMatcherURI := fmt.Sprintf("http://127.0.0.1:%v", constants.SUPER_MATCHER_PORT)

	// Init matchers
	matcher1 := matcher.NewMatcher(assetHolders, adj, appAddr, onchain, constants.KEY_MATCHER_1, superMatcherURI, clientNode, constants.CHAIN_ID, token)
	matcher1.Register()

	// matcher2 := matcher.NewMatcher(assetHolders, adj, appAddr, onchain, constants.KEY_MATCHER_2, superMatcherURI, clientNode, constants.CHAIN_ID, token)
	// matcher2.Register()

	// Init Bob
	bob := user.NewUser(constants.KEY_BOB)
	busBob := matcher1.SetupClient(bob.ID)
	bob.SetupClient(busBob, constants.CHAIN_URL, matcher1.Adjudicator, matcher1.AssetHolders, matcher1.App, matcher1.Stakes, token)
	if ok := matcher1.OpenAppChannel(bob.ID, bob.AppClient.WireAddress()); !ok {
		log.Fatalln("OpenAppChannel Failed")
	}
	bob.AcceptedChannel()

	// Init Alice
	alice := user.NewUser(constants.KEY_ALICE)
	busAlice := matcher1.SetupClient(alice.ID)
	alice.SetupClient(busAlice, constants.CHAIN_URL, matcher1.Adjudicator, matcher1.AssetHolders, matcher1.App, matcher1.Stakes, token)
	if ok := matcher1.OpenAppChannel(alice.ID, alice.AppClient.WireAddress()); !ok {
		log.Fatalln("OpenAppChannel Failed")
	}
	alice.AcceptedChannel()

	// Ask orders
	askOrders := []*app.Order{}
	for price := 5; price < 10; price++ {
		order := app.NewOrder(big.NewInt(int64(price)), big.NewInt(5), constants.ASK, bob.AppClient.EthWalletAddress())
		if err := order.Sign(bob.PrivateKey); err != nil {
			_logger.Error("Sign order got error, err: %v\n", err)
		}
		askOrders = append(askOrders, order)

	}

	// Bid orders
	bidOrders := []*app.Order{}
	for price := 1; price < 6; price++ {
		order := app.NewOrder(big.NewInt(int64(price)), big.NewInt(5), constants.BID, alice.AppClient.EthWalletAddress())
		if err := order.Sign(alice.PrivateKey); err != nil {
			_logger.Error("Sign order got error err: %v\n", err)
		}
		bidOrders = append(bidOrders, order)
	}

	// Send orders
	bob.SendNewOrders(askOrders)
	<-time.After(time.Second * 10)
	alice.SendNewOrders(bidOrders)

	{
		tokenInstance, err := TOKEN.NewToken(token, clientNode)
		if err != nil {
			log.Fatal(err)
		}
		matcherBalance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, matcher1.Address)
		if err != nil {
			log.Fatal(err)
		}
		bobBalance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, bob.AppClient.WalletAddress())
		if err != nil {
			log.Fatal(err)
		}
		aliceBalance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, alice.AppClient.WalletAddress())
		if err != nil {
			log.Fatal(err)
		}
		_logger.Info("matcher's balance: %v\n", matcherBalance)
		_logger.Info("bob's balance: %v\n", bobBalance)
		_logger.Info("alice's balance: %v\n", aliceBalance)
	}

	<-time.After(time.Second * 10)

	{
		// Create Final Order
		order, err := app.EndOrder(constants.KEY_ALICE)
		if err != nil {
			_logger.Error("create an end order is fail, err: %v\n", err)
		}
		alice.SendNewOrders([]*app.Order{order})
	}

	{
		// Create Final Order
		order, err := app.EndOrder(constants.KEY_BOB)
		if err != nil {
			_logger.Error("create an end order is fail, err: %v\n", err)
		}
		bob.SendNewOrders([]*app.Order{order})
	}

	// Payout.
	_logger.Info("Settle\n")
	alice.Settle()
	matcher1.Settle(alice.ID)
	bob.Settle()
	matcher1.Settle(bob.ID)

	// Cleanup.
	_logger.Info("Shutdown\n")
	alice.Shutdown()
	matcher1.Shutdown(alice.ID)
	bob.Shutdown()
	matcher1.Shutdown(bob.ID)

	{
		tokenInstance, err := TOKEN.NewToken(token, clientNode)
		if err != nil {
			log.Fatal(err)
		}
		matcherBalance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, matcher1.Address)
		if err != nil {
			log.Fatal(err)
		}
		bobBalance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, bob.AppClient.WalletAddress())
		if err != nil {
			log.Fatal(err)
		}
		aliceBalance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, alice.AppClient.WalletAddress())
		if err != nil {
			log.Fatal(err)
		}
		_logger.Info("matcher's balance: %v\n", matcherBalance)
		_logger.Info("bob's balance: %v\n", bobBalance)
		_logger.Info("alice's balance: %v\n", aliceBalance)
	}

	<-time.After(time.Second * 100)
}
