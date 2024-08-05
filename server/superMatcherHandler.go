package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) sm_getBatchesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type __batch struct {
		BatchID                uuid.UUID
		Orders                 []uuid.UUID
		IsValidmatcher         bool
		IsValidSignature       bool
		IsValidOrdersSignature []bool
		DuplicateOrders        []uuid.UUID
		IsSent                 bool
	}
	type __data struct {
		Batches []__batch
	}
	data := __data{}
	for _, batch := range s.superMatcher.BatchResults {
		_validOrders := []bool{}
		for _, order := range batch.OrdersStatus {
			_validOrders = append(_validOrders, order.IsValidSignature)
		}
		data.Batches = append(data.Batches, __batch{
			BatchID:                batch.BatchID,
			Orders:                 batch.Orders,
			IsValidmatcher:         batch.BatchStatus.IsValidMatcher,
			IsValidSignature:       batch.BatchStatus.IsValidSignature,
			IsValidOrdersSignature: _validOrders,
			DuplicateOrders:        batch.DuplicateOrders,
			IsSent:                 s.superMatcher.IsSent[batch.BatchID],
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

func (s *Server) sm_sendBatchesHandler(w http.ResponseWriter, r *http.Request) {
	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.superMatcher.SendBatchDemo()

	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Successfully!",
	})
	w.Write(jsonData)
}
