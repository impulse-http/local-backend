package requests

import (
	"github.com/impulse-http/local-backend/pkg/service"
	"net/http"
)

type GetRequestResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func MakeGetRequestHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, req *http.Request) {

	}
}
