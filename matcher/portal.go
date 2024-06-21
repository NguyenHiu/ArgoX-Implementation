package matcher

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (m *Matcher) SendBatch(batch *Batch) {
	m.Batches[batch.BatchID] = batch

	data, err := batch.Encode_TranferBatching()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/batch", m.SuperMatcherURI), bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error when sending batches to Super Matcher, err: ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		_logger.Error("[%v] SendBatch got error, code: %v\n", m.ID.String()[:6], res.StatusCode)
	}
}
