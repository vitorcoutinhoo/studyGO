package role

import (
	"context"
	"errors"
	"plantao/internal/domain/log"

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
	log        log.Logger
}

func NewRoleService(repository RoleRepository, log log.Logger) *RoleService {
	return &RoleService{repository: repository, log: log}
}

func (s *RoleService) GetAll(ctx context.Context) ([]Role, error) {
	s.log.Info("buscando todas as roles")

	roles, err := s.repository.FindAll(ctx)
	if err != nil {
		s.log.Error("erro ao buscar roles", "error", err)
		return nil, err
	}

	s.log.Info("roles encontradas com sucesso", "quantidade", len(roles))
	return roles, nil
}
