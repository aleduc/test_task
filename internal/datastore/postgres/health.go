package postgres

import (
	"context"

	"gorm.io/gorm"
)

type Health struct {
	pgClient *gorm.DB
}

func NewHealth(pgClient *gorm.DB) *Health {
	return &Health{pgClient: pgClient}
}

func (h *Health) Health(ctx context.Context) error {
	db, err := h.pgClient.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}
