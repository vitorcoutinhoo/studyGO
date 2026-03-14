package postgres

import (
	"context"
	"errors"
	"fmt"
	"plantao/internal/domain/comunicacao"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ModeloRepository struct {
	pool *pgxpool.Pool
}

func NewModeloRepository(pool *pgxpool.Pool) *ModeloRepository {
	return &ModeloRepository{pool: pool}
}

func (r *ModeloRepository) Store(ctx context.Context, com *comunicacao.Comunicacao) (*comunicacao.Comunicacao, error) {
	query := `
		INSERT INTO modelos_comunicacao (
			nome,
			tipo,
			assunto,
			corpo,
			ativo
		)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, created_at, updated_at
	`

	var createdAt time.Time
	var updatedAt time.Time
	var id uuid.UUID

	err := r.pool.QueryRow(
		ctx,
		query,
		com.Nome,
		string(com.TipoComunicacao),
		com.Assunto,
		com.Corpo,
		statusToDBModelo(com.Ativo),
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("erro ao inserir novo modelo de comunicação: %w", err)
	}

	com.Id = id
	com.CreatedAt = &createdAt
	com.UpdatedAt = &updatedAt

	return com, nil
}

func (r *ModeloRepository) Update(ctx context.Context, com *comunicacao.Comunicacao) error {
	query := `
	UPDATE modelos_comunicacao
	SET
		nome = $1,
		tipo = $2,
		assunto = $3,
		corpo = $4,
		ativo = $5,
		updated_at = NOW()
	WHERE id = $6 AND ativo = 'Y'
	`

	result, err := r.pool.Exec(
		ctx,
		query,
		com.Nome,
		string(com.TipoComunicacao),
		com.Assunto,
		com.Corpo,
		statusToDBModelo(com.Ativo),
		com.Id,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar modelo de comunicação: %w", err)
	}

	if result.RowsAffected() == 0 {
		return comunicacao.ErrorModeloComunicacaoNotFound
	}

	return nil
}

func (r *ModeloRepository) Disable(ctx context.Context, modeloId uuid.UUID) error {
	query := `
	UPDATE modelos_comunicacao
	SET
		tipo = 'INATIVO_' || id || '_' || tipo,
		ativo = 'N',
		updated_at = NOW()
	WHERE id = $1 AND ativo = 'Y'
	`

	result, err := r.pool.Exec(ctx, query, modeloId)

	if err != nil {
		return fmt.Errorf("erro ao desativar modelo de comunicação: %w", err)
	}

	if result.RowsAffected() == 0 {
		return comunicacao.ErrorModeloComunicacaoNotFound
	}

	return nil
}

func (r *ModeloRepository) FindById(ctx context.Context, modeloId uuid.UUID) (*comunicacao.Comunicacao, error) {
	query := `
	SELECT
		id,
		nome,
		tipo,
		assunto,
		corpo,
		ativo,
		created_at,
		updated_at
	FROM modelos_comunicacao
	WHERE id = $1 AND ativo = 'Y'
	`

	var com comunicacao.Comunicacao
	var tipo string
	var ativo string
	var createdAt time.Time
	var updatedAt time.Time

	err := r.pool.QueryRow(ctx, query, modeloId).Scan(
		&com.Id,
		&com.Nome,
		&tipo,
		&com.Assunto,
		&com.Corpo,
		&ativo,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("nenhum modelo de comunicação encontrado com ID [%s]", modeloId)
		}
		return nil, fmt.Errorf("erro ao procurar modelo de comunicação por ID: %w", err)
	}

	com.TipoComunicacao = comunicacao.TipoComunicacao(tipo)
	com.Ativo = dbToStatusModelo(ativo)
	com.CreatedAt = &createdAt
	com.UpdatedAt = &updatedAt

	return &com, nil
}

func (r *ModeloRepository) FindByTipo(ctx context.Context, tipoComunicacao string) (*comunicacao.Comunicacao, error) {
	query := `
	SELECT
		id,
		nome,
		tipo,
		assunto,
		corpo,
		ativo,
		created_at,
		updated_at
	FROM modelos_comunicacao
	WHERE tipo = $1 AND ativo = 'Y'
	`

	var com comunicacao.Comunicacao
	var tipo string
	var ativo string
	var createdAt time.Time
	var updatedAt time.Time

	err := r.pool.QueryRow(ctx, query, tipoComunicacao).Scan(
		&com.Id,
		&com.Nome,
		&tipo,
		&com.Assunto,
		&com.Corpo,
		&ativo,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("nenhum modelo de comunicação encontrado com tipo [%s]", tipoComunicacao)
		}
		return nil, fmt.Errorf("erro ao procurar modelo de comunicação por tipo: %w", err)
	}

	com.TipoComunicacao = comunicacao.TipoComunicacao(tipo)
	com.Ativo = dbToStatusModelo(ativo)
	com.CreatedAt = &createdAt
	com.UpdatedAt = &updatedAt

	return &com, nil
}

func (r *ModeloRepository) FindAll(ctx context.Context) ([]*comunicacao.Comunicacao, error) {
	query := `
	SELECT
		id,
		nome,
		tipo,
		assunto,
		corpo,
		ativo,
		created_at,
		updated_at
	FROM modelos_comunicacao
	WHERE ativo = 'Y'
	`

	rows, err := r.pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("erro ao obter modelos de comunicação: %w", err)
	}

	defer rows.Close()

	var modelos []*comunicacao.Comunicacao

	for rows.Next() {
		var com comunicacao.Comunicacao
		var tipo string
		var ativo string
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(
			&com.Id,
			&com.Nome,
			&tipo,
			&com.Assunto,
			&com.Corpo,
			&ativo,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao escanear modelo de comunicação: %w", err)
		}

		com.TipoComunicacao = comunicacao.TipoComunicacao(tipo)
		com.Ativo = dbToStatusModelo(ativo)
		com.CreatedAt = &createdAt
		com.UpdatedAt = &updatedAt

		modelos = append(modelos, &com)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar modelos de comunicação: %w", err)
	}

	return modelos, nil
}

func (r *ModeloRepository) ExistsTipo(ctx context.Context, tipo string) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM modelos_comunicacao
		WHERE tipo = $1 AND ativo = 'Y'
	)
	`

	var exists bool

	err := r.pool.QueryRow(ctx, query, tipo).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("erro ao verificar a existência do tipo de modelo comunicação: %w", err)
	}

	return exists, nil
}

func (r *ModeloRepository) ExistsTipoExcludingId(
	ctx context.Context,
	tipo string,
	id uuid.UUID,
) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM modelos_comunicacao
		WHERE tipo = $1
		AND id <> $2
		AND ativo = 'Y'
	)
	`

	var exists bool

	err := r.pool.QueryRow(ctx, query, tipo, id).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("erro ao verificar a existência do tipo de modelo comunicação: %w", err)
	}

	return exists, nil
}

func statusToDBModelo(status comunicacao.StatusModeloComunicacao) string {
	switch status {
	case comunicacao.StatusAtivo:
		return "Y"
	case comunicacao.StatusInativo:
		return "N"
	default:
		return "N"
	}
}

func dbToStatusModelo(dbValue string) comunicacao.StatusModeloComunicacao {
	switch dbValue {
	case "Y":
		return comunicacao.StatusAtivo
	case "N":
		return comunicacao.StatusInativo
	default:
		return comunicacao.StatusInativo
	}
}
