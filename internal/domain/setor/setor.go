package setor

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrorSetorNotFound      = errors.New("setor não encontrado")
	ErrorSetorAlreadyExists = errors.New("setor já existe")
	ErrorSetorNomeInvalido  = errors.New("nome do setor inválido")
)

type Setor struct {
	Id   uuid.UUID
	Nome string
}

type SetorRepository interface {
	Store(ctx context.Context, nome string) (*Setor, error)
	FindAll(ctx context.Context) ([]Setor, error)
	FindById(ctx context.Context, id uuid.UUID) (*Setor, error)
	ExistsNome(ctx context.Context, nome string) (bool, error)
	Disable(ctx context.Context, id uuid.UUID) error
}

type SetorService struct {
	repository SetorRepository
}

func NewSetorService(repository SetorRepository) *SetorService {
	return &SetorService{repository: repository}
}

func (s *SetorService) Create(ctx context.Context, nome string) (*Setor, error) {
	if len(nome) < 2 {
		return nil, ErrorSetorNomeInvalido
	}

	exists, err := s.repository.ExistsNome(ctx, nome)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrorSetorAlreadyExists
	}

	return s.repository.Store(ctx, nome)
}

func (s *SetorService) GetAll(ctx context.Context) ([]Setor, error) {
	return s.repository.FindAll(ctx)
}

func (s *SetorService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repository.FindById(ctx, id)
	if err != nil {
		return ErrorSetorNotFound
	}

	return s.repository.Disable(ctx, id)
}
