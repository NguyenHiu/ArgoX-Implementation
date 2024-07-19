package main

import (
	"log"
	"math/big"
	"time"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	"github.com/NguyenHiu/lightning-exchange/data"
	"github.com/NguyenHiu/lightning-exchange/deploy"
	"github.com/NguyenHiu/lightning-exchange/listener"
	"github.com/NguyenHiu/lightning-exchange/logger"
	"github.com/NguyenHiu/lightning-exchange/matcher"
	"github.com/NguyenHiu/lightning-exchange/orderApp"
	"github.com/NguyenHiu/lightning-exchange/reporter"
	"github.com/NguyenHiu/lightning-exchange/user"
	"github.com/NguyenHiu/lightning-exchange/util"
	"github.com/NguyenHiu/lightning-exchange/worker"
	"github.com/ethereum/go-ethereum/ethclient"
	"perun.network/go-perun/backend/ethereum/wallet"
)

var _logger = logger.NewLogger("Main", logger.Red, logger.Bold)

func main() {
	START_TIME := time.Now()

	// Get ethclient
	clientNode, err := ethclient.Dial(constants.CHAIN_URL)
	if err != nil {
		log.Fatal(err)
	}

	_logger.Debug("Deploy smart contracts...\n")
	// Deploy contracts
	deploy.DeployContracts()
	_token, _onchain, adj, assetHolders, appAddr := getContracts()

	_logger.Debug("Deposit ETH to the exchange...\n")
	// Deposit
	onchainInstance, _ := onchain.NewOnchain(_onchain, clientNode)
	util.DepositETH(onchainInstance, clientNode, constants.KEY_ALICE)
	util.DepositETH(onchainInstance, clientNode, constants.KEY_BOB)
	util.DepositETH(onchainInstance, clientNode, constants.KEY_DEPLOYER)

	_logger.Debug("Start Listener...\n")
	// Listener
	go listener.StartListener(_onchain)

	_logger.Debug("Start Reporter...\n")
	// Reporter
	rp, err := reporter.NewReporter(_onchain, constants.KEY_REPORTER, constants.CHAIN_ID)
	if err != nil {
		_logger.Error("Create reporter error, err: %v\n", err)
	}
	rp.Listening() // Start the reporter
	rp.Reporting() // Start reporting

	_logger.Debug("Start Worker...\n")
	// Worker
	worker := worker.NewWorker(_onchain, constants.KEY_WORKER, clientNode)
	worker.Listening()

	_logger.Debug("Start Super Matcher...\n")
	// Super Matcher
	sm := SetupSuperMatcher(_onchain, clientNode, constants.KEY_SUPER_MATCHER, constants.SUPER_MATCHER_PORT)

	_logger.Debug("Init matchers...\n")
	// Slice of matchers
	matchers := []*matcher.Matcher{}
	for i := 0; i < constants.NO_MATCHER; i++ {
		matcher := matcher.NewMatcher(
			assetHolders,
			adj,
			appAddr,
			_onchain,
			constants.MATCHER_PRVKEYS[i],
			clientNode,
			constants.CHAIN_ID,
			_token,
			sm,
		)
		matcher.Register()
		matchers = append(matchers, matcher)
	}

	_logger.Debug("Send Orders...")
	// Alice: Initialize connections to all matchers
	alice := user.NewUser(constants.KEY_ALICE)
	for i := 0; i < constants.NO_MATCHER; i++ {
		orderBus, tradeBus := matchers[i].SetupClient(alice.ID)
		alice.SetupClient(
			matchers[i].ID, orderBus,
			tradeBus,
			constants.CHAIN_URL,
			matchers[i].Adjudicator,
			matchers[i].AssetHolders,
			matchers[i].OrderApp,
			matchers[i].TradeApp,
			matchers[i].Stakes,
			matchers[i].EmptyStakes,
			_token,
		)

		if ok := matchers[i].OpenAppChannel(
			alice.ID,
			alice.Connections[matchers[i].ID].OrderAppClient.WireAddress(),
		); !ok {
			log.Fatalln("OpenAppChannel Failed")
		}
		alice.AcceptedChannel(matchers[i].ID)
		<-time.After(time.Millisecond * 500)
	}

	for _, matcher := range matchers {
		PrintBalances(_token, clientNode, matcher.Address)
	}

	orders, _ := data.LoadOrders("./data/_orders.json")
	for _, order := range orders {
		newOrder := orderApp.NewOrder(
			big.NewInt(int64(order.Price)),
			big.NewInt(int64(order.Amount)),
			order.Side,
			wallet.AsWalletAddr(alice.Address),
		)
		newOrder.Sign(alice.PrivateKey)
		j := 0
		_check := time.Now()
		for j < constants.SEND_TO {
			if time.Since(_check).Milliseconds() > 500*constants.NO_MATCHER {
				_logger.Error("Send Order Loop runs forever!\n")
				break
			}
			for _id, _conn := range alice.Connections {
				if j >= constants.SEND_TO {
					break
				}

				_conn.Mux.Lock()
				isBlocked := _conn.IsBlocked
				_conn.Mux.Unlock()
				if !isBlocked {
					alice.SendNewOrders(_id, []*orderApp.Order{newOrder})
					_conn.IsBlocked = true
					go func(conn *user.Connection) {
						<-time.After(time.Millisecond * 500)
						conn.Mux.Lock()
						defer conn.Mux.Unlock()
						conn.IsBlocked = false
					}(_conn)
					j += 1
					_check = time.Now()
				}
			}
		}
	}

	/*
	 *		Random new orders
	 */
	// _newOrders := []*data.OrderData{}
	// for i := 0; i < 1000; i++ {
	// 	aliceOrderData := data.RandomOrders(15, 30, 1, 10, 1)
	// 	_newOrders = append(_newOrders, aliceOrderData...)
	// 	aliceOrders := FromData(aliceOrderData, alice.PrivateKey)
	// 	for _, newOrder := range aliceOrders {
	// 		j := 0
	// 		_check := time.Now()
	// 		for j < constants.SEND_TO {
	// 			if time.Since(_check).Milliseconds() > 500*constants.NO_MATCHER {
	// 				_logger.Error("Send Order Loop runs forever!\n")
	// 				break
	// 			}
	// 			for _id, _conn := range alice.Connections {
	// 				_conn.Mux.Lock()
	// 				isBlocked := _conn.IsBlocked
	// 				_conn.Mux.Unlock()
	// 				if !isBlocked {
	// 					alice.SendNewOrders(_id, []*orderApp.Order{newOrder})
	// 					_conn.IsBlocked = true
	// 					go func(conn *user.Connection) {
	// 						<-time.After(time.Millisecond * 500)
	// 						conn.Mux.Lock()
	// 						defer conn.Mux.Unlock()
	// 						conn.IsBlocked = false
	// 					}(_conn)
	// 					j += 1
	// 					_check = time.Now()
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// data.SaveOrders(_newOrders, "./data/_orders.json")

	_logger.Debug("waiting for end orders\n")
	<-time.After(time.Second * 10)

	// Create Final Order
	_finalOrder, err := orderApp.EndOrder(constants.KEY_ALICE)
	if err != nil {
		_logger.Error("create an end order is fail, err: %v\n", err)
	}
	for _id := range alice.Connections {
		alice.SendNewOrders(_id, []*orderApp.Order{_finalOrder})
	}
	<-time.After(time.Second * 20)

	// Payout.
	_logger.Info("Settle\n")
	_logger.Debug("Settle All\n")
	alice.SettleAll()
	for _, matcher := range matchers {
		_logger.Debug("Settle Matcher::%v\n", matcher.Address)
		if _, _ok := alice.Connections[matcher.ID]; _ok {
			matcher.Settle(alice.ID)
		}
	}

	// Cleanup.
	_logger.Info("Shutdown\n")
	_logger.Info("Shutdown All\n")
	alice.ShutdownAll()
	for _, matcher := range matchers {
		_logger.Info("Shutdown Matcher::%v\n", matcher.Address)
		matcher.Shutdown(alice.ID)
	}

	// Export Data

	// _totalMatchedAmount := new(big.Int)
	// _totalTime := 0
	// _totalGas := util.CalculateTotalUsedGas(alice.Address) + util.CalculateTotalUsedGas(sm.Address)
	// for _, matcher := range matchers {
	// 	_totalGas += util.CalculateTotalUsedGas(matcher.Address)
	// 	_totalTime += int(matcher.TotalTime)
	// 	_totalMatchedAmount.Add(_totalMatchedAmount, matcher.TotalMatchedAmount)
	// }
	// _logger.Debug("total gas: %v\n", _totalGas)
	// _logger.Debug("supermatcher: %v\n", util.CalculateTotalUsedGas(sm.Address))
	// _logger.Debug("total matched amount: %v\n", _totalMatchedAmount)
	// _logger.Debug("total matched time: %v\n", _totalTime)
	_logger.Debug("Running time: %v seconds\n", time.Since(START_TIME).Seconds())

	<-time.After(time.Second * 2)
}
