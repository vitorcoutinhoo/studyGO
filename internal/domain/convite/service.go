package convite

import (
	"context"
	"fmt"
	"plantao/internal/domain/comunicacao"
	domainlog "plantao/internal/domain/log"
	"plantao/internal/infra/config"

	"github.com/google/uuid"
)

type ConviteService struct {
	respository  ConviteRepository
	envioService *comunicacao.EnvioService
	urlServer    string
	log          domainlog.Logger
}

func NewConviteService(respository ConviteRepository, envioService *comunicacao.EnvioService, cfg *config.Config, log domainlog.Logger) *ConviteService {
	return &ConviteService{
		respository:  respository,
		envioService: envioService,
		urlServer:    cfg.Frontend.URL + "/cadastro?token=",
		log:          log,
	}
}

func (s *ConviteService) CreateConvite(ctx context.Context, colaboradorId, colaboradorNome, colaboradorEmail string) error {
	s.log.Info("iniciando criação de convite", "id_colaborador", colaboradorId, "email", colaboradorEmail)

	id, err := uuid.Parse(colaboradorId)
	if err != nil {
		s.log.Warn("UUID do colaborador inválido para criação de convite", "id_colaborador", colaboradorId, "error", err)
		return fmt.Errorf("UUID inválido: %v", err)
	}

	s.log.Info("salvando convite no banco de dados", "id_colaborador", id)
	conv, err := s.respository.Store(ctx, id)
	if err != nil {
		s.log.Error("erro ao salvar convite", "id_colaborador", id, "error", err)
		return fmt.Errorf("erro ao salvar convite: %v", err)
	}

	s.log.Info("preparando envio de email do convite", "id_colaborador", id, "email", colaboradorEmail)
	go func() {
		data := map[string]any{
			string(comunicacao.Nome):  colaboradorNome,
			string(comunicacao.Email): colaboradorEmail,
			string(comunicacao.Link):  s.urlServer + conv.Token.String(),
		}

		err := s.envioService.SendEmailComunicacao(
			context.Background(),
			comunicacao.ColaboradorCadastrado,
			colaboradorEmail,
			id,
			data,
		)

		if err != nil {
			s.log.Warn("erro ao enviar email de convite", "id_colaborador", id, "email", colaboradorEmail, "error", err)
		}
	}()

	if conv != nil {
		s.log.Info("convite criado com sucesso", "id_colaborador", id, "expira_em", conv.ExpiraEm)
	} else {
		s.log.Info("convite criado com sucesso", "id_colaborador", id)
	}
	return nil
}

func (s *ConviteService) GetAllConvites(ctx context.Context) ([]Convite, error) {
	s.log.Info("buscando todos os convites")

	convites, err := s.respository.FindAll(ctx)
	if err != nil {
		s.log.Error("erro ao buscar convites", "error", err)
		return nil, err
	}

	s.log.Info("convites encontrados com sucesso", "quantidade", len(convites))
	return convites, nil
}

func (s *ConviteService) DisableConvite(ctx context.Context, token string) error {
	s.log.Info("iniciando desativação de convite")

	tk, err := uuid.Parse(token)
	if err != nil {
		s.log.Warn("token de convite inválido para desativação", "error", err)
		return fmt.Errorf("UUID inválido: %v", err)
	}

	if err := s.respository.Disable(ctx, tk); err != nil {
		s.log.Error("erro ao desativar convite", "error", err)
		return err
	}

	s.log.Info("convite desativado com sucesso")
	return nil
}
