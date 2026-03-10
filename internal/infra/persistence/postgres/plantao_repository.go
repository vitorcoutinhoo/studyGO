package postgres

import (
	"context"
	"fmt"
	"plantao/internal/domain/plantao"
	"plantao/internal/domain/shared"
	"strconv"

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
		INSERT INTO plantoes (id, id_colaborador, data_inicio, data_fim, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.pool.Exec(ctx, query,
		plantao.Id,
		plantao.ColaboradorId,
		plantao.Periodo.Inicio,
		plantao.Periodo.Fim,
		strconv.Itoa(int(plantao.Status)),
		plantao.CreatedAt,
		plantao.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to store plantao: %w", err)
	}

	return nil
}

func (r *PlantaoRepository) Update(ctx context.Context, plantao *plantao.Plantao) error {
	query := `
		UPDATE plantoes
		SET id_colaborador = $2, data_inicio = $3, data_fim = $4, status = $5, created_at = $6, updated_at = $7
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		plantao.Id,
		plantao.ColaboradorId,
		plantao.Periodo.Inicio,
		plantao.Periodo.Fim,
		strconv.Itoa(int(plantao.Status)),
		plantao.CreatedAt,
		plantao.UpdatedAt,
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
		SELECT id, id_colaborador, data_inicio, data_fim, status, created_at, updated_at
		FROM plantoes
		WHERE id = $1
	`
	var p plantao.Plantao
	var statusStr string
	row := r.pool.QueryRow(ctx, query, plantaoId)
	p.Periodo = &shared.Periodo{}

	err := row.Scan(
		&p.Id,
		&p.ColaboradorId,
		&p.Periodo.Inicio,
		&p.Periodo.Fim,
		&statusStr,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan plantao: %w", err)
	}

	statusInt, err := strconv.Atoi(statusStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse status: %w", err)
	}
	p.Status = plantao.StatusPlantao(statusInt)

	return &p, nil
}

func (r *PlantaoRepository) Find(
	ctx context.Context,
	filtro *plantao.Filtro,
) ([]plantao.Plantao, error) {

	query := `
		SELECT id, id_colaborador, data_inicio, data_fim, status, created_at, updated_at
		FROM plantoes
		WHERE 1=1
	`

	args := []any{}
	arg := 1

	// Filtro por colaborador
	if filtro != nil && filtro.ColaboradorID != "" {
		query += fmt.Sprintf(" AND id_colaborador = $%d", arg)
		args = append(args, filtro.ColaboradorID)
		arg++
	}

	// Filtro por período
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

	// Filtro por status
	if filtro != nil && filtro.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", arg)
		args = append(args, strconv.Itoa(int(*filtro.Status)))
		arg++
	}

	// Paginação
	if filtro != nil && filtro.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", arg)
		args = append(args, *filtro.Limit)
		arg++
	}

	if filtro != nil && filtro.Offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", arg)
		args = append(args, *filtro.Offset)
		arg++
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plantoes []plantao.Plantao
	for rows.Next() {
		var p plantao.Plantao
		var statusStr string
		p.Periodo = &shared.Periodo{}

		if err := rows.Scan(
			&p.Id,
			&p.ColaboradorId,
			&p.Periodo.Inicio,
			&p.Periodo.Fim,
			&statusStr,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}

		statusInt, err := strconv.Atoi(statusStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse status: %w", err)
		}
		p.Status = plantao.StatusPlantao(statusInt)

		plantoes = append(plantoes, p)
	}

	return plantoes, nil
}
