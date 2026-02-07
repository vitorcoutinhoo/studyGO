package colaborador

import (
	"context"
	"fmt"
	"plantao/internal/api/dto"
	"plantao/utils"
	"strings"

	"github.com/google/uuid"
)

// Serviço para gerenciar colaboradores
type ColaboradorService struct {
	repository ColaboradorRepository
}

// Cria uma nova instância do serviço de colaborador
func NewColaboradorService(repository ColaboradorRepository) *ColaboradorService {
	return &ColaboradorService{
		repository: repository,
	}
} // Fim NewColaboradorService

// Cria um novo colaborador com validações e armazenamento
func (s *ColaboradorService) CreateColaborador(ctx context.Context, colaboradorDTO *dto.CreateColaboradorRequest) (*Colaborador, error) {
	colaborador, err := createColaboradorDtoToDomain(colaboradorDTO)

	if err != nil {
		return nil, err
	}

	if s.repository.ExistsEmail(ctx, colaborador.Email) {
		return nil, ErrorInvalidEmail
	}

	col, err := NewColaborador(
		colaborador.Nome,
		colaborador.Email,
		colaborador.Telefone,
		colaborador.Cargo,
		colaborador.Setor,
		colaborador.Foto,
		colaborador.DataAdmissao,
		colaborador.DataDesligamento,
	)

	if err != nil {
		return nil, err
	}

	if err := s.repository.Store(ctx, col); err != nil {
		return nil, err
	}

	return col, nil
} // Fim CreateColaborador

// Atualiza um colaborador existente com novas informações
func (s *ColaboradorService) UpdateColaborador(ctx context.Context, colaboradorDTO *dto.UpdateColaboradorRequest, colaboradorId string) (*Colaborador, error) {
	id, err := uuid.Parse(colaboradorId)

	if err != nil {
		return nil, fmt.Errorf("UUID inválido: %v", err)
	}

	colaborador, err := s.repository.FindById(ctx, id)

	if err != nil {
		return nil, err
	}

	if colaborador == nil {
		return nil, ErrorColaboradorNotFound
	}

	if s.repository.ExistsEmailExcludingId(ctx, colaborador.Email, colaborador.Id) {
		return nil, ErrorInvalidEmail
	}

	sts, err := statusColaboradorFromDTO(colaboradorDTO.Status)

	if err != nil {
		return nil, err
	}

	err = colaborador.UpdateDados(
		colaboradorDTO.Email,
		colaboradorDTO.Telefone,
		colaboradorDTO.Cargo,
		colaboradorDTO.Setor,
		colaboradorDTO.Foto,
		sts,
	)

	if err != nil {
		return nil, err
	}

	if err := s.repository.Update(ctx, colaborador); err != nil {
		return nil, err
	}

	return colaborador, nil
} // Fim UpdateColaborador

// Desativa um colaborador pelo ID
func (s *ColaboradorService) DisableColaborador(ctx context.Context, colaboradorId string) error {
	id, err := uuid.Parse(colaboradorId)

	if err != nil {
		return fmt.Errorf("UUID inválido: %v", err)
	}

	if !s.repository.ExistsId(ctx, id) {
		return ErrorColaboradorNotFound
	}

	if err := s.repository.Disable(ctx, id); err != nil {
		return err
	}

	return nil
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
// Permite filtrar por nome, email, telefone, cargo, departamento e data de admissão
func (s *ColaboradorService) GetColaboradorByFilter(ctx context.Context, filterReq dto.GetColaboradoresByFilterRequest) ([]Colaborador, error) {
	filter := filterDtoToFilterDomain(filterReq)

	fmt.Println(filter)

	colaboradores, err := s.repository.FindByFilter(ctx, filter)

	if err != nil {
		return nil, err
	}

	return colaboradores, nil
} // Fim GetVolaboradorByFilter

// Converte a requisição de criação de colaborador para o domínio
func createColaboradorDtoToDomain(r *dto.CreateColaboradorRequest) (*Colaborador, error) {
	dataAdmissao, err := utils.ParseBrToUsDate(r.DataAdmissao)

	if err != nil {
		return nil, err
	}

	dataDesligamento, err := utils.ParseBrToUsDate(r.DataDesligamento)

	if err != nil {
		return nil, err
	}

	return &Colaborador{
		Nome:             r.Nome,
		Email:            r.Email,
		Telefone:         r.Telefone,
		Cargo:            r.Cargo,
		Setor:            r.Setor,
		Foto:             r.Foto,
		DataAdmissao:     dataAdmissao,
		DataDesligamento: dataDesligamento,
	}, nil
} // Fim toDomain

// Converte o status do colaborador a partir do DTO
func statusColaboradorFromDTO(value string) (StatusColaborador, error) {
	switch strings.ToLower(value) {
	case "ativo":
		return StatusAtivo, nil
	case "inativo":
		return StatusInativo, nil
	default:
		return 0, ErrorInvalidStatus
	}
} // Fim StatusColaboradorFromDTO

// Converte a requisição de filtro de colaboradores para o domínio
func filterDtoToFilterDomain(filterReq dto.GetColaboradoresByFilterRequest) ColaboradorFilter {
	filter := ColaboradorFilter{
		Nome:         filterReq.Nome,
		Email:        filterReq.Email,
		Telefone:     filterReq.Telefone,
		Cargo:        filterReq.Cargo,
		Departamento: filterReq.Setor,
	}

	return filter
} // Fim filterDtoToFilterDomain
