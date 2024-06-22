package main

import (
	"fmt"
	"log"
	"math/big"
	"time"

	App "github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/client"
	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/data"
	"github.com/NguyenHiu/lightning-exchange/deploy"
	"github.com/NguyenHiu/lightning-exchange/listener"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/NguyenHiu/lightning-exchange/matcher"
	"github.com/NguyenHiu/lightning-exchange/reporter"
	"github.com/NguyenHiu/lightning-exchange/supermatcher"
	"github.com/NguyenHiu/lightning-exchange/user"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

	// Start Super Matcher
	go StartSuperMatcher(onchain, clientNode, constants.KEY_SUPER_MATCHER, constants.SUPER_MATCHER_PORT)

	superMatcherURI := fmt.Sprintf("http://127.0.0.1:%v", constants.SUPER_MATCHER_PORT)

	// Init matcher
	matcher1 := matcher.NewMatcher(assetHolders, adj, appAddr, onchain, constants.KEY_MATCHER_1, superMatcherURI, clientNode, constants.CHAIN_ID, token)
	matcher1.Register()
	// Init Alice
	alice := user.NewUser(constants.KEY_ALICE)
	bus_1, adj_1, ahs_1, app_1, stakes_1 := matcher1.SetupClient(alice.ID)
	alice.SetupClient(bus_1, constants.CHAIN_URL, adj_1, ahs_1, app_1, stakes_1, token)
	_logger.Info("Opening channel.\n")
	ok := matcher1.OpenAppChannel(alice.ID, alice.AppClient.WireAddress())
	if !ok {
		log.Fatalln("OpenAppChannel Failed")
	}
	alice.AcceptedChannel()

	// Init matcher
	matcher2 := matcher.NewMatcher(assetHolders, adj, appAddr, onchain, constants.KEY_MATCHER_2, superMatcherURI, clientNode, constants.CHAIN_ID, token)
	matcher2.Register()
	// Init Bob
	bob := user.NewUser(constants.KEY_BOB)
	bus_2, adj_2, ahs_2, app_2, stakes_2 := matcher2.SetupClient(bob.ID)
	bob.SetupClient(bus_2, constants.CHAIN_URL, adj_2, ahs_2, app_2, stakes_2, token)
	ok = matcher2.OpenAppChannel(bob.ID, bob.AppClient.WireAddress())
	if !ok {
		log.Fatalln("OpenAppChannel Failed")
	}
	bob.AcceptedChannel()

	// Create Order 1
	order_1 := App.NewOrder(client.EthToWei(big.NewFloat(5)), big.NewInt(5), constants.ASK, alice.AppClient.WalletAddressAsEthwallet(), "P")
	alicePrvKey, err := crypto.HexToECDSA(constants.KEY_ALICE)
	if err != nil {
		panic(err)
	}
	order_1.Sign(*alicePrvKey)
	alice.SendNewOrder(&order_1)

	// Create Order 2
	order_2 := App.NewOrder(client.EthToWei(big.NewFloat(5)), big.NewInt(5), constants.BID, bob.AppClient.WalletAddressAsEthwallet(), "P")
	bobPrvKey, err := crypto.HexToECDSA(constants.KEY_BOB)
	if err != nil {
		panic(err)
	}
	order_2.Sign(*bobPrvKey)
	bob.SendNewOrder(&order_2)

	order_3 := App.NewOrder(client.EthToWei(big.NewFloat(7)), big.NewInt(5), constants.BID, bob.AppClient.WalletAddressAsEthwallet(), "P")
	order_3.Sign(*bobPrvKey)
	bob.SendNewOrder(&order_3)

	order_4 := App.NewOrder(client.EthToWei(big.NewFloat(6)), big.NewInt(5), constants.BID, bob.AppClient.WalletAddressAsEthwallet(), "P")
	order_4.Sign(*bobPrvKey)
	bob.SendNewOrder(&order_4)

	// <-time.After(time.Second * 10)

	// <-time.After(time.Second * 100)

	// Create Final Order
	lastOrder_1 := App.NewOrder(&big.Int{}, &big.Int{}, constants.BID, alice.AppClient.WalletAddressAsEthwallet(), "F")
	lastOrder_1.Sign(*alicePrvKey)
	alice.SendNewOrder(&lastOrder_1)

	// Create Final Order
	lastOrder_2 := App.NewOrder(&big.Int{}, &big.Int{}, constants.BID, bob.AppClient.WalletAddressAsEthwallet(), "F")
	lastOrder_2.Sign(*bobPrvKey)
	bob.SendNewOrder(&lastOrder_2)

	// listener.LogOrderBookOverview(rp.OnchainInstance)

	// Payout.
	_logger.Info("Settle\n")
	alice.Settle()
	matcher1.Settle(alice.ID)
	bob.Settle()
	matcher2.Settle(bob.ID)

	// Cleanup.
	_logger.Info("Shutdown\n")
	alice.Shutdown()
	matcher1.Shutdown(alice.ID)
	bob.Shutdown()
	matcher2.Shutdown(bob.ID)

	<-time.After(time.Second * 100)
}
