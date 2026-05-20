package usuario

import (
	"context"
	"fmt"
	"plantao/internal/domain/colaborador"
	"plantao/internal/domain/convite"
	"plantao/internal/domain/log"

	"github.com/google/uuid"
)

// Serviço para gerenciar usuários
type UsuarioService struct {
	repository            UsuarioRepository
	colaboradorRepository colaborador.ColaboradorRepository
	passwordHasher        PasswordHasher
	conviteRepository     convite.ConviteRepository
	log                   log.Logger
}

// Cria uma nova instância do serviço de usuário
func NewUsuarioService(repository UsuarioRepository, colaboradorRepository colaborador.ColaboradorRepository, passwordHasher PasswordHasher, conviteRepository convite.ConviteRepository, log log.Logger) *UsuarioService {
	return &UsuarioService{
		repository:            repository,
		colaboradorRepository: colaboradorRepository,
		passwordHasher:        passwordHasher,
		conviteRepository:     conviteRepository,
		log:                   log,
	}
} // Fim NewUsuarioService

// Cria novo usuário utilizando token
func (s *UsuarioService) CreateUsuarioByToken(ctx context.Context, tokenStr, senha string) (*Usuario, error) {
	s.log.Info("iniciando criação de usuário por token")

	token, err := uuid.Parse(tokenStr)
	if err != nil {
		s.log.Warn("token inválido para criação de usuário", "error", err)
		return nil, fmt.Errorf("token inválido")
	}

	s.log.Info("buscando convite para criação de usuário")
	conv, err := s.conviteRepository.FindByToken(ctx, token)
	if err != nil {
		s.log.Warn("convite não encontrado para criação de usuário", "error", err)
		return nil, err
	}

	if err := conv.Validate(); err != nil {
		s.log.Warn("convite inválido para criação de usuário", "id_colaborador", conv.IdColaborador, "error", err)
		return nil, err
	}

	s.log.Info("buscando colaborador vinculado ao convite", "id_colaborador", conv.IdColaborador)
	col, err := s.colaboradorRepository.FindById(ctx, conv.IdColaborador)
	if err != nil || col == nil {
		s.log.Warn("colaborador vinculado ao convite não encontrado", "id_colaborador", conv.IdColaborador, "error", err)
		return nil, colaborador.ErrorColaboradorNotFound
	}

	s.log.Info("verificando existência de usuário por email", "email", col.Email)
	exists, err := s.repository.ExistsEmail(ctx, col.Email)
	if err != nil {
		s.log.Error("erro ao verificar existência de usuário por email", "email", col.Email, "error", err)
		return nil, err
	}

	if exists {
		s.log.Warn("usuário já existe para email do colaborador", "email", col.Email, "id_colaborador", col.Id)
		return nil, ErrorEmailAlreadyExists
	}

	s.log.Info("criando objeto de usuário", "id_colaborador", conv.IdColaborador, "email", col.Email)
	newUsuario, err := NewUsuario(conv.IdColaborador, col.Email, senha, RoleColaborador, StatusAtivo)
	if err != nil {
		s.log.Warn("erro ao validar dados do usuário", "id_colaborador", conv.IdColaborador, "email", col.Email, "error", err)
		return nil, err
	}

	s.log.Info("gerando hash da senha do usuário", "id_colaborador", conv.IdColaborador)
	hashedPassword, err := s.passwordHasher.HashPassword(senha)
	if err != nil {
		s.log.Error("erro ao hashear senha do usuário", "id_colaborador", conv.IdColaborador, "error", err)
		return nil, fmt.Errorf("erro ao hashear senha: %w", err)
	}
	newUsuario.Senha = hashedPassword

	s.log.Info("salvando usuário no banco de dados", "id_colaborador", conv.IdColaborador, "email", col.Email)
	result, err := s.repository.Store(ctx, newUsuario)
	if err != nil {
		s.log.Error("erro ao salvar usuário no banco de dados", "id_colaborador", conv.IdColaborador, "email", col.Email, "error", err)
		return nil, err
	}

	if err := s.conviteRepository.MarkAsUsed(ctx, token); err != nil {
		s.log.Warn("erro ao marcar convite como utilizado", "id_colaborador", conv.IdColaborador, "error", err)
	}

	if result != nil {
		s.log.Info("usuário criado com sucesso", "id_usuario", result.Id, "id_colaborador", result.IdColaborador, "email", result.Email)
	} else {
		s.log.Info("usuário criado com sucesso")
	}
	return result, nil
} // Fim CreateUsuarioByToken

