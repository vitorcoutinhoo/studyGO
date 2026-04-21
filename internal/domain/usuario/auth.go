package usuario

import (
	"context"
	"errors"
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

func (s *AuthService) Authenticate(ctx context.Context, email, senha string) (*Usuario, *string, error) {
	user, err := s.repository.FindByEmail(ctx, email)

	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if user == nil {
		return nil, nil, ErrInvalidCredentials
	}

	if !s.password.ComparePassword(user.Senha, senha) {
		return nil, nil, ErrInvalidCredentials
	}

	token, err := s.tokenGenerator.GenerateToken(user.Id.String(), string(user.Role))
	if err != nil {
		return nil, nil, err
	}

	return &Usuario{
		Id:            user.Id,
		IdColaborador: user.IdColaborador,
		Email:         user.Email,
		Senha:         "",
		Role:          user.Role,
		Ativo:         user.Ativo,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, &token, nil
}
