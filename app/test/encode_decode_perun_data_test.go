package test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/NguyenHiu/lightning-exchange/app"
	"github.com/NguyenHiu/lightning-exchange/constants"
	"github.com/ethereum/go-ethereum/crypto"
	"perun.network/go-perun/backend/ethereum/wallet"
)

func Test(t *testing.T) {
	prvKey, _ := crypto.HexToECDSA(constants.KEY_ALICE)
	alice := wallet.AsWalletAddr(crypto.PubkeyToAddress(prvKey.PublicKey))
	order := app.NewOrder(big.NewInt(5), big.NewInt(5), constants.BID, alice)
	order.Sign(constants.KEY_ALICE)
	order_1 := app.NewOrder(big.NewInt(5), big.NewInt(5), constants.BID, alice)
	order_1.Sign(constants.KEY_ALICE)
	verifyApp := app.NewVerifyApp(alice)
	sampleData := verifyApp.InitData()

	sampleData.Orders[order.OrderID] = order
	sampleData.Orders[order_1.OrderID] = order_1
	sampleData.Msgs[order.OrderID] = []*app.Message{}
	sampleData.Msgs[order_1.OrderID] = []*app.Message{}

	buf := new(bytes.Buffer)
	if err := sampleData.Encode(buf); err != nil {
		t.Errorf("encode sample data fail, err: %v", err)
		return
	}

	decodedData, err := verifyApp.DecodeData(buf)
	if err != nil {
		t.Errorf("decode sample data fail, err: %v", err)
		return
	}

	d := decodedData.(*app.VerifyAppData)
	if len(d.Orders) != len(sampleData.Orders) {
		t.Error("missing order")
		return
	}

	if len(d.Msgs) != len(sampleData.Msgs) {
		t.Errorf("incorrect the numer of message, expect %v, got %v", len(sampleData.Msgs), len(d.Msgs))
		return
	}

	for k, v := range d.Orders {
		if k != app.EndID && !order.Equal(v) {
			t.Error("invalid decoded order")
			return
		}
		if len(d.Msgs[k]) != 0 {
			t.Errorf("incorrect the number of message, expect 0, got %v", len(d.Msgs[k]))
			return
		}
	}
}
