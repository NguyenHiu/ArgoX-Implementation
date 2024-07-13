package main

import (
	"context"
	"log"
	"math/big"
	"math/rand"
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

func SetupSuperMatcher(onchainAddr common.Address, client *ethclient.Client, privateKeyHex string, port int) *supermatcher.SuperMatcher {
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

func PrintBalances(tokenAddr common.Address, clientNode bind.ContractBackend, addrs ...common.Address) {
	tokenInstance, err := TOKEN.NewToken(tokenAddr, clientNode)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(addrs); i++ {
		bal, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: context.Background()}, addrs[i])
		if err != nil {
			log.Fatal(err)
		}
		_logger.Info("[%v] gvn token: %v\n", addrs[i].String()[:5], bal)
	}
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
	sm := SetupSuperMatcher(_onchain, clientNode, constants.KEY_SUPER_MATCHER, constants.SUPER_MATCHER_PORT)

	// Init matcher 1
	matcher1 := matcher.NewMatcher(assetHolders, adj, appAddr, _onchain, constants.KEY_MATCHER_1, clientNode, constants.CHAIN_ID, _token, sm)
	matcher1.Register()

	PrintBalances(_token, clientNode, matcher1.Address)

	// // Init matcher 2
	// matcher2 := matcher.NewMatcher(assetHolders, adj, appAddr, _onchain, constants.KEY_MATCHER_2, clientNode, constants.CHAIN_ID, _token, sm)
	// matcher2.Register()

	// Init Alice
	alice := user.NewUser(constants.KEY_ALICE)
	busAliceOrder, busAliceTrade := matcher1.SetupClient(alice.ID)
	alice.SetupClient(busAliceOrder, busAliceTrade, constants.CHAIN_URL, matcher1.Adjudicator, matcher1.AssetHolders, matcher1.OrderApp, matcher1.TradeApp, matcher1.Stakes, matcher1.EmptyStakes, _token)
	if ok := matcher1.OpenAppChannel(alice.ID, alice.OrderAppClient.WireAddress()); !ok {
		log.Fatalln("OpenAppChannel Failed")
	}
	alice.AcceptedChannel()
	PrintBalances(_token, clientNode, matcher1.Address, alice.Address)

	// Init Bob
	bob := user.NewUser(constants.KEY_BOB)
	busBobOrder, busBobTrade := matcher1.SetupClient(bob.ID)
	bob.SetupClient(busBobOrder, busBobTrade, constants.CHAIN_URL, matcher1.Adjudicator, matcher1.AssetHolders, matcher1.OrderApp, matcher1.TradeApp, matcher1.Stakes, matcher1.EmptyStakes, _token)
	if ok := matcher1.OpenAppChannel(bob.ID, bob.OrderAppClient.WireAddress()); !ok {
		log.Fatalln("OpenAppChannel Failed")
	}
	bob.AcceptedChannel()

	PrintBalances(_token, clientNode, matcher1.Address, bob.Address)

	/* */
	orders, _ := data.LoadOrders("./data/orders.json")
	for _, order := range orders {
		if rand.Int()%2 == 0 {
			newOrder := orderApp.NewOrder(
				big.NewInt(int64(order.Price)),
				big.NewInt(int64(order.Amount)),
				order.Side,
				wallet.AsWalletAddr(alice.Address),
			)
			newOrder.Sign(alice.PrivateKey)
			alice.SendNewOrders([]*orderApp.Order{newOrder})
			<-time.After(time.Millisecond * 500)
		} else {
			newOrder := orderApp.NewOrder(
				big.NewInt(int64(order.Price)),
				big.NewInt(int64(order.Amount)),
				order.Side,
				wallet.AsWalletAddr(bob.Address),
			)
			newOrder.Sign(bob.PrivateKey)
			bob.SendNewOrders([]*orderApp.Order{newOrder})
			<-time.After(time.Millisecond * 500)
		}
	}

	/*
	 *		Random new orders
	 */
	// _bobOrders := []*data.OrderData{}
	// _aliceOrders := []*data.OrderData{}
	// for i := 0; i < 100; i++ {
	// 	bobOrderData := data.RandomOrders(1, 20, 1, 10, 1)
	// 	// _bobOrders = append(_bobOrders, bobOrderData...)
	// 	bobOrders := FromData(bobOrderData, bob.PrivateKey)
	// 	bob.SendNewOrders(bobOrders)
	// 	<-time.After(time.Millisecond * 500)

	// 	aliceOrderData := data.RandomOrders(15, 30, 1, 10, 1)
	// 	// _aliceOrders = append(_aliceOrders, aliceOrderData...)
	// 	aliceOrders := FromData(aliceOrderData, alice.PrivateKey)
	// 	alice.SendNewOrders(aliceOrders)
	// 	<-time.After(time.Millisecond * 500)
	// }

	// matcher1.ExportOrdersData("./data/orders.json")

	// data.SaveOrders(_bobOrders, "./data/bobOrders.json")
	// data.SaveOrders(_aliceOrders, "./data/aliceOrders.json")

	PrintBalances(_token, clientNode, matcher1.Address, alice.Address, bob.Address)

	_logger.Debug("waiting for end orders\n")
	<-time.After(time.Second * 5)

	{
		// Create Final Order
		order, err := orderApp.EndOrder(constants.KEY_ALICE)
		if err != nil {
			_logger.Error("create an end order is fail, err: %v\n", err)
		}
		alice.SendNewOrders([]*orderApp.Order{order})
		bob.SendNewOrders([]*orderApp.Order{order})
	}

	// {
	// 	// Create Final Order
	// 	order, err := orderApp.EndOrder(constants.KEY_BOB)
	// 	if err != nil {
	// 		_logger.Error("create an end order is fail, err: %v\n", err)
	// 	}
	// }

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

	PrintBalances(_token, clientNode, matcher1.Address, alice.Address, bob.Address)
	_logger.Debug("total gas: %v\n",
		util.CalculateTotalUsedGas(alice.Address)+
			util.CalculateTotalUsedGas(bob.Address)+
			util.CalculateTotalUsedGas(matcher1.Address)+
			util.CalculateTotalUsedGas(matcher1.SuperMatcherInstance.Address),
	)

	_logger.Debug("supermatcher: %v\n",
		util.CalculateTotalUsedGas(matcher1.SuperMatcherInstance.Address))

	<-time.After(time.Second * 2)
}
