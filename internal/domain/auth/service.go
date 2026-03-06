package auth

import (
	"context"
	"errors"
	"plantao/internal/api/dto"
	"plantao/internal/domain/usuario"
	"plantao/internal/infra/security"
)

type AuthService struct {
	repository usuario.UsuarioRepository
	password   usuario.PasswordHasher
	jwtService *security.JWTService
}

var (
	ErrInvalidCredentials = errors.New("Senha ou email inválidos")
)

func NewAuthService(repository usuario.UsuarioRepository, password usuario.PasswordHasher, jwtService *security.JWTService) *AuthService {
	return &AuthService{
		repository: repository,
		password:   password,
		jwtService: jwtService,
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

	token, err := s.jwtService.GenerateToken(user.Id.String(), string(user.Role))

	if err != nil {
		return nil, err
	}

	return &token, nil
}
