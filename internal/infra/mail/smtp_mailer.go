package mail

import (
	"net/smtp"
	"plantao/internal/infra/config"
)

type SMTPMailer struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewSMTPMailer(conf *config.Config) *SMTPMailer {
	return &SMTPMailer{
		Host:     conf.SMTP.Host,
		Port:     conf.SMTP.Port,
		Username: conf.SMTP.Username,
		Password: conf.SMTP.Password,
		From:     conf.SMTP.From,
	}
}

func (m *SMTPMailer) SendEmail(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	return smtp.SendMail(
		m.Host+":"+m.Port,
		auth,
		m.From,
		[]string{to},
		msg,
	)
}
