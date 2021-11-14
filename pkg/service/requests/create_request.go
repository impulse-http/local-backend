package requests

import (
	"github.com/impulse-http/local-backend/pkg"
	"github.com/impulse-http/local-backend/pkg/service"
	"net/http"
)

type NewRequestRequest struct {
	Name    string          `json:"name"`
	Request pkg.RequestType `json:"request"`
}

type NewRequestResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func MakeNewRequestHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
