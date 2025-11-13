package service

import (
	"database/sql"

	"github.com/google/uuid"
	"main.go/types"
)

type ColaboradorService struct {
	db *sql.DB
}

// CreateColaborador implements types.ColaboradorRepository.
func (c *ColaboradorService) CreateColaborador(colaborador types.ColaboradorRequest) (*types.ColaboradorResponse, error) {
	panic("unimplemented")
}

// DeleteColaboradorById implements types.ColaboradorRepository.
func (c *ColaboradorService) DeleteColaboradorById(id uuid.UUID) error {
	panic("unimplemented")
}

// GetColaboradores implements types.ColaboradorRepository.
func (c *ColaboradorService) GetColaboradores() ([]*types.ColaboradorResponse, error) {
	panic("unimplemented")
}

// GetColaboradoresById implements types.ColaboradorRepository.
func (c *ColaboradorService) GetColaboradoresById(id uuid.UUID) (*types.ColaboradorResponse, error) {
	panic("unimplemented")
}

// UpdateColaborador implements types.ColaboradorRepository.
func (c *ColaboradorService) UpdateColaborador(id uuid.UUID, colaborador types.ColaboradorRequest) error {
	panic("unimplemented")
}

func NewColaboradorService(db *sql.DB) *ColaboradorService {
	return &ColaboradorService{
		db: db,
	}
}
