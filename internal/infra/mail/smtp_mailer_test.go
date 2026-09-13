package mail

import (
	"context"
	"errors"
	"testing"

	configsmtp "plantao/internal/domain/smtp"
)

type smtpRepositoryFake struct {
	err      error
	chamadas int
}

func (r *smtpRepositoryFake) Get(context.Context) (*configsmtp.Configuracao, error) {
	r.chamadas++
	return nil, r.err
}
func (r *smtpRepositoryFake) Create(context.Context, *configsmtp.Configuracao) error { return nil }
func (r *smtpRepositoryFake) Update(context.Context, *configsmtp.Configuracao) error { return nil }

func TestSMTPMailerBuscaConfiguracaoNoRepositorio(t *testing.T) {
	repo := &smtpRepositoryFake{err: configsmtp.ErrorConfiguracaoSMTPNotFound}
	err := NewSMTPMailer(repo).SendEmail(context.Background(), "destinatario@example.com", "assunto", "corpo")
	if !errors.Is(err, configsmtp.ErrorConfiguracaoSMTPNotFound) || repo.chamadas != 1 {
		t.Fatalf("erro/chamadas = %v/%d", err, repo.chamadas)
	}
}
