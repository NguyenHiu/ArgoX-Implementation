package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"perun.network/go-perun/backend/ethereum/wallet"
)

var _logger = logger.NewLogger("Main", logger.Red, logger.Bold)

func main() {

	MAIN_START_TIME := time.Now()

	/// Get user's configurations
	if len(os.Args) != 8 {
		log.Fatal("Didnt provide enough params")
	}
	_port, _ := strconv.Atoi(os.Args[1])
	constants.CHAIN_URL = fmt.Sprintf("ws://127.0.0.1:%v", _port)
	_chainID, _ := strconv.Atoi(os.Args[2])
	constants.CHAIN_ID = int64(_chainID)
	_noMatcher, _ := strconv.Atoi(os.Args[3])
	constants.NO_MATCHER = _noMatcher
	_noSendTo, _ := strconv.Atoi(os.Args[4])
	constants.SEND_TO = _noSendTo

	_TYPE_ := os.Args[5]
	if _TYPE_ != "run" && _TYPE_ != "random" {
		log.Fatal("Invalia type")
	}
	_FILENAME_ := os.Args[6]
	if _, err := os.Stat(_FILENAME_); os.IsNotExist(err) {
		log.Fatalf("File %s does not exist", _FILENAME_)
	}
	_PRICE_CURVE_FOLDER_ := os.Args[7]
	if info, err := os.Stat(_PRICE_CURVE_FOLDER_); err != nil {
		if os.IsNotExist(err) {
			// Attempt to create the directory if it does not exist
			if mkdirErr := os.MkdirAll(_PRICE_CURVE_FOLDER_, 0755); mkdirErr != nil {
				log.Fatalf("Failed to create directory %s: %s", _PRICE_CURVE_FOLDER_, mkdirErr)
			} else {
				log.Printf("Directory %s created", _PRICE_CURVE_FOLDER_)
			}
		} else {
			log.Fatal(err) // Other errors
		}
	} else if !info.IsDir() {
		log.Fatalf("%s is not a directory", _PRICE_CURVE_FOLDER_)
	}

	// Get ethclient
	clientNode, err := ethclient.Dial(constants.CHAIN_URL)
	if err != nil {
		log.Fatal(err)
	}

	// Deploy contracts
	_logger.Debug("Deploy smart contracts...\n")
	deploy.DeployContracts()
	_token, _onchain, adj, assetHolders, appAddr := getContracts()

	// Deposit ETH to the smart contract
	_logger.Debug("Deposit ETH to the exchange...\n")
	onchainInstance, _ := onchain.NewOnchain(_onchain, clientNode)
	util.DepositETH(onchainInstance, clientNode, constants.KEY_ALICE)
	util.DepositETH(onchainInstance, clientNode, constants.KEY_BOB)
	util.DepositETH(onchainInstance, clientNode, constants.KEY_DEPLOYER)

	// Listener
	_logger.Debug("Start Listener...\n")
	_listenerInstance := listener.NewListener()
	go _listenerInstance.StartListener(_onchain)

	// Reporter
	_logger.Debug("Start Reporter...\n")
	rp, err := reporter.NewReporter(_onchain, constants.KEY_REPORTER, constants.CHAIN_ID)
	if err != nil {
		_logger.Error("Create reporter error, err: %v\n", err)
	}
	rp.Listening()
	rp.Reporting()

	// Worker
	_logger.Debug("Start Worker...\n")
	w := worker.NewWorker(_onchain, constants.KEY_WORKER, clientNode)
	w.Listening()

	// Super Matcher
	_logger.Debug("Start Super Matcher...\n")
	sm := SetupSuperMatcher(_onchain, clientNode, constants.KEY_SUPER_MATCHER, constants.SUPER_MATCHER_PORT)

	// Slice of matchers
	_logger.Debug("Init matchers...\n")
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

	// Start Getting Price Curve
	for _, matcher := range matchers {
		// Local
		matcher.IsGetPriceCurve = true
		go matcher.GetPriceCurve()
	}

	// DEBUG
	NoOrderSentToEachMatcher := map[uuid.UUID]int{}
	for _, matcher := range matchers {
		NoOrderSentToEachMatcher[matcher.ID] = 0
	}
	// DEBUG

	// Onchain
	_listenerInstance.IsGetPriceCurve = true
	go _listenerInstance.GetPriceCurve()

	SEND_ORDER_START_TIME := time.Now()

	// Send orders
	_logger.Debug("Send Orders...\n")
	if _TYPE_ == "run" {
		// Run existed orders in `filename` file
		orders, _ := data.LoadOrders(_FILENAME_)
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
			_isSent := make(map[uuid.UUID]bool)
			for j < constants.SEND_TO {
				if time.Since(_check).Milliseconds() > 500*(int64(constants.NO_MATCHER+2)) {
					_logger.Error("Send Order Loop runs forever!\n")
					break
				}

				/** DEBUG */
				// var _id uuid.UUID
				// _rand := rand.Intn(3)
				// if _rand == 0 {
				// 	_id = matchers[0].ID
				// } else if _rand == 1 {
				// 	_id = matchers[1].ID
				// } else {
				// 	_id = matchers[rand.Intn(len(alice.Connections)-2)+2].ID
				// }
				// NoOrderSentToEachMatcher[_id] += 1
				// alice.SendNewOrders(
				// 	_id,
				// 	[]*orderApp.Order{newOrder},
				// )
				// <-time.After(time.Millisecond * 50)
				// j += 1
				/** DEBUG */

				for _id, _conn := range alice.Connections {
					if j >= constants.SEND_TO {
						break
					}

					_conn.Mux.Lock()
					isBlocked := _conn.IsBlocked
					_conn.Mux.Unlock()
					if !isBlocked {
						if _, _ok := _isSent[_id]; _ok {
							continue
						} else {
							_isSent[_id] = true
						}

						alice.SendNewOrders(_id, []*orderApp.Order{newOrder})
						// DEBUG
						NoOrderSentToEachMatcher[_id] += 1
						// DEBUG
						_conn.IsBlocked = true
						go func(conn *user.Connection) {
							<-time.After(time.Millisecond * 50)
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
	} else if _TYPE_ == "random" {
		// Random new orders
		_newOrders := []*data.OrderData{}
		for i := 0; i < 1000; i++ {
			aliceOrderData := data.RandomOrders(15, 30, 1, 10, 1)
			_newOrders = append(_newOrders, aliceOrderData...)
			aliceOrders := FromData(aliceOrderData, alice.PrivateKey)
			for _, newOrder := range aliceOrders {
				j := 0
				_check := time.Now()
				for j < constants.SEND_TO {
					if time.Since(_check).Milliseconds() > 500*int64((constants.NO_MATCHER+2)) {
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
								<-time.After(time.Millisecond * 50)
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
		}
		data.SaveOrders(_newOrders, _FILENAME_)
	}

	// Create Final Order
	_logger.Debug("Done sending orders phase.\n")
	_logger.Debug("Waiting 10 seconds for sending end orders...\n")
	<-time.After(time.Second * 10)

	_finalOrder, err := orderApp.EndOrder(constants.KEY_ALICE)
	if err != nil {
		_logger.Error("create an end order is fail, err: %v\n", err)
	}
	for _id := range alice.Connections {
		alice.SendNewOrders(_id, []*orderApp.Order{_finalOrder})
	}

	SEND_ORDER_TIME_IN_SECONDS := time.Since(SEND_ORDER_START_TIME).Seconds()

	_logger.Debug("Waiting 10 seconds before closing...\n")
	<-time.After(time.Second * 10)

	// Stop collect price curves
	_listenerInstance.IsGetPriceCurve = false
	for _, matcher := range matchers {
		matcher.IsGetPriceCurve = false
	}

	// // Payout.
	// _logger.Info("Settle\n")
	// for _, matcher := range matchers {
	// 	if _, _ok := alice.Connections[matcher.ID]; _ok {
	// 		alice.Settle(matcher.ID)
	// 		matcher.Settle(alice.ID)
	// 	}
	// }

	// // Cleanup.
	// _logger.Info("Shutdown\n")
	// for _, matcher := range matchers {
	// 	alice.Shutdown(matcher.ID)
	// 	matcher.Shutdown(alice.ID)
	// }

	/// Export results
	// Match Amount
	_totalOnchainMatchAmount := _listenerInstance.TotalMatchedAmountOnchain
	_totalLocalMatchAmount := new(big.Int)
	_numberOfMatchedORderLocal := 0

	// Match Time
	_totalLocalTime := 0
	_totalOnchainTime := _listenerInstance.TotalTimeOnchain

	// Profit
	_totalProfitLocal := 0
	_totalProfitOnchain := _listenerInstance.TotalProfitOnchain

	// Gas
	_address := []common.Address{
		alice.Address,
		sm.Address,
		rp.Address,
		w.Address,
	}
	for _, matcher := range matchers {
		_address = append(_address, matcher.Address)
	}
	_gasUsed := util.CalculateTotalUsedGas(_address)

	// Get data from matchers
	_totalMatcherGasUsed := 0
	for idx, matcher := range matchers {
		_numberOfMatchedORderLocal += int(matcher.NumberOfMatchedOrder)
		_totalMatcherGasUsed += _gasUsed[matcher.Address]
		_totalLocalTime += int(matcher.TotalTimeLocal)
		_totalLocalMatchAmount.Add(_totalLocalMatchAmount, matcher.TotalMatchedAmountLocal)
		_totalProfitLocal += int(matcher.TotalProfitLocal.Int64())
		ExportPriceCurve(matcher.PriceCurveLocal, fmt.Sprintf("%v/local_curve_%v.json", _PRICE_CURVE_FOLDER_, idx))
	}
	_totalGas := _gasUsed[alice.Address] + _gasUsed[sm.Address] + _gasUsed[w.Address] + _gasUsed[rp.Address] + _totalMatcherGasUsed
	_logger.Debug("Total Gas: %v\n", _totalGas)
	_logger.Debug("  > Super Matcher: %v\n", _gasUsed[sm.Address])
	_logger.Debug("  > Total Matcher: %v\n", _totalGas)
	_logger.Debug("  > Reporter: %v\n", _gasUsed[rp.Address])
	_logger.Debug("  > Worker: %v\n", _gasUsed[w.Address])
	_logger.Debug("Total Match Amount: %v\n", new(big.Int).Add(
		_totalLocalMatchAmount,
		_totalOnchainMatchAmount,
	))
	_logger.Debug("  > Total Local Match Amount: %v\n", _totalLocalMatchAmount)
	_logger.Debug("  > Total Onchain Match Amount: %v\n", _totalOnchainMatchAmount)
	_logger.Debug("Total Match Time: %v\n", _totalLocalTime+_totalOnchainTime)
	_logger.Debug("  > Total Local Match Time: %v\n", _totalLocalTime)
	_logger.Debug("  > Total Onchain Match Time: %v\n", _totalOnchainTime)
	_logger.Debug("Sending Order Time: %v seconds\n", SEND_ORDER_TIME_IN_SECONDS)

	_logger.Debug("Exporting Price Curves...\n")
	ExportPriceCurve(_listenerInstance.PriceCurveOnchain, fmt.Sprintf("%v/onchain_curve.json", _PRICE_CURVE_FOLDER_))

	_logger.Debug("Exporting Run Logs...\n")
	ExportRunLogs(
		_gasUsed[alice.Address],
		_gasUsed[sm.Address],
		_gasUsed[rp.Address],
		_gasUsed[w.Address],
		_totalMatcherGasUsed,
		_totalLocalMatchAmount,
		_totalOnchainMatchAmount,
		_totalLocalTime,
		_totalOnchainTime,
		_numberOfMatchedORderLocal,
		int(_listenerInstance.NumberOfMatchedOrder),
		_totalProfitLocal,
		int(_totalProfitOnchain.Int64()),
		sm.NoBatches,
		fmt.Sprintf("%v/logs.json", _PRICE_CURVE_FOLDER_),
	)

	_logger.Debug("Total Running Time: %v seconds\n", time.Since(MAIN_START_TIME).Seconds())

	for _id, _num := range NoOrderSentToEachMatcher {
		_logger.Debug("Number of orders sent to matcher %v: %v\n", _id, _num)
	}

}
