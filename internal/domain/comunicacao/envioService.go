package comunicacao

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/google/uuid"

	"plantao/internal/domain/log"
	"plantao/internal/utils"
)

type EnvioService struct {
	envioRepository  EnvioComunicacaoRepository
	emailRepository  Mailer
	modeloRepository ModeloComunicaRepository
	log              log.Logger
}

func NewEnvioService(envioRepository EnvioComunicacaoRepository, emailRepository Mailer, modeloRepository ModeloComunicaRepository, log log.Logger) *EnvioService {
	return &EnvioService{
		envioRepository:  envioRepository,
		emailRepository:  emailRepository,
		modeloRepository: modeloRepository,
		log:              log,
	}
}

func (s *EnvioService) SendEmailComunicacao(
	ctx context.Context,
	tipoComunicacao TipoComunicacao,
	destinatario string,
	idColaborador uuid.UUID,
	data map[string]any,
) error {
	s.log.Info("iniciando envio de comunicação por email", "tipo_comunicacao", tipoComunicacao, "destinatario", destinatario, "id_colaborador", idColaborador)

	modelo, err := s.modeloRepository.FindByTipo(ctx, string(tipoComunicacao))

	if err != nil {
		s.log.Error("erro ao buscar modelo de comunicação para envio", "tipo_comunicacao", tipoComunicacao, "destinatario", destinatario, "error", err)
		return err
	}

	s.log.Info("renderizando template de comunicação", "id_modelo", modelo.Id, "tipo_comunicacao", tipoComunicacao)
	body, err := renderTemplate(modelo.Corpo, data)

	if err != nil {
		s.log.Error("erro ao renderizar template de comunicação", "id_modelo", modelo.Id, "tipo_comunicacao", tipoComunicacao, "error", err)
		return err
	}

	s.log.Info("enviando email de comunicação", "tipo_comunicacao", tipoComunicacao, "destinatario", destinatario)
	err = s.emailRepository.SendEmail(
		ctx,
		destinatario,
		modelo.Assunto,
		body,
	)

	emailLog := "Envio feito com sucesso!"
	statusEnvio := Enviado

	if err != nil {
		s.log.Warn("erro ao enviar email de comunicação", "tipo_comunicacao", tipoComunicacao, "destinatario", destinatario, "error", err)
		emailLog = err.Error()
		statusEnvio = Erro
	} else {
		s.log.Info("email de comunicação enviado com sucesso", "tipo_comunicacao", tipoComunicacao, "destinatario", destinatario)
	}

	newEnvio, err := NewEnvio(
		modelo.Id,
		idColaborador,
		tipoComunicacao,
		destinatario,
		emailLog,
		statusEnvio,
	)

	if err != nil {
		s.log.Error("erro ao criar registro de envio de comunicação", "tipo_comunicacao", tipoComunicacao, "destinatario", destinatario, "status_envio", statusEnvio, "error", err)
		return err
	}

	if err := s.envioRepository.Store(ctx, newEnvio); err != nil {
		s.log.Error("erro ao registrar envio de comunicação", "tipo_comunicacao", tipoComunicacao, "destinatario", destinatario, "status_envio", statusEnvio, "error", err)
		return err
	}

	s.log.Info("envio de comunicação registrado com sucesso", "tipo_comunicacao", tipoComunicacao, "destinatario", destinatario, "status_envio", statusEnvio)
	return nil
}

func renderTemplate(htmlBody string, data map[string]any) (string, error) {
	tmpl, err := template.New("email").Parse(htmlBody)

	if err != nil {
		return "", fmt.Errorf("erro ao processar template do email: %w", err)
	}

	var result bytes.Buffer

	err = tmpl.Execute(&result, formatTemplateDates(data))

	if err != nil {
		return "", fmt.Errorf("erro ao reenderizar corpo do email: %w", err)
	}

	return result.String(), nil
}

func formatTemplateDates(data map[string]any) map[string]any {
	formattedData := make(map[string]any, len(data))
	for key, value := range data {
		formattedData[key] = value
	}

	for _, tag := range []TagBody{DataInicio, DataFim, DataAtual} {
		key := string(tag)
		value, ok := formattedData[key]
		if !ok {
			continue
		}

		if formattedDate, ok := formatTemplateDate(value); ok {
			formattedData[key] = formattedDate
		}
	}

	return formattedData
}

func formatTemplateDate(value any) (string, bool) {
	var date time.Time

	switch typedValue := value.(type) {
	case time.Time:
		date = typedValue
	case *time.Time:
		if typedValue == nil {
			return "", false
		}
		date = *typedValue
	case string:
		parsedDate, ok := parseTemplateDate(typedValue)
		if !ok {
			return "", false
		}
		date = parsedDate
	default:
		return "", false
	}

	formattedDate, err := utils.ParseUsToBrDate(&date, nil)
	return formattedDate, err == nil
}

func parseTemplateDate(value string) (time.Time, bool) {
	layouts := []string{
		time.RFC3339Nano,
		"2006-01-02 15:04:05.999999999Z07:00",
		"2006-01-02 15:04:05.999999999Z0700",
		"2006-01-02 15:04:05.999999999Z07",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02",
		"02/01/2006",
	}

	for _, layout := range layouts {
		if date, err := time.Parse(layout, strings.TrimSpace(value)); err == nil {
			return date, true
		}
	}

	return time.Time{}, false
}
