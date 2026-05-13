package postgres

import (
	"context"
	"fmt"
	"plantao/internal/domain/role"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository struct {
	pool *pgxpool.Pool
}

func NewRoleRepository(pool *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{pool: pool}
}

func (r *RoleRepository) FindAll(ctx context.Context) ([]role.Role, error) {
	query := `SELECT id, nome FROM roles ORDER BY nome`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar roles: %w", err)
	}
	defer rows.Close()

	var roles []role.Role
	for rows.Next() {
		var ro role.Role
		if err := rows.Scan(&ro.Id, &ro.Nome); err != nil {
			return nil, fmt.Errorf("erro ao escanear role: %w", err)
		}
		roles = append(roles, ro)
	}

	return roles, nil
}
