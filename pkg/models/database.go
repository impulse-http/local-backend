package models

import (
	"time"
)

type RequestHistoryEntry struct {
	Id        int64        `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	Request   RequestType  `json:"request"`
	Response  ResponseType `json:"response"`
}

type RequestEntry struct {
	Id        int64       `json:"id"`
	Name      string      `json:"name"`
	CreatedAt time.Time   `json:"created_at"`
	Request   RequestType `json:"request"`
}
