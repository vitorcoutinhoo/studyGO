package usuario

import (
	"context"
	"fmt"
	"plantao/internal/domain/colaborador"

	"github.com/google/uuid"
)

// Serviço para gerenciar usuários
type UsuarioService struct {
	repository            UsuarioRepository
	colaboradorRepository colaborador.ColaboradorRepository
	passwordHasher        PasswordHasher
}

// Cria uma nova instância do serviço de usuário
func NewUsuarioService(repository UsuarioRepository, colaboradorRepository colaborador.ColaboradorRepository, passwordHasher PasswordHasher) *UsuarioService {
	return &UsuarioService{
		repository:            repository,
		colaboradorRepository: colaboradorRepository,
		passwordHasher:        passwordHasher,
	}
} // Fim NewUsuarioService

// Cria um novo usuário com validações e armazenamento
func (s *UsuarioService) CreateUsuario(ctx context.Context, email, senha, colaboradorId string) (*Usuario, error) {
	colaboradorUUID, err := uuid.Parse(colaboradorId)

	if err != nil {
		return nil, fmt.Errorf("UUID do colaborador inválido: %v", err)
	}

	exists, err := s.colaboradorRepository.ExistsId(ctx, colaboradorUUID)

	if err != nil {
		return nil, fmt.Errorf("erro ao verificar existência de ID do colaborador: %w", err)
	}

	if !exists {
		return nil, colaborador.ErrorColaboradorNotFound
	}

	newUsuario, err := NewUsuario(colaboradorUUID, email, senha, RoleColaborador, StatusAtivo)

	if err != nil {
		return nil, err
	}

	exists, err = s.repository.ExistsEmail(ctx, email)

	if err != nil {
		return nil, fmt.Errorf("erro ao verificar existência de email: %w", err)
	}

	if exists {
		return nil, ErrorEmailAlreadyExists
	}

	hashedPassword, err := s.passwordHasher.HashPassword(senha)

	if err != nil {
		return nil, fmt.Errorf("erro ao hashear senha: %w", err)
	}

	newUsuario.Senha = hashedPassword

	return s.repository.Store(ctx, newUsuario)
} // Fim CreateUsuario

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
func (s *UsuarioService) DisableUsuario(ctx context.Context, usuarioId string) error {
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

	return s.repository.Disable(ctx, existingUsuario.Id)
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
