package models

import (
	"time"

	"github.com/google/uuid"
)

type Feriados struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Data      time.Time `gorm:"type:date;not null" json:"data"`
	Nome      string    `gorm:"type:varchar(255);not null" json:"nome"`
	Descricao string    `gorm:"type:text" json:"descricao"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
