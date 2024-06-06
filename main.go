package main

import (
	"fmt"
	"log"
	"math/big"
	"time"

	App "github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/client"
	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/matcher"
	"github.com/NguyenHiu/lightning-exchange/user"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Init matcher
	matcherIns := matcher.NewMatcher(constants.CHAIN_URL, constants.CHAIN_ID, constants.KEY_DEPLOYER)

	// Init Alice
	alice := user.NewUser(constants.KEY_ALICE)
	bus_1, adj_1, ahs_1, app_1, stakes_1 := matcherIns.SetupClient(alice.ID)
	alice.SetupClient(bus_1, constants.CHAIN_URL, adj_1, ahs_1, app_1, stakes_1)
	log.Println("Opening channel.")
	ok := matcherIns.OpenAppChannel(alice.ID, alice.AppClient.WireAddress())
	if !ok {
		log.Fatalln("OpenAppChannel Failed")
	}
	alice.AcceptedChannel()

	// Init Bob
	bob := user.NewUser(constants.KEY_BOB)
	bus_2, adj_2, ahs_2, app_2, stakes_2 := matcherIns.SetupClient(bob.ID)
	bob.SetupClient(bus_2, constants.CHAIN_URL, adj_2, ahs_2, app_2, stakes_2)
	ok = matcherIns.OpenAppChannel(bob.ID, bob.AppClient.WireAddress())
	if !ok {
		log.Fatalln("OpenAppChannel Failed")
	}
	bob.AcceptedChannel()

	// Create Order 1
	order_1 := App.NewOrder(client.EthToWei(big.NewFloat(5)).Int64(), 5, constants.BID, alice.AppClient.WalletAddressAsEthwallet(), "P")
	alicePrvKey, err := crypto.HexToECDSA(constants.KEY_ALICE)
	if err != nil {
		panic(err)
	}
	order_1.Sign(*alicePrvKey)
	alice.SendNewOrder(&order_1)

	// Create Order 2
	order_2 := App.NewOrder(client.EthToWei(big.NewFloat(6)).Int64(), 5, constants.BID, bob.AppClient.WalletAddressAsEthwallet(), "P")
	bobPrvKey, err := crypto.HexToECDSA(constants.KEY_BOB)
	if err != nil {
		panic(err)
	}
	order_2.Sign(*bobPrvKey)
	bob.SendNewOrder(&order_2)

	order_3 := App.NewOrder(client.EthToWei(big.NewFloat(7)).Int64(), 5, constants.BID, bob.AppClient.WalletAddressAsEthwallet(), "P")
	order_3.Sign(*bobPrvKey)
	bob.SendNewOrder(&order_3)

	order_4 := App.NewOrder(client.EthToWei(big.NewFloat(6)).Int64(), 5, constants.BID, bob.AppClient.WalletAddressAsEthwallet(), "P")
	order_4.Sign(*bobPrvKey)
	bob.SendNewOrder(&order_4)

	// Create Final Order
	lastOrder_1 := App.NewOrder(0, 0, constants.BID, alice.AppClient.WalletAddressAsEthwallet(), "F")
	lastOrder_1.Sign(*alicePrvKey)
	alice.SendNewOrder(&lastOrder_1)

	// Create Final Order
	lastOrder_2 := App.NewOrder(0, 0, constants.BID, bob.AppClient.WalletAddressAsEthwallet(), "F")
	lastOrder_2.Sign(*bobPrvKey)
	bob.SendNewOrder(&lastOrder_2)

	<-time.After(200 * time.Second)
	log.Println("DONE")

	// Payout.
	fmt.Println("Settle")
	alice.Settle()
	matcherIns.Settle(alice.ID)
	bob.Settle()
	matcherIns.Settle(bob.ID)

	// Cleanup.
	fmt.Println("Shutdown")
	alice.Shutdown()
	matcherIns.Shutdown(alice.ID)
	bob.Shutdown()
	matcherIns.Shutdown(bob.ID)

}
