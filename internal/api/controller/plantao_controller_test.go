package controller

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthenticatedUserIDUsaIdentidadeDoContexto(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set("userId", "usuario-autenticado")

	id, ok := authenticatedUserID(ctx)
	if !ok || id != "usuario-autenticado" {
		t.Fatalf("identidade = %q, ok = %v", id, ok)
	}
}

func TestParseDateTimeInterpretaDataCivilEmSaoPaulo(t *testing.T) {
	data, err := parseDateTime("2026-07-06")
	if err != nil {
		t.Fatal(err)
	}
	if data.Location().String() != "America/Sao_Paulo" ||
		data.Format("2006-01-02 15:04") != "2026-07-06 00:00" {
		t.Fatalf("data interpretada incorretamente: %s (%s)", data, data.Location())
	}
}
