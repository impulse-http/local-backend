package requests

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/impulse-http/local-backend/pkg/models"
	"github.com/impulse-http/local-backend/pkg/service"
)

func MakeNewRequestHandler(s *service.Service) service.Handler {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			service.WriteJSONError(w, "failed to read body", 500)
			return
		}

		var requestJson models.NewRequestRequest
		if err := json.Unmarshal(body, &requestJson); err != nil {
			log.Println(err)
			service.WriteJSONError(w, "error while unmarshall json", 500)
			return
		}

		requestId, err := s.DB.CreateRequest(ctx, &requestJson)
		if err != nil {
			log.Println(err)
			service.WriteJSONError(w, "error while creating request "+err.Error(), 500)
			return
		}
		service.WriteJsonResponse(w, models.NewRequestResponse{
			Id:   int(requestId),
			Name: requestJson.Name,
		})
	}
}
