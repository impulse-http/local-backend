package handlers

import (
	"github.com/gorilla/mux"
	"github.com/impulse-http/local-backend/pkg/models"
	"github.com/impulse-http/local-backend/pkg/service"
	"net/http"
	"strconv"
)

func MakeUpdateCollectionHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			service.WriteJSONError(writer, "Couldn't parse request parameters", 500)
			return
		}
		col := models.StoredCollection{}
		if err := service.ReadJsonRequest(req, &col); err != nil {
			service.WriteJSONError(writer, "Couldn't read request body", 500)
			return
		}
		resp, err := s.DB.UpdateCollection(req.Context(), id, col)
		if err != nil {
			service.WriteJSONError(writer, "Couldn't update collection", 500)
			return
		}
		service.WriteJsonResponse(writer, resp)
	}
}
