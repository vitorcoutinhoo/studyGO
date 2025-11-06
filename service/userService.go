package service

import (
	"fmt"

	"github.com/google/uuid"
	"main.go/db"
	"main.go/dto"
	"main.go/models"
)

func CreateUser(userDto dto.UserCreateDTO) (dto.UserResponseDTO, error) {
	user := models.User{
		ID:        uuid.New(),
		Email:     userDto.Email,
		SenhaHash: userDto.SenhaHash,
		Role:      "USER",
	}

	result := db.DB.Create(&user)

	if result.Error != nil {
		return dto.UserResponseDTO{}, fmt.Errorf("failed to create user: %v", result.Error)
	}

	return dto.UserResponseDTO{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func GetAllUsers() ([]dto.UserResponseDTO, error) {
	var userResponseDTO []dto.UserResponseDTO
	result := db.DB.Model(&models.User{}).Find(&userResponseDTO)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get users: %v", result.Error)
	}

	return userResponseDTO, nil
}
