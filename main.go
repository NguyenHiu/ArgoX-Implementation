package main

import (
	"fmt"
	"log"
	"math/big"

	App "github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/client"
	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/matcher"
	"github.com/NguyenHiu/lightning-exchange/user"
	utils "github.com/NguyenHiu/lightning-exchange/utils"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	matcherIns := matcher.NewMatcher(constants.CHAIN_URL, constants.CHAIN_ID, constants.KEY_DEPLOYER)
	alice := user.NewUser(constants.KEY_ALICE)

	bus, adj, ahs, app, stakes := matcherIns.SetupClient(alice.ID)
	alice.SetupClient(bus, constants.CHAIN_URL, adj, ahs, app, stakes)

	// Print balances before transactions.
	l := utils.NewBalanceLogger(constants.CHAIN_URL)
	l.LogBalances(matcherIns.ClientConfigs[alice.ID].AppClient, alice.AppClient)

	// Open app channel and play.
	log.Println("Opening channel.")
	ok := matcherIns.OpenAppChannel(alice.ID, alice.AppClient.WireAddress())
	if !ok {
		log.Fatalln("OpenAppChannel Failed")
	}
	alice.AcceptedChannel()

	newOrder := App.NewOrder(client.EthToWei(big.NewFloat(5)).Int64(), 5, constants.BID, alice.AppClient.WalletAddressAsEthwallet(), "P")
	alicePrvKey, err := crypto.HexToECDSA(constants.KEY_ALICE)
	if err != nil {
		panic(err)
	}
	newOrder.Sign(*alicePrvKey)
	alice.SendNewOrder(&newOrder)

	lastOrder := App.NewOrder(0, 0, constants.BID, alice.AppClient.WalletAddressAsEthwallet(), "F")
	lastOrder.Sign(*alicePrvKey)
	alice.SendNewOrder(&lastOrder)

	// Payout.
	fmt.Println("Settle")
	alice.Settle()
	matcherIns.Settle(alice.ID)

	// Print balances after transactions.
	fmt.Println("LogBalances")
	l.LogBalances(matcherIns.ClientConfigs[alice.ID].AppClient, alice.AppClient)

	// Cleanup.
	fmt.Println("Shutdown")
	alice.Shutdown()
	matcherIns.Shutdown(alice.ID)

}
