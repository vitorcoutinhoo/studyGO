package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"plantao/internal/domain/financeiro"
	"plantao/internal/domain/plantao"
	"plantao/internal/domain/shared"
	"plantao/internal/domain/usuario"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlantaoRepository struct {
	pool *pgxpool.Pool
}

func NewPlantaoRepository(pool *pgxpool.Pool) *PlantaoRepository {
	return &PlantaoRepository{pool: pool}
}

func (r *PlantaoRepository) Store(ctx context.Context, p *plantao.Plantao) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin plantao creation transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var colaboradorID string
	err = tx.QueryRow(ctx, `
		SELECT id
		FROM colaboradores
		WHERE id = $1
		FOR UPDATE
	`, p.ColaboradorId).Scan(&colaboradorID)
	if errors.Is(err, pgx.ErrNoRows) {
		return plantao.ErrorColaboradorNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to lock plantao colaborador: %w", err)
	}

	var sobreposto bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM plantoes
			WHERE id_colaborador = $1
			  AND status <> $4
			  AND (
				(data_inicio < $3 AND data_fim > $2)
				OR (data_inicio = data_fim AND $2 = $3 AND data_inicio = $2)
			  )
		)
	`, p.ColaboradorId, p.Periodo.Inicio, p.Periodo.Fim, strconv.Itoa(int(plantao.StatusPlantaoCancelado))).Scan(&sobreposto)
	if err != nil {
		return fmt.Errorf("failed to check overlapping plantoes: %w", err)
	}
	if sobreposto {
		return plantao.ErrorExistingPlantao
	}

	query := `
		INSERT INTO plantoes (id, id_colaborador, data_inicio, data_fim, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	tag, err := tx.Exec(ctx, query,
		p.Id,
		p.ColaboradorId,
		p.Periodo.Inicio,
		p.Periodo.Fim,
		strconv.Itoa(int(p.Status)),
		p.CreatedAt,
		p.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to store plantao: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("failed to store plantao: unexpected affected rows")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit plantao creation transaction: %w", err)
	}
	return nil
}

