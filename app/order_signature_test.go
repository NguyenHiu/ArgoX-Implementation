package app

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/ethereum/go-ethereum/crypto"
	"perun.network/go-perun/backend/ethereum/wallet"
)

func TestSignAndVerifyOrderSignature(t *testing.T) {
	price := big.NewInt(10)
	amount := big.NewInt(5)
	side := constants.BID
	k, _ := crypto.GenerateKey()
	p, _ := k.Public().(*ecdsa.PublicKey)
	owner := wallet.AsWalletAddr(crypto.PubkeyToAddress(*p))
	newOrder := NewOrder(price, amount, side, owner, "P")
	t.Logf("OrderID: %v\n", newOrder.OrderID)

	t.Logf("Before signing, signature: %v\n", newOrder.OwnerSignture)
	err := newOrder.Sign(*k)
	if err != nil {
		t.Errorf("Error when signing, error: %v", err)
	}
	t.Logf("After signing, signature: %v\n", newOrder.OwnerSignture)

	// verify signature
	if !newOrder.IsValidSignature() {
		t.Errorf("Invalid signature, error: %v", err)
	} else {
		t.Log("Valid signature")
	}
}
