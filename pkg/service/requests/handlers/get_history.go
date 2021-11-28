package handlers

import (
	"fmt"
	"github.com/impulse-http/local-backend/pkg/models"
	"github.com/impulse-http/local-backend/pkg/service"
	"log"
	"net/http"
)

type RequestHistoryResponse struct {
	Entries []models.RequestHistoryEntry `json:"history"`
}

func MakeGetHistoryRequestHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		history, err := s.DB.GetHistory(ctx)
		if err != nil {
			log.Println("Couldn't get history" + fmt.Sprint(err))
			service.WriteJSONError(writer, "error while reading from db", 500)
		}
		service.WriteJsonResponse(writer, RequestHistoryResponse{Entries: history})
	}
}
