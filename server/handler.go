package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

func (s *Server) getUserDataHandler(w http.ResponseWriter, r *http.Request) {
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

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	w.Write(jsonData)
}

func (s *Server) getMatcherDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type _matcher_order struct {
		Price  int
		Amount int
		Side   bool
	}
	type _data struct {
		Name      string
		ID        uuid.UUID
		Address   common.Address
		BidOrders []_matcher_order
		AskOrders []_matcher_order
	}
	data := []_data{}
	for idx, matcher := range s.matchers {
		_bOrder := []_matcher_order{}
		_aOrder := []_matcher_order{}
		for _, order := range matcher.BidOrders {
			_bOrder = append(_bOrder, _matcher_order{
				Price:  int(order.Data.Price.Int64()),
				Amount: int(order.Data.Amount.Int64()),
				Side:   order.Data.Side,
			})
		}
		for _, order := range matcher.AskOrders {
			_aOrder = append(_aOrder, _matcher_order{
				Price:  int(order.Data.Price.Int64()),
				Amount: int(order.Data.Amount.Int64()),
				Side:   order.Data.Side,
			})
		}
		data = append(data, _data{
			Name:      fmt.Sprintf("Matcher %v", idx+1),
			ID:        matcher.ID,
			Address:   matcher.Address,
			BidOrders: _bOrder,
			AskOrders: _aOrder,
		})
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	w.Write(jsonData)
}

func (s *Server) sendOrderAndEndHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Access form data values
	// orderID := r.Form.Get("order_id")
	// quantity := r.Form.Get("quantity")

	// Process the form data
	// ...

	fmt.Fprintf(w, "Order received and processed successfully")
}
