package matcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (m *Matcher) SendBatch(batch *Batch) {
	m.Batches[batch.BatchID] = batch

	jsonData, err := json.Marshal(batch)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post(fmt.Sprintf("%v/batch/submit", m.SuperMatcherURI), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error when sending batches to Super Matcher, err: ", err)
	}

	fmt.Printf("SendBatch, res: %v\n", res)
}
