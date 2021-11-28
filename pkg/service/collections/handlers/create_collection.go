package handlers

import (
	"github.com/impulse-http/local-backend/pkg/models"
	"github.com/impulse-http/local-backend/pkg/service"
	"log"
	"net/http"
)

func MakeCreateCollectionHandler(s *service.Service) service.Handler {
	return func(writer http.ResponseWriter, req *http.Request) {
		collection := models.Collection{}
		err := service.ReadJsonRequest(req, &collection)
		if err != nil {
			log.Println("Error reading body" + err.Error())
			service.WriteJSONError(writer, "Couldn't read request body", 500)
			return
		}
		ctx := req.Context()
		c, err := s.DB.CreateCollection(ctx, collection)
		if err != nil {
			service.WriteJSONError(writer, "Error creating new collection", 500)
			return
		}
		service.WriteJsonResponse(writer, c)
	}
}
