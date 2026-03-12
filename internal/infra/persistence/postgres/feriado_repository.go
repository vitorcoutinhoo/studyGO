package postgres

import (
	"context"
	"fmt"
	"time"

	"plantao/internal/domain/financeiro"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FeriadoRepository struct {
	pool *pgxpool.Pool
}

func NewFeriadoRepository(pool *pgxpool.Pool) *FeriadoRepository {
	return &FeriadoRepository{pool: pool}
}

func (r *FeriadoRepository) FindById(ctx context.Context, id uuid.UUID) (*financeiro.Feriado, error) {
	query := `
		SELECT id, data, nome, descricao, created_at
		FROM feriados
		WHERE id = $1
	`

	var f financeiro.Feriado
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&f.Id,
		&f.Data,
		&f.Nome,
		&f.Descricao,
		&f.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find feriado: %w", err)
	}

	return &f, nil
}

func (r *FeriadoRepository) FindByAno(ctx context.Context, ano int) ([]financeiro.Feriado, error) {
	query := `
		SELECT id, data, nome, descricao, created_at
		FROM feriados
		WHERE EXTRACT(YEAR FROM data) = $1
		ORDER BY data ASC
	`

	rows, err := r.pool.Query(ctx, query, ano)
	if err != nil {
		return nil, fmt.Errorf("failed to find feriados: %w", err)
	}
	defer rows.Close()

	var feriados []financeiro.Feriado
	for rows.Next() {
		var f financeiro.Feriado
		if err := rows.Scan(&f.Id, &f.Data, &f.Nome, &f.Descricao, &f.CreatedAt); err != nil {
			return nil, err
		}
		feriados = append(feriados, f)
	}

	return feriados, nil
}

func (r *FeriadoRepository) UpdateData(ctx context.Context, id uuid.UUID, novaData time.Time) error {
	query := `UPDATE feriados SET data = $2 WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id, novaData)
	if err != nil {
		return fmt.Errorf("failed to update feriado data: %w", err)
	}

	return nil
}
