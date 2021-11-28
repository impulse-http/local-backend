package handlers

import (
	"bytes"
	"github.com/impulse-http/local-backend/pkg/models"
	"github.com/impulse-http/local-backend/pkg/service"
	"log"
	"net/http"
	"strings"
)

type CallRequest struct {
	Request models.RequestType `json:"request"`
}

type CallResponse struct {
	Id       int64               `json:"id"`
	Response models.ResponseType `json:"response"`
}

// MakeRequestSendHandler create a handler that make a call to api and save it to history
func MakeRequestSendHandler(s *service.Service) service.Handler {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		requestJson := CallRequest{}
		if err := service.ReadJsonRequest(req, &requestJson); err != nil {
			service.WriteJSONError(w, "failed to read body", 500)
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
		res := models.ResponseType{
			Headers: response.Header,
			Body:    buf.String(),
		}

		id, err := s.DB.CreateHistoryEntry(ctx, &requestJson.Request, &res)
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
