package usuario

import (
	"context"
	"errors"
	"plantao/internal/api/dto"
)

var ErrInvalidCredentials = errors.New("Senha ou email inválidos")

type AuthService struct {
	repository     UsuarioRepository
	password       PasswordHasher
	tokenGenerator TokenGenerator
}

func NewAuthService(repository UsuarioRepository, password PasswordHasher, tokenGenerator TokenGenerator) *AuthService {
	return &AuthService{
		repository:     repository,
		password:       password,
		tokenGenerator: tokenGenerator,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, login *dto.LoginRequestDTO) (*string, error) {
	user, err := s.repository.FindByEmail(ctx, login.Email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !s.password.ComparePassword(user.Senha, login.Senha) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokenGenerator.GenerateToken(user.Id.String(), string(user.Role))

	if err != nil {
		return nil, err
	}

	return &token, nil
}
