package cargo

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrorCargoNotFound    = errors.New("cargo não encontrado")
	ErrorCargoAlreadyExists = errors.New("cargo já existe")
	ErrorCargoNomeInvalido  = errors.New("nome do cargo inválido")
)

type Cargo struct {
	Id   uuid.UUID
	Nome string
}

type CargoRepository interface {
	Store(ctx context.Context, nome string) (*Cargo, error)
	FindAll(ctx context.Context) ([]Cargo, error)
	FindById(ctx context.Context, id uuid.UUID) (*Cargo, error)
	ExistsNome(ctx context.Context, nome string) (bool, error)
	Disable(ctx context.Context, id uuid.UUID) error
}

type CargoService struct {
	repository CargoRepository
}

func NewCargoService(repository CargoRepository) *CargoService {
	return &CargoService{repository: repository}
}

func (s *CargoService) Create(ctx context.Context, nome string) (*Cargo, error) {
	if len(nome) < 2 {
		return nil, ErrorCargoNomeInvalido
	}

	exists, err := s.repository.ExistsNome(ctx, nome)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrorCargoAlreadyExists
	}

	return s.repository.Store(ctx, nome)
}

func (s *CargoService) GetAll(ctx context.Context) ([]Cargo, error) {
	return s.repository.FindAll(ctx)
}

func (s *CargoService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repository.FindById(ctx, id)
	if err != nil {
		return ErrorCargoNotFound
	}

	return s.repository.Disable(ctx, id)
}
