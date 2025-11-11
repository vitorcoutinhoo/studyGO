package service

import (
	"fmt"

	"github.com/google/uuid"
	"main.go/db"
	"main.go/dto"
	"main.go/models"
)

func CreateColaboradores(colaboradorDto dto.ColabotadoresRequestDTO) (dto.ColabotadoresResponseDTO, error) {
	colaborador := models.Colaboradores{
		ID:           uuid.New(),
		Nome:         colaboradorDto.Nome,
		Email:        colaboradorDto.Email,
		Telefone:     colaboradorDto.Telefone,
		Cargo:        colaboradorDto.Cargo,
		Departamento: colaboradorDto.Departamento,
		FotoURL:      colaboradorDto.FotoURL,
		DataAdmissao: colaboradorDto.DataAdmissao,
	}

	result := db.DB.Create(&colaborador)

	if result.Error != nil {
		return dto.ColabotadoresResponseDTO{}, fmt.Errorf("Erro ao criar colaborador: %v", result.Error)
	}

	return dto.ColabotadoresResponseDTO{
		ID:               colaborador.ID.String(),
		Nome:             colaborador.Nome,
		Email:            colaborador.Email,
		Telefone:         colaborador.Telefone,
		Cargo:            colaborador.Cargo,
		Departamento:     colaborador.Departamento,
		FotoURL:          colaborador.FotoURL,
		DataAdmissao:     colaborador.DataAdmissao,
		DataDesligamento: colaborador.DataDesligamento,
	}, nil
}

func GetAllColaboradores() ([]dto.ColabotadoresResponseDTO, error) {
	var colaboradoresDto []dto.ColabotadoresResponseDTO
	result := db.DB.Model(&models.Colaboradores{}).Find(&colaboradoresDto)

	if result.Error != nil {
		return nil, fmt.Errorf("Error ao consultar colaboradores: %v", result.Error)
	}

	return colaboradoresDto, nil
}

func GetColaboradorById(id string) (*dto.ColabotadoresResponseDTO, error) {
	idConverted, err := uuid.Parse(id)

	if err != nil {
		return nil, fmt.Errorf("Erro ao converter id do colaborador: %v", err)
	}

	var colaboradorDto dto.ColabotadoresResponseDTO
	result := db.DB.Model(&models.Colaboradores{}).First(&colaboradorDto, idConverted)

	if result.Error != nil {
		return nil, fmt.Errorf("Erro ao consultar colaborador com id[%v]: %v", idConverted, result.Error)
	}

	return &colaboradorDto, nil
}

func UpdateColaborador(id string, colaboradorDto dto.ColabotadoresRequestDTO) (*dto.ColabotadoresResponseDTO, error) {
	idConverted, err := uuid.Parse(id)

	if err != nil {
		return nil, fmt.Errorf("Erro ao converter id do colaborador: %v", err)
	}

	var colaborador models.Colaboradores
	result := db.DB.First(&colaborador, idConverted)

	if result.Error != nil {
		return nil, fmt.Errorf("Colaborador com id[%v] não encontrado: %v", idConverted, result.Error)
	}

	colaborador.Nome = colaboradorDto.Nome
	colaborador.Email = colaboradorDto.Email
	colaborador.Telefone = colaboradorDto.Telefone
	colaborador.Cargo = colaboradorDto.Cargo
	colaborador.Departamento = colaboradorDto.Departamento
	colaborador.FotoURL = colaboradorDto.FotoURL
	colaborador.DataAdmissao = colaboradorDto.DataAdmissao
	colaborador.DataDesligamento = colaboradorDto.DataDesligamento

	db.DB.Save(&colaborador)

	return &dto.ColabotadoresResponseDTO{
		ID:               colaborador.ID.String(),
		Nome:             colaborador.Nome,
		Email:            colaborador.Email,
		Telefone:         colaborador.Telefone,
		Cargo:            colaborador.Cargo,
		Departamento:     colaborador.Departamento,
		FotoURL:          colaborador.FotoURL,
		DataAdmissao:     colaborador.DataAdmissao,
		DataDesligamento: colaborador.DataDesligamento,
	}, nil
}

func DeleteColaboradorById(id string) error {
	idConverted, err := uuid.Parse(id)

	if err != nil {
		return fmt.Errorf("Erro ao converter id do colaborador: %v", err)
	}

	var colaborador models.Colaboradores
	result := db.DB.First(&colaborador, idConverted)

	if result.Error != nil {
		return fmt.Errorf("Colaborador com id[%v] não encontrado: %v", idConverted, result.Error)
	}

	db.DB.Delete(&colaborador)

	return nil
}
