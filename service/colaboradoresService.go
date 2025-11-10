package service

import (
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
		return dto.ColabotadoresResponseDTO{}, result.Error
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
		return nil, result.Error
	}

	return colaboradoresDto, nil
}

func GetColaboradorById(id string) (*dto.ColabotadoresResponseDTO, error) {
	idConverted, err := uuid.Parse(id)

	if err != nil {
		return nil, err
	}

	var colaboradorDto dto.ColabotadoresResponseDTO
	result := db.DB.Model(&models.Colaboradores{}).First(&colaboradorDto, idConverted)

	if result.Error != nil {
		return nil, result.Error
	}

	return &colaboradorDto, nil
}

func UpdateColaborador(id string, colaboradorDto dto.ColabotadoresRequestDTO) (*dto.ColabotadoresResponseDTO, error) {
	idConverted, err := uuid.Parse(id)

	if err != nil {
		return nil, err
	}

	var colaborador models.Colaboradores
	result := db.DB.First(&colaborador, idConverted)

	if result.Error != nil {
		return nil, result.Error
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
		return err
	}

	var colaborador models.Colaboradores
	result := db.DB.First(&colaborador, idConverted)

	if result.Error != nil {
		return result.Error
	}

	db.DB.Delete(&colaborador)

	return nil
}
