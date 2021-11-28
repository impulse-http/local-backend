package handlers

import (
	"github.com/gorilla/mux"
	"github.com/impulse-http/local-backend/pkg/service"
	"net/http"
	"strconv"
)

func MakeDeleteCollectionHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		id, err := strconv.ParseInt(params["id"], 10, 64)
		if err != nil {
			service.WriteJSONError(writer, "Can't parse url param", 500)
			return
		}
		ctx := req.Context()
		if err := s.DB.DeleteCollection(ctx, id); err != nil {
			service.WriteJSONError(writer, "Can't delete from db", 500)
			return
		}
		writer.WriteHeader(200)
		return
	}
}
