package models

import (
	"time"

	"github.com/google/uuid"
)

type Colaboradores struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Nome             string     `gorm:"type:varchar(255);not null" json:"nome"`
	Email            string     `gorm:"type:varchar(255);unique;not null" json:"email"`
	Telefone         *string    `gorm:"type:varchar(50)" json:"telefone"`
	Cargo            *string    `gorm:"type:VARCHAR(100)" json:"cargo"`
	Departamento     *string    `gorm:"type:VARCHAR(100)" json:"departamento"`
	FotoURL          *string    `gorm:"type:text" json:"foto_url"`
	Ativo            string     `gorm:"char(1);default:'Y'" json:"ativo"`
	DataAdmissao     *time.Time `gorm:"type:date" json:"data_admissao"`
	DataDesligamento *time.Time `gorm:"type:date" json:"data_desligamento"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	Plantoes []Plantoes `gorm:"foreignKey:IDColaborador;references:ID"`
}
