package colaborador

import (
	"context"
	"fmt"
	"plantao/internal/api/dto"
	"plantao/utils"
	"strings"
	"time"

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
func (s *ColaboradorService) CreateColaborador(ctx context.Context, colaboradorDTO *dto.CreateColaboradorRequest) (*dto.ColaboradorResponse, error) {
	colaborador, err := createColaboradorDtoToDomain(colaboradorDTO)

	if err != nil {
		return nil, err
	}

	if s.repository.ExistsEmail(ctx, colaborador.Email) {
		return nil, ErrorEmailAlreadyExists
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
		colaborador.Status,
		colaborador.AtivoPlantao,
	)

	if err != nil {
		return nil, err
	}

	colaboradorReturn, err := s.repository.Store(ctx, col)

	if err != nil {
		return nil, err
	}

	colaboradorResponse, err := colaboradorToResponse(colaboradorReturn)

	if err != nil {
		return nil, err
	}

	return colaboradorResponse, nil
} // Fim CreateColaborador

// Atualiza um colaborador existente com novas informações
func (s *ColaboradorService) UpdateColaborador(ctx context.Context, colaboradorDTO *dto.UpdateColaboradorRequest, colaboradorId string) error {
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

	if s.repository.ExistsEmailExcludingId(ctx, colaborador.Email, colaborador.Id) {
		return ErrorEmailAlreadyExists
	}

	colaboradorDomain, err := updateColaboradorDtoToDomain(colaboradorDTO)

	if err != nil {
		return nil
	}

	err = colaborador.UpdateDados(
		&colaboradorDomain.Nome,
		&colaboradorDomain.Email,
		&colaboradorDomain.Telefone,
		&colaboradorDomain.Cargo,
		&colaboradorDomain.Setor,
		&colaboradorDomain.Foto,
		colaboradorDomain.DataAdmissao,
		colaboradorDomain.DataDesligamento,
		&colaboradorDomain.Status,
		&colaboradorDomain.AtivoPlantao,
	)

	if err != nil {
		return err
	}

	if err := s.repository.Update(ctx, colaborador); err != nil {
		return err
	}

	return nil
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
func (s *ColaboradorService) GetColaboradorById(ctx context.Context, colaboradorId string) (*dto.ColaboradorResponse, error) {
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

	colaboradorResponse, err := colaboradorToResponse(colaborador)

	if err != nil {
		return nil, err
	}

	return colaboradorResponse, nil
} // Fim GetColaboradorById

// Recupera colaboradores com base em filtros opcionais
// Permite filtrar por nome, email, telefone, cargo, departamento e data de admissão
func (s *ColaboradorService) GetColaboradorByFilter(ctx context.Context, filterReq dto.GetColaboradoresByFilterRequest) ([]dto.ColaboradorResponse, error) {
	filter := filterDtoToFilterDomain(filterReq)

	colaboradores, err := s.repository.FindByFilter(ctx, filter)

	if err != nil {
		return nil, err
	}

	responses := make([]dto.ColaboradorResponse, 0, len(colaboradores))

	for i := range colaboradores {
		resp, err := colaboradorToResponse(&colaboradores[i])
		if err != nil {
			return nil, err
		}

		responses = append(responses, *resp)
	}

	return responses, nil
} // Fim GetVolaboradorByFilter

// Converte a requisição de criação de colaborador para o domínio
func createColaboradorDtoToDomain(r *dto.CreateColaboradorRequest) (*Colaborador, error) {
	dataAdmissao, err := utils.ParseBrToUsDate(&r.DataAdmissao)

	if err != nil {
		return nil, err
	}

	var dataDesligamento *time.Time

	if r.DataDesligamento != nil {
		dataTemp, err := utils.ParseBrToUsDate(r.DataDesligamento)

		if err != nil {
			return nil, err
		}

		dataDesligamento = dataTemp
	} else {
		dataDesligamento = nil
	}

	ativo, err := statusColaboradorFromDTO(r.Status)

	if err != nil {
		return nil, err
	}

	ativoPlantao, err := statusColaboradorFromDTO(r.AtivoPlantao)

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
		Status:           ativo,
		AtivoPlantao:     ativoPlantao,
		DataAdmissao:     dataAdmissao,
		DataDesligamento: dataDesligamento,
	}, nil
} // Fim toDomain

// Conversor do DTO de UPDATE para uma entidade do DOMÍNIO
func updateColaboradorDtoToDomain(r *dto.UpdateColaboradorRequest) (*Colaborador, error) {
	dataAdmissao, err := utils.ParseBrToUsDate(r.DataAdmissao)

	if err != nil {
		return nil, err
	}

	var dataDesligamento *time.Time

	if r.DataDesligamento != nil {
		dataTemp, err := utils.ParseBrToUsDate(r.DataDesligamento)

		if err != nil {
			return nil, err
		}

		dataDesligamento = dataTemp
	} else {
		dataDesligamento = nil
	}

	ativo, err := statusColaboradorFromDTO(*r.Status)

	if err != nil {
		return nil, err
	}

	ativoPlantao, err := statusColaboradorFromDTO(*r.AtivoPlantao)

	if err != nil {
		return nil, err
	}

	return &Colaborador{
		Nome:             *r.Nome,
		Email:            *r.Email,
		Telefone:         *r.Telefone,
		Cargo:            *r.Cargo,
		Setor:            *r.Setor,
		Foto:             *r.Foto,
		Status:           ativo,
		AtivoPlantao:     ativoPlantao,
		DataAdmissao:     dataAdmissao,
		DataDesligamento: dataDesligamento,
	}, nil
} // Fim updateColaboradorDtoToDomain

// Converte um colaborador do domínio para um DTO
func colaboradorToResponse(c *Colaborador) (*dto.ColaboradorResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("colaborador vazio ou nulo")
	}

	status, err := statusColaboradorToDTO(c.Status)

	if err != nil {
		return nil, err
	}

	statusPlantao, err := statusColaboradorToDTO(c.AtivoPlantao)

	if err != nil {
		return nil, err
	}

	dataAdmissao, err := utils.ParseUsToBrDate(c.DataAdmissao)

	if err != nil {
		return nil, err
	}

	var dataDeligamento string

	if c.DataDesligamento != nil {
		dataTemp, err := utils.ParseUsToBrDate(c.DataDesligamento)

		if err != nil {
			return nil, err
		}

		dataDeligamento = dataTemp
	}

	return &dto.ColaboradorResponse{
		Id:               c.Id.String(),
		Nome:             c.Nome,
		Email:            c.Email,
		Telefone:         c.Telefone,
		Cargo:            c.Cargo,
		Setor:            c.Setor,
		Foto:             c.Foto,
		Status:           status,
		AtivoPlantao:     statusPlantao,
		DataAdmissao:     dataAdmissao,
		DataDesligamento: dataDeligamento,
	}, nil
} // Fim colaboradorToResponse

// Converte um status do domínio para um status do DTO(Texto)
func statusColaboradorToDTO(status StatusColaborador) (string, error) {
	switch status {
	case 1:
		return "ativo", nil
	case 2:
		return "inativo", nil
	default:
		return "", ErrorInvalidStatus
	}
} // Fim statusColaboradorToDTO

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
