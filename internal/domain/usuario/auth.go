package usuario

import (
	"context"
	"errors"
	"plantao/internal/domain/log"
)

var ErrInvalidCredentials = errors.New("Senha ou email inválidos")

type AuthService struct {
	repository     UsuarioRepository
	password       PasswordHasher
	tokenGenerator TokenGenerator
	log            log.Logger
}

func NewAuthService(repository UsuarioRepository, password PasswordHasher, tokenGenerator TokenGenerator, log log.Logger) *AuthService {
	return &AuthService{
		repository:     repository,
		password:       password,
		tokenGenerator: tokenGenerator,
		log:            log,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, email, senha string) (*Usuario, *string, error) {
	s.log.Info("iniciando autenticação de usuário", "email", email)

	user, err := s.repository.FindByEmail(ctx, email)

	if err != nil {
		s.log.Warn("credenciais inválidas na autenticação", "email", email, "error", err)
		return nil, nil, ErrInvalidCredentials
	}

	if user == nil {
		s.log.Warn("usuário não encontrado na autenticação", "email", email)
		return nil, nil, ErrInvalidCredentials
	}

	if !s.password.ComparePassword(user.Senha, senha) {
		s.log.Warn("senha inválida na autenticação", "id_usuario", user.Id, "email", email)
		return nil, nil, ErrInvalidCredentials
	}

	s.log.Info("gerando token de autenticação", "id_usuario", user.Id, "role", user.Role)
	token, err := s.tokenGenerator.GenerateToken(user.Id.String(), string(user.Role))
	if err != nil {
		s.log.Error("erro ao gerar token de autenticação", "id_usuario", user.Id, "error", err)
		return nil, nil, err
	}

	s.log.Info("usuário autenticado com sucesso", "id_usuario", user.Id, "email", email)
	return &Usuario{
		Id:            user.Id,
		IdColaborador: user.IdColaborador,
		Email:         user.Email,
		Senha:         "",
		Role:          user.Role,
		Ativo:         user.Ativo,
		Auditoria:     user.Auditoria,
	}, &token, nil
}
