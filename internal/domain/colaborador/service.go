package colaborador

import (
	"context"
	"plantao/internal/api/controller"
	"plantao/utils"
	"strings"
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
func (s *ColaboradorService) CreateColaborador(ctx context.Context, colaboradorDTO *controller.CreateColaboradorRequest) (*Colaborador, error) {
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
func (s *ColaboradorService) UpdateColaborador(ctx context.Context, colaboradorDTO *controller.UpdateColaboradorRequest) (*Colaborador, error) {
	colaborador, err := s.repository.FindById(ctx, colaboradorDTO.Id)

	if err != nil {
		return nil, err
	}

	if colaborador == nil {
		return nil, ErrorColaboradorNotFound
	}

	if s.repository.ExistsEmailExcludingId(ctx, colaborador.Email, colaborador.Id) {
		return nil, ErrorInvalidEmail
	}

	sts, err := StatusColaboradorFromDTO(colaboradorDTO.Status)

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

// Deleta um colaborador pelo ID
func (s *ColaboradorService) DeleteColaborador(ctx context.Context, colaboradorId string) error {
	if !s.repository.ExistsId(ctx, colaboradorId) {
		return ErrorColaboradorNotFound
	}

	if err := s.repository.Delete(ctx, colaboradorId); err != nil {
		return err
	}

	return nil
} // Fim DeleteColaborador

// Recupera um colaborador pelo ID
func (s *ColaboradorService) GetColaboradorById(ctx context.Context, colaboradorId string) (*Colaborador, error) {
	colaborador, _ := s.repository.FindById(ctx, colaboradorId)

	if colaborador == nil {
		return nil, ErrorColaboradorNotFound
	}

	return colaborador, nil
} // Fim GetColaboradorById

// Recupera todos os colaboradores
func (s *ColaboradorService) GetAllColaboradores(ctx context.Context) ([]Colaborador, error) {
	colaboradores, err := s.repository.Find(ctx)

	if err != nil {
		return nil, err
	}

	return colaboradores, nil
} // Fim GetAllColaboradores

// Converte a requisição de criação de colaborador para o domínio
func createColaboradorDtoToDomain(r *controller.CreateColaboradorRequest) (*Colaborador, error) {
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

// Converte a requisição de atualização de colaborador para o domínio
func updateColaboradorDtoToDomain(r *controller.UpdateColaboradorRequest) (*Colaborador, error) {
	dataAdmissao, err := utils.ParseBrToUsDate(r.DataAdmissao)

	if err != nil {
		return nil, err
	}

	dataDesligamento, err := utils.ParseBrToUsDate(r.DataDesligamento)

	if err != nil {
		return nil, err
	}

	sts, err := StatusColaboradorFromDTO(r.Status)

	if err != nil {
		return nil, err
	}

	return &Colaborador{
		Id:               r.Id,
		Email:            r.Email,
		Telefone:         r.Telefone,
		Cargo:            r.Cargo,
		Setor:            r.Setor,
		Foto:             r.Foto,
		Status:           sts,
		DataAdmissao:     dataAdmissao,
		DataDesligamento: dataDesligamento,
	}, nil
} // Fim toDomainUpdate

// Converte o status do colaborador a partir do DTO
func StatusColaboradorFromDTO(value string) (StatusColaborador, error) {
	switch strings.ToLower(value) {
	case "ativo":
		return StatusAtivo, nil
	case "inativo":
		return StatusInativo, nil
	default:
		return 0, ErrorInvalidStatus
	}
} // Fim StatusColaboradorFromDTO
