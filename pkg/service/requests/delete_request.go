package requests

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/impulse-http/local-backend/pkg/models"
	"github.com/impulse-http/local-backend/pkg/service"
)

func MakeDeleteRequestHandler(s *service.Service) service.Handler {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			service.WriteJSONError(w, "failed to read body", 500)
			return
		}

		var requestJson models.DeleteRequestRequest
		if err := json.Unmarshal(body, &requestJson); err != nil {
			log.Println(err)
			service.WriteJSONError(w, "error while unmarshall json", 500)
			return
		}

		if err := s.DB.DeleteRequest(ctx, requestJson.Id); err != nil {
			log.Println(err)
			service.WriteJSONError(w, "error while querying db", 500)
			return
		}
		w.WriteHeader(http.StatusOK)
		service.WriteJsonResponse(w, "")
	}
}
