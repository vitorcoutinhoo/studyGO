package colaborador

import (
	"context"
	"fmt"
	"strings"

	"plantao/internal/domain/comunicacao"

	"github.com/google/uuid"
)

// Serviço para gerenciar colaboradores
type ColaboradorService struct {
	repository   ColaboradorRepository
	envioService *comunicacao.EnvioService
}

// Cria uma nova instância do serviço de colaborador
func NewColaboradorService(repository ColaboradorRepository, envioService *comunicacao.EnvioService) *ColaboradorService {
	return &ColaboradorService{
		repository:   repository,
		envioService: envioService,
	}
} // Fim NewColaboradorService

// Cria um novo colaborador com validações e armazenamento
func (s *ColaboradorService) CreateColaborador(ctx context.Context, col *Colaborador) (*Colaborador, error) {
	exists, err := s.repository.ExistsEmail(ctx, col.Email)

	if err != nil {
		return nil, fmt.Errorf("erro ao verificar existência de email: %w", err)
	}

	if exists {
		return nil, ErrorEmailAlreadyExists
	}

	colaborador, err := NewColaborador(
		col.Nome,
		col.Email,
		col.Telefone,
		col.Foto,
		col.DataAdmissao,
		col.DataDesligamento,
		col.Status,
		col.AtivoPlantao,
		col.Cargo,
		col.Setor,
	)

	if err != nil {
		return nil, err
	}

	colaboradorReturn, err := s.repository.Store(ctx, colaborador)

	if err != nil {
		return nil, err
	}

	// err = s.envioService.SendEmailComunicacao(
	// 	ctx,
	// 	"Boas Vindas",
	// 	colaboradorReturn.Id.String(),
	// 	colaboradorReturn.Email,
	// 	colaboradorReturn.Nome,
	// )

	if err != nil {
		return nil, err
	}

	return colaboradorReturn, nil
} // Fim CreateColaborador

// Atualiza um colaborador existente com novas informações
func (s *ColaboradorService) UpdateColaborador(ctx context.Context, col *Colaborador, colaboradorId string) error {
	id, err := uuid.Parse(colaboradorId)

	if err != nil {
		return fmt.Errorf("UUID inválido: %v", err)
	}

	colaborador, err := s.repository.FindById(ctx, id)

	if err != nil {
		return err
	}

	if colaborador == nil {
		return ErrorColaboradorNotFound
	}

	exists, err := s.repository.ExistsEmailExcludingId(ctx, colaborador.Email, colaborador.Id)

	if err != nil {
		return fmt.Errorf("erro ao verificar existência de email excluindo ID: %w", err)
	}

	if exists {
		return ErrorEmailAlreadyExists
	}

	err = colaborador.UpdateDados(
		&col.Nome,
		&col.Email,
		&col.Telefone,
		&col.Foto,
		col.DataAdmissao,
		col.DataDesligamento,
		&col.Status,
		&col.AtivoPlantao,
		&col.Cargo,
		&col.Setor,
	)

	if err != nil {
		return err
	}

	return s.repository.Update(ctx, colaborador)
} // Fim UpdateColaborador

// Desativa um colaborador pelo ID
func (s *ColaboradorService) DisableColaborador(ctx context.Context, colaboradorId string) error {
	id, err := uuid.Parse(colaboradorId)

	if err != nil {
		return fmt.Errorf("UUID inválido: %v", err)
	}

	exists, err := s.repository.ExistsId(ctx, id)

	if err != nil {
		return fmt.Errorf("erro ao verificar existência de ID: %w", err)
	}

	if !exists {
		return ErrorColaboradorNotFound
	}

	return s.repository.Disable(ctx, id)
} // Fim DisableColaborador

// Recupera um colaborador pelo ID
func (s *ColaboradorService) GetColaboradorById(ctx context.Context, colaboradorId string) (*Colaborador, error) {
	id, err := uuid.Parse(colaboradorId)

	if err != nil {
		return nil, fmt.Errorf("UUID inválido: %v", err)
	}

	colaborador, err := s.repository.FindById(ctx, id)

	if colaborador == nil {
		return nil, ErrorColaboradorNotFound
	}

	if err != nil {
		return nil, err
	}

	return colaborador, nil
} // Fim GetColaboradorById

// Recupera colaboradores com base em filtros opcionais
func (s *ColaboradorService) GetColaboradorByFilter(ctx context.Context, filter ColaboradorFilter) ([]Colaborador, error) {
	return s.repository.FindByFilter(ctx, filter)
} // Fim GetColaboradorByFilter

// Converte string para StatusColaborador
func ParseStatusColaborador(value string) (StatusColaborador, error) {
	switch strings.ToLower(value) {
	case "ativo":
		return StatusAtivo, nil
	case "inativo":
		return StatusInativo, nil
	default:
		return 0, ErrorInvalidStatus
	}
} // Fim ParseStatusColaborador

// Converte StatusColaborador para string
func StatusColaboradorString(status StatusColaborador) (string, error) {
	switch status {
	case StatusAtivo:
		return "ativo", nil
	case StatusInativo:
		return "inativo", nil
	default:
		return "", ErrorInvalidStatus
	}
} // Fim StatusColaboradorString

func ParseCargoColaborador(s string) (CargoColaborador, error) {
	switch s {
	case string(CargoAnalista):
		return CargoAnalista, nil
	case string(CargoGerente):
		return CargoGerente, nil
	case string(CargoConsultor):
		return CargoConsultor, nil
	case string(CargoTecnico):
		return CargoTecnico, nil
	case string(CargoOutro):
		return CargoOutro, nil
	case string(CargoDesenvolvedorFrontend):
		return CargoDesenvolvedorFrontend, nil
	case string(CargoDesenvolvedorBackend):
		return CargoDesenvolvedorBackend, nil
	case string(CargoDesenvolvedorFullstack):
		return CargoDesenvolvedorFullstack, nil
	default:
		return "", fmt.Errorf("cargo inválido: %s", s)
	}
}

func ParseSetorColaborador(s string) (SetorColaborador, error) {
	switch s {
	case string(SetorRH):
		return SetorRH, nil
	case string(SetorTI):
		return SetorTI, nil
	case string(SetorFinanceiro):
		return SetorFinanceiro, nil
	case string(SetorSuporte):
		return SetorSuporte, nil
	case string(SetorDesenvolvimento):
		return SetorDesenvolvimento, nil
	case string(SetorDiretoria):
		return SetorDiretoria, nil
	default:
		return "", fmt.Errorf("setor inválido: %s", s)
	}
}
