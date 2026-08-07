package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

const plantaoUpdateID = "20000000-0000-0000-0000-000000000060"

func executarUpdatePlantaoInvalido(t *testing.T, plantaoID, body string, autenticado bool) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = gin.Params{{Key: "id", Value: plantaoID}}
	if autenticado {
		ctx.Set("userId", "usuario-1")
	}
	ctx.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/plantoes/"+plantaoID, bytes.NewBufferString(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	(&PlantaoController{}).UpdatePlantao(ctx)
	return recorder
}

func TestUpdatePlantaoValidaContratoAntesDoServico(t *testing.T) {
	tests := []struct {
		nome        string
		plantaoID   string
		body        string
		autenticado bool
		status      int
	}{
		{"id do plantão inválido", "invalido", `{}`, true, http.StatusBadRequest},
		{"sem autenticação", plantaoUpdateID, `{}`, false, http.StatusUnauthorized},
		{"corpo vazio", plantaoUpdateID, `{}`, true, http.StatusBadRequest},
		{"json inválido", plantaoUpdateID, `{`, true, http.StatusBadRequest},
		{"colaborador nulo", plantaoUpdateID, `{"colaborador_id":null}`, true, http.StatusBadRequest},
		{"colaborador inválido", plantaoUpdateID, `{"colaborador_id":"invalido"}`, true, http.StatusBadRequest},
		{"início nulo", plantaoUpdateID, `{"data_inicio":null}`, true, http.StatusBadRequest},
		{"fim nulo", plantaoUpdateID, `{"data_fim":null}`, true, http.StatusBadRequest},
		{"início inválido", plantaoUpdateID, `{"data_inicio":"01/09/2026"}`, true, http.StatusBadRequest},
		{"fim inválido", plantaoUpdateID, `{"data_fim":"amanhã"}`, true, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			recorder := executarUpdatePlantaoInvalido(t, tt.plantaoID, tt.body, tt.autenticado)
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, esperado %d, body = %s", recorder.Code, tt.status, recorder.Body.String())
			}
		})
	}
}
