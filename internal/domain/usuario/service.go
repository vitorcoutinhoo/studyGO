package usuario

import (
	"context"
	"fmt"
	"plantao/internal/api/dto"
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
// Recebe como parâmetros o repositório de usuário e o repositório de colaborador para realizar as operações necessárias
func NewUsuarioService(repository UsuarioRepository, colaboradorRepository colaborador.ColaboradorRepository, passwordHasher PasswordHasher) *UsuarioService {
	return &UsuarioService{
		repository:            repository,
		colaboradorRepository: colaboradorRepository,
		passwordHasher:        passwordHasher,
	}
} // Fim NewUsuarioService

// Cria um novo usuário com validações e armazenamento
func (s *UsuarioService) CreateUsuario(ctx context.Context, usuario *dto.UsuarioRequestDTO, colaboradorId string) (*dto.UsuarioResponseDTO, error) {
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

	newUsuario, err := NewUsuario(colaboradorUUID, usuario.Email, usuario.Senha, RoleColaborador, StatusAtivo)

	if err != nil {
		return nil, err
	}

	exists, err = s.repository.ExistsEmail(ctx, usuario.Email)

	if err != nil {
		return nil, fmt.Errorf("erro ao verificar existência de email: %w", err)
	}

	if exists {
		return nil, ErrorEmailAlreadyExists
	}

	hashedPassword, err := s.passwordHasher.HashPassword(usuario.Senha)

	if err != nil {
		return nil, fmt.Errorf("erro ao hashear senha: %w", err)
	}

	newUsuario.Senha = hashedPassword

	createdUsuario, err := s.repository.Store(ctx, newUsuario)

	if err != nil {
		return nil, err
	}

	usuarioResponse, err := usuarioToResponse(createdUsuario)

	if err != nil {
		return nil, err
	}

	return usuarioResponse, nil
} // Fim CreateUsuario

// Atualiza um usuário existente com novas informações
func (s *UsuarioService) UpdateUsuario(ctx context.Context, usuario *dto.UsuarioRequestDTO, usuarioId string) error {
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

	err = existingUsuario.UpdateUsuario(usuario.Email, usuario.Senha, nil)

	if err != nil {
		return err
	}

	exists, err := s.repository.ExistsEmailExcludingId(ctx, usuario.Email, existingUsuario.Id)

	if err != nil {
		return fmt.Errorf("erro ao verificar existência de email excluindo ID: %w", err)
	}

	if exists {
		return ErrorEmailAlreadyExists
	}

	if usuario.Senha != "" {
		hashedPassword, err := s.passwordHasher.HashPassword(usuario.Senha)

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

// Recupera um usuário pelo ID, retornando suas informações em um formato adequado para resposta
func (s *UsuarioService) GetUsuarioById(ctx context.Context, usuarioId string) (*dto.UsuarioResponseDTO, error) {
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

	return usuarioToResponse(existingUsuario)
} // Fim GetUsuarioById

func (s *UsuarioService) GetUsuarioByEmail(ctx context.Context, email string) (*dto.UsuarioResponseDTO, error) {
	existingUsuario, err := s.repository.FindByEmail(ctx, email)

	if err != nil {
		return nil, ErrorUserNotFound
	}

	response, err := usuarioToResponse(existingUsuario)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// Converte um objeto de domínio de usuário para um DTO de resposta, formatando os campos conforme necessário
func usuarioToResponse(usuario *Usuario) (*dto.UsuarioResponseDTO, error) {
	a, err := statusToString(usuario.Ativo)

	if err != nil {
		return nil, err
	}

	return &dto.UsuarioResponseDTO{
		Id:            usuario.Id.String(),
		IdColaborador: usuario.IdColaborador.String(),
		Email:         usuario.Email,
		Role:          string(usuario.Role),
		Ativo:         a,
	}, nil
} // Fim usuarioToResponse

// Converte o status do usuário para uma string legível, verificando se o status é válido
func statusToString(status StatusUsuario) (string, error) {
	switch status {
	case StatusAtivo:
		return "ativo", nil
	case StatusInativo:
		return "inativo", nil
	default:
		return "desconecido", ErrorInvalidStatus
	}
} // Fim statusToString
