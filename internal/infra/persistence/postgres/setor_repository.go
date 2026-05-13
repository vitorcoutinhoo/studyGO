package postgres

import (
	"context"
	"errors"
	"fmt"
	"plantao/internal/domain/setor"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SetorRepository struct {
	pool *pgxpool.Pool
}

func NewSetorRepository(pool *pgxpool.Pool) *SetorRepository {
	return &SetorRepository{pool: pool}
}

func (r *SetorRepository) Store(ctx context.Context, nome string) (*setor.Setor, error) {
	query := `INSERT INTO setores (nome) VALUES ($1) RETURNING id, nome`

	var s setor.Setor
	err := r.pool.QueryRow(ctx, query, nome).Scan(&s.Id, &s.Nome)
	if err != nil {
		return nil, fmt.Errorf("erro ao inserir setor: %w", err)
	}

	return &s, nil
}

func (r *SetorRepository) FindAll(ctx context.Context) ([]setor.Setor, error) {
	query := `SELECT id, nome FROM setores WHERE ativo = 'Y' ORDER BY nome`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar setores: %w", err)
	}
	defer rows.Close()

	var setores []setor.Setor
	for rows.Next() {
		var s setor.Setor
		if err := rows.Scan(&s.Id, &s.Nome); err != nil {
			return nil, fmt.Errorf("erro ao escanear setor: %w", err)
		}
		setores = append(setores, s)
	}

	return setores, nil
}

func (r *SetorRepository) FindById(ctx context.Context, id uuid.UUID) (*setor.Setor, error) {
	query := `SELECT id, nome FROM setores WHERE id = $1 AND ativo = 'Y'`

	var s setor.Setor
	err := r.pool.QueryRow(ctx, query, id).Scan(&s.Id, &s.Nome)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, setor.ErrorSetorNotFound
		}
		return nil, fmt.Errorf("erro ao buscar setor: %w", err)
	}

	return &s, nil
}

func (r *SetorRepository) ExistsNome(ctx context.Context, nome string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM setores WHERE nome ILIKE $1 AND ativo = 'Y')`

	var exists bool
	err := r.pool.QueryRow(ctx, query, nome).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar setor: %w", err)
	}

	return exists, nil
}

func (r *SetorRepository) Disable(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE setores SET ativo = 'N' WHERE id = $1 AND ativo = 'Y'`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao desativar setor: %w", err)
	}

	if result.RowsAffected() == 0 {
		return setor.ErrorSetorNotFound
	}

	return nil
}
