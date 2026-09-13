package mail

import (
	"context"
	"errors"
	"fmt"
	gosmtp "net/smtp"

	configsmtp "plantao/internal/domain/smtp"
)

var ErrorEnvioEmail = errors.New("não foi possível enviar email")

type SMTPMailer struct {
	repository configsmtp.Repository
}

func NewSMTPMailer(repository configsmtp.Repository) *SMTPMailer {
	return &SMTPMailer{repository: repository}
}

func (m *SMTPMailer) SendEmail(ctx context.Context, to string, subject string, body string) error {
	configuracao, err := m.repository.Get(ctx)
	if err != nil {
		return err
	}

	auth := gosmtp.PlainAuth("", configuracao.Usuario, configuracao.Senha, configuracao.Host)

	msg := []byte("From: " + configuracao.Remetente + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	err = gosmtp.SendMail(
		configuracao.Host+":"+fmt.Sprint(configuracao.Porta),
		auth,
		configuracao.Remetente,
		[]string{to},
		msg,
	)
	if err != nil {
		return ErrorEnvioEmail
	}
	return nil
}
