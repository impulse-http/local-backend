package requests

import (
	"github.com/impulse-http/local-backend/pkg/models"
	"github.com/impulse-http/local-backend/pkg/service"
	"net/http"
)

type UpdateRequestRequest struct {
	Name    string             `json:"name"`
	Request models.RequestType `json:"request"`
}

type UpdateRequestResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func MakeUpdateRequestHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}