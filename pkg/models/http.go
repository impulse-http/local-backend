package models

type NewRequestRequest struct {
	Name    string      `json:"name"`
	Request RequestType `json:"request"`
}

type NewRequestResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
