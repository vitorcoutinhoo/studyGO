package postgres

import (
	"context"
	"errors"
	"fmt"
	"plantao/internal/domain/cargo"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CargoRepository struct {
	pool *pgxpool.Pool
}

func NewCargoRepository(pool *pgxpool.Pool) *CargoRepository {
	return &CargoRepository{pool: pool}
}

func (r *CargoRepository) Store(ctx context.Context, nome string) (*cargo.Cargo, error) {
	query := `INSERT INTO cargos (nome) VALUES ($1) RETURNING id, nome`

	var c cargo.Cargo
	err := r.pool.QueryRow(ctx, query, nome).Scan(&c.Id, &c.Nome)
	if err != nil {
		return nil, fmt.Errorf("erro ao inserir cargo: %w", err)
	}

	return &c, nil
}

func (r *CargoRepository) FindAll(ctx context.Context) ([]cargo.Cargo, error) {
	query := `SELECT id, nome FROM cargos WHERE ativo = 'Y' ORDER BY nome`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cargos: %w", err)
	}
	defer rows.Close()

	var cargos []cargo.Cargo
	for rows.Next() {
		var c cargo.Cargo
		if err := rows.Scan(&c.Id, &c.Nome); err != nil {
			return nil, fmt.Errorf("erro ao escanear cargo: %w", err)
		}
		cargos = append(cargos, c)
	}

	return cargos, nil
}

func (r *CargoRepository) FindById(ctx context.Context, id uuid.UUID) (*cargo.Cargo, error) {
	query := `SELECT id, nome FROM cargos WHERE id = $1 AND ativo = 'Y'`

	var c cargo.Cargo
	err := r.pool.QueryRow(ctx, query, id).Scan(&c.Id, &c.Nome)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, cargo.ErrorCargoNotFound
		}
		return nil, fmt.Errorf("erro ao buscar cargo: %w", err)
	}

	return &c, nil
}

func (r *CargoRepository) ExistsNome(ctx context.Context, nome string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM cargos WHERE nome ILIKE $1 AND ativo = 'Y')`

	var exists bool
	err := r.pool.QueryRow(ctx, query, nome).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar cargo: %w", err)
	}

	return exists, nil
}

func (r *CargoRepository) Disable(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE cargos SET ativo = 'N' WHERE id = $1 AND ativo = 'Y'`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao desativar cargo: %w", err)
	}

	if result.RowsAffected() == 0 {
		return cargo.ErrorCargoNotFound
	}

	return nil
}
