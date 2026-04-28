package convite

import (
	"context"
	"fmt"
	"log"
	"plantao/internal/domain/comunicacao"
	"plantao/internal/infra/config"

	"github.com/google/uuid"
)

type ConviteService struct {
	respository  ConviteRepository
	envioService *comunicacao.EnvioService
	urlServer    string
}

func NewConviteService(respository ConviteRepository, envioService *comunicacao.EnvioService, cfg *config.Config) *ConviteService {
	return &ConviteService{
		respository:  respository,
		envioService: envioService,
		urlServer:    "http://" + cfg.Server.Host + ":" + cfg.Server.Port + "/api/v1/usuarios/cadastro?token=",
	}
}

func (s *ConviteService) CreateConvite(ctx context.Context, colaboradorId, colaboradorNome, colaboradorEmail string) error {
	id, err := uuid.Parse(colaboradorId)
	if err != nil {
		return fmt.Errorf("UUID inválido: %v", err)
	}

	conv, err := s.respository.Store(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao salvar convite: %v", err)
	}

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
			data,
		)

		if err != nil {
			log.Printf("[email] erro ao enviar boas-vindas para %s: %v", colaboradorEmail, err)
		}
	}()

	return nil
}

func (s *ConviteService) GetAllConvites(ctx context.Context) ([]Convite, error) {
	return s.respository.FindAll(ctx)
}

func (s *ConviteService) DisableConvite(ctx context.Context, token string) error {
	tk, err := uuid.Parse(token)
	if err != nil {
		return fmt.Errorf("UUID inválido: %v", err)
	}

	return s.respository.Disable(ctx, tk)
}
