package test

import (
	"math/big"
	"testing"

	"github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/ethereum/go-ethereum/crypto"
	"perun.network/go-perun/backend/ethereum/wallet"
)

func TestOrderSignature(t *testing.T) {

	prvKey, _ := crypto.HexToECDSA(constants.KEY_ALICE)
	alice := wallet.AsWalletAddr(crypto.PubkeyToAddress(prvKey.PublicKey))
	order := app.NewOrder(big.NewInt(5), big.NewInt(5), constants.BID, alice)

	t.Logf("OrderID: %v\n", order.OrderID)

	t.Logf("Before signing, signature: %v\n", order.Signature)
	err := order.Sign(constants.KEY_ALICE)
	if err != nil {
		t.Errorf("Error when signing, error: %v", err)
	}
	t.Logf("After signing, signature: %v\n", order.Signature)

	// verify signature
	if !order.IsValidSignature() {
		t.Errorf("Invalid signature, error: %v", err)
	} else {
		t.Log("Valid signature")
	}
}
