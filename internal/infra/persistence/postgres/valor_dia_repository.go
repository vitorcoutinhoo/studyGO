package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"plantao/internal/domain/financeiro"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *ValorDiaRepository) WithTransaction(ctx context.Context, fn func(financeiro.ValorDiaTransaction) error) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("failed to begin valor dia transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(&valorDiaTransaction{tx: tx}); err != nil {
		return translateValorDiaTransactionError(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return translateValorDiaTransactionError(err)
	}
	return nil
}

func translateValorDiaTransactionError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "55P03", "40001", "40P01":
			return financeiro.ErrorConflitoValorDia
		}
	}
	return err
}

type valorDiaTransaction struct {
	tx pgx.Tx
}

func (t *valorDiaTransaction) LockByTipoDia(ctx context.Context, tipoDia financeiro.TipoDia) (*financeiro.ValorDia, error) {
	const query = `
		SELECT id, tipo_dia, valor, descricao, vigencia_inicio, vigencia_fim, created_at, updated_at
		FROM config_valores_dia
		WHERE tipo_dia = $1
		FOR UPDATE NOWAIT
	`
	var valor financeiro.ValorDia
	err := t.tx.QueryRow(ctx, query, tipoDia).Scan(
		&valor.Id,
		&valor.TipoDia,
		&valor.Valor,
		&valor.Descricao,
		&valor.VigenciaInicio,
		&valor.VigenciaFim,
		&valor.CreatedAt,
		&valor.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, financeiro.ErrorValorDiaNotFound
	}
	if err != nil {
		return nil, err
	}
	return &valor, nil
}

func (t *valorDiaTransaction) Update(ctx context.Context, valor *financeiro.ValorDia) (*financeiro.ValorDia, error) {
	tag, err := t.tx.Exec(ctx, `
		UPDATE config_valores_dia
		SET valor = $2,
		    descricao = $3,
		    vigencia_inicio = $4,
		    vigencia_fim = $5,
		    updated_at = NOW()
		WHERE id = $1 AND tipo_dia = $6
	`, valor.Id, valor.Valor, valor.Descricao, valor.VigenciaInicio, valor.VigenciaFim, valor.TipoDia)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() != 1 {
		return nil, financeiro.ErrorConflitoValorDia
	}

	const query = `
		SELECT id, tipo_dia, valor, descricao, vigencia_inicio, vigencia_fim, created_at, updated_at
		FROM config_valores_dia
		WHERE id = $1
	`
	var atualizado financeiro.ValorDia
	err = t.tx.QueryRow(ctx, query, valor.Id).Scan(
		&atualizado.Id,
		&atualizado.TipoDia,
		&atualizado.Valor,
		&atualizado.Descricao,
		&atualizado.VigenciaInicio,
		&atualizado.VigenciaFim,
		&atualizado.CreatedAt,
		&atualizado.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, financeiro.ErrorConflitoValorDia
	}
	if err != nil {
		return nil, err
	}
	return &atualizado, nil
}
