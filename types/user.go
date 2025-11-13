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
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Ativo     string    `json:"ativo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	CreateUser(user UserRequest) (*UserResponse, error)
	GetUsers() ([]*UserResponse, error)
	GetUserById(id uuid.UUID) (*UserResponse, error)
	UpdateUser(id uuid.UUID, user UserRequest) error
	DeletUserById(id uuid.UUID) error
}
