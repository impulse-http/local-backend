package handlers

import (
	"github.com/gorilla/mux"
	"github.com/impulse-http/local-backend/pkg/models"
	"github.com/impulse-http/local-backend/pkg/service"
	"net/http"
	"strconv"
)

func MakeUpdateRequestHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		ctx := request.Context()
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			service.WriteJSONError(writer, "Couldn't parse request id from url", 500)
			return
		}
		req := &models.Request{}
		if err = service.ReadJsonRequest(request, req); err != nil {
			service.WriteJSONError(writer, "Coulnd't parse json body", 500)
			return
		}
		res, err := s.DB.UpdateRequest(ctx, id, req)
		if err != nil {
			service.WriteJSONError(writer, "Couldn't update request", 500)
			return
		}
		service.WriteJsonResponse(writer, res)
		return
	}
}
