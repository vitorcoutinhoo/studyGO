package postgres

import (
	"context"
	"errors"
	"fmt"

	"plantao/internal/domain/smtp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SMTPRepository struct{ pool *pgxpool.Pool }

func NewSMTPRepository(pool *pgxpool.Pool) *SMTPRepository { return &SMTPRepository{pool: pool} }

func (r *SMTPRepository) Get(ctx context.Context) (*smtp.Configuracao, error) {
	const query = `SELECT id, smtp_host, smtp_port, smtp_username, smtp_password, smtp_from, created_at, updated_at FROM config_smtp WHERE singleton = TRUE`
	var configuracao smtp.Configuracao
	err := r.pool.QueryRow(ctx, query).Scan(&configuracao.Id, &configuracao.Host, &configuracao.Porta, &configuracao.Usuario, &configuracao.Senha, &configuracao.Remetente, &configuracao.CreatedAt, &configuracao.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, smtp.ErrorConfiguracaoSMTPNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get SMTP configuration: %w", err)
	}
	return &configuracao, nil
}

func (r *SMTPRepository) Create(ctx context.Context, configuracao *smtp.Configuracao) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO config_smtp (id, smtp_host, smtp_port, smtp_username, smtp_password, smtp_from) VALUES ($1, $2, $3, $4, $5, $6)`, configuracao.Id, configuracao.Host, configuracao.Porta, configuracao.Usuario, configuracao.Senha, configuracao.Remetente)
	if err != nil {
		return translateSMTPStoreError(err)
	}
	return nil
}

func (r *SMTPRepository) Update(ctx context.Context, configuracao *smtp.Configuracao) error {
	const query = `UPDATE config_smtp SET smtp_host = $2, smtp_port = $3, smtp_username = $4, smtp_password = $5, smtp_from = $6, updated_at = NOW() WHERE id = $1 AND singleton = TRUE RETURNING updated_at`
	err := r.pool.QueryRow(ctx, query, configuracao.Id, configuracao.Host, configuracao.Porta, configuracao.Usuario, configuracao.Senha, configuracao.Remetente).Scan(&configuracao.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return smtp.ErrorConfiguracaoSMTPNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update SMTP configuration: %w", err)
	}
	return nil
}

func translateSMTPStoreError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return smtp.ErrorConfiguracaoSMTPAlreadyExists
	}
	return fmt.Errorf("failed to store SMTP configuration: %w", err)
}