func (r *PlantaoRepository) Update(ctx context.Context, plantao *plantao.Plantao) error {
	query := `
		UPDATE plantoes
		SET id_colaborador = $2, data_inicio = $3, data_fim = $4, status = $5, created_at = $6, updated_at = $7
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		plantao.Id,
		plantao.ColaboradorId,
		plantao.Periodo.Inicio,
		plantao.Periodo.Fim,
		strconv.Itoa(int(plantao.Status)),
		plantao.CreatedAt,
		plantao.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update plantao: %w", err)
	}

	return nil
}

func (r *PlantaoRepository) Delete(
	ctx context.Context,
	plantaoID string,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	queries := []string{
		`DELETE FROM status_plantao WHERE id_plantao = $1`,
		`DELETE FROM pagamentos WHERE id_plantao = $1`,
		`DELETE FROM plantoes_detalhes WHERE id_plantao = $1`,
	}

	for _, query := range queries {
		if _, err := tx.Exec(ctx, query, plantaoID); err != nil {
			return fmt.Errorf(
				"failed to delete plantao dependencies: %w",
				err,
			)
		}
	}

	result, err := tx.Exec(
		ctx,
		`DELETE FROM plantoes WHERE id = $1`,
		plantaoID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete plantao: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("plantao not found")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PlantaoRepository) FindById(ctx context.Context, plantaoId string) (*plantao.Plantao, error) {
	query := `
		SELECT id, id_colaborador, data_inicio, data_fim, status, valor_total, observacoes, created_at, updated_at
		FROM plantoes
		WHERE id = $1
	`
	var p plantao.Plantao
	var statusStr string
	row := r.pool.QueryRow(ctx, query, plantaoId)
	p.Periodo = &shared.Periodo{}

	err := row.Scan(
		&p.Id,
		&p.ColaboradorId,
		&p.Periodo.Inicio,
		&p.Periodo.Fim,
		&statusStr,
		&p.ValorTotal,
		&p.Observacoes,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, plantao.ErrorPlantaoNotFinded
		}
		return nil, fmt.Errorf("failed to scan plantao: %w", err)
	}

	statusInt, err := strconv.Atoi(statusStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse status: %w", err)
	}
	p.Status = plantao.StatusPlantao(statusInt)

	return &p, nil
}

func (r *PlantaoRepository) Find(
	ctx context.Context,
	filtro *plantao.Filtro,
) ([]plantao.Plantao, error) {

	query := `
		SELECT id, id_colaborador, data_inicio, data_fim, status, valor_total, observacoes, created_at, updated_at
		FROM plantoes
		WHERE 1=1
	`

	args := []any{}
	arg := 1

	// Filtro por colaborador
	if filtro != nil && filtro.ColaboradorID != "" {
		query += fmt.Sprintf(" AND id_colaborador = $%d", arg)
		args = append(args, filtro.ColaboradorID)
		arg++
	}

	// Filtro por período
	if filtro != nil && filtro.Periodo != nil {
		query += fmt.Sprintf(
			" AND data_inicio <= $%d AND data_fim >= $%d",
			arg,
			arg+1,
		)
		args = append(args,
			filtro.Periodo.Fim,
			filtro.Periodo.Inicio,
		)
		arg += 2
	}

	// Filtro por status
	if filtro != nil && filtro.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", arg)
		args = append(args, strconv.Itoa(int(*filtro.Status)))
		arg++
	}

	// Paginação
	if filtro != nil && filtro.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", arg)
		args = append(args, *filtro.Limit)
		arg++
	}

	if filtro != nil && filtro.Offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", arg)
		args = append(args, *filtro.Offset)
		arg++
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plantoes []plantao.Plantao
	for rows.Next() {
		var p plantao.Plantao
		var statusStr string
		p.Periodo = &shared.Periodo{}

		if err := rows.Scan(
			&p.Id,
			&p.ColaboradorId,
			&p.Periodo.Inicio,
			&p.Periodo.Fim,
			&statusStr,
			&p.ValorTotal,
			&p.Observacoes,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}

		statusInt, err := strconv.Atoi(statusStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse status: %w", err)
		}
		p.Status = plantao.StatusPlantao(statusInt)

		plantoes = append(plantoes, p)
	}

	return plantoes, nil
}

func (r *PlantaoRepository) FindRelatorio(ctx context.Context, filtro *plantao.RelatorioFiltro) ([]plantao.RelatorioItem, error) {
	query := `
		WITH plantoes_filtrados AS (
			SELECT c.id AS colaborador_id, p.id AS plantao_id, p.status,
			       p.data_inicio, p.data_fim, c.nome AS nome_colaborador,
			       p.valor_total AS valor_total_persistido, p.observacoes
		FROM plantoes p
		JOIN colaboradores c ON p.id_colaborador = c.id
		WHERE (p.data_inicio AT TIME ZONE 'America/Sao_Paulo')::date <= $2::date
		  AND (p.data_fim AT TIME ZONE 'America/Sao_Paulo')::date >= $1::date
	`
	args := []any{filtro.DataInicio, filtro.DataFim}
	arg := 3

	if filtro.ColaboradorID != "" {
		query += fmt.Sprintf(" AND p.id_colaborador = $%d", arg)
		args = append(args, filtro.ColaboradorID)
		arg++
	}
	if filtro.Status != nil {
		query += fmt.Sprintf(" AND p.status = $%d", arg)
		args = append(args, strconv.Itoa(int(*filtro.Status)))
	}

	query += `
		),
		dias_plantao AS (
			SELECT pf.*, dia.data::date AS data
			FROM plantoes_filtrados pf
		CROSS JOIN LATERAL generate_series(
			(pf.data_inicio AT TIME ZONE 'America/Sao_Paulo')::date::timestamp,
			(pf.data_fim AT TIME ZONE 'America/Sao_Paulo')::date::timestamp,
			INTERVAL '1 day'
		) AS dia(data)
		),
		valores_dia AS (
			SELECT d.*, pd.id AS detalhe_id, COALESCE(pd.valor, cv.valor) AS valor_dia
			FROM dias_plantao d
		LEFT JOIN plantoes_detalhes pd
			ON pd.id_plantao = d.plantao_id AND pd.data = d.data
		LEFT JOIN feriados f ON f.data = d.data
		CROSS JOIN LATERAL (
			SELECT CASE
				WHEN f.id IS NOT NULL THEN 'FERIADO'
				WHEN EXTRACT(DOW FROM d.data) = 6 THEN 'SABADO'
				WHEN EXTRACT(DOW FROM d.data) = 0 THEN 'DOMINGO'
				ELSE 'UTIL'
			END AS tipo_dia
		) tipo
		LEFT JOIN config_valores_dia cv
			ON cv.tipo_dia = tipo.tipo_dia
		),
		relatorio AS (
			SELECT v.*,
			       CASE
				   WHEN COALESCE(v.valor_total_persistido, 0) = 0 THEN
				       CASE
					   WHEN COUNT(v.valor_dia) OVER (PARTITION BY v.plantao_id) =
					        COUNT(*) OVER (PARTITION BY v.plantao_id)
					   THEN SUM(v.valor_dia) OVER (PARTITION BY v.plantao_id)
					   ELSE NULL
				       END
				   ELSE v.valor_total_persistido
			       END AS valor_total_relatorio
			FROM valores_dia v
		)
		SELECT colaborador_id, plantao_id, status, data_inicio, data_fim, data,
		       nome_colaborador, valor_total_relatorio, valor_dia, observacoes
		FROM relatorio
		WHERE data BETWEEN $1::date AND $2::date
		ORDER BY data, nome_colaborador, plantao_id, detalhe_id
	`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query plantao report: %w", err)
	}
	defer rows.Close()

	items := make([]plantao.RelatorioItem, 0)
	for rows.Next() {
		var item plantao.RelatorioItem
		var status string
		var valorTotal *float64
		var valor *float64
		if err := rows.Scan(
			&item.ColaboradorID,
			&item.PlantaoID,
			&status,
			&item.DataInicio,
			&item.DataFim,
			&item.Data,
			&item.NomeColaborador,
			&valorTotal,
			&valor,
			&item.Observacoes,
		); err != nil {
			return nil, fmt.Errorf("failed to scan plantao report: %w", err)
		}
		if valorTotal == nil || valor == nil {
			return nil, financeiro.ErrorValorDiaNotFound
		}
		item.ValorTotal = *valorTotal
		item.Valor = *valor

		statusInt, err := strconv.Atoi(status)
		if err != nil {
			return nil, fmt.Errorf("failed to parse plantao report status: %w", err)
		}
		item.Status = plantao.StatusPlantao(statusInt)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read plantao report: %w", err)
	}

	return items, nil
}

func (r *PlantaoRepository) WithTransaction(ctx context.Context, fn func(plantao.PlantaoTransaction) error) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(&plantaoTransaction{tx: tx}); err != nil {
		return translateTransactionError(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return translateTransactionError(err)
	}
	return nil
}

func translateTransactionError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "55P03", "40001", "40P01":
			return plantao.ErrorConflitoConcorrencia
		}
	}
	return err
}

type plantaoTransaction struct {
	tx pgx.Tx
}

func (t *plantaoTransaction) LockPlantao(ctx context.Context, plantaoID string) (*plantao.Plantao, error) {
	const query = `
		SELECT id, id_colaborador, data_inicio, data_fim, status, valor_total, observacoes, created_at, updated_at
		FROM plantoes
		WHERE id = $1
		FOR UPDATE NOWAIT
	`
	var p plantao.Plantao
	var status string
	p.Periodo = &shared.Periodo{}
	err := t.tx.QueryRow(ctx, query, plantaoID).Scan(
		&p.Id, &p.ColaboradorId, &p.Periodo.Inicio, &p.Periodo.Fim, &status,
		&p.ValorTotal, &p.Observacoes, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, plantao.ErrorPlantaoNotFinded
	}
	if err != nil {
		return nil, err
	}
	statusInt, err := strconv.Atoi(status)
	if err != nil {
		return nil, plantao.ErrorInvalidStatusPlantao
	}
	p.Status = plantao.StatusPlantao(statusInt)
	return &p, nil
}

func (t *plantaoTransaction) FindActor(ctx context.Context, usuarioID string) (*plantao.Actor, error) {
	const query = `
		SELECT id, id_colaborador, role
		FROM usuarios_login
		WHERE id = $1 AND ativo = 'Y'
		FOR SHARE
	`
	var actor plantao.Actor
	err := t.tx.QueryRow(ctx, query, usuarioID).Scan(&actor.UsuarioID, &actor.ColaboradorID, &actor.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, usuario.ErrorUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &actor, nil
}

func (t *plantaoTransaction) ColaboradorExists(ctx context.Context, colaboradorID string) (bool, error) {
	var exists bool
	err := t.tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM colaboradores WHERE id = $1)`, colaboradorID).Scan(&exists)
	return exists, err
}

