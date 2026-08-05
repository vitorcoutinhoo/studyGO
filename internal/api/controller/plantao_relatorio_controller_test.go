package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"plantao/internal/api/dto"
	"plantao/internal/domain/plantao"

	"github.com/gin-gonic/gin"
)

const colaboradorRelatorioID = "10000000-0000-0000-0000-000000000010"

type relatorioRepositoryStub struct {
	items  []plantao.RelatorioItem
	filtro *plantao.RelatorioFiltro
	err    error
}

func (*relatorioRepositoryStub) Store(context.Context, *plantao.Plantao) error  { return nil }
func (*relatorioRepositoryStub) Update(context.Context, *plantao.Plantao) error { return nil }
func (*relatorioRepositoryStub) Delete(context.Context, string) error           { return nil }
func (*relatorioRepositoryStub) FindById(context.Context, string) (*plantao.Plantao, error) {
	return nil, nil
}
func (*relatorioRepositoryStub) Find(context.Context, *plantao.Filtro) ([]plantao.Plantao, error) {
	return nil, nil
}
func (r *relatorioRepositoryStub) FindRelatorio(_ context.Context, filtro *plantao.RelatorioFiltro) ([]plantao.RelatorioItem, error) {
	r.filtro = filtro
	return r.items, r.err
}
func (*relatorioRepositoryStub) WithTransaction(context.Context, func(plantao.PlantaoTransaction) error) error {
	return nil
}

func novoControllerRelatorio(repo *relatorioRepositoryStub) *PlantaoController {
	return NewPlantaoController(plantao.NewPlantaoService(repo, nil, nil))
}

func executarRelatorio(t *testing.T, controller *PlantaoController, target string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, target, nil)
	controller.GetRelatorio(ctx)
	return recorder
}

func TestGetRelatorioAceitaCombinacoesDeFiltros(t *testing.T) {
	status := plantao.StatusPlantaoConcluido
	tests := []struct {
		nome                string
		query               string
		colaboradorEsperado string
		statusEsperado      *plantao.StatusPlantao
	}{
		{"somente período", "", "", nil},
		{"com colaborador", "&colaborador_id=" + colaboradorRelatorioID, colaboradorRelatorioID, nil},
		{"com status", "&status=2", "", &status},
		{"com colaborador e status", "&colaborador_id=" + colaboradorRelatorioID + "&status=2", colaboradorRelatorioID, &status},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			repo := &relatorioRepositoryStub{items: []plantao.RelatorioItem{{
				ColaboradorID:   colaboradorRelatorioID,
				PlantaoID:       "20000000-0000-0000-0000-000000000010",
				Status:          plantao.StatusPlantaoConcluido,
				DataInicio:      time.Date(2026, 8, 1, 3, 0, 0, 0, time.UTC),
				DataFim:         time.Date(2026, 8, 3, 3, 0, 0, 0, time.UTC),
				Data:            time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC),
				NomeColaborador: "Maria Silva",
				ValorTotal:      450,
				Valor:           150,
			}}}
			recorder := executarRelatorio(t, novoControllerRelatorio(repo), "/api/v1/plantoes/relatorio?data_inicio=2026-08-01&data_fim=2026-08-31"+tt.query)

			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
			}
			if repo.filtro == nil || repo.filtro.ColaboradorID != tt.colaboradorEsperado {
				t.Fatalf("filtro de colaborador = %+v", repo.filtro)
			}
			if (repo.filtro.Status == nil) != (tt.statusEsperado == nil) ||
				(repo.filtro.Status != nil && *repo.filtro.Status != *tt.statusEsperado) {
				t.Fatalf("filtro de status = %+v", repo.filtro.Status)
			}

			var resposta []dto.RelatorioPlantaoResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &resposta); err != nil {
				t.Fatal(err)
			}
			if len(resposta) != 1 || resposta[0].Data != "2026-08-02" || resposta[0].NomeColaborador != "Maria Silva" {
				t.Fatalf("resposta inesperada: %+v", resposta)
			}
		})
	}
}

func TestGetRelatorioRetornaArrayVazio(t *testing.T) {
	recorder := executarRelatorio(
		t,
		novoControllerRelatorio(&relatorioRepositoryStub{}),
		"/api/v1/plantoes/relatorio?data_inicio=2026-08-01&data_fim=2026-08-31",
	)
	if recorder.Code != http.StatusOK || recorder.Body.String() != "[]" {
		t.Fatalf("status/body = %d/%s", recorder.Code, recorder.Body.String())
	}
}

func TestGetRelatorioValidaQueryParameters(t *testing.T) {
	tests := []struct {
		nome   string
		target string
		status int
	}{
		{"datas ausentes", "/api/v1/plantoes/relatorio", http.StatusBadRequest},
		{"início inválido", "/api/v1/plantoes/relatorio?data_inicio=01-08-2026&data_fim=2026-08-31", http.StatusBadRequest},
		{"fim inválido", "/api/v1/plantoes/relatorio?data_inicio=2026-08-01&data_fim=31-08-2026", http.StatusBadRequest},
		{"uuid inválido", "/api/v1/plantoes/relatorio?data_inicio=2026-08-01&data_fim=2026-08-31&colaborador_id=invalido", http.StatusBadRequest},
		{"status não numérico", "/api/v1/plantoes/relatorio?data_inicio=2026-08-01&data_fim=2026-08-31&status=pago", http.StatusBadRequest},
		{"fim anterior", "/api/v1/plantoes/relatorio?data_inicio=2026-08-31&data_fim=2026-08-01", http.StatusUnprocessableEntity},
		{"status fora do intervalo", "/api/v1/plantoes/relatorio?data_inicio=2026-08-01&data_fim=2026-08-31&status=5", http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			repo := &relatorioRepositoryStub{}
			recorder := executarRelatorio(t, novoControllerRelatorio(repo), tt.target)
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, esperado %d, body = %s", recorder.Code, tt.status, recorder.Body.String())
			}
			if repo.filtro != nil {
				t.Fatal("repositório foi consultado com parâmetros inválidos")
			}
		})
	}
}

func TestGetRelatorioNaoExpoeErroInterno(t *testing.T) {
	repo := &relatorioRepositoryStub{err: errors.New("detalhe interno do banco")}
	recorder := executarRelatorio(t, novoControllerRelatorio(repo), "/api/v1/plantoes/relatorio?data_inicio=2026-08-01&data_fim=2026-08-31")
	if recorder.Code != http.StatusInternalServerError || recorder.Body.String() != `{"code":"INTERNAL_ERROR","message":"erro interno do servidor"}` {
		t.Fatalf("status/body = %d/%s", recorder.Code, recorder.Body.String())
	}
}
