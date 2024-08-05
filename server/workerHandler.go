package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) getBatchBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type __batch struct {
		BatchID uuid.UUID
		Price   int
		Amount  int
		Side    bool
	}
	type __data struct {
		BidBatches []__batch
		AskBatches []__batch
	}
	data := __data{
		BidBatches: []__batch{},
		AskBatches: []__batch{},
	}
	for _, batch := range s.worker.BidBatches {
		data.BidBatches = append(data.BidBatches, __batch{
			BatchID: batch.BatchID,
			Price:   batch.Price,
			Amount:  batch.Amount,
			Side:    batch.Side,
		})
	}
	for _, batch := range s.worker.AskBatches {
		data.AskBatches = append(data.AskBatches, __batch{
			BatchID: batch.BatchID,
			Price:   batch.Price,
			Amount:  batch.Amount,
			Side:    batch.Side,
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

func (s *Server) matchBatches(w http.ResponseWriter, r *http.Request) {
	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.worker.Matching()

	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Successfully!",
	})
	w.Write(jsonData)
}
