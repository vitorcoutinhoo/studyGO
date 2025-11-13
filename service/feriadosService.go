package service

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"main.go/types"
)

type FeriadosService struct {
	db *sql.DB
}

func NewFeriadosService(db *sql.DB) *FeriadosService {
	return &FeriadosService{
		db: db,
	}
}

func (fs *FeriadosService) CreateFeriado(feriado types.FeriadoRequest) (*types.FeriadoResponse, error) {
	var flag bool
	checkFeriado := "SELECT true FROM feriados WHERE data = $1 LIMIT 1"
	err := fs.db.QueryRow(checkFeriado, feriado.Data).Scan(&flag)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if flag {
		return nil, fmt.Errorf("feriado na data[%v] já cadastrado", feriado.Data)
	}

	sqlStatement := `
		INSERT INTO feriados (data, descricao) 
		VALUES ($1, $2)
		RETURNING id, data, descricao, created_at, updated_at
	`
	f := new(types.FeriadoResponse)
	err = fs.db.QueryRow(sqlStatement, feriado.Data, feriado.Descricao).Scan(
		&f.ID,
		&f.Data,
		&f.Descricao,
		&f.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (fs *FeriadosService) GetFeriados() ([]types.FeriadoResponse, error) {
	sqlStatement := `
		SELECT id, data, descricao, created_at, updated_at
		FROM feriados
		ORDER BY data
	`

	rows, err := fs.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feriados []types.FeriadoResponse
	for rows.Next() {
		f := types.FeriadoResponse{}
		err := rows.Scan(
			&f.ID,
			&f.Data,
			&f.Descricao,
			&f.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		feriados = append(feriados, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return feriados, nil
}

func (fs *FeriadosService) GetFeriadoByID(id uuid.UUID) (*types.FeriadoResponse, error) {
	sqlStatement := `
		SELECT id, data, descricao, created_at, updated_at
		FROM feriados
		WHERE id = $1
	`
	f := new(types.FeriadoResponse)
	err := fs.db.QueryRow(sqlStatement, id).Scan(
		&f.ID,
		&f.Data,
		&f.Descricao,
		&f.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("feriado com id[%v] não encontrado", id)
		}

		return nil, err
	}

	return f, nil
}

func (fs *FeriadosService) DeleteFeriadobyID(id uuid.UUID) error {
	sqlStatement := `
		DELETE FROM feriados
		WHERE id = $1
	`

	var deletedID uuid.UUID

	err := fs.db.QueryRow(sqlStatement, id).Scan(&deletedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("feriado com id[%v] não encontrado", id)
		}
		return err
	}

	return nil
}

func (fs *FeriadosService) DeleteAllFeriados() error {
	sqlStatement := `
		DELETE FROM feriados
	`
	_, err := fs.db.Exec(sqlStatement)
	if err != nil {
		return err
	}

	fmt.Println("Todos os feriados foram deletados com sucesso.")
	return nil
}
