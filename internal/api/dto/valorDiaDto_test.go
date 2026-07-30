package dto

import (
	"encoding/json"
	"testing"
)

func TestUpdateValorDiaRequestDistingueOmitidoDeNull(t *testing.T) {
	var omitido UpdateValorDiaRequest
	if err := json.Unmarshal([]byte(`{"valor": 175.50}`), &omitido); err != nil {
		t.Fatal(err)
	}
	if !omitido.Valor.Set || omitido.Valor.Value == nil ||
		omitido.Descricao.Set || omitido.VigenciaInicio.Set || omitido.VigenciaFim.Set {
		t.Fatal("campos omitidos foram marcados como presentes")
	}

	var nulos UpdateValorDiaRequest
	if err := json.Unmarshal([]byte(`{"descricao": null, "vigencia_fim": null}`), &nulos); err != nil {
		t.Fatal(err)
	}
	if !nulos.Descricao.Set || nulos.Descricao.Value != nil ||
		!nulos.VigenciaFim.Set || nulos.VigenciaFim.Value != nil {
		t.Fatalf("null não foi preservado: %+v", nulos)
	}
}

func TestUpdateValorDiaRequestPreservaNullEmCamposObrigatorios(t *testing.T) {
	var req UpdateValorDiaRequest
	if err := json.Unmarshal([]byte(`{"valor": null, "vigencia_inicio": null}`), &req); err != nil {
		t.Fatal(err)
	}
	if !req.Valor.Set || req.Valor.Value != nil ||
		!req.VigenciaInicio.Set || req.VigenciaInicio.Value != nil {
		t.Fatalf("presença de null não foi preservada: %+v", req)
	}
}

func TestUpdateValorDiaRequestRejeitaTipoInvalidoEmCampoString(t *testing.T) {
	var req UpdateValorDiaRequest
	if err := json.Unmarshal([]byte(`{"descricao": 123}`), &req); err == nil {
		t.Fatal("descrição numérica deveria ser rejeitada")
	}
}
