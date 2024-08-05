package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NguyenHiu/lightning-exchange/listener"
	"github.com/NguyenHiu/lightning-exchange/matcher"
	"github.com/NguyenHiu/lightning-exchange/reporter"
	"github.com/NguyenHiu/lightning-exchange/supermatcher"
	"github.com/NguyenHiu/lightning-exchange/user"
	"github.com/NguyenHiu/lightning-exchange/worker"
)

type Server struct {
	port int

	user         *user.User
	matchers     []*matcher.Matcher
	superMatcher *supermatcher.SuperMatcher
	listener     *listener.Listener
	reporter     *reporter.Reporter
	worker       *worker.Worker
}

func (s *Server) Start() {
	// http.HandleFunc("/user/get-profile", s.getUserProfileHandler)
	http.HandleFunc("/user/send-order", s.sendOrderHandler)
	http.HandleFunc("/user/get-orders", s.getOrdersHandler)

	// http.HandleFunc("/matcher/get-matchers-data", s.getMatcherProfileHandler)
	http.HandleFunc("/matcher/get-local-order-book", s.getLocalOrderBookHandler)
	http.HandleFunc("/matcher/get-batches", s.getBatchesHandler)
	http.HandleFunc("/matcher/matching", s.matchingHandler)
	http.HandleFunc("/matcher/batching", s.batchingHandler)
	http.HandleFunc("/matcher/send-batches", s.sendBatchesHandler)
	http.HandleFunc("/matcher/send-batch-details", s.sendBatchDetailsHandler)

	http.HandleFunc("/super-matcher/get-batches", s.sm_getBatchesHandler)
	http.HandleFunc("/super-matcher/send-batches", s.sm_sendBatchesHandler)

	http.HandleFunc("/searcher/get-batch-book", s.getBatchBook)
	http.HandleFunc("/searcher/match-batches", s.matchBatches)

	http.HandleFunc("/reporter/get-matched-batches", s.getMatchedBatches)
	http.HandleFunc("/reporter/report-batch", s.reportBatch)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Server listening on port %d", s.port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func NewServer(port int,
	user *user.User,
	matchers []*matcher.Matcher,
	superMatcher *supermatcher.SuperMatcher,
	reporter *reporter.Reporter,
	listener *listener.Listener,
	worker *worker.Worker,
) *Server {
	return &Server{
		port:         port,
		user:         user,
		matchers:     matchers,
		superMatcher: superMatcher,
		listener:     listener,
		reporter:     reporter,
		worker:       worker,
	}
}

// func main() {
// 	server := NewServer(7000)
// 	server.Start()
// }
