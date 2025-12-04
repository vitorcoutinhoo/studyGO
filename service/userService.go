package service

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"main.go/types"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// CreateUser implements types.UserRepository.
func (s *UserService) CreateUser(colaboradorId uuid.UUID, user types.UserRequest) (*types.UserResponse, error) {
	sqlStatement := `
        INSERT INTO usuarios_login (id_colaborador, email, senha_hash) 
        VALUES ($1, $2, $3)
        RETURNING 
            id, id_colaborador, email, role, ativo, created_at, updated_at
    `

	u := new(types.UserResponse)

	passwordHash, errPassword := BcryptHashPassword(user.SenhaHash)

	if errPassword != nil {
		return nil, errPassword
	}

	err := s.db.QueryRow(
		sqlStatement,
		colaboradorId,
		user.Email,
		passwordHash,
	).Scan(
		&u.ID,
		&u.IDColaborador,
		&u.Email,
		&u.Role,
		&u.Ativo,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				switch pqErr.Constraint {
				case "usuarios_login_id_colaborador_key":
					return nil, fmt.Errorf("o colaborador (ID: %s) já possui um cadastro de usuário", colaboradorId)
				case "usuarios_login_email_key":
					return nil, fmt.Errorf("o email '%s' já está em uso", user.Email)
				}
			case "23503":
				return nil, fmt.Errorf("colaborador com ID '%s' não encontrado. Não é possível criar o login", colaboradorId)
			}
		}

		return nil, err
	}

	return u, nil
}

// GetUsers implements types.UserRepository.
func (s *UserService) GetUsers() ([]*types.UserResponse, error) {
	sqlStatement := `
        SELECT id, id_colaborador, email, role, ativo, created_at, updated_at 
        FROM usuarios_login
        ORDER BY created_at DESC
    `

	rows, err := s.db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*types.UserResponse

	for rows.Next() {
		u := new(types.UserResponse)

		err := rows.Scan(
			&u.ID,
			&u.IDColaborador,
			&u.Email,
			&u.Role,
			&u.Ativo,
			&u.CreatedAt,
			&u.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserById implements types.UserRepository.
func (s *UserService) GetUserById(id uuid.UUID) (*types.UserResponse, error) {
	sqlStatement := `
        SELECT id, id_colaborador, email, role, ativo, created_at, updated_at 
        FROM usuarios_login
        WHERE id = $1
    `

	u := new(types.UserResponse)

	err := s.db.QueryRow(sqlStatement, id).Scan(
		&u.ID,
		&u.IDColaborador,
		&u.Email,
		&u.Role,
		&u.Ativo,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuário com ID '%s' não encontrado", id)
		}

		return nil, err
	}

	return u, nil
}

// UpdateUser implements types.UserRepository.
func (s *UserService) UpdateUser(id uuid.UUID, user types.UserRequest) error {
	sqlStatement := `
        UPDATE usuarios_login 
        SET email = $1, senha_hash = $2, updated_at = NOW()
        WHERE id = $3
        RETURNING id 
    `

	var returnedID uuid.UUID

	err := s.db.QueryRow(sqlStatement, user.Email, user.SenhaHash, id).Scan(&returnedID)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("usuário com ID '%s' não encontrado para atualizar", id)
		}

		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return fmt.Errorf("o email '%s' já está em uso por outro usuário", user.Email)
			}
		}

		return err
	}

	return nil
}

// DeletUserById implements types.UserRepository.
func (s *UserService) DeleteUserById(id uuid.UUID) error {
	sqlStatement := `
        DELETE FROM usuarios_login
        WHERE id = $1
        RETURNING id
    `

	var deletedID uuid.UUID

	err := s.db.QueryRow(sqlStatement, id).Scan(&deletedID)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("usuário com ID '%s' não encontrado para deletar", id)
		}

		return err
	}

	return nil
}
