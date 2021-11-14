package requests

import (
	"fmt"
	"github.com/impulse-http/local-backend/pkg/database"
	"github.com/impulse-http/local-backend/pkg/service"
	"log"
	"net/http"
)

type RequestHistoryResponse struct {
	Entries []database.RequestHistoryEntry `json:"history"`
}

func MakeGetHistoryRequestHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, req *http.Request) {
		history, err := s.DB.GetHistory()
		if err != nil {
			log.Println("Couldn't get history" + fmt.Sprint(err))
			service.WriteJSONError(writer, "error while reading from db", 500)
		}
		service.WriteJsonResponse(writer, RequestHistoryResponse{Entries: history})
	}
}
