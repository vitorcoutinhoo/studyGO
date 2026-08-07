package setor

import (
	"context"
	"errors"
	"plantao/internal/domain/log"

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
	log        log.Logger
}

func NewSetorService(repository SetorRepository, log log.Logger) *SetorService {
	return &SetorService{repository: repository, log: log}
}

func (s *SetorService) Create(ctx context.Context, nome string) (*Setor, error) {
	s.log.Info("iniciando criação de setor", "nome", nome)

	if len(nome) < 2 {
		s.log.Warn("nome do setor inválido", "nome", nome)
		return nil, ErrorSetorNomeInvalido
	}

	s.log.Info("verificando existência de setor no banco de dados", "nome", nome)
	exists, err := s.repository.ExistsNome(ctx, nome)
	if err != nil {
		s.log.Error("erro ao verificar existência de setor", "nome", nome, "error", err)
		return nil, err
	}
	if exists {
		s.log.Warn("setor já está cadastrado", "nome", nome)
		return nil, ErrorSetorAlreadyExists
	}

	setor, err := s.repository.Store(ctx, nome)
	if err != nil {
		s.log.Error("erro ao salvar setor no banco de dados", "nome", nome, "error", err)
		return nil, err
	}

	if setor != nil {
		s.log.Info("setor criado com sucesso", "id_setor", setor.Id, "nome", setor.Nome)
	} else {
		s.log.Info("setor criado com sucesso")
	}
	return setor, nil
}

func (s *SetorService) GetAll(ctx context.Context) ([]Setor, error) {
	s.log.Info("buscando todos os setores")

	setores, err := s.repository.FindAll(ctx)
	if err != nil {
		s.log.Error("erro ao buscar setores", "error", err)
		return nil, err
	}

	s.log.Info("setores encontrados com sucesso", "quantidade", len(setores))
	return setores, nil
}

func (s *SetorService) Delete(ctx context.Context, id uuid.UUID) error {
	s.log.Info("iniciando desativação de setor", "id_setor", id)

	_, err := s.repository.FindById(ctx, id)
	if err != nil {
		s.log.Warn("setor não encontrado para desativação", "id_setor", id, "error", err)
		return ErrorSetorNotFound
	}

	if err := s.repository.Disable(ctx, id); err != nil {
		s.log.Error("erro ao desativar setor", "id_setor", id, "error", err)
		return err
	}

	s.log.Info("setor desativado com sucesso", "id_setor", id)
	return nil
}
