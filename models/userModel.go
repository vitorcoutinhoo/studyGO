package models

import (
	"time"

	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	SenhaHash string    `gorm:"type:text;not null" json:"senha"`
	Role      string    `gorm:"type:varchar(50);default:'colaborador'" json:"role"`
	Ativo     string    `gorm:"type:char(1);default:'Y'" json:"ativo"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
