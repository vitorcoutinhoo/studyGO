package dto

import (
	"bytes"
	"encoding/json"
	"time"
)

type CreateSMTPRequest struct {
	Host      string `json:"smtp_host" binding:"required"`
	Porta     int    `json:"smtp_port" binding:"required"`
	Usuario   string `json:"smtp_username" binding:"required"`
	Senha     string `json:"smtp_password" binding:"required"`
	Remetente string `json:"smtp_from" binding:"required"`
}

type OptionalInt struct {
	Set   bool
	Value *int
}

func (o *OptionalInt) UnmarshalJSON(data []byte) error {
	o.Set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return nil
	}
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.Value = &value
	return nil
}

type UpdateSMTPRequest struct {
	Host      OptionalString `json:"smtp_host"`
	Porta     OptionalInt    `json:"smtp_port"`
	Usuario   OptionalString `json:"smtp_username"`
	Senha     OptionalString `json:"smtp_password"`
	Remetente OptionalString `json:"smtp_from"`
}

func (r UpdateSMTPRequest) HasFields() bool {
	return r.Host.Set || r.Porta.Set || r.Usuario.Set || r.Senha.Set || r.Remetente.Set
}

func (r UpdateSMTPRequest) HasNullField() bool {
	return (r.Host.Set && r.Host.Value == nil) || (r.Porta.Set && r.Porta.Value == nil) ||
		(r.Usuario.Set && r.Usuario.Value == nil) || (r.Senha.Set && r.Senha.Value == nil) ||
		(r.Remetente.Set && r.Remetente.Value == nil)
}

type SMTPResponse struct {
	Id        string    `json:"id"`
	Host      string    `json:"smtp_host"`
	Porta     int       `json:"smtp_port"`
	Usuario   string    `json:"smtp_username"`
	Remetente string    `json:"smtp_from"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