func (t *plantaoTransaction) CountDetalhes(ctx context.Context, plantaoID string) (int, error) {
	var total int
	err := t.tx.QueryRow(ctx, `SELECT COUNT(*) FROM plantoes_detalhes WHERE id_plantao = $1`, plantaoID).Scan(&total)
	return total, err
}

func (t *plantaoTransaction) CountPagamentos(ctx context.Context, plantaoID string) (int, error) {
	var total int
	err := t.tx.QueryRow(ctx, `SELECT COUNT(*) FROM pagamentos WHERE id_plantao = $1`, plantaoID).Scan(&total)
	return total, err
}

func (t *plantaoTransaction) FindFeriados(ctx context.Context, inicio, fim time.Time) (map[string]bool, error) {
	rows, err := t.tx.Query(ctx, `SELECT data FROM feriados WHERE data BETWEEN $1 AND $2`, inicio, fim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]bool)
	for rows.Next() {
		var data time.Time
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
		result[data.Format("2006-01-02")] = true
	}
	return result, rows.Err()
}

func (t *plantaoTransaction) FindValorDiaCentavos(ctx context.Context, tipoDia financeiro.TipoDia) (int64, error) {
	const query = `
		SELECT
			CASE
				WHEN valor > 0 AND valor = ROUND(valor, 2) THEN (valor * 100)::bigint
				ELSE NULL
			END
		FROM config_valores_dia
		WHERE tipo_dia = $1
		FOR SHARE
	`
	var centavos *int64
	err := t.tx.QueryRow(ctx, query, tipoDia).Scan(&centavos)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, financeiro.ErrorValorDiaNotFound
	}
	if err != nil {
		return 0, err
	}

	if centavos == nil || *centavos <= 0 {
		return 0, financeiro.ErrorValorDiaInvalido
	}
	return *centavos, nil
}

