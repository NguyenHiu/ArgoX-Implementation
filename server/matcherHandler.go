package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

func (s *Server) getMatcherProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type __matcher_order struct {
		Price  int
		Amount int
		Side   bool
	}
	type __data struct {
		Name      string
		ID        uuid.UUID
		Address   common.Address
		BidOrders []__matcher_order
		AskOrders []__matcher_order
	}
	data := []__data{}
	for idx, matcher := range s.matchers {
		_bOrder := []__matcher_order{}
		_aOrder := []__matcher_order{}
		for _, order := range matcher.BidOrders {
			_bOrder = append(_bOrder, __matcher_order{
				Price:  int(order.Data.Price.Int64()),
				Amount: int(order.Data.Amount.Int64()),
				Side:   order.Data.Side,
			})
		}
		for _, order := range matcher.AskOrders {
			_aOrder = append(_aOrder, __matcher_order{
				Price:  int(order.Data.Price.Int64()),
				Amount: int(order.Data.Amount.Int64()),
				Side:   order.Data.Side,
			})
		}
		data = append(data, __data{
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

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *Server) getLocalOrderBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type __matcher_order struct {
		ID     uuid.UUID
		Price  int
		Amount int
		Side   bool
	}
	type __data struct {
		MatcherID uuid.UUID
		BidOrders []__matcher_order
		AskOrders []__matcher_order
	}
	data := []__data{}
	for _, matcher := range s.matchers {
		_bOrder := []__matcher_order{}
		_aOrder := []__matcher_order{}
		for _, order := range matcher.BidOrders {
			_bOrder = append(_bOrder, __matcher_order{
				ID:     order.Data.From,
				Price:  int(order.Data.Price.Int64()),
				Amount: int(order.Data.Amount.Int64()),
				Side:   order.Data.Side,
			})
		}
		for _, order := range matcher.AskOrders {
			_aOrder = append(_aOrder, __matcher_order{
				ID:     order.Data.From,
				Price:  int(order.Data.Price.Int64()),
				Amount: int(order.Data.Amount.Int64()),
				Side:   order.Data.Side,
			})
		}
		data = append(data, __data{
			MatcherID: matcher.ID,
			BidOrders: _bOrder,
			AskOrders: _aOrder,
		})
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *Server) matchingHandler(w http.ResponseWriter, r *http.Request) {
	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	s.matchers[0].Matching()

	jsonData, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Successfully!",
	})
	w.Write(jsonData)
}

func (s *Server) batchingHandler(w http.ResponseWriter, r *http.Request) {
	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	type __batching struct {
		NoOrder int
	}
	var b __batching
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusInternalServerError)
		return
	}
	s.matchers[0].Batching(b.NoOrder)

	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Successfully!",
	})
	w.Write(jsonData)

}

func (s *Server) getBatchesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type __order struct {
		OrderID uuid.UUID
		Price   int
		Amount  int
		Side    bool
	}

	type __batch struct {
		BatchID uuid.UUID
		Price   int
		Amount  int
		Side    bool
		Orders  []__order
		Status  string
	}
	type __data struct {
		MatcherID uuid.UUID
		Batches   []__batch
	}
	data := []__data{}
	for _, matcher := range s.matchers {
		batch := []__batch{}
		for _, b := range matcher.BatchesList {
			orders := []__order{}
			for _, ord := range b.Orders {
				orders = append(orders, __order{
					OrderID: ord.OriginalOrder.OrderID,
					Price:   int(ord.ShadowOrder.Price.Int64()),
					Amount:  int(ord.ShadowOrder.Amount.Int64()),
					Side:    ord.OriginalOrder.Side,
				})
			}
			batch = append(batch, __batch{
				BatchID: b.BatchID,
				Price:   int(b.Price.Int64()),
				Amount:  int(b.Amount.Int64()),
				Side:    b.Side,
				Orders:  orders,
				Status:  matcher.BatchStatusMapping[b.BatchID],
			})
		}
		data = append(data, __data{
			MatcherID: matcher.ID,
			Batches:   batch,
		})
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *Server) sendBatchesHandler(w http.ResponseWriter, r *http.Request) {
	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	for _, matcher := range s.matchers {
		for _, batch := range matcher.BatchesList {
			if _isSent, ok := matcher.IsSent[batch.BatchID]; ok && !_isSent {
				matcher.IsSent[batch.BatchID] = true
				matcher.SendBatch(batch)
			} else if !ok {
				matcher.IsSent[batch.BatchID] = true
				matcher.SendBatch(batch)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Successfully!",
	})
	w.Write(jsonData)
}

func (s *Server) sendBatchDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	type __info struct {
		MatcherID uuid.UUID
		BatchID   uuid.UUID
	}
	var b __info = __info{
		MatcherID: uuid.UUID{},
		BatchID:   uuid.UUID{},
	}
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusInternalServerError)
		return
	}

	for _, matcher := range s.matchers {
		if matcher.ID == b.MatcherID {
			matcher.SendBatchDetails(b.BatchID)
			w.Header().Set("Content-Type", "application/json")
			jsonData, _ := json.Marshal(struct {
				Message string `json:"message"`
			}{
				Message: "Successfully!",
			})
			w.Write(jsonData)
		}
	}

}
