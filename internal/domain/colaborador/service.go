package colaborador

import (
	"context"
	"fmt"
	"io"
	"strings"

	"plantao/internal/domain/comunicacao"
	"plantao/internal/domain/convite"
	"plantao/internal/domain/log"
	"plantao/internal/infra/config"

	"github.com/google/uuid"
)

// Serviço para gerenciar colaboradores
type ColaboradorService struct {
	repository        ColaboradorRepository
	envioService      *comunicacao.EnvioService
	storageImage      FileStorage
	conviteRepository convite.ConviteRepository
	urlServer         string
	log               log.Logger
}

// Cria uma nova instância do serviço de colaborador
func NewColaboradorService(repository ColaboradorRepository, envioService *comunicacao.EnvioService, storageImage FileStorage, conviteRepository convite.ConviteRepository, cfg *config.Config, log log.Logger) *ColaboradorService {
	return &ColaboradorService{
		repository:        repository,
		envioService:      envioService,
		storageImage:      storageImage,
		conviteRepository: conviteRepository,
		urlServer:         cfg.Frontend.URL + "/cadastro?token=",
		log:               log,
	}
} // Fim NewColaboradorService

// Cria um novo colaborador com validações e armazenamento
func (s *ColaboradorService) CreateColaborador(ctx context.Context, col *Colaborador, file io.ReadSeeker) (*Colaborador, error) {
	s.log.Info("verificando existência de email do colaborador no banco de dados")
	exists, err := s.repository.ExistsEmail(ctx, col.Email)

	if err != nil {
		s.log.Error("erro ao verificar existência do email do colaborador", "email", col.Email, "error", err)
		return nil, fmt.Errorf("erro ao verificar existência de email: %w", err)
	}

	if exists {
		s.log.Warn("email do colaborador já está cadastrado", "email", col.Email)
		return nil, ErrorEmailAlreadyExists
	}

	s.log.Info("criando novo abjeto de colaborador")
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
		s.log.Error("erro ao criar objeto de colaborador", "error", err)
		return nil, err
	}

	s.log.Info("verificando a existência do arquivo de imagen do colaborador")
	if file != nil {
		s.log.Info("salvando aquivo de imagem do colaborador")
		url, err := s.storageImage.Save(file, colaborador.Foto)

		if err != nil {
			s.log.Error("erro ao salvar arquivo de imagem do colaborador", "error", err)
			return nil, err
		}
		s.log.Info("arquivo de imagem salva com sucesso")

		colaborador.Foto = url
	}

	s.log.Info("salvando colaborador no banco de dados")
	colaboradorReturn, err := s.repository.Store(ctx, colaborador)

	if err != nil {
		s.log.Error("erro ao salvar colaborador no banco de dados", "error", err)
		return nil, err
	}

	s.log.Info("salvando convite para colaborador no banco de dados")
	conv, err := s.conviteRepository.Store(ctx, colaboradorReturn.Id)
	if err != nil {
		s.log.Error("erro ao salvar convite para colaborador", "id_colaborador", colaboradorReturn.Id, "error", err)
		return nil, err
	}

	s.log.Info("preparando envio de email para colaboraro", "id_colaborador", colaboradorReturn.Id)
	go func() {
		data := map[string]any{
			string(comunicacao.Nome):  colaboradorReturn.Nome,
			string(comunicacao.Email): colaboradorReturn.Email,
			string(comunicacao.Link):  s.urlServer + conv.Token.String(),
		}

		s.log.Info("enviando email para colaborador", "id_colaborador", colaboradorReturn.Id)
		err := s.envioService.SendEmailComunicacao(
			context.Background(),
			comunicacao.ColaboradorCadastrado,
			colaboradorReturn.Email,
			colaboradorReturn.Id,
			data,
		)

		if err != nil {
			s.log.Warn("erro ao enviar email de boas-vindas para colaborador", "email", colaboradorReturn.Email, "error", err)
		}
	}()

	s.log.Info("colaborador inserido com sucesso", "id_colaborador", colaboradorReturn.Id)
	return colaboradorReturn, nil
} // Fim CreateColaborador

