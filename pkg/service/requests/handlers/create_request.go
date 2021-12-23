package handlers

import (
	"github.com/gorilla/mux"
	"github.com/impulse-http/local-backend/pkg/models"
	"github.com/impulse-http/local-backend/pkg/service"
	"log"
	"net/http"
	"strconv"
)

func MakeNewRequestHandler(s *service.Service) service.Handler {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		var colId int64 = 0
		sId := vars["id"]

		if sId != "" {
			var err error
			colId, err = strconv.ParseInt(vars["id"], 10, 64)
			if err != nil {
				service.WriteJSONError(w, "Couldn't parse url parameter id", 500)
				return
			}
		}

		ctx := req.Context()
		request := &models.Request{}
		if err := service.ReadJsonRequest(req, request); err != nil {
			log.Println(err)
			service.WriteJSONError(w, "failed to read body", 500)
			return
		}

		requestId, err := s.DB.CreateRequest(ctx, request, colId)
		if err != nil {
			log.Println(err)
			service.WriteJSONError(w, "error while creating request "+err.Error(), 500)
			return
		}
		service.WriteJsonResponse(w, models.StoredRequest{
			Id:      int(requestId),
			Name:    request.Name,
			Request: request.Request,
		})
	}
}
