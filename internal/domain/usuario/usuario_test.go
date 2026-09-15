package usuario

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestNewUsuarioAceitaRoleFinanceiro(t *testing.T) {
	u, err := NewUsuario(uuid.New(), "financeiro@example.com", "senha123", RoleFinanceiro, StatusAtivo)
	if err != nil {
		t.Fatal(err)
	}
	if u.Role != RoleFinanceiro {
		t.Fatalf("role = %q, esperado %q", u.Role, RoleFinanceiro)
	}
}

func TestNewUsuarioRejeitaRoleDesconhecida(t *testing.T) {
	_, err := NewUsuario(uuid.New(), "usuario@example.com", "senha123", Role("desconhecida"), StatusAtivo)
	if !errors.Is(err, ErrorInvalidRole) {
		t.Fatalf("erro = %v, esperado %v", err, ErrorInvalidRole)
	}
}
