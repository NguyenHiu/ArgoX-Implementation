package test

import (
	"testing"
)

func Test(t *testing.T) {
	// prvKey, _ := crypto.HexToECDSA(constants.KEY_ALICE)
	// alice := wallet.AsWalletAddr(crypto.PubkeyToAddress(prvKey.PublicKey))
	// order := app.NewOrder(big.NewInt(5), big.NewInt(5), constants.BID, alice)
	// order.Sign(constants.KEY_ALICE)
	// order_1 := app.NewOrder(big.NewInt(5), big.NewInt(5), constants.BID, alice)
	// order_1.Sign(constants.KEY_ALICE)
	// verifyApp := app.NewVerifyApp(alice)
	// sampleData := verifyApp.InitData()

	// sampleData.Orders[order.OrderID] = order
	// sampleData.Orders[order_1.OrderID] = order_1

	// buf := new(bytes.Buffer)
	// if err := sampleData.Encode(buf); err != nil {
	// 	t.Errorf("encode sample data fail, err: %v", err)
	// 	return
	// }

	// decodedData, err := verifyApp.DecodeData(buf)
	// if err != nil {
	// 	t.Errorf("decode sample data fail, err: %v", err)
	// 	return
	// }

	// d := decodedData.(*app.VerifyAppData)
	// if len(d.Orders) != len(sampleData.Orders) {
	// 	t.Error("missing order")
	// 	return
	// }

	// for k, v := range d.Orders {
	// 	if k != app.EndID && !order.Equal(v) {
	// 		t.Error("invalid decoded order")
	// 		return
	// 	}
	// }
}
