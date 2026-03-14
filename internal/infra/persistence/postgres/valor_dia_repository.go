package postgres

import (
	"context"
	"fmt"
	"time"

	"plantao/internal/domain/financeiro"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ValorDiaRepository struct {
	pool *pgxpool.Pool
}

func NewValorDiaRepository(pool *pgxpool.Pool) *ValorDiaRepository {
	return &ValorDiaRepository{pool: pool}
}

func (r *ValorDiaRepository) FindVigentes(ctx context.Context) ([]financeiro.ValorDia, error) {
	query := `
		SELECT id, tipo_dia, valor, vigencia_inicio, vigencia_fim
		FROM config_valores_dia
		WHERE vigencia_fim IS NULL
		ORDER BY tipo_dia ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find valores vigentes: %w", err)
	}
	defer rows.Close()

	var valores []financeiro.ValorDia
	for rows.Next() {
		var v financeiro.ValorDia
		if err := rows.Scan(&v.Id, &v.TipoDia, &v.Valor, &v.VigenciaInicio, &v.VigenciaFim); err != nil {
			return nil, err
		}
		valores = append(valores, v)
	}

	return valores, nil
}

func (r *ValorDiaRepository) FindVigenteByTipoDia(ctx context.Context, tipoDia financeiro.TipoDia) (*financeiro.ValorDia, error) {
	query := `
		SELECT id, tipo_dia, valor, vigencia_inicio, vigencia_fim
		FROM config_valores_dia
		WHERE tipo_dia = $1 AND vigencia_fim IS NULL
	`

	var v financeiro.ValorDia
	err := r.pool.QueryRow(ctx, query, tipoDia).Scan(
		&v.Id, &v.TipoDia, &v.Valor, &v.VigenciaInicio, &v.VigenciaFim,
	)
	if err != nil {
		// nenhum registro vigente — não é erro
		return nil, nil
	}

	return &v, nil
}

func (r *ValorDiaRepository) FindVigenteByData(ctx context.Context, data time.Time) (map[financeiro.TipoDia]float64, error) {
	query := `
		SELECT tipo_dia, valor
		FROM config_valores_dia
		WHERE vigencia_inicio <= $1
		  AND (vigencia_fim IS NULL OR vigencia_fim >= $1)
	`

	rows, err := r.pool.Query(ctx, query, data)
	if err != nil {
		return nil, fmt.Errorf("failed to find valores dia: %w", err)
	}
	defer rows.Close()

	valores := make(map[financeiro.TipoDia]float64)
	for rows.Next() {
		var tipoDia string
		var valor float64
		if err := rows.Scan(&tipoDia, &valor); err != nil {
			return nil, err
		}
		valores[financeiro.TipoDia(tipoDia)] = valor
	}

	return valores, nil
}

func (r *ValorDiaRepository) Store(ctx context.Context, v *financeiro.ValorDia) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO config_valores_dia (id, tipo_dia, valor, vigencia_inicio)
		 VALUES ($1, $2, $3, $4)`,
		v.Id, v.TipoDia, v.Valor, v.VigenciaInicio,
	)
	if err != nil {
		return fmt.Errorf("failed to store valor dia: %w", err)
	}
	return nil
}

func (r *ValorDiaRepository) CloseVigencia(ctx context.Context, id uuid.UUID, vigenciaFim time.Time) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE config_valores_dia SET vigencia_fim = $2, updated_at = NOW() WHERE id = $1`,
		id, vigenciaFim,
	)
	if err != nil {
		return fmt.Errorf("failed to close vigencia: %w", err)
	}
	return nil
}
