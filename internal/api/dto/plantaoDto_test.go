package dto

import (
	"encoding/json"
	"testing"
)

func TestUpdatePlantaoRequestDistingueCamposOmitidosENulos(t *testing.T) {
	var parcial UpdatePlantaoRequest
	if err := json.Unmarshal([]byte(`{"data_inicio":"2026-09-01"}`), &parcial); err != nil {
		t.Fatal(err)
	}
	if !parcial.DataInicio.Set || parcial.DataInicio.Value == nil ||
		parcial.DataFim.Set || parcial.ColaboradorID.Set || !parcial.HasFields() {
		t.Fatalf("presença dos campos incorreta: %+v", parcial)
	}

	var nulo UpdatePlantaoRequest
	if err := json.Unmarshal([]byte(`{"colaborador_id":null}`), &nulo); err != nil {
		t.Fatal(err)
	}
	if !nulo.ColaboradorID.Set || nulo.ColaboradorID.Value != nil {
		t.Fatalf("null não foi preservado: %+v", nulo)
	}

	var vazio UpdatePlantaoRequest
	if err := json.Unmarshal([]byte(`{}`), &vazio); err != nil {
		t.Fatal(err)
	}
	if vazio.HasFields() {
		t.Fatal("corpo vazio foi considerado atualização")
	}
}

func TestUpdatePlantaoRequestRejeitaTiposInvalidos(t *testing.T) {
	for _, body := range []string{
		`{"colaborador_id":123}`,
		`{"data_inicio":123}`,
		`{"data_fim":false}`,
	} {
		var req UpdatePlantaoRequest
		if err := json.Unmarshal([]byte(body), &req); err == nil {
			t.Fatalf("corpo deveria ser rejeitado: %s", body)
		}
	}
}
