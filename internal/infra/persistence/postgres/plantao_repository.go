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
	p.Periodo = &plantao.Periodo{}
	
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

func (r *PlantaoRepository) Find(
	ctx context.Context,
	filtro *plantao.Filtro,
) ([]plantao.Plantao, error) {

	query := `
		SELECT id, colaborador_id, inicio, fim
		FROM plantoes
		WHERE 1=1
	`

	args := []any{}
	arg := 1

	// Filtro por colaborador
	if filtro != nil && filtro.ColaboradorID != "" {
		query += fmt.Sprintf(" AND colaborador_id = $%d", arg)
		args = append(args, filtro.ColaboradorID)
		arg++
	}

	// Filtro por período (ESSENCIAL validar Periodo != nil)
	if filtro != nil && filtro.Periodo != nil {
		query += fmt.Sprintf(
			" AND data_inicio >= $%d AND data_fim <= $%d",
			arg,
			arg+1,
		)
		args = append(args,
			filtro.Periodo.Inicio,
			filtro.Periodo.Fim,
		)
		arg += 2
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plantoes []plantao.Plantao
	for rows.Next() {
		var p plantao.Plantao
		p.Periodo = &plantao.Periodo{}

		if err := rows.Scan(
			&p.Id,
			&p.ColaboradorId,
			&p.Periodo.Inicio,
			&p.Periodo.Fim,
		); err != nil {
			return nil, err
		}
		plantoes = append(plantoes, p)
	}

	return plantoes, nil
}
