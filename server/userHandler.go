package server

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/NguyenHiu/lightning-exchange/orderApp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"perun.network/go-perun/backend/ethereum/wallet"
)

func (s *Server) getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data := struct {
		ID      uuid.UUID
		Address common.Address
	}{
		ID:      s.user.ID,
		Address: s.user.Address,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *Server) sendOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	type __order struct {
		Price  int  `json:"price"`
		Amount int  `json:"amount"`
		Side   bool `json:"side"`
	}

	var order __order
	fmt.Println("r.Body: ", r.Body)
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	newOrder := orderApp.NewOrder(
		big.NewInt(int64(order.Price)),
		big.NewInt(int64(order.Amount)),
		order.Side,
		wallet.AsWalletAddr(s.user.Address),
	)
	newOrder.Sign(s.user.PrivateKey)
	s.user.SendNewOrders(s.matchers[0].ID, []*orderApp.Order{newOrder})

	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Successfully!",
	})
	w.Write(jsonData)
}

func (s *Server) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type __order_status struct {
		OrderID uuid.UUID
		Status  string
	}

	type __matcher_order_status struct {
		MatcherID uuid.UUID
		Orders    []__order_status
	}

	_matcher_order_status := []__matcher_order_status{}
	for _, matcher := range s.matchers {
		_order_status := []__order_status{}
		for _, orderID := range s.user.SentOrderList {
			_status := matcher.GetOrderStatus(orderID)
			if _status != "" {
				_order_status = append(_order_status, __order_status{
					OrderID: orderID,
					Status:  _status,
				})
			}
		}
		_matcher_order_status = append(_matcher_order_status, __matcher_order_status{
			MatcherID: matcher.ID,
			Orders:    _order_status,
		})
	}

	jsonData, err := json.Marshal(_matcher_order_status)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
