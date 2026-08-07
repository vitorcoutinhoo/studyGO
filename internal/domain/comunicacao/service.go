package comunicacao

import (
	"context"
	"fmt"
	"plantao/internal/domain/log"

	"github.com/google/uuid"
)

type ModeloComunicacaoService struct {
	repository ModeloComunicaRepository
	log        log.Logger
}

func NewModeloComunicacaoService(repository ModeloComunicaRepository, log log.Logger) *ModeloComunicacaoService {
	return &ModeloComunicacaoService{
		repository: repository,
		log:        log,
	}
}

func (s *ModeloComunicacaoService) CreateModeloComunicacao(ctx context.Context, nome, tipoComunicacao, assunto, corpo string) (*Comunicacao, error) {
	s.log.Info("iniciando criação de modelo de comunicação", "nome", nome, "tipo_comunicacao", tipoComunicacao)

	s.log.Info("verificando existência de modelo de comunicação por tipo", "tipo_comunicacao", tipoComunicacao)
	exists, err := s.repository.ExistsTipo(ctx, tipoComunicacao)

	if err != nil {
		s.log.Error("erro ao verificar existência de modelo de comunicação", "tipo_comunicacao", tipoComunicacao, "error", err)
		return nil, err
	}

	if exists {
		s.log.Warn("modelo de comunicação já existe", "tipo_comunicacao", tipoComunicacao)
		return nil, ErrorModeloComunicacaoAlreadyExists
	}

	sts, err := ParseTipoComunicacao(tipoComunicacao)

	if err != nil {
		s.log.Warn("tipo de comunicação inválido", "tipo_comunicacao", tipoComunicacao, "error", err)
		return nil, err
	}

	s.log.Info("criando objeto de modelo de comunicação", "tipo_comunicacao", tipoComunicacao)
	newModelo, err := NewComunicacao(nome, assunto, corpo, StatusAtivo, sts)

	if err != nil {
		s.log.Warn("erro ao validar modelo de comunicação", "tipo_comunicacao", tipoComunicacao, "error", err)
		return nil, err
	}

	modelo, err := s.repository.Store(ctx, newModelo)
	if err != nil {
		s.log.Error("erro ao salvar modelo de comunicação", "tipo_comunicacao", tipoComunicacao, "error", err)
		return nil, err
	}

	if modelo != nil {
		s.log.Info("modelo de comunicação criado com sucesso", "id_modelo", modelo.Id, "tipo_comunicacao", modelo.TipoComunicacao)
	} else {
		s.log.Info("modelo de comunicação criado com sucesso")
	}
	return modelo, nil
}

func (s *ModeloComunicacaoService) UpdateModeloComunicacao(ctx context.Context, id, nome, tipoComunicacao, assunto, corpo, ativo string) error {
	s.log.Info("iniciando atualização de modelo de comunicação", "id_modelo", id, "tipo_comunicacao", tipoComunicacao)

	modeloID, err := uuid.Parse(id)

	if err != nil {
		s.log.Warn("UUID do modelo de comunicação inválido", "id_modelo", id, "error", err)
		return err
	}

	modelo, err := s.repository.FindById(ctx, modeloID)

	if err != nil {
		s.log.Error("erro ao buscar modelo de comunicação para atualização", "id_modelo", modeloID, "error", err)
		return err
	}

	if modelo == nil {
		s.log.Warn("modelo de comunicação não encontrado para atualização", "id_modelo", modeloID)
		return ErrorModeloComunicacaoNotFound
	}

	exists, err := s.repository.ExistsTipoExcludingId(ctx, tipoComunicacao, modeloID)

	if err != nil {
		s.log.Error("erro ao verificar existência de tipo de comunicação excluindo modelo", "id_modelo", modeloID, "tipo_comunicacao", tipoComunicacao, "error", err)
		return err
	}

	if exists {
		s.log.Warn("tipo de comunicação já utilizado por outro modelo", "id_modelo", modeloID, "tipo_comunicacao", tipoComunicacao)
		return ErrorModeloComunicacaoAlreadyExists
	}

	tipo, err := ParseTipoComunicacao(tipoComunicacao)

	if err != nil {
		s.log.Warn("tipo de comunicação inválido para atualização", "id_modelo", modeloID, "tipo_comunicacao", tipoComunicacao, "error", err)
		return err
	}

	status, err := ParseStatusModeloComunicacao(ativo)

	if err != nil {
		s.log.Warn("status de modelo de comunicação inválido", "id_modelo", modeloID, "status", ativo, "error", err)
		return err
	}

	err = modelo.UpdateComunicacao(&nome, &assunto, &corpo, &status, &tipo)

	if err != nil {
		s.log.Warn("erro ao validar atualização do modelo de comunicação", "id_modelo", modeloID, "error", err)
		return err
	}

	if err := s.repository.Update(ctx, modelo); err != nil {
		s.log.Error("erro ao atualizar modelo de comunicação no banco de dados", "id_modelo", modeloID, "error", err)
		return err
	}

	s.log.Info("modelo de comunicação atualizado com sucesso", "id_modelo", modeloID)
	return nil
}