func (t *plantaoTransaction) InsertDetalhes(ctx context.Context, plantaoID string, detalhes []plantao.PlantaoDetalhe) error {
	for _, detalhe := range detalhes {
		tag, err := t.tx.Exec(ctx, `
			INSERT INTO plantoes_detalhes (id, id_plantao, data, tipo_dia, valor)
			VALUES ($1, $2, $3, $4, $5::numeric / 100)
		`, uuid.NewString(), plantaoID, detalhe.Data, detalhe.TipoDia, detalhe.ValorCentavos)
		if err != nil {
			return err
		}
		if tag.RowsAffected() != 1 {
			return plantao.ErrorDetalhesInconsistentes
		}
	}
	return nil
}

func (t *plantaoTransaction) ClosePlantao(ctx context.Context, plantaoID string, valorTotalCentavos int64, observacoes *string) error {
	tag, err := t.tx.Exec(ctx, `
		UPDATE plantoes
		SET status = $2, valor_total = $3::numeric / 100, observacoes = $4, updated_at = NOW()
		WHERE id = $1 AND status = $5
	`, plantaoID, strconv.Itoa(int(plantao.StatusPlantaoConcluido)), valorTotalCentavos, observacoes, strconv.Itoa(int(plantao.StatusPlantaoEmAndamento)))
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return plantao.ErrorConflitoConcorrencia
	}
	return nil
}

func (t *plantaoTransaction) UpdateStatus(ctx context.Context, plantaoID string, status plantao.StatusPlantao) error {
	tag, err := t.tx.Exec(ctx, `UPDATE plantoes SET status = $2, updated_at = NOW() WHERE id = $1`, plantaoID, strconv.Itoa(int(status)))
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return plantao.ErrorPlantaoNotFinded
	}
	return nil
}

