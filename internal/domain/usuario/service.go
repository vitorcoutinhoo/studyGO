package usuario

import (
	"context"
	"fmt"
	"plantao/internal/domain/colaborador"
	"plantao/internal/domain/convite"

	"github.com/google/uuid"
)

// Serviço para gerenciar usuários
type UsuarioService struct {
	repository            UsuarioRepository
	colaboradorRepository colaborador.ColaboradorRepository
	passwordHasher        PasswordHasher
	conviteRepository     convite.ConviteRepository
}

// Cria uma nova instância do serviço de usuário
func NewUsuarioService(repository UsuarioRepository, colaboradorRepository colaborador.ColaboradorRepository, passwordHasher PasswordHasher, conviteRepository convite.ConviteRepository) *UsuarioService {
	return &UsuarioService{
		repository:            repository,
		colaboradorRepository: colaboradorRepository,
		passwordHasher:        passwordHasher,
		conviteRepository:     conviteRepository,
	}
} // Fim NewUsuarioService

// Cria novo usuário utilizando token
func (s *UsuarioService) CreateUsuarioByToken(ctx context.Context, tokenStr, email, senha string) (*Usuario, error) {
	token, err := uuid.Parse(tokenStr)
	if err != nil {
		return nil, fmt.Errorf("token inválido")
	}

	convite, err := s.conviteRepository.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if err := convite.Validate(); err != nil {
		return nil, err
	}

	exists, err := s.colaboradorRepository.ExistsId(ctx, convite.IdColaborador)
	if err != nil || !exists {
		return nil, colaborador.ErrorColaboradorNotFound
	}

	exists, err = s.repository.ExistsEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrorEmailAlreadyExists
	}

	newUsuario, err := NewUsuario(convite.IdColaborador, email, senha, RoleColaborador, StatusAtivo)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := s.passwordHasher.HashPassword(senha)
	if err != nil {
		return nil, fmt.Errorf("erro ao hashear senha: %w", err)
	}
	newUsuario.Senha = hashedPassword

	result, err := s.repository.Store(ctx, newUsuario)
	if err != nil {
		return nil, err
	}

	_ = s.conviteRepository.MarkAsUsed(ctx, token)

	return result, nil
} // Fim CreateUsuarioByToken

// Atualiza um usuário existente com novas informações
func (s *UsuarioService) UpdateUsuario(ctx context.Context, email, senha, usuarioId string) error {
	usuarioUUID, err := uuid.Parse(usuarioId)

	if err != nil {
		return fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	existingUsuario, err := s.repository.FindById(ctx, usuarioUUID)

	if existingUsuario == nil {
		return ErrorUserNotFound
	}

	if err != nil {
		return err
	}

	if existingUsuario.Ativo == StatusInativo {
		return colaborador.ErrorInactiveColaborador
	}

	err = existingUsuario.UpdateUsuario(email, senha, nil)

	if err != nil {
		return err
	}

	exists, err := s.repository.ExistsEmailExcludingId(ctx, email, existingUsuario.Id)

	if err != nil {
		return fmt.Errorf("erro ao verificar existência de email excluindo ID: %w", err)
	}

	if exists {
		return ErrorEmailAlreadyExists
	}

	if senha != "" {
		hashedPassword, err := s.passwordHasher.HashPassword(senha)

		if err != nil {
			return fmt.Errorf("erro ao hashear senha: %w", err)
		}

		existingUsuario.Senha = hashedPassword
	}

	return s.repository.Update(ctx, existingUsuario)
} // Fim UpdateUsuario

// Desativa um usuário existente, marcando-o como inativo
func (s *UsuarioService) DeleteUsuario(ctx context.Context, usuarioId string) error {
	usuarioUUID, err := uuid.Parse(usuarioId)

	if err != nil {
		return fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	existingUsuario, err := s.repository.FindById(ctx, usuarioUUID)

	if existingUsuario == nil {
		return ErrorUserNotFound
	}

	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, existingUsuario.Id)
} // Fim DisableUsuario

// Recupera um usuário pelo ID
func (s *UsuarioService) GetUsuarioById(ctx context.Context, usuarioId string) (*Usuario, error) {
	usuarioUUID, err := uuid.Parse(usuarioId)

	if err != nil {
		return nil, fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	existingUsuario, err := s.repository.FindById(ctx, usuarioUUID)

	if existingUsuario == nil {
		return nil, ErrorUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return existingUsuario, nil
} // Fim GetUsuarioById

func (s *UsuarioService) GetUsuarioByEmail(ctx context.Context, email string) (*Usuario, error) {
	existingUsuario, err := s.repository.FindByEmail(ctx, email)

	if err != nil {
		return nil, ErrorUserNotFound
	}

	return existingUsuario, nil
}

func (s *UsuarioService) GetAll(ctx context.Context) (*[]Usuario, error) {
	u, err := s.repository.FindAll(ctx)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UsuarioService) ExistsUsuarioById(ctx context.Context, id string) error {
	usuarioUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	exists, err := s.repository.ExistsId(ctx, usuarioUUID)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("usuário não existe")
	}

	return nil
}

// Converte o status do usuário para string
func StatusUsuarioString(status StatusUsuario) (string, error) {
	switch status {
	case StatusAtivo:
		return "ativo", nil
	case StatusInativo:
		return "inativo", nil
	default:
		return "desconecido", ErrorInvalidStatus
	}
} // Fim StatusUsuarioString
