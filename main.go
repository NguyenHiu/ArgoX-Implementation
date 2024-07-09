package main

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	TOKEN "github.com/NguyenHiu/lightning-exchange/contracts/generated/token"
	"github.com/NguyenHiu/lightning-exchange/data"
	"github.com/NguyenHiu/lightning-exchange/deploy"
	"github.com/NguyenHiu/lightning-exchange/listener"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/NguyenHiu/lightning-exchange/matcher"
	"github.com/NguyenHiu/lightning-exchange/orderApp"
	"github.com/NguyenHiu/lightning-exchange/reporter"
	"github.com/NguyenHiu/lightning-exchange/supermatcher"
	"github.com/NguyenHiu/lightning-exchange/user"
	"github.com/NguyenHiu/lightning-exchange/util"
	"github.com/NguyenHiu/lightning-exchange/worker"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"perun.network/go-perun/backend/ethereum/wallet"
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

func SetupMatcher(onchainAddr common.Address, client *ethclient.Client, privateKeyHex string, port int) *supermatcher.SuperMatcher {
	onchainInstance, err := onchain.NewOnchain(onchainAddr, client)
	if err != nil {
		log.Fatal(err)
	}

	sm, err := supermatcher.NewSuperMatcher(onchainInstance, privateKeyHex, port, constants.CHAIN_ID)
	if err != nil {
		log.Fatal(err)
	}

	return sm
}

func FromData(ordersData []*data.OrderData, ownerPrvKey string) []*orderApp.Order {
	prvkey, _ := crypto.HexToECDSA(ownerPrvKey)
	address := crypto.PubkeyToAddress(prvkey.PublicKey)

	orders := []*orderApp.Order{}
	for _, order := range ordersData {
		id, _ := uuid.NewRandom()
		appOrder := &orderApp.Order{
			OrderID: id,
			Price:   big.NewInt(int64(order.Price)),
			Amount:  big.NewInt(int64(order.Amount)),
			Side:    order.Side,
			Owner:   wallet.AsWalletAddr(address),
		}
		if err := appOrder.Sign(ownerPrvKey); err != nil {
			log.Fatal(err)
		}
		orders = append(orders, appOrder)
	}
	return orders
}