func (t *plantaoTransaction) InsertHistorico(ctx context.Context, plantaoID string, statusAntigo, statusNovo plantao.StatusPlantao, usuarioID string, observacoes *string) error {
	tag, err := t.tx.Exec(ctx, `
		INSERT INTO status_plantao (id, id_plantao, status_antigo, status_novo, id_usuario, observacoes)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, uuid.NewString(), plantaoID, strconv.Itoa(int(statusAntigo)), strconv.Itoa(int(statusNovo)), usuarioID, observacoes)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return plantao.ErrorConflitoConcorrencia
	}
	return nil
}

func (t *plantaoTransaction) InsertPagamentoPendente(ctx context.Context, plantaoID, colaboradorID string, valorTotalCentavos int64) error {
	tag, err := t.tx.Exec(ctx, `
		INSERT INTO pagamentos (id, id_plantao, id_colaborador, valor_total, status)
		VALUES ($1, $2, $3, $4::numeric / 100, 'pendente')
	`, uuid.NewString(), plantaoID, colaboradorID, valorTotalCentavos)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return plantao.ErrorPagamentoExistente
	}
	return nil
}

func (t *plantaoTransaction) FindPagamento(ctx context.Context, plantaoID string) (*plantao.Pagamento, error) {
	rows, err := t.tx.Query(ctx, `
		SELECT id, id_plantao, id_colaborador, (valor_total * 100)::bigint, status
		FROM pagamentos
		WHERE id_plantao = $1
		FOR UPDATE
	`, plantaoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pagamentos []plantao.Pagamento
	for rows.Next() {
		var p plantao.Pagamento
		if err := rows.Scan(&p.ID, &p.PlantaoID, &p.ColaboradorID, &p.ValorTotalCentavos, &p.Status); err != nil {
			return nil, err
		}
		pagamentos = append(pagamentos, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(pagamentos) == 0 {
		return nil, plantao.ErrorPagamentoNotFound
	}
	if len(pagamentos) != 1 {
		return nil, plantao.ErrorPagamentoInconsistente
	}
	return &pagamentos[0], nil
}

func (t *plantaoTransaction) SummarizeDetalhes(ctx context.Context, plantaoID string) (*plantao.DetalhesResumo, error) {
	var resumo plantao.DetalhesResumo
	err := t.tx.QueryRow(ctx, `
		SELECT COUNT(*), COUNT(DISTINCT data), COALESCE((SUM(valor) * 100)::bigint, 0)
		FROM plantoes_detalhes
		WHERE id_plantao = $1
	`, plantaoID).Scan(&resumo.Quantidade, &resumo.DatasDistintas, &resumo.ValorTotalCentavos)
	if err != nil {
		return nil, err
	}
	return &resumo, nil
}

func (t *plantaoTransaction) PayPagamento(ctx context.Context, pagamentoID string, dataPagamento time.Time, observacoes *string) error {
	tag, err := t.tx.Exec(ctx, `
		UPDATE pagamentos
		SET status = 'pago', data_pagamento = $2, observacoes = $3, updated_at = NOW()
		WHERE id = $1 AND status = 'pendente'
	`, pagamentoID, dataPagamento, observacoes)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return plantao.ErrorPlantaoJaPago
	}
	return nil
}

func (t *plantaoTransaction) LockPlantoesAgendadosAte(ctx context.Context, instante time.Time, limite int) ([]*plantao.Plantao, error) {
	rows, err := t.tx.Query(ctx, `
		SELECT id, id_colaborador, data_inicio, data_fim, status, valor_total, observacoes, created_at, updated_at
		FROM plantoes
		WHERE status = $1 AND data_inicio <= $2
		ORDER BY data_inicio, id
		LIMIT $3
		FOR UPDATE SKIP LOCKED
	`, strconv.Itoa(int(plantao.StatusPlantaoAgendado)), instante, limite)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plantoes := make([]*plantao.Plantao, 0)
	for rows.Next() {
		var p plantao.Plantao
		var status string
		p.Periodo = &shared.Periodo{}
		if err := rows.Scan(
			&p.Id,
			&p.ColaboradorId,
			&p.Periodo.Inicio,
			&p.Periodo.Fim,
			&status,
			&p.ValorTotal,
			&p.Observacoes,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		statusInt, err := strconv.Atoi(status)
		if err != nil {
			return nil, plantao.ErrorInvalidStatusPlantao
		}
		p.Status = plantao.StatusPlantao(statusInt)
		plantoes = append(plantoes, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return plantoes, nil
}

func (t *plantaoTransaction) StartPlantao(ctx context.Context, plantaoID string) error {
	tag, err := t.tx.Exec(ctx, `
		UPDATE plantoes
		SET status = $2, updated_at = NOW()
		WHERE id = $1 AND status = $3
	`,
		plantaoID,
		strconv.Itoa(int(plantao.StatusPlantaoEmAndamento)),
		strconv.Itoa(int(plantao.StatusPlantaoAgendado)),
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return plantao.ErrorConflitoConcorrencia
	}
	return nil
}

func (t *plantaoTransaction) InsertHistoricoAutomatico(ctx context.Context, plantaoID string, statusAntigo, statusNovo plantao.StatusPlantao, observacoes string) error {
	tag, err := t.tx.Exec(ctx, `
		INSERT INTO status_plantao (id, id_plantao, status_antigo, status_novo, id_usuario, observacoes)
		VALUES ($1, $2, $3, $4, NULL, $5)
	`,
		uuid.NewString(),
		plantaoID,
		strconv.Itoa(int(statusAntigo)),
		strconv.Itoa(int(statusNovo)),
		observacoes,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return plantao.ErrorConflitoConcorrencia
	}
	return nil
}