// Atualiza um usuário existente com novas informações
func (s *UsuarioService) UpdateUsuario(ctx context.Context, email, senha, usuarioId string) error {
	s.log.Info("iniciando atualização de usuário", "id_usuario", usuarioId, "email", email)

	usuarioUUID, err := uuid.Parse(usuarioId)

	if err != nil {
		s.log.Warn("UUID do usuário inválido para atualização", "id_usuario", usuarioId, "error", err)
		return fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	s.log.Info("buscando usuário para atualização", "id_usuario", usuarioUUID)
	existingUsuario, err := s.repository.FindById(ctx, usuarioUUID)

	if existingUsuario == nil {
		s.log.Warn("usuário não encontrado para atualização", "id_usuario", usuarioUUID)
		return ErrorUserNotFound
	}

	if err != nil {
		s.log.Error("erro ao buscar usuário para atualização", "id_usuario", usuarioUUID, "error", err)
		return err
	}

	if existingUsuario.Ativo == StatusInativo {
		s.log.Warn("tentativa de atualizar usuário inativo", "id_usuario", usuarioUUID)
		return colaborador.ErrorInactiveColaborador
	}

	err = existingUsuario.UpdateUsuario(email, senha, nil)

	if err != nil {
		s.log.Warn("erro ao validar atualização de usuário", "id_usuario", usuarioUUID, "error", err)
		return err
	}

	s.log.Info("verificando existência de email em outro usuário", "id_usuario", usuarioUUID, "email", email)
	exists, err := s.repository.ExistsEmailExcludingId(ctx, email, existingUsuario.Id)

	if err != nil {
		s.log.Error("erro ao verificar existência de email excluindo usuário", "id_usuario", usuarioUUID, "email", email, "error", err)
		return fmt.Errorf("erro ao verificar existência de email excluindo ID: %w", err)
	}

	if exists {
		s.log.Warn("email já utilizado por outro usuário", "id_usuario", usuarioUUID, "email", email)
		return ErrorEmailAlreadyExists
	}

	if senha != "" {
		s.log.Info("gerando hash da nova senha do usuário", "id_usuario", usuarioUUID)
		hashedPassword, err := s.passwordHasher.HashPassword(senha)

		if err != nil {
			s.log.Error("erro ao hashear nova senha do usuário", "id_usuario", usuarioUUID, "error", err)
			return fmt.Errorf("erro ao hashear senha: %w", err)
		}

		existingUsuario.Senha = hashedPassword
	}

	if err := s.repository.Update(ctx, existingUsuario); err != nil {
		s.log.Error("erro ao atualizar usuário no banco de dados", "id_usuario", usuarioUUID, "error", err)
		return err
	}

	s.log.Info("usuário atualizado com sucesso", "id_usuario", usuarioUUID)
	return nil
} // Fim UpdateUsuario

// Desativa um usuário existente, marcando-o como inativo
func (s *UsuarioService) DeleteUsuario(ctx context.Context, usuarioId string) error {
	s.log.Info("iniciando desativação de usuário", "id_usuario", usuarioId)

	usuarioUUID, err := uuid.Parse(usuarioId)

	if err != nil {
		s.log.Warn("UUID do usuário inválido para desativação", "id_usuario", usuarioId, "error", err)
		return fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	s.log.Info("buscando usuário para desativação", "id_usuario", usuarioUUID)
	existingUsuario, err := s.repository.FindById(ctx, usuarioUUID)

	if existingUsuario == nil {
		s.log.Warn("usuário não encontrado para desativação", "id_usuario", usuarioUUID)
		return ErrorUserNotFound
	}

	if err != nil {
		s.log.Error("erro ao buscar usuário para desativação", "id_usuario", usuarioUUID, "error", err)
		return err
	}

	if err := s.repository.Delete(ctx, existingUsuario.Id); err != nil {
		s.log.Error("erro ao desativar usuário", "id_usuario", existingUsuario.Id, "error", err)
		return err
	}

	s.log.Info("usuário desativado com sucesso", "id_usuario", existingUsuario.Id)
	return nil
} // Fim DisableUsuario

// Recupera um usuário pelo ID
func (s *UsuarioService) GetUsuarioById(ctx context.Context, usuarioId string) (*Usuario, error) {
	s.log.Info("buscando usuário por ID", "id_usuario", usuarioId)

	usuarioUUID, err := uuid.Parse(usuarioId)

	if err != nil {
		s.log.Warn("UUID do usuário inválido para busca", "id_usuario", usuarioId, "error", err)
		return nil, fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	existingUsuario, err := s.repository.FindById(ctx, usuarioUUID)

	if existingUsuario == nil {
		s.log.Warn("usuário não encontrado", "id_usuario", usuarioUUID)
		return nil, ErrorUserNotFound
	}

	if err != nil {
		s.log.Error("erro ao buscar usuário por ID", "id_usuario", usuarioUUID, "error", err)
		return nil, err
	}

	s.log.Info("usuário encontrado com sucesso", "id_usuario", usuarioUUID)
	return existingUsuario, nil
} // Fim GetUsuarioById

func (s *UsuarioService) GetUsuarioByEmail(ctx context.Context, email string) (*Usuario, error) {
	s.log.Info("buscando usuário por email", "email", email)

	existingUsuario, err := s.repository.FindByEmail(ctx, email)

	if err != nil {
		s.log.Warn("usuário não encontrado por email", "email", email, "error", err)
		return nil, ErrorUserNotFound
	}

	s.log.Info("usuário encontrado por email com sucesso", "email", email)
	return existingUsuario, nil
}

func (s *UsuarioService) GetAll(ctx context.Context) (*[]Usuario, error) {
	s.log.Info("buscando todos os usuários")

	u, err := s.repository.FindAll(ctx)

	if err != nil {
		s.log.Error("erro ao buscar usuários", "error", err)
		return nil, err
	}

	quantidade := 0
	if u != nil {
		quantidade = len(*u)
	}

	s.log.Info("usuários encontrados com sucesso", "quantidade", quantidade)
	return u, nil
}

func (s *UsuarioService) UpdateRole(ctx context.Context, usuarioId string, novaRole Role) error {
	s.log.Info("iniciando atualização de role do usuário", "id_usuario", usuarioId, "nova_role", novaRole)

	id, err := uuid.Parse(usuarioId)
	if err != nil {
		s.log.Warn("UUID do usuário inválido para atualização de role", "id_usuario", usuarioId, "error", err)
		return fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	u, err := s.repository.FindById(ctx, id)
	if err != nil || u == nil {
		s.log.Warn("usuário não encontrado para atualização de role", "id_usuario", id, "error", err)
		return ErrorUserNotFound
	}

	if !isRoleValid(novaRole) {
		s.log.Warn("role inválida para atualização de usuário", "id_usuario", id, "nova_role", novaRole)
		return ErrorInvalidRole
	}

	u.Role = novaRole
	if err := s.repository.Update(ctx, u); err != nil {
		s.log.Error("erro ao atualizar role do usuário", "id_usuario", id, "nova_role", novaRole, "error", err)
		return err
	}

	s.log.Info("role do usuário atualizada com sucesso", "id_usuario", id, "nova_role", novaRole)
	return nil
}

func (s *UsuarioService) ExistsUsuarioById(ctx context.Context, id string) error {
	s.log.Info("verificando existência de usuário por ID", "id_usuario", id)

	usuarioUUID, err := uuid.Parse(id)
	if err != nil {
		s.log.Warn("UUID do usuário inválido para verificação de existência", "id_usuario", id, "error", err)
		return fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	exists, err := s.repository.ExistsId(ctx, usuarioUUID)
	if err != nil {
		s.log.Error("erro ao verificar existência de usuário", "id_usuario", usuarioUUID, "error", err)
		return err
	}

	if !exists {
		s.log.Warn("usuário não existe", "id_usuario", usuarioUUID)
		return fmt.Errorf("usuário não existe")
	}

	s.log.Info("usuário existe", "id_usuario", usuarioUUID)
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
