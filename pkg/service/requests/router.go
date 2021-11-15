package requests

import (
	"github.com/gorilla/mux"
	"github.com/impulse-http/local-backend/pkg/service"
)

func AddRequestHandlers(s *service.Service, r *mux.Router) {
	r.HandleFunc("/requests/bare/make/", MakeRequestSendHandler(s)).Methods("POST")
	r.HandleFunc("/requests/", MakeNewRequestHandler(s)).Methods("POST")
	r.HandleFunc("/requests/{id:[0-9]+}/", MakeGetRequestHandler(s)).Methods("GET")
	r.HandleFunc("/requests/{id:[0-9]+}/", MakeUpdateRequestHandler(s)).Methods("PUT")
	r.HandleFunc("/requests/history/", MakeGetHistoryRequestHandler(s)).Methods("GET")
}
