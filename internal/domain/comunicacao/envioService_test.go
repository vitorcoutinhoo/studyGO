package comunicacao

import (
	"testing"
	"time"
)

func TestRenderTemplateFormataTodasAsTagsDeData(t *testing.T) {
	dataFim := time.Date(2026, time.September, 16, 23, 59, 0, 0, time.FixedZone("UTC-3", -3*60*60))
	body, err := renderTemplate(
		`<p>{{.dataInicio}}|{{.dataFim}}|{{.dataAtual}}</p>`,
		map[string]any{
			string(DataInicio): "2026-09-14T08:30:00Z",
			string(DataFim):    &dataFim,
			string(DataAtual):  "2026-09-15 08:30:00",
		},
	)
	if err != nil {
		t.Fatalf("erro ao renderizar template: %v", err)
	}

	expected := `<p>14/09/2026|16/09/2026|15/09/2026</p>`
	if body != expected {
		t.Fatalf("corpo = %q, esperado %q", body, expected)
	}
}

func TestRenderTemplateFormataDataEmTemplateTexto(t *testing.T) {
	body, err := renderTemplate(
		`Plantão: {{.dataInicio}} até {{.dataFim}}`,
		map[string]any{
			string(DataInicio): time.Date(2026, time.September, 14, 8, 30, 0, 0, time.UTC),
			string(DataFim):    "2026-09-15",
		},
	)
	if err != nil {
		t.Fatalf("erro ao renderizar template: %v", err)
	}

	expected := `Plantão: 14/09/2026 até 15/09/2026`
	if body != expected {
		t.Fatalf("corpo = %q, esperado %q", body, expected)
	}
}

func TestRenderTemplatePreservaDadosOriginaisEValoresNaoReconhecidos(t *testing.T) {
	data := map[string]any{
		string(DataInicio): "data não reconhecida",
		string(Nome):       "Ana",
	}

	body, err := renderTemplate(`{{.nome}}: {{.dataInicio}}`, data)
	if err != nil {
		t.Fatalf("erro ao renderizar template: %v", err)
	}
	if body != "Ana: data não reconhecida" {
		t.Fatalf("corpo inesperado: %q", body)
	}
	if data[string(DataInicio)] != "data não reconhecida" {
		t.Fatalf("dados de entrada foram alterados: %#v", data)
	}
}
