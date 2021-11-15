package service

import (
	"encoding/json"
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter, d interface{}) {
	data, err := json.Marshal(d)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		return
	}
}
