package main

import (
	"fmt"
	"log"
	"math/big"
	"time"

	App "github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/client"
	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/NguyenHiu/lightning-exchange/data"
	"github.com/NguyenHiu/lightning-exchange/deploy"
	"github.com/NguyenHiu/lightning-exchange/matcher"
	"github.com/NguyenHiu/lightning-exchange/user"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

func main() {
	deploy.DeployContracts()
	token, onchain, adj, assetHolders, appAddr := getContracts()

	clientNode := "http://127.0.0.1:8545"

	superMatcherURI := ""

	// Init matcher
	matcher1 := matcher.NewMatcher(assetHolders, adj, appAddr, onchain, constants.KEY_MATCHER, superMatcherURI, clientNode, constants.CHAIN_ID, token)
	matcher1.Register()
	// Init Alice
	alice := user.NewUser(constants.KEY_ALICE)
	bus_1, adj_1, ahs_1, app_1, stakes_1 := matcher1.SetupClient(alice.ID)
	alice.SetupClient(bus_1, constants.CHAIN_URL, adj_1, ahs_1, app_1, stakes_1, token)
	log.Println("Opening channel.")
	ok := matcher1.OpenAppChannel(alice.ID, alice.AppClient.WireAddress())
	if !ok {
		log.Fatalln("OpenAppChannel Failed")
	}
	alice.AcceptedChannel()

	// Init matcher
	matcher2 := matcher.NewMatcher(assetHolders, adj, appAddr, onchain, constants.KEY_MATCHER, superMatcherURI, clientNode, constants.CHAIN_ID, token)
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
	order_1 := App.NewOrder(client.EthToWei(big.NewFloat(5)), big.NewInt(5), constants.BID, alice.AppClient.WalletAddressAsEthwallet(), "P")
	alicePrvKey, err := crypto.HexToECDSA(constants.KEY_ALICE)
	if err != nil {
		panic(err)
	}
	order_1.Sign(*alicePrvKey)
	alice.SendNewOrder(&order_1)

	// Create Order 2
	order_2 := App.NewOrder(client.EthToWei(big.NewFloat(6)), big.NewInt(5), constants.BID, bob.AppClient.WalletAddressAsEthwallet(), "P")
	bobPrvKey, err := crypto.HexToECDSA(constants.KEY_BOB)
	if err != nil {
		panic(err)
	}
	order_2.Sign(*bobPrvKey)
	bob.SendNewOrder(&order_2)

	// order_3 := App.NewOrder(client.EthToWei(big.NewFloat(7)), big.NewInt(5), constants.BID, bob.AppClient.WalletAddressAsEthwallet(), "P")
	// order_3.Sign(*bobPrvKey)
	// bob.SendNewOrder(&order_3)

	// order_4 := App.NewOrder(client.EthToWei(big.NewFloat(6)), big.NewInt(6), constants.BID, bob.AppClient.WalletAddressAsEthwallet(), "P")
	// order_4.Sign(*bobPrvKey)
	// bob.SendNewOrder(&order_4)

	// Create Final Order
	lastOrder_1 := App.NewOrder(&big.Int{}, &big.Int{}, constants.BID, alice.AppClient.WalletAddressAsEthwallet(), "F")
	lastOrder_1.Sign(*alicePrvKey)
	alice.SendNewOrder(&lastOrder_1)

	// Create Final Order
	lastOrder_2 := App.NewOrder(&big.Int{}, &big.Int{}, constants.BID, bob.AppClient.WalletAddressAsEthwallet(), "F")
	lastOrder_2.Sign(*bobPrvKey)
	bob.SendNewOrder(&lastOrder_2)

	<-time.After(20 * time.Second)
	log.Println("DONE")

	// Payout.
	fmt.Println("Settle")
	alice.Settle()
	matcher1.Settle(alice.ID)
	bob.Settle()
	matcher2.Settle(bob.ID)

	// Cleanup.
	fmt.Println("Shutdown")
	alice.Shutdown()
	matcher1.Shutdown(alice.ID)
	bob.Shutdown()
	matcher2.Shutdown(bob.ID)

}
