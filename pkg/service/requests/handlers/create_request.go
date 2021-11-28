package handlers

import (
	"github.com/impulse-http/local-backend/pkg/models"
	"github.com/impulse-http/local-backend/pkg/service"
	"log"
	"net/http"
)

func MakeNewRequestHandler(s *service.Service) service.Handler {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		request := &models.Request{}
		if err := service.ReadJsonRequest(req, request); err != nil {
			log.Println(err)
			service.WriteJSONError(w, "failed to read body", 500)
			return
		}

		requestId, err := s.DB.CreateRequest(ctx, request)
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
