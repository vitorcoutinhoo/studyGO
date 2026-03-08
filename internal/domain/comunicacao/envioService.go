package comunicacao

import (
	"context"
	"plantao/internal/domain/email"

	"github.com/google/uuid"
)

type EnvioService struct {
	envioRepository  EnvioComunicacaoRepository
	emailRepository  email.EmailSender
	modeloRepository ModeloComunicaRepository
}

func NewEnvioService(envioRepository EnvioComunicacaoRepository, emailRepository email.EmailSender, modeloRepository ModeloComunicaRepository) *EnvioService {
	return &EnvioService{
		envioRepository:  envioRepository,
		emailRepository:  emailRepository,
		modeloRepository: modeloRepository,
	}
}

func (s *EnvioService) SendEmailComunicacao(
	ctx context.Context,
	nomeModelo string,
	idColaborador string,
	destinatario string,
	nomeColaborador string,
) error {
	modelo, err := s.modeloRepository.FindByNome(ctx, nomeModelo)

	if err != nil {
		return err
	}

	idCol, err := uuid.Parse(idColaborador)

	if err != nil {
		return err
	}

	err = s.emailRepository.SendEmail(
		destinatario,
		modelo.Assunto,
		modelo.Corpo,
	)

	emailLog := "Envio feito com sucesso!"
	statusEnvio := Enviado

	if err != nil {
		emailLog = err.Error()
		statusEnvio = Erro
	}

	newEnvio, err := NewEnvio(
		modelo.Id,
		idCol,
		Email,
		destinatario,
		modelo.Assunto,
		modelo.Corpo,
		emailLog,
		statusEnvio,
	)

	if err != nil {
		return err
	}

	return s.envioRepository.Store(ctx, newEnvio)
}
