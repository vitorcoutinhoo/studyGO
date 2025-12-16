package service

import "database/sql"

type ConfigValoresDiaService struct {
	db *sql.DB
}

func NewConfigValoresDiaService(db *sql.DB) *ConfigValoresDiaService {
	return &ConfigValoresDiaService{
		db: db,
	}
}
