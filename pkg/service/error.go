package service

import (
	"encoding/json"
	"net/http"
)

type ErrorMessage struct {
	Message string            `json:"message"`
	Details map[string]string `json:"details"`
}

func WriteJSONError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(ErrorMessage{Message: message})
	if err != nil {

	}
	w.Write(data)
}
