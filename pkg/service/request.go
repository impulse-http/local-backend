package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ReadJsonRequest(req *http.Request, i interface{}) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, i); err != nil {
		return err
	}
	return nil
}
