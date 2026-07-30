package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUpdateValorVigenteRejeitaTipoAusente(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/admin/config-valores", bytes.NewBufferString(`{"valor": 10}`))
	ctx.Request.Header.Set("Content-Type", "application/json")

	(&ValorDiaController{}).UpdateValorVigente(ctx)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, esperado 400", recorder.Code)
	}
}

func TestUpdateValorVigenteRejeitaPatchVazioEDataInvalida(t *testing.T) {
	tests := []struct {
		nome string
		body string
	}{
		{"patch vazio", `{}`},
		{"data inválida", `{"vigencia_inicio":"30/07/2026"}`},
		{"valor nulo", `{"valor":null}`},
		{"início nulo", `{"vigencia_inicio":null}`},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Params = gin.Params{{Key: "tipo_dia", Value: "UTIL"}}
			ctx.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/admin/config-valores/UTIL", bytes.NewBufferString(tt.body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			(&ValorDiaController{}).UpdateValorVigente(ctx)
			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, esperado 400", recorder.Code)
			}
		})
	}
}

func TestParseValorDiaDate(t *testing.T) {
	data, err := parseValorDiaDate("2026-07-30")
	if err != nil || data.Format("2006-01-02") != "2026-07-30" {
		t.Fatalf("data/erro = %v/%v", data, err)
	}
}
