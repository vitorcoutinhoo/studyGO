package cargo

import (
	"context"
	"errors"
	"plantao/internal/domain/log"

	"github.com/google/uuid"
)

var (
	ErrorCargoNotFound      = errors.New("cargo não encontrado")
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
	log        log.Logger
}

func NewCargoService(repository CargoRepository, log log.Logger) *CargoService {
	return &CargoService{repository: repository, log: log}
}

func (s *CargoService) Create(ctx context.Context, nome string) (*Cargo, error) {
	s.log.Info("iniciando criação de cargo", "nome", nome)

	if len(nome) < 2 {
		s.log.Warn("nome do cargo inválido", "nome", nome)
		return nil, ErrorCargoNomeInvalido
	}

	s.log.Info("verificando existência de cargo no banco de dados", "nome", nome)
	exists, err := s.repository.ExistsNome(ctx, nome)
	if err != nil {
		s.log.Error("erro ao verificar existência de cargo", "nome", nome, "error", err)
		return nil, err
	}
	if exists {
		s.log.Warn("cargo já está cadastrado", "nome", nome)
		return nil, ErrorCargoAlreadyExists
	}

	cargo, err := s.repository.Store(ctx, nome)
	if err != nil {
		s.log.Error("erro ao salvar cargo no banco de dados", "nome", nome, "error", err)
		return nil, err
	}

	if cargo != nil {
		s.log.Info("cargo criado com sucesso", "id_cargo", cargo.Id, "nome", cargo.Nome)
	} else {
		s.log.Info("cargo criado com sucesso")
	}
	return cargo, nil
}

func (s *CargoService) GetAll(ctx context.Context) ([]Cargo, error) {
	s.log.Info("buscando todos os cargos")

	cargos, err := s.repository.FindAll(ctx)
	if err != nil {
		s.log.Error("erro ao buscar cargos", "error", err)
		return nil, err
	}

	s.log.Info("cargos encontrados com sucesso", "quantidade", len(cargos))
	return cargos, nil
}

func (s *CargoService) Delete(ctx context.Context, id uuid.UUID) error {
	s.log.Info("iniciando desativação de cargo", "id_cargo", id)

	_, err := s.repository.FindById(ctx, id)
	if err != nil {
		s.log.Warn("cargo não encontrado para desativação", "id_cargo", id, "error", err)
		return ErrorCargoNotFound
	}

	if err := s.repository.Disable(ctx, id); err != nil {
		s.log.Error("erro ao desativar cargo", "id_cargo", id, "error", err)
		return err
	}

	s.log.Info("cargo desativado com sucesso", "id_cargo", id)
	return nil
}
