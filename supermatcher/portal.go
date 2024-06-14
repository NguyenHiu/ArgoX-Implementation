package supermatcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (sm *SuperMatcher) SetupHTTPServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		batch := Batch{}

		if err = json.Unmarshal(data, &batch); err != nil {
			log.Fatal(err)
		}

		if batch.Equal(&Batch{}) {
			fmt.Fprintln(w, "missing param(s)")
		}

		fmt.Println("Got batch: ", batch)
		sm.AddBatch(&batch)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%v", sm.Port), nil); err != nil {
		log.Fatal(err)
	}
}
