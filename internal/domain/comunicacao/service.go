package comunicacao

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ModeloComunicacaoService struct {
	repository ModeloComunicaRepository
}

func NewModeloComunicacaoService(repository ModeloComunicaRepository) *ModeloComunicacaoService {
	return &ModeloComunicacaoService{
		repository: repository,
	}
}

func (s *ModeloComunicacaoService) CreateModeloComunicacao(ctx context.Context, nome, tipoComunicacao, assunto, corpo string) (*Comunicacao, error) {
	exists, err := s.repository.ExistsName(ctx, nome)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrorModeloComunicacaoAlreadyExists
	}

	sts, err := ParseTipoComunicacao(tipoComunicacao)

	if err != nil {
		return nil, err
	}

	nomeEnvio, err := ParseNomeEnvio(nome)

	if err != nil {
		return nil, err
	}

	newModelo, err := NewComunicacao(assunto, corpo, StatusAtivo, sts, nomeEnvio)

	if err != nil {
		return nil, err
	}

	return s.repository.Store(ctx, newModelo)
}

func (s *ModeloComunicacaoService) UpdateModeloComunicacao(ctx context.Context, id, nome, tipoComunicacao, assunto, corpo, ativo string) error {
	modeloID, err := uuid.Parse(id)

	if err != nil {
		return err
	}

	modelo, err := s.repository.FindById(ctx, modeloID)

	if err != nil {
		return err
	}

	if modelo == nil {
		return ErrorModeloComunicacaoNotFound
	}

	exists, err := s.repository.ExistsNameExcludingId(ctx, nome, modeloID)

	if err != nil {
		return err
	}

	if exists {
		return ErrorModeloComunicacaoAlreadyExists
	}

	tipo, err := ParseTipoComunicacao(tipoComunicacao)

	if err != nil {
		return err
	}

	status, err := ParseStatusModeloComunicacao(ativo)

	if err != nil {
		return err
	}

	nomeEnvio, err := ParseNomeEnvio(nome)

	if err != nil {
		return err
	}

	err = modelo.UpdateComunicacao(assunto, corpo, &status, &tipo, &nomeEnvio)

	if err != nil {
		return err
	}

	return s.repository.Update(ctx, modelo)
}

func (s *ModeloComunicacaoService) DisableModeloComunicacao(ctx context.Context, id string) error {
	modeloID, err := uuid.Parse(id)

	if err != nil {
		return err
	}

	modelo, err := s.repository.FindById(ctx, modeloID)

	if err != nil {
		return err
	}

	if modelo == nil {
		return ErrorModeloComunicacaoNotFound
	}

	return s.repository.Disable(ctx, modeloID)
}

func (s *ModeloComunicacaoService) GetModeloComunicacaoById(ctx context.Context, id string) (*Comunicacao, error) {
	modeloID, err := uuid.Parse(id)

	if err != nil {
		return nil, err
	}

	modelo, err := s.repository.FindById(ctx, modeloID)

	if err != nil {
		return nil, err
	}

	if modelo == nil {
		return nil, ErrorModeloComunicacaoNotFound
	}

	return modelo, nil
}

func (s *ModeloComunicacaoService) GetAllModelosComunicacao(ctx context.Context) ([]Comunicacao, error) {
	modelos, err := s.repository.FindAll(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]Comunicacao, 0, len(modelos))

	for _, m := range modelos {
		result = append(result, *m)
	}

	return result, nil
}

func ParseStatusModeloComunicacao(s string) (StatusModeloComunicacao, error) {
	switch s {
	case "ATIVO":
		return StatusAtivo, nil
	case "INATIVO":
		return StatusInativo, nil
	}

	return 0, ErrorInvalidStatus
}

func ParseStatusModeloComunicacaoString(s StatusModeloComunicacao) string {
	switch s {
	case StatusAtivo:
		return "ATIVO"
	case StatusInativo:
		return "INATIVO"
	}

	return "UNKNOWN"
}

func ParseTipoComunicacao(s string) (TipoComunicacao, error) {
	switch s {
	case "EMAIL":
		return Email, nil
	case "SMS":
		return SMS, nil
	}

	return "", ErrorInvalidTipoComunicacao
}

func ParseNomeEnvio(s string) (NomeEnvio, error) {
	switch s {
	case "Boas Vindas":
		return BoasVindas, nil
	case "Novo Plantão":
		return NovoPlantao, nil
	case "Cadastro Atualizado":
		return CadastroAtualizado, nil
	case "Cadastro Excluido":
		return CadastroExcluido, nil
	default:
		return "", fmt.Errorf("tipo de envio inválido: %s", s)
	}
}
