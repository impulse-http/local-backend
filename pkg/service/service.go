package service

import (
	"github.com/impulse-http/local-backend/pkg/database"
)

type Service struct {
	DB *database.Database
}

func NewService(db *database.Database) *Service {
	return &Service{DB: db}
}
