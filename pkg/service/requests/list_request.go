package requests

import (
	"fmt"
	"log"
	"net/http"

	"github.com/impulse-http/local-backend/pkg/models"
	"github.com/impulse-http/local-backend/pkg/service"
)

type ResponseListRequests struct {
	Entries []*models.NewRequestRequest
}

func MakeListRequestHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		requests, err := s.DB.GetListRequests(ctx)
		if err != nil {
			log.Println("Couldn't get list requests" + fmt.Sprint(err))
			service.WriteJSONError(writer, "error while reading from db", 500)
			return
		}
		service.WriteJsonResponse(writer, ResponseListRequests{Entries: requests})
	}
}
