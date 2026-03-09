package postgres

import (
	"context"
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
		id_colaborador,
		tipo,
		destinatario,
		assunto,
		corpo,
		status,
		data_envio,
		erro_log
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		com.IdModelo,
		com.IdColaborador,
		string(com.TipoComunicacao),
		com.Destinatario,
		com.Assunto,
		com.Corpo,
		string(com.Status),
		com.DataEnvio,
		com.ErroLog,
	)

	if err != nil {
		return err
	}

	return nil
}
