package postgres

import (
	"context"
	"errors"
	"fmt"
	"plantao/internal/domain/colaborador"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ColaboradorRepository implementa a interface colaborador.ColaboradorRepository usando PostgreSQL.
type ColaboradorRepository struct {
	pool *pgxpool.Pool
}

// NewColaboradorRepository cria uma nova instância de ColaboradorRepository.
func NewColaboradorRepository(pool *pgxpool.Pool) *ColaboradorRepository {
	return &ColaboradorRepository{pool: pool}
} // Fim NewColaboradorRepository

// Salva um novo colaborador no banco de dados.
func (r *ColaboradorRepository) Store(ctx context.Context, colaborador *colaborador.Colaborador) error {
	query := `
		INSERT INTO colaboradores (
			nome,
			email,
			telefone,
			cargo,
			departamento,
			foto_url,
			ativo,
			data_admissao,
			data_desligamento
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`
	ativo := statusColaboradorToDB(colaborador.Status)

	_, err := r.pool.Exec(
		ctx,
		query,
		colaborador.Nome,
		colaborador.Email,
		colaborador.Telefone,
		colaborador.Cargo,
		colaborador.Setor,
		colaborador.Foto,
		ativo,
		colaborador.DataAdmissao,
		colaborador.DataDesligamento,
	)

	if err != nil {
		return fmt.Errorf("erro ao salvar colaborador: %w", err)
	}

	return nil
} // Fim Store

// Atualiza os dados de um colaborador existente no banco de dados.
func (r *ColaboradorRepository) Update(ctx context.Context, colaborador *colaborador.Colaborador) error {
	query := `
		UPDATE colaboradores
		SET
			nome = $1,
			email = $2,
			telefone = $3,
			cargo = $4,
			departamento = $5,
			foto_url = $6,
			ativo = $7,
			data_admissao = $8,
			data_desligamento = $9,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $10
	`
	ativo := statusColaboradorToDB(colaborador.Status)

	_, err := r.pool.Exec(
		ctx,
		query,
		colaborador.Nome,
		colaborador.Email,
		colaborador.Telefone,
		colaborador.Cargo,
		colaborador.Setor,
		colaborador.Foto,
		ativo,
		colaborador.DataAdmissao,
		colaborador.DataDesligamento,
		colaborador.Id,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar colaborador: %w", err)
	}

	return nil
} // Fim Update

// Desativa um colaborador no banco de dados, marcando-o como inativo.
func (r *ColaboradorRepository) Disable(ctx context.Context, colaboradorId string) error {
	query := `
		UPDATE colaboradores
		SET
			ativo = 'N',
			data_desligamento = CURRENT_DATE,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		  AND ativo = 'Y'
	`

	result, err := r.pool.Exec(ctx, query, colaboradorId)

	if err != nil {
		return fmt.Errorf("erro ao deletar colaborador: %w", err)
	}

	if result.RowsAffected() == 0 {
		return colaborador.ErrorColaboradorNotFound
	}

	return nil
} // Fim Disable

// Busca um colaborador no banco de dados usando seu ID.
func (r *ColaboradorRepository) FindById(ctx context.Context, colaboradorId string) (*colaborador.Colaborador, error) {
	query := `
		SELECT id, nome, email, telefone, cargo, departamento, foto_url, ativo, data_admissao, data_desligamento
		FROM colaboradores
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, colaboradorId)

	var c colaborador.Colaborador
	var ativo string

	err := row.Scan(
		&c.Id,
		&c.Nome,
		&c.Email,
		&c.Telefone,
		&c.Cargo,
		&c.Setor,
		&c.Foto,
		&ativo,
		&c.DataAdmissao,
		&c.DataDesligamento,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, colaborador.ErrorColaboradorNotFound
		}

		return nil, fmt.Errorf("erro ao buscar colaborador por ID: %w", err)
	}

	c.Status = statusColaboradorFromDB(ativo)

	return &c, nil
} // Fim FindById

// Busca um colaborador no banco de dados usando seu email.
func (r *ColaboradorRepository) FindByEmail(ctx context.Context, email string) (*colaborador.Colaborador, error) {
	query := `
		SELECT id, nome, email, telefone, cargo, departamento, foto_url, ativo, data_admissao, data_desligamento
		FROM colaboradores
		WHERE email = $1
	`

	row := r.pool.QueryRow(ctx, query, email)

	var c colaborador.Colaborador
	var ativo string

	err := row.Scan(
		&c.Id,
		&c.Nome,
		&c.Email,
		&c.Telefone,
		&c.Cargo,
		&c.Setor,
		&c.Foto,
		&ativo,
		&c.DataAdmissao,
		&c.DataDesligamento,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, colaborador.ErrorColaboradorNotFound
		}

		return nil, fmt.Errorf("erro ao buscar colaborador por email: %w", err)
	}

	c.Status = statusColaboradorFromDB(ativo)

	return &c, nil
} // Fim FindByEmail

// Busca colaboradores no banco de dados com base em filtros opcionais.
// Permite filtrar por nome, email, telefone, cargo, departamento e data de admissão.
func (r *ColaboradorRepository) FindByFilter(ctx context.Context, filter colaborador.ColaboradorFilter) ([]colaborador.Colaborador, error) {
	query := `
		SELECT id, nome, email, telefone, cargo, departamento, foto_url, ativo, data_admissao, data_desligamento
		FROM colaboradores
	`

	var conditions []string
	var args []any
	argPos := 1

	if filter.Nome != nil {
		conditions = append(conditions, fmt.Sprintf("nome ILIKE $%d", argPos))
		args = append(args, "%"+*filter.Nome+"%")
		argPos++
	}

	if filter.Email != nil {
		conditions = append(conditions, fmt.Sprintf("email ILIKE $%d", argPos))
		args = append(args, "%"+*filter.Email+"%")
		argPos++
	}

	if filter.Telefone != nil {
		conditions = append(conditions, fmt.Sprintf("telefone ILIKE $%d", argPos))
		args = append(args, "%"+*filter.Telefone+"%")
		argPos++
	}

	if filter.Cargo != nil {
		conditions = append(conditions, fmt.Sprintf("cargo ILIKE $%d", argPos))
		args = append(args, "%"+*filter.Cargo+"%")
		argPos++
	}

	if filter.Departamento != nil {
		conditions = append(conditions, fmt.Sprintf("departamento ILIKE $%d", argPos))
		args = append(args, "%"+*filter.Departamento+"%")
		argPos++
	}

	if filter.DataAdmissao != nil {
		conditions = append(conditions, fmt.Sprintf("data_admissao = $%d", argPos))
		args = append(args, *filter.DataAdmissao)
		argPos++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := r.pool.Query(ctx, query, args...)

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar colaboradores: %w", err)
	}

	defer rows.Close()

	var colaboradores []colaborador.Colaborador

	for rows.Next() {
		var c colaborador.Colaborador
		var ativo string

		err := rows.Scan(
			&c.Id,
			&c.Nome,
			&c.Email,
			&c.Telefone,
			&c.Cargo,
			&c.Setor,
			&c.Foto,
			&ativo,
			&c.DataAdmissao,
			&c.DataDesligamento,
		)

		if err != nil {
			return nil, err
		}

		c.Status = statusColaboradorFromDB(ativo)
		colaboradores = append(colaboradores, c)
	}

	return colaboradores, nil
} // Fim FindByFilter

// Verifica se um email já existe no banco de dados.
func (r *ColaboradorRepository) ExistsEmail(ctx context.Context, email string) bool {
	query := `
		SELECT COUNT(1)
		FROM colaboradores
		WHERE email = $1
	`

	var count int

	err := r.pool.QueryRow(ctx, query, email).Scan(&count)

	if err != nil {
		return false
	}

	return count > 0
} // Fim ExistsEmail

// Verifica se o ID de um colaborador existe no banco de dados.
func (r *ColaboradorRepository) ExistsId(ctx context.Context, colaboradorId string) bool {
	query := `
		SELECT COUNT(1)
		FROM colaboradores
		WHERE id = $1
	`

	var count int

	err := r.pool.QueryRow(ctx, query, colaboradorId).Scan(&count)

	if err != nil {
		return false
	}

	return count > 0
} // Fim ExistsId

// Verifica se um email já existe no banco de dados, excluindo um colaborador específico pelo ID.
func (r *ColaboradorRepository) ExistsEmailExcludingId(ctx context.Context, email, colaboradorId string) bool {
	query := `
		SELECT COUNT(1)
		FROM colaboradores
		WHERE email = $1 AND id != $2
	`

	var count int

	err := r.pool.QueryRow(ctx, query, email, colaboradorId).Scan(&count)

	if err != nil {
		return false
	}

	return count > 0
} // Fim ExistsEmailExcludingId

// Converte o status do colaborador para o formato do banco de dados.
func statusColaboradorToDB(status colaborador.StatusColaborador) string {
	switch status {
	case colaborador.StatusAtivo:
		return "Y"
	case colaborador.StatusInativo:
		return "N"
	default:
		return "N"
	}
} // Fim statusColaboradorToDB

// Converte o status do colaborador do formato do banco de dados para o formato da aplicação.
func statusColaboradorFromDB(ativo string) colaborador.StatusColaborador {
	switch ativo {
	case "Y":
		return colaborador.StatusAtivo
	case "N":
		return colaborador.StatusInativo
	default:
		return colaborador.StatusInativo
	}
} // Fim statusColaboradorFromDB
