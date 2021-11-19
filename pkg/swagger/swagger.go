package swagger

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/go-openapi/spec"
	"github.com/impulse-http/local-backend/pkg/models"
)

func ParseFromCollection() {

	data := []byte(`{
  "person": {
    "name": {
      "first": "Leonid",
      "last": "Bugaev",
      "fullName": "Leonid Bugaev"
    },
    "github": {
      "handle": "buger",
      "followers": 109
    },
    "avatars": [
      { "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
    ]
  },
  "company": {
    "name": "Acme"
  }
}`)

	swagger := spec.Swagger{
		VendorExtensible: spec.VendorExtensible{},
		SwaggerProps:     spec.SwaggerProps{Paths: &spec.Paths{Paths: map[string]spec.PathItem{}}},
	}

	request := []models.RequestType{
		{
			Url:    "http://google.com",
			Method: "POST",
			Body:   "{\"test\": \"body\"}",
		},
	}

	for _, r := range request {
		pathItem := spec.PathItem{
			PathItemProps: spec.PathItemProps{Post: &spec.Operation{
				OperationProps: spec.OperationProps{Description: "test"},
			}},
		}
		jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			fmt.Printf("Key: '%s'\n Value: '%s'\n Type: %s\n", string(key), string(value), dataType)
			return nil
		}, "person", "name")

		swagger.Paths.Paths[r.Url] = pathItem
	}
}
