package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/contracts/generated/onchain"
	TOKEN "github.com/NguyenHiu/lightning-exchange/contracts/generated/token"
	"github.com/NguyenHiu/lightning-exchange/data"
	"github.com/NguyenHiu/lightning-exchange/orderApp"
	"github.com/NguyenHiu/lightning-exchange/supermatcher"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"perun.network/go-perun/backend/ethereum/wallet"
)

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

	sm, err := supermatcher.NewSuperMatcher(onchainInstance, privateKeyHex, port, int(constants.CHAIN_ID))
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
		// _logger.Info("[%v] gvn token: %v\n", addrs[i].String()[:5], bal)
		fmt.Printf("[%v] gvn token: %v\n", addrs[i].String()[:5], bal)
	}
}

func ExportPriceCurve(priceCurve []*big.Int, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("cannot open file: %v", filename)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(&priceCurve); err != nil {
		return fmt.Errorf("cannot write data into file: %v", filename)
	}

	return nil
}

func ExportRunLogs(
	aliceGas,
	smGas,
	rGas,
	wGas,
	mGas int,
	localMatchedAmount,
	onchainMatchedAmount *big.Int,
	localMatchTime,
	onchainMatchTime int,
	filename string,
) error {
	data := struct {
		AliceGas             int
		SuperMatcherGas      int
		ReporterGas          int
		WorkerGas            int
		MatcherGas           int
		LocalMatchedAmount   *big.Int
		OnchainMatchedAmount *big.Int
		LocalMatchTime       int
		OnchainMatchTime     int
	}{
		AliceGas:             aliceGas,
		SuperMatcherGas:      smGas,
		ReporterGas:          rGas,
		WorkerGas:            wGas,
		MatcherGas:           mGas,
		LocalMatchedAmount:   localMatchedAmount,
		OnchainMatchedAmount: onchainMatchedAmount,
		LocalMatchTime:       localMatchTime,
		OnchainMatchTime:     onchainMatchTime,
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("cannot open file: %v", filename)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(&data); err != nil {
		return fmt.Errorf("cannot write data into file: %v", filename)
	}

	return nil
}
