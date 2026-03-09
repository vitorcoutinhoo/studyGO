package comunicacao

import (
	"context"
	"fmt"
	"plantao/internal/api/dto"

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

func (s *ModeloComunicacaoService) CreateModeloComunicacao(ctx context.Context, com *dto.ModeloComunicacaoRequestDTO) (*dto.ModeloComunicacaoResponseDTO, error) {
	exists, err := s.repository.ExistsName(ctx, com.Nome)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrorModeloComunicacaoAlreadyExists
	}

	sts, err := parseTipoComunicacao(com.TipoComunicacao)

	if err != nil {
		return nil, err
	}

	nomeEnvio, err := parseNomeEnvio(com.Nome)

	if err != nil {
		return nil, err
	}

	newModelo, err := NewComunicacao(com.Assunto, com.Corpo, StatusAtivo, sts, nomeEnvio)

	if err != nil {
		return nil, err
	}

	createdModelo, err := s.repository.Store(ctx, newModelo)

	if err != nil {
		return nil, err
	}

	response, err := domainToResponse(*createdModelo)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *ModeloComunicacaoService) UpdateModeloComunicacao(ctx context.Context, id string, com *dto.ModeloComunicacaoUpdateRequestDTO) error {
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

	exists, err := s.repository.ExistsNameExcludingId(ctx, com.Nome, modeloID)

	if err != nil {
		return err
	}

	if exists {
		return ErrorModeloComunicacaoAlreadyExists
	}

	tipo, err := parseTipoComunicacao(com.TipoComunicacao)

	if err != nil {
		return err
	}

	status, err := parseStatusModeloComunicacao(com.Ativo)

	if err != nil {
		return err
	}

	nomeEnvio, err := parseNomeEnvio(com.Nome)

	if err != nil {
		return err
	}

	err = modelo.UpdateComunicacao(
		com.Assunto,
		com.Corpo,
		&status,
		&tipo,
		&nomeEnvio,
	)

	if err != nil {
		return err
	}

	err = s.repository.Update(ctx, modelo)

	if err != nil {
		return err
	}

	return nil
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

	err = s.repository.Disable(ctx, modeloID)

	if err != nil {
		return err
	}

	return nil
}

func (s *ModeloComunicacaoService) GetModeloComunicacaoById(ctx context.Context, id string) (*dto.ModeloComunicacaoResponseDTO, error) {
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

	return domainToResponse(*modelo)
}

func (s *ModeloComunicacaoService) GetAllModelosComunicacao(ctx context.Context) (*[]dto.ModeloComunicacaoResponseDTO, error) {
	modelos, err := s.repository.FindAll(ctx)

	if err != nil {
		return nil, err
	}

	response := make([]dto.ModeloComunicacaoResponseDTO, 0, len(modelos))

	for _, modelo := range modelos {

		dtoResp, err := domainToResponse(*modelo)

		if err != nil {
			return nil, err
		}

		response = append(response, *dtoResp)
	}

	return &response, nil
}

func domainToResponse(m Comunicacao) (*dto.ModeloComunicacaoResponseDTO, error) {
	return &dto.ModeloComunicacaoResponseDTO{
		Id:              m.Id.String(),
		Nome:            string(m.Nome),
		TipoComunicacao: string(m.TipoComunicacao),
		Assunto:         m.Assunto,
		Corpo:           m.Corpo,
		Ativo:           parseStatusModeloComunicacaoString(m.Ativo),
	}, nil
}

func parseStatusModeloComunicacao(s string) (StatusModeloComunicacao, error) {
	switch s {
	case "ATIVO":
		return StatusAtivo, nil
	case "INATIVO":
		return StatusInativo, nil
	}

	return 0, ErrorInvalidStatus
}

func parseStatusModeloComunicacaoString(s StatusModeloComunicacao) string {
	switch s {
	case StatusAtivo:
		return "ATIVO"
	case StatusInativo:
		return "INATIVO"
	}

	return "UNKNOWN"
}

func parseTipoComunicacao(s string) (TipoComunicacao, error) {
	switch s {
	case "EMAIL":
		return Email, nil
	case "SMS":
		return SMS, nil
	}

	return "", ErrorInvalidTipoComunicacao
}

func parseNomeEnvio(s string) (NomeEnvio, error) {
	switch s {
	case "Boas Vindas":
		return NomeEnvio(BoasVindas), nil
	case "Novo Plantão":
		return NomeEnvio(NovoPlantao), nil
	case "Cadastro Atualizado":
		return NomeEnvio(CadastroAtualizado), nil
	case "Cadastro Excluido":
		return NomeEnvio(CadastroExcluido), nil
	default:
		return "", fmt.Errorf("tipo de envio inválido: %s", s)
	}
}
