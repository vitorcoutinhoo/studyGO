package service

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"main.go/types"
)

type ColaboradorService struct {
	db *sql.DB
}

func NewColaboradorService(db *sql.DB) *ColaboradorService {
	return &ColaboradorService{
		db: db,
	}
}

// CreateColaborador implements types.ColaboradorRepository.
func (c *ColaboradorService) CreateColaborador(colaborador types.ColaboradorRequest) (*types.ColaboradorResponse, error) {
	sqlStatement := `
        INSERT INTO colaboradores (
            nome, email, telefone, cargo, departamento, foto_url, data_admissao
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING 
            id, nome, email, telefone, cargo, departamento, foto_url, 
            ativo, data_admissao, data_desligamento, created_at, updated_at
    `

	cResponse := new(types.ColaboradorResponse)

	err := c.db.QueryRow(
		sqlStatement,
		colaborador.Nome,
		colaborador.Email,
		colaborador.Telefone,
		colaborador.Cargo,
		colaborador.Departamento,
		colaborador.FotoURL,
		colaborador.DataAdmissao,
	).Scan(
		&cResponse.ID,
		&cResponse.Nome,
		&cResponse.Email,
		&cResponse.Telefone,
		&cResponse.Cargo,
		&cResponse.Departamento,
		&cResponse.FotoURL,
		&cResponse.Ativo,
		&cResponse.DataAdmissao,
		&cResponse.DataDesligamento,
		&cResponse.CreatedAt,
		&cResponse.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return nil, fmt.Errorf("o email '%s' já está cadastrado", colaborador.Email)
			}
		}
		return nil, err
	}

	return cResponse, nil
}

// GetColaboradores implements types.ColaboradorRepository.
func (c *ColaboradorService) GetColaboradores() ([]*types.ColaboradorResponse, error) {
	sqlStatement := `
        SELECT 
            id, nome, email, telefone, cargo, departamento, foto_url, 
            ativo, data_admissao, data_desligamento, created_at, updated_at
        FROM colaboradores
        ORDER BY nome ASC
    `

	rows, err := c.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var colaboradores []*types.ColaboradorResponse

	for rows.Next() {
		colab := new(types.ColaboradorResponse)

		err := rows.Scan(
			&colab.ID,
			&colab.Nome,
			&colab.Email,
			&colab.Telefone,
			&colab.Cargo,
			&colab.Departamento,
			&colab.FotoURL,
			&colab.Ativo,
			&colab.DataAdmissao,
			&colab.DataDesligamento,
			&colab.CreatedAt,
			&colab.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		colaboradores = append(colaboradores, colab)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return colaboradores, nil
}

// GetColaboradoresById implements types.ColaboradorRepository.
func (c *ColaboradorService) GetColaboradoresById(id uuid.UUID) (*types.ColaboradorResponse, error) {
	sqlStatement := `
        SELECT 
            id, nome, email, telefone, cargo, departamento, foto_url, 
            ativo, data_admissao, data_desligamento, created_at, updated_at
        FROM colaboradores
        WHERE id = $1
    `

	colab := new(types.ColaboradorResponse)

	err := c.db.QueryRow(sqlStatement, id).Scan(
		&colab.ID,
		&colab.Nome,
		&colab.Email,
		&colab.Telefone,
		&colab.Cargo,
		&colab.Departamento,
		&colab.FotoURL,
		&colab.Ativo,
		&colab.DataAdmissao,
		&colab.DataDesligamento,
		&colab.CreatedAt,
		&colab.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("colaborador com ID '%s' não encontrado", id)
		}

		return nil, err
	}

	return colab, nil
}

// UpdateColaborador implements types.ColaboradorRepository.
func (c *ColaboradorService) UpdateColaborador(id uuid.UUID, colaborador types.ColaboradorRequest) error {
	panic("unimplemented")
}

// DeleteColaboradorById implements types.ColaboradorRepository.
func (c *ColaboradorService) DeleteColaboradorById(id uuid.UUID) error {
	panic("unimplemented")
}
