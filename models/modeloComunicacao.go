package models

import (
	"time"

	"github.com/google/uuid"
)

type ModeloComunicacao struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Nome      string    `gorm:"type:varchar(255);not null;unique" json:"nome"`
	Tipo      string    `gorm:"type:varchar(100);not null" json:"tipo"`
	Assunto   string    `gorm:"type:varchar(255);not null" json:"assunto"`
	Corpo     string    `gorm:"type:text;not null" json:"corpo"`
	Ativo     string    `gorm:"type:char(1);default:'Y'" json:"ativo"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
