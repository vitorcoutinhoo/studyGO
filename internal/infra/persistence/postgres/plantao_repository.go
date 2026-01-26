package postgres

import (
	"context"
	"fmt"
	"plantao/internal/domain/plantao"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PlantaoRepository struct {
	pool *pgxpool.Pool
}

func NewPlantaoRepository(pool *pgxpool.Pool) *PlantaoRepository {
	return &PlantaoRepository{pool: pool}
}

func (r *PlantaoRepository) Store(ctx context.Context, plantao *plantao.Plantao) error {
	query := `
		INSERT INTO plantoes (id, colaborador_id, inicio, fim, status)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(ctx, query,
		plantao.Id,
		plantao.ColaboradorId,
		plantao.Periodo.Inicio,
		plantao.Periodo.Fim,
		plantao.Status,
	)

	if err != nil {
		return fmt.Errorf("failed to store plantao: %w", err)
	}

	return nil
}

func (r *PlantaoRepository) Update(ctx context.Context, plantao *plantao.Plantao) error {
	query := `
		UPDATE plantoes
		SET colaborador_id = $1, inicio = $2, fim = $3, status = $4
		WHERE id = $5
	`

	_, err := r.pool.Exec(ctx, query,
		plantao.Id,
		plantao.ColaboradorId,
		plantao.Periodo.Inicio,
		plantao.Periodo.Fim,
		plantao.Status,
	)

	if err != nil {
		return fmt.Errorf("failed to update plantao: %w", err)
	}

	return nil
}

func (r *PlantaoRepository) Delete(ctx context.Context, plantaoId string) error {
	query := `
		DELETE FROM plantoes
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, plantaoId)
	if err != nil {
		return fmt.Errorf("failed to delete plantao: %w", err)
	}

	return nil
}

func (r *PlantaoRepository) FindById(ctx context.Context, plantaoId string) (*plantao.Plantao, error) {
	query := `
		SELECT id, colaborador_id, inicio, fim, status
		FROM plantoes
		WHERE id = $1
	`
	var p plantao.Plantao
	row := r.pool.QueryRow(ctx, query, plantaoId)

	err := row.Scan(
		&p.Id,
		&p.ColaboradorId,
		&p.Periodo.Inicio,
		&p.Periodo.Fim,
		&p.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan plantao: %w", err)
	}

	return &p, nil
}

func (r *PlantaoRepository) Find(ctx context.Context, filtro *plantao.Filtro) ([]*plantao.Plantao, error) {
	query := `
		SELECT id, colaborador_id, inicio, fim, status
		FROM plantoes
		WHERE ($1 is null or colaborador_id = $1)
		AND ($2 is null or (inicio >= $2 and fim <= $3))
		AND ($4 is null or status = $4)
		LIMIT $5 OFFSET $6
	`

	row, err := r.pool.Query(ctx, query,
		filtro.ColaboradorID,
		filtro.Periodo.Inicio,
		filtro.Periodo.Fim,
		filtro.Status,
		filtro.Limit,
		filtro.Offset,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find plantao: %w", err)
	}

	defer row.Close()
	var plantoes []*plantao.Plantao

	for row.Next() {
		var p plantao.Plantao
		err := row.Scan(
			&p.Id,
			&p.ColaboradorId,
			&p.Periodo.Inicio,
			&p.Periodo.Fim,
			&p.Status,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan plantao: %w", err)
		}

		plantoes = append(plantoes, &p)
	}

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return plantoes, nil
}
