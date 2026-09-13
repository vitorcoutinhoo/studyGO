package postgres

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"plantao/internal/domain/smtp"
)

func TestTranslateSMTPStoreErrorMapeiaDuplicidade(t *testing.T) {
	err := translateSMTPStoreError(&pgconn.PgError{Code: "23505"})
	if !errors.Is(err, smtp.ErrorConfiguracaoSMTPAlreadyExists) {
		t.Fatalf("erro = %v", err)
	}
}