// Atualiza um colaborador existente com novas informações
func (s *ColaboradorService) UpdateColaborador(ctx context.Context, col *Colaborador, colaboradorId string, file io.ReadSeeker) error {
	s.log.Info("iniciando atualização de colaborador", "id_colaborador", colaboradorId)

	id, err := uuid.Parse(colaboradorId)

	if err != nil {
		s.log.Warn("UUID do colaborador inválido para atualização", "id_colaborador", colaboradorId, "error", err)
		return fmt.Errorf("UUID inválido: %v", err)
	}

	s.log.Info("buscando colaborador para atualização", "id_colaborador", id)
	colaborador, err := s.repository.FindById(ctx, id)

	if err != nil {
		s.log.Error("erro ao buscar colaborador para atualização", "id_colaborador", id, "error", err)
		return err
	}

	if colaborador == nil {
		s.log.Warn("colaborador não encontrado para atualização", "id_colaborador", id)
		return ErrorColaboradorNotFound
	}

	s.log.Info("verificando existência de email excluindo colaborador", "id_colaborador", colaborador.Id, "email", colaborador.Email)
	exists, err := s.repository.ExistsEmailExcludingId(ctx, colaborador.Email, colaborador.Id)

	if err != nil {
		s.log.Error("erro ao verificar existência de email excluindo colaborador", "id_colaborador", colaborador.Id, "email", colaborador.Email, "error", err)
		return fmt.Errorf("erro ao verificar existência de email excluindo ID: %w", err)
	}

	if exists {
		s.log.Warn("email do colaborador já está cadastrado em outro registro", "id_colaborador", colaborador.Id, "email", colaborador.Email)
		return ErrorEmailAlreadyExists
	}

	imagemAnterior := colaborador.Foto

	s.log.Info("atualizando dados do colaborador em memória", "id_colaborador", colaborador.Id)
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
		s.log.Warn("erro ao validar atualização do colaborador", "id_colaborador", colaborador.Id, "error", err)
		return err
	}

	if file != nil {
		s.log.Info("salvando nova imagem do colaborador", "id_colaborador", colaborador.Id)
		url, err := s.storageImage.Save(file, colaborador.Foto)

		if err != nil {
			s.log.Error("erro ao salvar nova imagem do colaborador", "id_colaborador", colaborador.Id, "error", err)
			return err
		}

		colaborador.Foto = url

		s.log.Info("removendo imagem anterior do colaborador", "id_colaborador", colaborador.Id, "imagem_anterior", imagemAnterior)
		err = s.storageImage.Delete(imagemAnterior)

		if err != nil {
			s.log.Error("erro ao remover imagem anterior do colaborador", "id_colaborador", colaborador.Id, "imagem_anterior", imagemAnterior, "error", err)
			return err
		}
	}

	s.log.Info("salvando atualização do colaborador no banco de dados", "id_colaborador", colaborador.Id)
	if err := s.repository.Update(ctx, colaborador); err != nil {
		s.log.Error("erro ao atualizar colaborador no banco de dados", "id_colaborador", colaborador.Id, "error", err)
		return err
	}

	s.log.Info("colaborador atualizado com sucesso", "id_colaborador", colaborador.Id)
	return nil
} // Fim UpdateColaborador

// Desativa um colaborador pelo ID
func (s *ColaboradorService) DisableColaborador(ctx context.Context, colaboradorId string) error {
	s.log.Info("iniciando desativação de colaborador", "id_colaborador", colaboradorId)

	id, err := uuid.Parse(colaboradorId)

	if err != nil {
		s.log.Warn("UUID do colaborador inválido para desativação", "id_colaborador", colaboradorId, "error", err)
		return fmt.Errorf("UUID inválido: %v", err)
	}

	s.log.Info("verificando existência de colaborador para desativação", "id_colaborador", id)
	exists, err := s.repository.ExistsId(ctx, id)

	if err != nil {
		s.log.Error("erro ao verificar existência de colaborador por ID", "id_colaborador", id, "error", err)
		return fmt.Errorf("erro ao verificar existência de ID: %w", err)
	}

	if !exists {
		s.log.Warn("colaborador não encontrado para desativação", "id_colaborador", id)
		return ErrorColaboradorNotFound
	}

	// Buscar colaborador
	s.log.Info("buscando colaborador para desativação", "id_colaborador", id)
	col, err := s.repository.FindById(ctx, id)
	if err != nil {
		s.log.Error("erro ao buscar colaborador para desativação", "id_colaborador", id, "error", err)
		return fmt.Errorf("erro ao buscar colaborador: %w", err)
	}

	// Se existir foto, deletar
	if col.Foto != "" {
		s.log.Info("deletando foto do colaborador antes da desativação", "id_colaborador", id, "foto", col.Foto)
		err = s.storageImage.Delete(col.Foto)
		if err != nil {
			s.log.Error("erro ao deletar foto do colaborador", "id_colaborador", id, "foto", col.Foto, "error", err)
			return fmt.Errorf("erro ao deletar foto: %w", err)
		}
	}

	if err := s.repository.Disable(ctx, id); err != nil {
		s.log.Error("erro ao desativar colaborador", "id_colaborador", id, "error", err)
		return err
	}

	s.log.Info("colaborador desativado com sucesso", "id_colaborador", id)
	return nil
} // Fim DisableColaborador

// Recupera um colaborador pelo ID
func (s *ColaboradorService) GetColaboradorById(ctx context.Context, colaboradorId string) (*Colaborador, error) {
	s.log.Info("buscando colaborador por ID", "id_colaborador", colaboradorId)

	id, err := uuid.Parse(colaboradorId)

	if err != nil {
		s.log.Warn("UUID do colaborador inválido para busca", "id_colaborador", colaboradorId, "error", err)
		return nil, fmt.Errorf("UUID inválido: %v", err)
	}

	colaborador, err := s.repository.FindById(ctx, id)

	if colaborador == nil {
		s.log.Warn("colaborador não encontrado", "id_colaborador", id)
		return nil, ErrorColaboradorNotFound
	}

	if err != nil {
		s.log.Error("erro ao buscar colaborador por ID", "id_colaborador", id, "error", err)
		return nil, err
	}

	s.log.Info("colaborador encontrado com sucesso", "id_colaborador", id)
	return colaborador, nil
} // Fim GetColaboradorById

// Recupera colaboradores com base em filtros opcionais
func (s *ColaboradorService) GetColaboradorByFilter(ctx context.Context, filter ColaboradorFilter) ([]Colaborador, error) {
	s.log.Info("buscando colaboradores por filtro", "filter", filter)

	colaboradores, err := s.repository.FindByFilter(ctx, filter)
	if err != nil {
		s.log.Error("erro ao buscar colaboradores por filtro", "filter", filter, "error", err)
		return nil, err
	}

	s.log.Info("colaboradores encontrados com sucesso", "quantidade", len(colaboradores))
	return colaboradores, nil
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