// TODO: Simulation
func main() {

	clientNode, err := ethclient.Dial(constants.CHAIN_URL)
	if err != nil {
		log.Fatal(err)
	}

	// Deploy contracts
	deploy.DeployContracts()
	_token, _onchain, adj, assetHolders, appAddr := getContracts()

	// Listen events from onchain contract
	go listener.StartListener(_onchain)

	// Deposit
	onchainInstance, _ := onchain.NewOnchain(_onchain, clientNode)
	util.DepositETH(onchainInstance, clientNode, constants.KEY_ALICE)
	util.DepositETH(onchainInstance, clientNode, constants.KEY_BOB)
	util.DepositETH(onchainInstance, clientNode, constants.KEY_DEPLOYER)

	// Create a Reporter
	rp, err := reporter.NewReporter(_onchain, constants.KEY_REPORTER, constants.CHAIN_ID)
	if err != nil {
		_logger.Error("Create reporter error, err: %v\n", err)
	}
	// Start the reporter
	rp.Listening()
	rp.Reporting()

	worker := worker.NewWorker(_onchain, constants.KEY_WORKER, clientNode)
	worker.Listening()

	// Start Super Matcher
	sm := SetupMatcher(_onchain, clientNode, constants.KEY_SUPER_MATCHER, constants.SUPER_MATCHER_PORT)

	// Init matcher 1
	matcher1 := matcher.NewMatcher(assetHolders, adj, appAddr, _onchain, constants.KEY_MATCHER_1, clientNode, constants.CHAIN_ID, _token, sm)
	matcher1.Register()

	// Init Bob
	bob := user.NewUser(constants.KEY_BOB)
	busBobOrder, busBobTrade := matcher1.SetupClient(bob.ID)
	bob.SetupClient(busBobOrder, busBobTrade, constants.CHAIN_URL, matcher1.Adjudicator, matcher1.AssetHolders, matcher1.OrderApp, matcher1.TradeApp, matcher1.Stakes, _token)
	if ok := matcher1.OpenAppChannel(bob.ID, bob.OrderAppClient.WireAddress()); !ok {
		log.Fatalln("OpenAppChannel Failed")
	}
	bob.AcceptedChannel()

	// // Init matcher 2
	// matcher2 := matcher.NewMatcher(assetHolders, adj, appAddr, _onchain, constants.KEY_MATCHER_2, clientNode, constants.CHAIN_ID, _token, sm)
	// matcher2.Register()

	// Init Alice
	alice := user.NewUser(constants.KEY_ALICE)
	busAliceOrder, busAliceTrade := matcher1.SetupClient(alice.ID)
	alice.SetupClient(busAliceOrder, busAliceTrade, constants.CHAIN_URL, matcher1.Adjudicator, matcher1.AssetHolders, matcher1.OrderApp, matcher1.TradeApp, matcher1.Stakes, _token)
	if ok := matcher1.OpenAppChannel(alice.ID, alice.OrderAppClient.WireAddress()); !ok {
		log.Fatalln("OpenAppChannel Failed")
	}
	alice.AcceptedChannel()

	{
		bobOrders := FromData(data.RandomOrders(1, 20, 1, 10, 100, "./data/bobOrders.json"), bob.PrivateKey)
		aliceOrders := FromData(data.RandomOrders(15, 30, 1, 10, 100, "./data/aliceOrders.json"), alice.PrivateKey)
		bob.SendNewOrders(bobOrders)
		alice.SendNewOrders(aliceOrders)
		<-time.After(time.Second * 10)
	}
	{
		bobOrders := FromData(data.RandomOrders(1, 20, 1, 10, 100, "./data/bobOrders.json"), bob.PrivateKey)
		aliceOrders := FromData(data.RandomOrders(15, 30, 1, 10, 100, "./data/aliceOrders.json"), alice.PrivateKey)
		bob.SendNewOrders(bobOrders)
		<-time.After(time.Second * 10)
		alice.SendNewOrders(aliceOrders)
	}

	{
		tokenInstance, err := TOKEN.NewToken(_token, clientNode)
		if err != nil {
			log.Fatal(err)
		}
		matcher1Balance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, matcher1.Address)
		if err != nil {
			log.Fatal(err)
		}
		// matcher2Balance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, matcher2.Address)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		bobBalance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, bob.OrderAppClient.WalletAddress())
		if err != nil {
			log.Fatal(err)
		}
		aliceBalance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, alice.OrderAppClient.WalletAddress())
		if err != nil {
			log.Fatal(err)
		}
		_logger.Info("matcher 1's balance: %v\n", matcher1Balance)
		// _logger.Info("matcher 2's balance: %v\n", matcher2Balance)
		_logger.Info("bob's balance: %v\n", bobBalance)
		_logger.Info("alice's balance: %v\n", aliceBalance)
	}

	// _logger.Debug("waiting for end orders\n")
	<-time.After(time.Second * 15)

	{
		// Create Final Order
		order, err := orderApp.EndOrder(constants.KEY_ALICE)
		if err != nil {
			_logger.Error("create an end order is fail, err: %v\n", err)
		}
		alice.SendNewOrders([]*orderApp.Order{order})
	}

	{
		// Create Final Order
		order, err := orderApp.EndOrder(constants.KEY_BOB)
		if err != nil {
			_logger.Error("create an end order is fail, err: %v\n", err)
		}
		bob.SendNewOrders([]*orderApp.Order{order})
	}

	<-time.After(time.Second * 5)

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
		tokenInstance, err := TOKEN.NewToken(_token, clientNode)
		if err != nil {
			log.Fatal(err)
		}
		matcher1Balance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, matcher1.Address)
		if err != nil {
			log.Fatal(err)
		}
		// matcher2Balance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, matcher2.Address)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		bobBalance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, bob.OrderAppClient.WalletAddress())
		if err != nil {
			log.Fatal(err)
		}
		aliceBalance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, alice.OrderAppClient.WalletAddress())
		if err != nil {
			log.Fatal(err)
		}
		_logger.Info("matcher 1's balance: %v\n", matcher1Balance)
		// _logger.Info("matcher 2's balance: %v\n", matcher2Balance)
		_logger.Info("bob's balance: %v\n", bobBalance)
		_logger.Info("alice's balance: %v\n", aliceBalance)
	}

	<-time.After(time.Second * 2)
}
