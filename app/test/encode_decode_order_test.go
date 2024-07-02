package test

import (
	"testing"
)

func TestEncodeDecodeOrder(t *testing.T) {

	// prvKey, _ := crypto.HexToECDSA(constants.KEY_ALICE)
	// alice := wallet.AsWalletAddr(crypto.PubkeyToAddress(prvKey.PublicKey))

	// order := app.NewOrder(big.NewInt(5), big.NewInt(5), constants.BID, alice)
	// err := order.Sign(constants.KEY_ALICE)
	// if err != nil {
	// 	t.Errorf("Sign order error, err: %v\n", err)
	// }

	// encodedData := order.Encode_TransferLightning()
	// t.Logf("encoded order: %v\n", encodedData)

	// decodedOrder, err := app.Order_Decode_TransferLightning(encodedData)
	// if err != nil {
	// 	t.Errorf("Can not decode order")
	// }
	// t.Logf("decoded order: %v\n", decodedOrder)

	// if order.OrderID != decodedOrder.OrderID {
	// 	t.Errorf("Expected OrderID: %v, got %v", order.OrderID, decodedOrder.OrderID)
	// }

	// if order.Price.Cmp(decodedOrder.Price) != 0 {
	// 	t.Errorf("Expected Price: %v, got %v", order.Price, decodedOrder.Price)
	// }

	// if order.Amount.Cmp(decodedOrder.Amount) != 0 {
	// 	t.Errorf("Expected Amount: %v, got %v", order.Amount, decodedOrder.Amount)
	// }

	// if order.Side != decodedOrder.Side {
	// 	t.Errorf("Expected Side: %v, got %v", order.Side, decodedOrder.Side)
	// }

	// if order.Owner.Cmp(decodedOrder.Owner) != 0 {
	// 	t.Errorf("Expected Owner: %v, got %v", order.Owner, decodedOrder.Owner)
	// }

	// if !bytes.Equal(order.Signature, decodedOrder.Signature) {
	// 	t.Errorf("Expected Signature: %v, got %v", order.Signature, decodedOrder.Signature)
	// }
}
