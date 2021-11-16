package handlers

import (
	"github.com/gorilla/mux"
	"github.com/impulse-http/local-backend/pkg/service"
	"net/http"
	"strconv"
)

func MakeGetRequestHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		vars := mux.Vars(req)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			service.WriteJSONError(writer, "Couldn't parse url parameter id", 500)
			return
		}
		response, err := s.DB.GetRequest(ctx, id)
		if err != nil {
			service.WriteJSONError(writer, "Couldn't get request", 500)
			return
		}
		service.WriteJsonResponse(writer, response)
		return
	}
}
