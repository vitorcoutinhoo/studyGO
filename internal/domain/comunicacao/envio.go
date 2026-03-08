package comunicacao

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

/*
CREATE TABLE envios_comunicacao (
    id UUID DEFAULT gen_random_uuid() NOT NULL,
    id_modelo UUID,
    id_colaborador UUID NOT NULL,
    tipo VARCHAR(100) NOT NULL,
    destinatario VARCHAR(255) NOT NULL,
    assunto VARCHAR(255),
    corpo TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'enviado',
    data_envio TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_leitura TIMESTAMP,
    erro_log TEXT,
    CONSTRAINT envios_comunicacao_pkey PRIMARY KEY (id),
    CONSTRAINT envios_comunicacao_id_modelo_fkey FOREIGN KEY (id_modelo)
        REFERENCES modelos_comunicacao (id),
    CONSTRAINT envios_comunicacao_id_colaborador_fkey FOREIGN KEY (id_colaborador)
        REFERENCES colaboradores (id)
);
*/

type StatusEnvio string

const (
	Enviado StatusEnvio = "ENVIADO"
	Erro    StatusEnvio = "ERRO"
)

type Envio struct {
	Id              uuid.UUID
	IdModelo        uuid.UUID
	IdColaborador   uuid.UUID
	TipoComunicacao TipoComunicacao
	Destinatario    string
	Assunto         string
	Corpo           string
	Status          StatusEnvio
	DataEnvio       time.Time
	ErroLog         string
}

var (
	ErrorInvalidDestinatario = errors.New("Destinatário invalido!")
	ErrorInvalidStatusEnvio  = errors.New("Status de envio invalido!")
)

func NewEnvio(idModelo, idColaborador uuid.UUID, tipoComunicacao TipoComunicacao, destinatario, assunto, corpo, erroLog string, status StatusEnvio) (*Envio, error) {
	if len(destinatario) < 1 {
		return nil, ErrorInvalidDestinatario
	}

	if !isStatusEnvioValid(status) {
		return nil, ErrorInvalidStatusEnvio
	}

	return &Envio{
		IdModelo:        idModelo,
		IdColaborador:   idColaborador,
		TipoComunicacao: tipoComunicacao,
		Destinatario:    destinatario,
		Assunto:         assunto,
		Corpo:           corpo,
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
