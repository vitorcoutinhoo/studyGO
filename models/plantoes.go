package models

import (
	"time"

	"github.com/google/uuid"
)

type Plantoes struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DataInicio  time.Time `gorm:"type:date;not null" json:"data_inicio"`
	DataFim     time.Time `gorm:"type:date;not null" json:"data_fim"`
	Status      string    `gorm:"type:varchar(255);default: 'Agendado'" json:"status"`
	ValorTotal  float64   `gorm:"type:decimal(10,2);not null" json:"valor_total"`
	Observacoes string    `gorm:"type:text" json:"observacoes"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Colaborador Colaboradores `gorm:"foreignKey:IDColaborador;references:ID"`
}
