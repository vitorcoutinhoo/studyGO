package postgres

import (
	"context"
	"fmt"
	"plantao/internal/infra/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(config *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), config.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	return pool, nil
}
