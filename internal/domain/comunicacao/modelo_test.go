package comunicacao

import (
	"errors"
	"testing"
)

func TestNewComunicacaoInformaTagsObrigatoriasAusentes(t *testing.T) {
	_, err := NewComunicacao(
		"Plantão pago",
		"Pagamento realizado",
		"<p>Olá {{.nome}}</p>",
		StatusAtivo,
		PlantaoPago,
	)
	if !errors.Is(err, ErrorInvalidCorpo) {
		t.Fatalf("erro = %v", err)
	}

	expected := `para o tipo de comunicação “Plantão Pago”, inclua as tags obrigatórias: {{.dataInicio}}, {{.dataFim}}, {{.valorPago}}: Corpo da comunicação invalido`
	if err.Error() != expected {
		t.Fatalf("mensagem = %q, esperado %q", err.Error(), expected)
	}
}

func TestNewComunicacaoAceitaTodasAsTagsObrigatorias(t *testing.T) {
	modelo, err := NewComunicacao(
		"Plantão pago",
		"Pagamento realizado",
		"<p>{{.nome}} {{.dataInicio}} {{.dataFim}} {{.valorPago}}</p>",
		StatusAtivo,
		PlantaoPago,
	)
	if err != nil || modelo == nil {
		t.Fatalf("modelo/erro = %+v/%v", modelo, err)
	}
}
