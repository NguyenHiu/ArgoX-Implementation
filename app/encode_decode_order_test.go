package app

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"perun.network/go-perun/backend/ethereum/wallet"

	"testing"
)

func TestEncodeDecodeOrder(t *testing.T) {
	prvkey, _ := crypto.GenerateKey()
	pubkey, _ := prvkey.Public().(*ecdsa.PublicKey)
	owner := wallet.AsWalletAddr(crypto.PubkeyToAddress(*pubkey))

	order := NewOrder(10.0, 5.0, BID, owner)
	err := order.Sign(*prvkey)
	if err != nil {
		t.Errorf("Sign order error, err: %v\n", err)
	}

	encodedData := order.EncodeOrder()
	fmt.Printf("encoded order: %v\n", encodedData)

	decodedOrder, _ := DecodeOrder(encodedData)
	fmt.Printf("decoded order: %v\n", decodedOrder)

	if order.OrderID != decodedOrder.OrderID {
		t.Errorf("Expected OrderID: %v, got %v", order.OrderID, decodedOrder.OrderID)
	}

	if order.Price != decodedOrder.Price {
		t.Errorf("Expected Price: %f, got %f", order.Price, decodedOrder.Price)
	}

	if order.Amount != decodedOrder.Amount {
		t.Errorf("Expected Amount: %f, got %f", order.Amount, decodedOrder.Amount)
	}

	if order.Side != decodedOrder.Side {
		t.Errorf("Expected Side: %v, got %v", order.Side, decodedOrder.Side)
	}

	if order.Owner.Cmp(decodedOrder.Owner) != 0 {
		t.Errorf("Expected Owner: %v, got %v", order.Owner, decodedOrder.Owner)
	}

	if order.Status != decodedOrder.Status {
		t.Errorf("Expected Status: %v, got %v", order.Status, decodedOrder.Status)
	}

	if order.MatchedAmoount != decodedOrder.MatchedAmoount {
		t.Errorf("Expected MatchedAmoount: %v, got %v", order.MatchedAmoount, decodedOrder.MatchedAmoount)
	}

}
