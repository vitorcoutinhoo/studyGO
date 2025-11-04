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

/*
CREATE TABLE usuarios_login (
    id RAW(16) DEFAULT SYS_GUID() NOT NULL,
    id_colaborador RAW(16) NOT NULL,
    email VARCHAR2(255) NOT NULL UNIQUE,
    senha_hash CLOB NOT NULL,
    role VARCHAR2(50) DEFAULT 'colaborador',
    ativo CHAR(1) DEFAULT 'Y',
    created_at DATE DEFAULT SYSDATE,
    updated_at DATE DEFAULT SYSDATE,
    CONSTRAINT usuarios_login_pkey PRIMARY KEY (id),
    CONSTRAINT usuarios_login_id_colaborador_fkey FOREIGN KEY (id_colaborador)
        REFERENCES colaboradores (id)
);
*/
