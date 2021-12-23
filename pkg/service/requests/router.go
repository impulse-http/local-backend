package requests

import (
	"github.com/gorilla/mux"
	"github.com/impulse-http/local-backend/pkg/service"
	"github.com/impulse-http/local-backend/pkg/service/requests/handlers"
)

func AddRequestHandlers(s *service.Service, r *mux.Router) {
	r.HandleFunc("/requests/bare/make/", handlers.MakeRequestSendHandler(s)).Methods("POST")
	r.HandleFunc("/requests/{id:[0-9]+}/", handlers.MakeGetRequestHandler(s)).Methods("GET")
	r.HandleFunc("/requests/{id:[0-9]+}/", handlers.MakeUpdateRequestHandler(s)).Methods("PUT")
	r.HandleFunc("/requests/{id:[0-9]+}/", handlers.MakeDeleteRequestHandler(s)).Methods("DELETE")
	r.HandleFunc("/requests/history/", handlers.MakeGetHistoryRequestHandler(s)).Methods("GET")
	r.HandleFunc("/collections/{id:[0-9]+}/requests/", handlers.MakeListRequestHandler(s)).Methods("GET")
	r.HandleFunc("/collections/{id:[0-9]+}/requests/", handlers.MakeNewRequestHandler(s)).Methods("POST")
}
