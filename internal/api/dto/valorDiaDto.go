package dto

import (
	"bytes"
	"encoding/json"
	"time"
)

type SetValorDiaRequest struct {
	TipoDia        string  `json:"tipo_dia" binding:"required"`
	Valor          float64 `json:"valor" binding:"required"`
	VigenciaInicio string  `json:"vigencia_inicio" binding:"required"`
}

type ValorDiaResponse struct {
	Id             string     `json:"id"`
	TipoDia        string     `json:"tipo_dia"`
	Valor          float64    `json:"valor"`
	VigenciaInicio time.Time  `json:"vigencia_inicio"`
	VigenciaFim    *time.Time `json:"vigencia_fim"`
}

type OptionalString struct {
	Set   bool
	Value *string
}

func (o *OptionalString) UnmarshalJSON(data []byte) error {
	o.Set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		o.Value = nil
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.Value = &value
	return nil
}

type OptionalFloat64 struct {
	Set   bool
	Value *float64
}

func (o *OptionalFloat64) UnmarshalJSON(data []byte) error {
	o.Set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		o.Value = nil
		return nil
	}

	var value float64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.Value = &value
	return nil
}

type UpdateValorDiaRequest struct {
	Valor          OptionalFloat64 `json:"valor"`
	Descricao      OptionalString  `json:"descricao"`
	VigenciaInicio OptionalString  `json:"vigencia_inicio"`
	VigenciaFim    OptionalString  `json:"vigencia_fim"`
}

func (r UpdateValorDiaRequest) HasFields() bool {
	return r.Valor.Set ||
		r.Descricao.Set ||
		r.VigenciaInicio.Set ||
		r.VigenciaFim.Set
}

type UpdateValorDiaResponse struct {
	Id             string     `json:"id"`
	TipoDia        string     `json:"tipo_dia"`
	Valor          float64    `json:"valor"`
	Descricao      *string    `json:"descricao"`
	VigenciaInicio time.Time  `json:"vigencia_inicio"`
	VigenciaFim    *time.Time `json:"vigencia_fim"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
