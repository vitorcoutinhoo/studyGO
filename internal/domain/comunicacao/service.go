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
	exists, err := s.repository.ExistsTipo(ctx, tipoComunicacao)

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

	newModelo, err := NewComunicacao(nome, assunto, corpo, StatusAtivo, sts)

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

	exists, err := s.repository.ExistsTipoExcludingId(ctx, tipoComunicacao, modeloID)

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

	err = modelo.UpdateComunicacao(&nome, &assunto, &corpo, &status, &tipo)

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
	case "Plantão Agendado":
		return PlantaoAgendado, nil
	case "Plantão Concluido":
		return PlantaoConluido, nil
	case "Plantão Ainda Está Aberto":
		return PlantaoAindaAberto, nil
	case "Plantão Pago":
		return PlantaoPago, nil
	case "Colaborador Cadastrado":
		return ColaboradorCadastrado, nil
	case "Colaborador Atualizado":
		return ColaboradorAtualizado, nil
	case "Colaborador Deletado":
		return ColaboradorDeletado, nil
	case "Usuário Cadastrado":
		return UsuarioCadastrado, nil
	case "Email do Usuário Atualizado":
		return EmailAtualizado, nil
	case "Senha do Usuário Atualizada":
		return SenhaAtualizada, nil
	case "Usuário Deletado":
		return UsuarioDeletado, nil
	default:
		return "", fmt.Errorf("tipo de comunicação inválido: %s", s)
	}
}
