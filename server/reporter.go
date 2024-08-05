package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) getMatchedBatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type __batch struct {
		BatchID     uuid.UUID
		MatchedTime int64
	}
	type __data struct {
		Batches []__batch
	}
	data := __data{
		Batches: []__batch{},
	}
	for k, v := range s.reporter.PendingBatches {
		data.Batches = append(data.Batches, __batch{
			BatchID:     k,
			MatchedTime: v,
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

func (s *Server) reportBatch(w http.ResponseWriter, r *http.Request) {
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
		BatchID uuid.UUID
	}
	var b __info
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusInternalServerError)
		return
	}

	s.reporter.ReportABatch(b.BatchID)

	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Successfully!",
	})

	w.Write(jsonData)
}