func (s *ModeloComunicacaoService) DisableModeloComunicacao(ctx context.Context, id string) error {
	s.log.Info("iniciando desativação de modelo de comunicação", "id_modelo", id)

	modeloID, err := uuid.Parse(id)

	if err != nil {
		s.log.Warn("UUID do modelo de comunicação inválido para desativação", "id_modelo", id, "error", err)
		return err
	}

	modelo, err := s.repository.FindById(ctx, modeloID)

	if err != nil {
		s.log.Error("erro ao buscar modelo de comunicação para desativação", "id_modelo", modeloID, "error", err)
		return err
	}

	if modelo == nil {
		s.log.Warn("modelo de comunicação não encontrado para desativação", "id_modelo", modeloID)
		return ErrorModeloComunicacaoNotFound
	}

	if err := s.repository.Disable(ctx, modeloID); err != nil {
		s.log.Error("erro ao desativar modelo de comunicação", "id_modelo", modeloID, "error", err)
		return err
	}

	s.log.Info("modelo de comunicação desativado com sucesso", "id_modelo", modeloID)
	return nil
}

func (s *ModeloComunicacaoService) GetModeloComunicacaoById(ctx context.Context, id string) (*Comunicacao, error) {
	s.log.Info("buscando modelo de comunicação por ID", "id_modelo", id)

	modeloID, err := uuid.Parse(id)

	if err != nil {
		s.log.Warn("UUID do modelo de comunicação inválido para busca", "id_modelo", id, "error", err)
		return nil, err
	}

	modelo, err := s.repository.FindById(ctx, modeloID)

	if err != nil {
		s.log.Error("erro ao buscar modelo de comunicação por ID", "id_modelo", modeloID, "error", err)
		return nil, err
	}

	if modelo == nil {
		s.log.Warn("modelo de comunicação não encontrado", "id_modelo", modeloID)
		return nil, ErrorModeloComunicacaoNotFound
	}

	s.log.Info("modelo de comunicação encontrado com sucesso", "id_modelo", modeloID)
	return modelo, nil
}

func (s *ModeloComunicacaoService) GetAllModelosComunicacao(ctx context.Context) ([]Comunicacao, error) {
	s.log.Info("buscando todos os modelos de comunicação")

	modelos, err := s.repository.FindAll(ctx)

	if err != nil {
		s.log.Error("erro ao buscar modelos de comunicação", "error", err)
		return nil, err
	}

	result := make([]Comunicacao, 0, len(modelos))

	for _, m := range modelos {
		result = append(result, *m)
	}

	s.log.Info("modelos de comunicação encontrados com sucesso", "quantidade", len(result))
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
