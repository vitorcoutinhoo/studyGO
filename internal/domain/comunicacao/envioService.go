package comunicacao

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
)

type EnvioService struct {
	envioRepository  EnvioComunicacaoRepository
	emailRepository  Mailer
	modeloRepository ModeloComunicaRepository
}

func NewEnvioService(envioRepository EnvioComunicacaoRepository, emailRepository Mailer, modeloRepository ModeloComunicaRepository) *EnvioService {
	return &EnvioService{
		envioRepository:  envioRepository,
		emailRepository:  emailRepository,
		modeloRepository: modeloRepository,
	}
}

func (s *EnvioService) SendEmailComunicacao(
	ctx context.Context,
	tipoComunicacao TipoComunicacao,
	idColaborador string,
	destinatario string,
	data map[string]any,
) error {
	modelo, err := s.modeloRepository.FindByTipo(ctx, string(tipoComunicacao))

	if err != nil {
		return err
	}

	body, err := renderTemplate(modelo.Corpo, data)

	if err != nil {
		return err
	}

	err = s.emailRepository.SendEmail(
		destinatario,
		modelo.Assunto,
		body,
	)

	emailLog := "Envio feito com sucesso!"
	statusEnvio := Enviado

	if err != nil {
		emailLog = err.Error()
		statusEnvio = Erro
	}

	newEnvio, err := NewEnvio(
		modelo.Id,
		tipoComunicacao,
		destinatario,
		emailLog,
		statusEnvio,
	)

	if err != nil {
		return err
	}

	return s.envioRepository.Store(ctx, newEnvio)
}

func renderTemplate(htmlBody string, data map[string]any) (string, error) {
	tmpl, err := template.New("email").Parse(htmlBody)

	if err != nil {
		return "", fmt.Errorf("erro ao processar template do email: %w", err)
	}

	var result bytes.Buffer

	err = tmpl.Execute(&result, data)

	if err != nil {
		return "", fmt.Errorf("erro ao reenderizar corpo do email: %w", err)
	}

	return result.String(), nil
}
