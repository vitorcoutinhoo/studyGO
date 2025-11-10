package service

import (
	"fmt"

	"github.com/google/uuid"
	"main.go/db"
	"main.go/dto"
	"main.go/models"
)

func CreateUser(userDto dto.UserRequestDTO) (dto.UserResponseDTO, error) {
	user := models.User{
		ID:        uuid.New(),
		Email:     userDto.Email,
		SenhaHash: userDto.SenhaHash,
		Role:      "USER",
	}

	result := db.DB.Create(&user)

	if result.Error != nil {
		return dto.UserResponseDTO{}, fmt.Errorf("Failed to create user: %v", result.Error)
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
		return nil, fmt.Errorf("Failed to get users: %v", result.Error)
	}

	return userResponseDTO, nil
}

func GetUserById(id string) (*dto.UserResponseDTO, error) {
	idConverted, err := uuid.Parse(id)

	if err != nil {
		return nil, fmt.Errorf("Invalid ID: %v", err)
	}

	var userResponseDTO dto.UserResponseDTO
	result := db.DB.Model(&models.User{}).First(&userResponseDTO, idConverted)

	if result.Error != nil {
		return nil, fmt.Errorf("Failed to get user: %v", result.Error)
	}

	return &userResponseDTO, nil
}

func UpdateUser(id string, userDto dto.UserRequestDTO) (*dto.UserResponseDTO, error) {
	idConverted, err := uuid.Parse(id)

	if err != nil {
		return nil, fmt.Errorf("Invalid ID: %v", err)
	}

	var user models.User
	result := db.DB.First(&user, idConverted)

	if result.Error != nil {
		return nil, fmt.Errorf("User not found: %v", result.Error)
	}

	user.Email = userDto.Email
	user.SenhaHash = userDto.SenhaHash
	db.DB.Save(&user)

	return &dto.UserResponseDTO{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func DeletUserById(id string) error {
	idConverted, err := uuid.Parse(id)

	if err != nil {
		return fmt.Errorf("Invalid ID: %v", err)
	}

	var user models.User
	result := db.DB.First(&user, idConverted)

	if result.Error != nil {
		return fmt.Errorf("User not found: %v", result.Error)
	}

	db.DB.Delete(&user)

	return nil
}
