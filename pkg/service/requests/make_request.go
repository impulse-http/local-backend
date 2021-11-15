package requests

import (
	"bytes"
	"encoding/json"
	"github.com/impulse-http/local-backend/pkg"
	"github.com/impulse-http/local-backend/pkg/service"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type CallRequest struct {
	Request pkg.RequestType `json:"request"`
}

type CallResponse struct {
	Id       int64            `json:"id"`
	Response pkg.ResponseType `json:"response"`
}

func MakeRequestSendHandler(s *service.Service) service.Handler {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			service.WriteJSONError(w, "failed to read body", 500)
			return
		}
		log.Println(string(body))
		var requestJson CallRequest
		if err := json.Unmarshal(body, &requestJson); err != nil {
			log.Println(err)
			service.WriteJSONError(w, "error while unmarshall json", 500)
			return
		}
		method := strings.ToUpper(requestJson.Request.Method)

		bodyReader := bytes.NewBufferString(requestJson.Request.Body)
		userRequest, err := http.NewRequest(method, requestJson.Request.Url, bodyReader)
		if err != nil {
			log.Println(err)
			service.WriteJSONError(w, "error while creating request", 500)
			return
		}
		client := &http.Client{}
		response, err := client.Do(userRequest)
		if err != nil {
			log.Println(err)
			service.WriteJSONError(w, "error while making request", 500)
			return
		}

		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(response.Body); err != nil {
			return
		}
		res := pkg.ResponseType{
			Headers: response.Header,
			Body:    buf.String(),
		}

		id, err := s.DB.CreateHistoryEntry(&requestJson.Request, &res)
		if err != nil {
			log.Println(err)
			service.WriteJSONError(w, "error while inserting to database", 500)
		}

		service.WriteJsonResponse(
			w,
			CallResponse{
				Id:       id,
				Response: res,
			},
		)
	}
}
