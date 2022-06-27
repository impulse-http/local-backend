package database

import (
	"context"
	"database/sql"

	"github.com/impulse-http/local-backend/pkg/models"
)

type DatabaseI interface {
	CreateRequest(context.Context, *models.NewRequestRequest) (int64, error)
	DeleteRequest(context.Context, int64) error
	GetListRequests(context.Context, int64) ([]*models.NewRequestRequest, error)
	GetHistory(context.Context) ([]models.RequestHistoryEntry, error)
	CreateHistoryEntry(context.Context, *models.RequestType, *models.ResponseType) (int64, error)
}

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{db: db}
}
