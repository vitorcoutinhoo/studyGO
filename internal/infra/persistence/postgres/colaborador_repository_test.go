package postgres

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"plantao/internal/domain/colaborador"
)

func TestTranslateColaboradorWriteErrorMapeiaCamposUnicos(t *testing.T) {
	for _, teste := range []struct {
		constraint string
		esperado   error
	}{
		{"colaboradores_telefone_key", colaborador.ErrorTelefoneAlreadyExists},
		{"colaboradores_email_key", colaborador.ErrorEmailAlreadyExists},
	} {
		err := translateColaboradorWriteError(&pgconn.PgError{Code: "23505", ConstraintName: teste.constraint})
		if !errors.Is(err, teste.esperado) {
			t.Fatalf("constraint %s: erro = %v", teste.constraint, err)
		}
	}
}
