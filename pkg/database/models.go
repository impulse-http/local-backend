package database

import (
	"github.com/impulse-http/local-backend/pkg"
	"time"
)

type RequestHistoryEntry struct {
	Id        int64            `json:"id"`
	CreatedAt time.Time        `json:"created_at"`
	Request   pkg.RequestType  `json:"request"`
	Response  pkg.ResponseType `json:"response"`
}
