package models

type Request struct {
	Name         string      `json:"name"`
	Request      RequestType `json:"request"`
	CollectionId int64       `json:"collection_id"`
}

type StoredRequest struct {
	Id      int         `json:"id"`
	Name    string      `json:"name"`
	Request RequestType `json:"request"`
}

type Collection struct {
	Name string `json:"name"`
}

type StoredCollection struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
