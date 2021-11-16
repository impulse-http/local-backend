package collections

import (
	"github.com/gorilla/mux"
	"github.com/impulse-http/local-backend/pkg/service"
	"github.com/impulse-http/local-backend/pkg/service/collections/handlers"
)

func AddCollectionsHandlers(s *service.Service, r *mux.Router) {
	r.HandleFunc("/collections/", handlers.MakeListCollectionHandler(s)).Methods("Get")
	r.HandleFunc("/collections/", handlers.MakeCreateCollectionHandler(s)).Methods("Post")
	r.HandleFunc("/collections/{id}/", handlers.MakeUpdateCollectionHandler(s)).Methods("Put")
	r.HandleFunc("/collections/{id}/", handlers.MakeDeleteCollectionHandler(s)).Methods("Delete")
}
