package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func ReadJsonRequest(req *http.Request, i interface{}) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print(err)
		return err
	}
	if err = json.Unmarshal(body, i); err != nil {
		log.Print(err)
		return err
	}
	return nil
}
