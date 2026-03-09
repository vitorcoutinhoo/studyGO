package comunicacao

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type StatusEnvio string

const (
	Enviado StatusEnvio = "ENVIADO"
	Erro    StatusEnvio = "ERRO"
)

type Envio struct {
	Id              uuid.UUID
	IdModelo        uuid.UUID
	TipoComunicacao TipoComunicacao
	Destinatario    string
	Status          StatusEnvio
	DataEnvio       time.Time
	ErroLog         string
}

var (
	ErrorInvalidDestinatario = errors.New("Destinatário invalido!")
	ErrorInvalidStatusEnvio  = errors.New("Status de envio invalido!")
)

func NewEnvio(idModelo uuid.UUID, tipoComunicacao TipoComunicacao, destinatario, erroLog string, status StatusEnvio) (*Envio, error) {
	if len(destinatario) < 1 {
		return nil, ErrorInvalidDestinatario
	}

	if !isStatusEnvioValid(status) {
		return nil, ErrorInvalidStatusEnvio
	}

	return &Envio{
		IdModelo:        idModelo,
		TipoComunicacao: tipoComunicacao,
		Destinatario:    destinatario,
		Status:          status,
		ErroLog:         erroLog,
	}, nil
}

func isStatusEnvioValid(t StatusEnvio) bool {
	switch t {
	case Enviado, Erro:
		return true
	}

	return false
}
