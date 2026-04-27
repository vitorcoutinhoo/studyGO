package postgres

import (
	"context"
	"fmt"
	"plantao/internal/domain/convite"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConviteRepository struct {
	pool *pgxpool.Pool
}

func NewConviteRepository(pool *pgxpool.Pool) *ConviteRepository {
	return &ConviteRepository{pool: pool}
}

func (r *ConviteRepository) Store(ctx context.Context, idColaborador uuid.UUID) (*convite.Convite, error) {
	query := `
        INSERT INTO convites (id_colaborador)
        VALUES ($1)
        RETURNING token, id_colaborador, expira_em, usado, created_at
    `
	var c convite.Convite
	err := r.pool.QueryRow(ctx, query, idColaborador).Scan(
		&c.Token, &c.IdColaborador, &c.ExpiraEm, &c.Usado, &c.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar convite: %w", err)
	}

	return &c, nil
}

func (r *ConviteRepository) FindByToken(ctx context.Context, token uuid.UUID) (*convite.Convite, error) {
	query := `
        SELECT token, id_colaborador, expira_em, usado, created_at
        FROM convites
        WHERE token = $1
    `
	var c convite.Convite
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&c.Token, &c.IdColaborador, &c.ExpiraEm, &c.Usado, &c.CreatedAt,
	)

	if err != nil {
		return nil, convite.ErrorConviteNotFound
	}

	return &c, nil
}

func (r *ConviteRepository) MarkAsUsed(ctx context.Context, token uuid.UUID) error {
	query := `UPDATE convites SET usado = TRUE WHERE token = $1`
	_, err := r.pool.Exec(ctx, query, token)
	return err
}

func (r *ConviteRepository) Disable(ctx context.Context, token uuid.UUID) error {
	query := `UPDATE convites 
	SET usado = TRUE, 
	expira_em = '2000-01-01 00:00:00'
	WHERE token = $1`

	_, err := r.pool.Exec(ctx, query, token)
	return err
}

// Na implementação — internal/infra/postgres/convite_repository.go
func (r *ConviteRepository) FindAll(ctx context.Context) ([]convite.Convite, error) {
	query := `
        SELECT token, id_colaborador, expira_em, usado, created_at
        FROM convites
        ORDER BY created_at DESC
    `

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar convites: %w", err)
	}

	defer rows.Close()

	var convites []convite.Convite
	for rows.Next() {
		var c convite.Convite
		err := rows.Scan(
			&c.Token,
			&c.IdColaborador,
			&c.ExpiraEm,
			&c.Usado,
			&c.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler convite: %w", err)
		}
		convites = append(convites, c)
	}

	return convites, nil
}
