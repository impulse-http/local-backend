package handlers

import (
	"github.com/gorilla/mux"
	"github.com/impulse-http/local-backend/pkg/service"
	"net/http"
	"strconv"
)

func MakeListRequestHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		vars := mux.Vars(req)

		var id int64 = 0
		sId := vars["id"]

		if sId != "" {
			var err error
			id, err = strconv.ParseInt(vars["id"], 10, 64)
			if err != nil {
				service.WriteJSONError(writer, "Couldn't parse url parameter id", 500)
				return
			}
		}

		response, err := s.DB.GetRequests(ctx, id)
		if err != nil {
			service.WriteJSONError(writer, "Couldn't get requests", 500)
			return
		}
		service.WriteJsonResponse(writer, response)
		return
	}
}
