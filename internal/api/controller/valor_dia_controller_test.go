package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUpdateValorRejeitaTipoAusente(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/admin/config-valores", bytes.NewBufferString(`{"valor":10}`))
	ctx.Request.Header.Set("Content-Type", "application/json")

	(&ValorDiaController{}).UpdateValor(ctx)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, esperado 400", recorder.Code)
	}
}

func TestUpdateValorRejeitaPatchInvalidoOuVigencia(t *testing.T) {
	tests := []struct {
		nome string
		body string
	}{
		{"patch vazio", `{}`},
		{"valor nulo", `{"valor":null}`},
		{"início removido", `{"vigencia_inicio":"2026-01-01"}`},
		{"fim removido", `{"vigencia_fim":null}`},
		{"valor acompanhado de vigência", `{"valor":100,"vigencia_inicio":"2026-01-01"}`},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Params = gin.Params{{Key: "tipo_dia", Value: "UTIL"}}
			ctx.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/admin/config-valores/UTIL", bytes.NewBufferString(tt.body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			(&ValorDiaController{}).UpdateValor(ctx)
			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, esperado 400", recorder.Code)
			}
		})
	}
}

func TestSetValorRejeitaCamposDeVigencia(t *testing.T) {
	for _, body := range []string{
		`{"tipo_dia":"UTIL","valor":100,"vigencia_inicio":"2026-01-01"}`,
		`{"tipo_dia":"UTIL","valor":100,"vigencia_fim":null}`,
	} {
		gin.SetMode(gin.TestMode)
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/config-valores", bytes.NewBufferString(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		(&ValorDiaController{}).SetValor(ctx)
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, esperado 400", recorder.Code)
		}
	}
}
