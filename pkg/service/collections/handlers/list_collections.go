package handlers

import (
	"github.com/impulse-http/local-backend/pkg/service"
	"net/http"
)

func MakeListCollectionHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		cols, err := s.DB.ListCollections(ctx)
		if err != nil {
			service.WriteJSONError(writer, "Couldn't get collections list", 500)
			return
		}
		service.WriteJsonResponse(writer, cols)
		return
	}
}
