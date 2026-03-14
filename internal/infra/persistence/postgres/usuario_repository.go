package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"plantao/internal/domain/usuario"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsuarioRepository struct {
	pool *pgxpool.Pool
}

func NewUsuarioRepository(pool *pgxpool.Pool) *UsuarioRepository {
	return &UsuarioRepository{pool: pool}
}

func (r *UsuarioRepository) Store(ctx context.Context, usuario *usuario.Usuario) (*usuario.Usuario, error) {
	query := `
	INSERT INTO usuarios_login (
		id_colaborador,
		email,
		senha_hash,
		role,
		ativo
	)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id, id_colaborador, email, role, ativo, created_at, updated_at
	`

	var ativoDB string

	err := r.pool.QueryRow(
		ctx,
		query,
		usuario.IdColaborador,
		usuario.Email,
		usuario.Senha,
		usuario.Role,
		statusToDBUsuario(usuario.Ativo),
	).Scan(
		&usuario.Id,
		&usuario.IdColaborador,
		&usuario.Email,
		&usuario.Role,
		&ativoDB,
		&usuario.CreatedAt,
		&usuario.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("erro ao inserir novo usuário: %w", err)
	}

	usuario.Ativo = dbToStatusUsuario(ativoDB)

	return usuario, nil
}

func (r *UsuarioRepository) Update(ctx context.Context, u *usuario.Usuario) error {
	query := `
	UPDATE usuarios_login
	SET
		email = $1,
		senha_hash = $2,
		ativo = $3,
		updated_at = NOW()
	WHERE id = $4 AND ativo = 'Y'
	`

	result, err := r.pool.Exec(
		ctx,
		query,
		u.Email,
		u.Senha,
		statusToDBUsuario(u.Ativo),
		u.Id,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	if result.RowsAffected() == 0 {
		return usuario.ErrorUserNotFound
	}

	return nil
}

func (r *UsuarioRepository) Delete(ctx context.Context, usuarioId uuid.UUID) error {
	query := `
		DELETE FROM usuarios_login WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, usuarioId)

	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("usuário não encontrado")
	}

	return nil
}

func (r *UsuarioRepository) FindById(ctx context.Context, usuarioId uuid.UUID) (*usuario.Usuario, error) {
	query := `
	SELECT
		id,
		id_colaborador,
		email,
		senha_hash,
		role,
		ativo,
		created_at,
		updated_at
	FROM usuarios_login
	WHERE id = $1 AND ativo = 'Y'
	`

	var u usuario.Usuario
	var ativoDB string

	err := r.pool.QueryRow(ctx, query, usuarioId).Scan(
		&u.Id,
		&u.IdColaborador,
		&u.Email,
		&u.Senha,
		&u.Role,
		&ativoDB,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, usuario.ErrorUserNotFound
		}
		return nil, fmt.Errorf("erro ao procurar usuário por id: %w", err)
	}

	sts := dbToStatusUsuario(ativoDB)

	u.Ativo = sts

	return &u, nil
}

func (r *UsuarioRepository) FindByEmail(ctx context.Context, email string) (*usuario.Usuario, error) {
	query := `
	SELECT
		id,
		id_colaborador,
		email,
		senha_hash,
		role,
		ativo,
		created_at,
		updated_at
	FROM usuarios_login
	WHERE email = $1 AND ativo = 'Y'
	`

	var u usuario.Usuario
	var ativoDB string

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.Id,
		&u.IdColaborador,
		&u.Email,
		&u.Senha,
		&u.Role,
		&ativoDB,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, usuario.ErrorUserNotFound
		}
		return nil, fmt.Errorf("erro ao procurar usuário por email: %w", err)
	}

	sts := dbToStatusUsuario(ativoDB)

	u.Ativo = sts

	return &u, nil
}

func (r *UsuarioRepository) FindAll(ctx context.Context) (*[]usuario.Usuario, error) {
	query := `
		SELECT id, id_colaborador, email, role, ativo, created_at, updated_at
		FROM usuarios_login
	`

	rows, err := r.pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %w", err)
	}

	defer rows.Close()

	var usuarios []usuario.Usuario
	var ativoDB string

	for rows.Next() {
		var u usuario.Usuario

		err := rows.Scan(
			&u.Id,
			&u.IdColaborador,
			&u.Email,
			&u.Role,
			&ativoDB,
			&u.CreatedAt,
			&u.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler usuário: %w", err)
		}

		sts := dbToStatusUsuario(ativoDB)

		u.Ativo = sts

		usuarios = append(usuarios, u)
	}

	return &usuarios, nil
}

func (r *UsuarioRepository) ExistsEmail(ctx context.Context, email string) (bool, error) {
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM usuarios_login
		WHERE email = $1 AND ativo = 'Y'
	)
	`

	var exists bool

	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência de email: %w", err)
	}

	return exists, nil
}

func (r *UsuarioRepository) ExistsId(ctx context.Context, usuarioId uuid.UUID) (bool, error) {
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM usuarios_login
		WHERE id = $1 AND ativo = 'Y'
	)
	`

	var exists bool

	err := r.pool.QueryRow(ctx, query, usuarioId).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência de id: %w", err)
	}

	return exists, nil
}

func (r *UsuarioRepository) ExistsEmailExcludingId(ctx context.Context, email string, usuarioId uuid.UUID) (bool, error) {
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM usuarios_login
		WHERE email = $1
		AND id <> $2 AND ativo = 'Y'
	)
	`

	var exists bool

	err := r.pool.QueryRow(ctx, query, email, usuarioId).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("erro ao verificar email: %w", err)
	}

	return exists, nil
}

func statusToDBUsuario(status usuario.StatusUsuario) string {
	switch status {
	case usuario.StatusAtivo:
		return "Y"
	case usuario.StatusInativo:
		return "N"
	default:
		return "N"
	}
}

func dbToStatusUsuario(dbValue string) usuario.StatusUsuario {
	switch dbValue {
	case "Y":
		return usuario.StatusAtivo
	case "N":
		return usuario.StatusInativo
	default:
		return usuario.StatusInativo
	}
}
