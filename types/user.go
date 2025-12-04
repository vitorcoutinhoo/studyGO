package types

import (
	"time"

	"github.com/google/uuid"
)

type UserRequest struct {
	Email     string `json:"email"`
	SenhaHash string `json:"senha"`
}

type UserResponse struct {
	ID            uuid.UUID `json:"id"`
	IDColaborador uuid.UUID `json:"id_colaborador"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	Ativo         string    `json:"ativo"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserRepository interface {
	CreateUser(colaboradorId uuid.UUID, user UserRequest) (*UserResponse, error)
	GetUsers() ([]*UserResponse, error)
	GetUserById(id uuid.UUID) (*UserResponse, error)
	UpdateUser(id uuid.UUID, user UserRequest) error
	DeleteUserById(id uuid.UUID) error
}
