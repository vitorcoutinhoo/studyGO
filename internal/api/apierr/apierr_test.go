package apierr

import (
	"net/http"
	"testing"

	"plantao/internal/domain/financeiro"
	"plantao/internal/domain/plantao"
)

func TestClassifyErrosDeFechamentoEPagamento(t *testing.T) {
	tests := []struct {
		err    error
		status int
		code   string
	}{
		{plantao.ErrorUsuarioSemPermissao, http.StatusForbidden, "FORBIDDEN"},
		{plantao.ErrorPlantaoJaFechado, http.StatusConflict, "CONFLICT"},
		{plantao.ErrorConflitoConcorrencia, http.StatusConflict, "CONFLICT"},
		{plantao.ErrorPagamentoNotFound, http.StatusNotFound, "NOT_FOUND"},
		{financeiro.ErrorValorDiaAlreadyExists, http.StatusConflict, "CONFLICT"},
		{financeiro.ErrorConflitoValorDia, http.StatusConflict, "CONFLICT"},
		{financeiro.ErrorAtualizacaoVazia, http.StatusBadRequest, "BAD_REQUEST"},
	}

	for _, tt := range tests {
		status, response := classify(tt.err)
		if status != tt.status || response.Code != tt.code {
			t.Fatalf("%v: status/code = %d/%s, esperado %d/%s", tt.err, status, response.Code, tt.status, tt.code)
		}
	}
}
