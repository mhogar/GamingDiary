package oauth

import (
	"fmt"
	"net/http"
)

type tokenServer struct {
	Code   chan string
	server *http.Server
}

func newTokenServer() tokenServer {
	s := tokenServer{
		Code: make(chan string),
	}

	s.server = &http.Server{Addr: ":8080", Handler: s}
	return s
}

func (s tokenServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<html><body><h1>Request Received</h1><p>Return to process...</p></body></html>")

	query := r.URL.Query()
	s.Code <- query["code"][0]
}

func (s tokenServer) Start() {
	go s.server.ListenAndServe()
}

func (s tokenServer) Close() {
	s.server.Close()
}

func (s tokenServer) URL() string {
	return "http://localhost" + s.server.Addr
}
