package role

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrorRoleNotFound = errors.New("role não encontrada")

type Role struct {
	Id   uuid.UUID
	Nome string
}

type RoleRepository interface {
	FindAll(ctx context.Context) ([]Role, error)
}

type RoleService struct {
	repository RoleRepository
}

func NewRoleService(repository RoleRepository) *RoleService {
	return &RoleService{repository: repository}
}

func (s *RoleService) GetAll(ctx context.Context) ([]Role, error) {
	return s.repository.FindAll(ctx)
}
