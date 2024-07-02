package supermatcher

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/NguyenHiu/lightning-exchange/matcher"
)

func (sm *SuperMatcher) SetupHTTPServer() {
	http.HandleFunc("/batch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		batch := &matcher.Batch{}
		if err := batch.Decode_TransferBatching(data); err != nil {
			log.Fatal(err)
		}

		if batch.Equal(&matcher.Batch{}) {
			fmt.Fprintln(w, "missing param(s)")
		}

		sm.AddBatch(batch)
	})

	_logger.Info("Super Matcher is listening...\n")
	if err := http.ListenAndServe(fmt.Sprintf(":%v", sm.Port), nil); err != nil {
		log.Fatal(err)
	}
}
