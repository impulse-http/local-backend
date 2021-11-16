package handlers

import (
	"github.com/gorilla/mux"
	"github.com/impulse-http/local-backend/pkg/service"
	"net/http"
	"strconv"
)

func MakeDeleteRequestHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		vars := mux.Vars(req)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			service.WriteJSONError(writer, "Couldn't parse url parameter id", 500)
			return
		}
		if err := s.DB.DeleteRequest(ctx, id); err != nil {
			service.WriteJSONError(writer, "Couldn't delete request", 500)
			return
		}
		writer.WriteHeader(200)
	}
}
