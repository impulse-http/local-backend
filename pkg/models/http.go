package models

type NewRequestRequest struct {
	Name    string      `json:"name"`
	Request RequestType `json:"request"`
}

type DeleteRequestRequest struct {
	Id int64 `json:"id"`
}

type NewRequestResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
