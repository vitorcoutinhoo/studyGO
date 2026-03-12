package postgres

import (
	"context"
	"fmt"
	"plantao/internal/domain/comunicacao"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EnvioRepository struct {
	pool *pgxpool.Pool
}

func NewEnvioRepository(pool *pgxpool.Pool) *EnvioRepository {
	return &EnvioRepository{pool: pool}
}

func (r *EnvioRepository) Store(ctx context.Context, com *comunicacao.Envio) error {
	query := `
	INSERT INTO envios_comunicacao (
		id_modelo,
		tipo,
		destinatario,
		status,
		erro_log
	)
	VALUES ($1,$2,$3,$4,$5)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		com.IdModelo,
		string(com.TipoComunicacao),
		com.Destinatario,
		string(com.Status),
		com.ErroLog,
	)

	if err != nil {
		return fmt.Errorf("erro ao inserir envio de comunicação: %w", err)
	}

	return nil
}
