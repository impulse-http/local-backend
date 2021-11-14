package requests

import (
	"github.com/impulse-http/local-backend/pkg/service"
	"net/http"
)

func MakeDeleteRequestHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, req *http.Request) {

	}
}
