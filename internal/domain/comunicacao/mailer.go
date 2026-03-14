package comunicacao

type Mailer interface {
	SendEmail(to string, subject string, body string) error
}