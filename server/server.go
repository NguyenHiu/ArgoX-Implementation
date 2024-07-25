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
	searcher     *worker.Worker
}

func (s *Server) Start() {
	http.HandleFunc("/get-user-data", s.getUserDataHandler)
	http.HandleFunc("/get-matchers-data", s.getMatcherDataHandler)
	http.HandleFunc("/send-order", s.sendOrderAndEndHandler)
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
	searcher *worker.Worker,
) *Server {
	return &Server{
		port:         port,
		user:         user,
		matchers:     matchers,
		superMatcher: superMatcher,
		listener:     listener,
		reporter:     reporter,
		searcher:     searcher,
	}
}

// func main() {
// 	server := NewServer(7000)
// 	server.Start()
// }
