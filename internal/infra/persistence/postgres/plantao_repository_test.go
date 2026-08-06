package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"plantao/internal/domain/plantao"
	"plantao/internal/domain/shared"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestTranslateTransactionErrorMapeiaConcorrencia(t *testing.T) {
	for _, code := range []string{"55P03", "40001", "40P01"} {
		err := translateTransactionError(&pgconn.PgError{Code: code})
		if !errors.Is(err, plantao.ErrorConflitoConcorrencia) {
			t.Fatalf("SQLSTATE %s não foi mapeado: %v", code, err)
		}
	}
}

func TestPlantaoTransactionRollbackELockConcorrente(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL não configurada")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	const (
		colaboradorID = "10000000-0000-0000-0000-000000000001"
		plantaoID     = "20000000-0000-0000-0000-000000000001"
	)
	_, _ = pool.Exec(ctx, `DELETE FROM plantoes_detalhes WHERE id_plantao = $1`, plantaoID)
	_, _ = pool.Exec(ctx, `DELETE FROM pagamentos WHERE id_plantao = $1`, plantaoID)
	_, _ = pool.Exec(ctx, `DELETE FROM status_plantao WHERE id_plantao = $1`, plantaoID)
	_, _ = pool.Exec(ctx, `DELETE FROM plantoes WHERE id = $1`, plantaoID)
	_, _ = pool.Exec(ctx, `DELETE FROM colaboradores WHERE id = $1`, colaboradorID)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM plantoes_detalhes WHERE id_plantao = $1`, plantaoID)
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM pagamentos WHERE id_plantao = $1`, plantaoID)
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM status_plantao WHERE id_plantao = $1`, plantaoID)
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM plantoes WHERE id = $1`, plantaoID)
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM colaboradores WHERE id = $1`, colaboradorID)
	})

	if _, err := pool.Exec(ctx, `
		INSERT INTO colaboradores (id, nome, email, telefone, cargo, departamento)
		VALUES ($1, 'Teste transação', 'transacao@example.test', '+550000000001', 'Teste', 'Teste')
	`, colaboradorID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO plantoes (id, id_colaborador, data_inicio, data_fim, status)
		VALUES ($1, $2, '2026-07-06T03:00:00Z', '2026-07-06T03:00:00Z', '1')
	`, plantaoID, colaboradorID); err != nil {
		t.Fatal(err)
	}

	repository := NewPlantaoRepository(pool)
	sentinel := errors.New("forçar rollback")
	err = repository.WithTransaction(ctx, func(tx plantao.PlantaoTransaction) error {
		if _, err := tx.LockPlantao(ctx, plantaoID); err != nil {
			return err
		}
		if err := tx.InsertDetalhes(ctx, plantaoID, []plantao.PlantaoDetalhe{{
			Data: time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC), TipoDia: "UTIL", ValorCentavos: 10000,
		}}); err != nil {
			return err
		}
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("erro = %v, esperado sentinel", err)
	}
	var detalhes int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM plantoes_detalhes WHERE id_plantao = $1`, plantaoID).Scan(&detalhes); err != nil {
		t.Fatal(err)
	}
	if detalhes != 0 {
		t.Fatalf("rollback deixou %d detalhes", detalhes)
	}

	lockObtido := make(chan struct{})
	liberarLock := make(chan struct{})
	primeiraFinalizada := make(chan error, 1)
	go func() {
		primeiraFinalizada <- repository.WithTransaction(ctx, func(tx plantao.PlantaoTransaction) error {
			if _, err := tx.LockPlantao(ctx, plantaoID); err != nil {
				return err
			}
			close(lockObtido)
			<-liberarLock
			return nil
		})
	}()
	<-lockObtido

	err = repository.WithTransaction(ctx, func(tx plantao.PlantaoTransaction) error {
		_, err := tx.LockPlantao(ctx, plantaoID)
		return err
	})
	if !errors.Is(err, plantao.ErrorConflitoConcorrencia) {
		t.Fatalf("segundo lock = %v, esperado conflito", err)
	}
	close(liberarLock)
	if err := <-primeiraFinalizada; err != nil {
		t.Fatalf("primeira transação: %v", err)
	}
}

func TestPlantaoTransactionInicioAutomatico(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL não configurada")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	const (
		colaboradorID = "10000000-0000-0000-0000-000000000002"
		plantaoID     = "20000000-0000-0000-0000-000000000002"
		futuroID      = "20000000-0000-0000-0000-000000000003"
	)
	cleanup := func(cleanupCtx context.Context) {
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM status_plantao WHERE id_plantao IN ($1, $2)`, plantaoID, futuroID)
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM plantoes WHERE id IN ($1, $2)`, plantaoID, futuroID)
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM colaboradores WHERE id = $1`, colaboradorID)
	}
	cleanup(ctx)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		cleanup(cleanupCtx)
	})

	if _, err := pool.Exec(ctx, `
		INSERT INTO colaboradores (id, nome, email, telefone, cargo, departamento)
		VALUES ($1, 'Teste worker', 'worker@example.test', '+550000000002', 'Teste', 'Teste')
	`, colaboradorID); err != nil {
		t.Fatal(err)
	}
	agora := time.Now()
	if _, err := pool.Exec(ctx, `
		INSERT INTO plantoes (id, id_colaborador, data_inicio, data_fim, status)
		VALUES
			($1, $3, $4, $4, '0'),
			($2, $3, $5, $5, '0')
	`, plantaoID, futuroID, colaboradorID, agora.Add(-time.Minute), agora.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	repository := NewPlantaoRepository(pool)
	lockObtido := make(chan struct{})
	liberarLock := make(chan struct{})
	primeiraFinalizada := make(chan error, 1)
	go func() {
		primeiraFinalizada <- repository.WithTransaction(ctx, func(tx plantao.PlantaoTransaction) error {
			selecionados, err := tx.LockPlantoesAgendadosAte(ctx, agora, 100)
			if err != nil {
				return err
			}
			if len(selecionados) != 1 || selecionados[0].Id != plantaoID {
				return errors.New("seleção inesperada na primeira transação")
			}
			close(lockObtido)
			<-liberarLock
			return nil
		})
	}()
	<-lockObtido

	err = repository.WithTransaction(ctx, func(tx plantao.PlantaoTransaction) error {
		selecionados, err := tx.LockPlantoesAgendadosAte(ctx, agora, 100)
		if err != nil {
			return err
		}
		if len(selecionados) != 0 {
			return errors.New("SKIP LOCKED selecionou plantão já bloqueado")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	close(liberarLock)
	if err := <-primeiraFinalizada; err != nil {
		t.Fatal(err)
	}

	sentinel := errors.New("forçar rollback automático")
	err = repository.WithTransaction(ctx, func(tx plantao.PlantaoTransaction) error {
		selecionados, err := tx.LockPlantoesAgendadosAte(ctx, agora, 100)
		if err != nil {
			return err
		}
		if len(selecionados) != 1 {
			return errors.New("plantão vencido não selecionado")
		}
		if err := tx.StartPlantao(ctx, plantaoID); err != nil {
			return err
		}
		if err := tx.InsertHistoricoAutomatico(ctx, plantaoID, plantao.StatusPlantaoAgendado, plantao.StatusPlantaoEmAndamento, "teste"); err != nil {
			return err
		}
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("erro = %v, esperado sentinel", err)
	}
	var status string
	var historicos int
	if err := pool.QueryRow(ctx, `SELECT status FROM plantoes WHERE id = $1`, plantaoID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM status_plantao WHERE id_plantao = $1`, plantaoID).Scan(&historicos); err != nil {
		t.Fatal(err)
	}
	if status != "0" || historicos != 0 {
		t.Fatalf("rollback deixou status/históricos = %s/%d", status, historicos)
	}

	err = repository.WithTransaction(ctx, func(tx plantao.PlantaoTransaction) error {
		selecionados, err := tx.LockPlantoesAgendadosAte(ctx, agora, 100)
		if err != nil {
			return err
		}
		if len(selecionados) != 1 || selecionados[0].Id != plantaoID {
			return errors.New("plantão vencido não selecionado após rollback")
		}
		if err := tx.StartPlantao(ctx, plantaoID); err != nil {
			return err
		}
		return tx.InsertHistoricoAutomatico(ctx, plantaoID, plantao.StatusPlantaoAgendado, plantao.StatusPlantaoEmAndamento, "teste")
	})
	if err != nil {
		t.Fatal(err)
	}
	var usuarioNulo bool
	if err := pool.QueryRow(ctx, `SELECT status FROM plantoes WHERE id = $1`, plantaoID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `
		SELECT COUNT(*), BOOL_AND(id_usuario IS NULL)
		FROM status_plantao
		WHERE id_plantao = $1
	`, plantaoID).Scan(&historicos, &usuarioNulo); err != nil {
		t.Fatal(err)
	}
	if status != "1" || historicos != 1 || !usuarioNulo {
		t.Fatalf("commit deixou status/históricos/usuário nulo = %s/%d/%t", status, historicos, usuarioNulo)
	}
}

func TestPlantaoRepositoryRelatorioIntegracao(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL não configurada")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	const (
		colaboradorA = "10000000-0000-0000-0000-000000000020"
		colaboradorB = "10000000-0000-0000-0000-000000000021"
		plantaoA     = "20000000-0000-0000-0000-000000000020"
		plantaoB     = "20000000-0000-0000-0000-000000000021"
		semDetalhes  = "20000000-0000-0000-0000-000000000022"
	)
	cleanup := func(cleanupCtx context.Context) {
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM plantoes_detalhes WHERE id_plantao IN ($1, $2, $3)`, plantaoA, plantaoB, semDetalhes)
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM plantoes WHERE id IN ($1, $2, $3)`, plantaoA, plantaoB, semDetalhes)
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM colaboradores WHERE id IN ($1, $2)`, colaboradorA, colaboradorB)
	}
	cleanup(ctx)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		cleanup(cleanupCtx)
	})
	var configuracoes int
	if err := pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM config_valores_dia
		WHERE tipo_dia IN ('UTIL', 'SABADO', 'DOMINGO', 'FERIADO')
	`).Scan(&configuracoes); err != nil {
		t.Fatal(err)
	}
	if configuracoes != 4 {
		t.Skip("configurações globais de todos os tipos de dia não disponíveis")
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO colaboradores (id, nome, email, telefone, cargo, departamento)
		VALUES
			($1, 'Ana Relatório', 'ana.relatorio@example.test', '+550000000020', 'Teste', 'Teste'),
			($2, 'Bruno Relatório', 'bruno.relatorio@example.test', '+550000000021', 'Teste', 'Teste')
	`, colaboradorA, colaboradorB); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO plantoes (id, id_colaborador, data_inicio, data_fim, status, valor_total, observacoes)
		VALUES
			($1, $4, '2026-08-01T03:00:00Z', '2026-08-03T03:00:00Z', '2', 0, 'plantão A'),
			($2, $5, '2026-08-15T03:00:00Z', '2026-08-15T03:00:00Z', '4', 999.00, NULL),
			($3, $4, '2026-08-10T03:00:00Z', '2026-08-10T03:00:00Z', '1', NULL, NULL)
	`, plantaoA, plantaoB, semDetalhes, colaboradorA, colaboradorB); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO plantoes_detalhes (id_plantao, data, tipo_dia, valor)
		VALUES
			($1, '2026-07-31', 'UTIL', 50.00),
			($1, '2026-08-01', 'UTIL', 100.00),
			($2, '2026-08-15', 'SABADO', 175.50),
			($1, '2026-08-03', 'UTIL', 200.00),
			($1, '2026-09-01', 'UTIL', 50.00)
	`, plantaoA, plantaoB); err != nil {
		t.Fatal(err)
	}

	repository := NewPlantaoRepository(pool)
	inicio := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	fim := time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC)
	statusPago := plantao.StatusPlantaoPago

	tests := []struct {
		nome         string
		filtro       *plantao.RelatorioFiltro
		esperados    int
		plantaoID    string
		primeiraData string
		ultimaData   string
	}{
		{
			nome:         "período inclusivo ordenado incluindo plantão sem detalhes",
			filtro:       &plantao.RelatorioFiltro{DataInicio: inicio, DataFim: fim},
			esperados:    5,
			primeiraData: "2026-08-01",
			ultimaData:   "2026-08-15",
		},
		{
			nome: "período recorta os dias do plantão",
			filtro: &plantao.RelatorioFiltro{
				DataInicio: time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC),
				DataFim:    time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC),
			},
			esperados:    1,
			plantaoID:    plantaoA,
			primeiraData: "2026-08-02",
			ultimaData:   "2026-08-02",
		},
		{
			nome:      "filtro por colaborador",
			filtro:    &plantao.RelatorioFiltro{DataInicio: inicio, DataFim: fim, ColaboradorID: colaboradorA},
			esperados: 4,
		},
		{
			nome:      "filtro por status",
			filtro:    &plantao.RelatorioFiltro{DataInicio: inicio, DataFim: fim, Status: &statusPago},
			esperados: 1,
			plantaoID: plantaoB,
		},
		{
			nome:      "filtros combinados",
			filtro:    &plantao.RelatorioFiltro{DataInicio: inicio, DataFim: fim, ColaboradorID: colaboradorB, Status: &statusPago},
			esperados: 1,
			plantaoID: plantaoB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			items, err := repository.FindRelatorio(ctx, tt.filtro)
			if err != nil {
				t.Fatal(err)
			}
			if len(items) != tt.esperados {
				t.Fatalf("itens = %d, esperado %d: %+v", len(items), tt.esperados, items)
			}
			if tt.plantaoID != "" {
				for _, item := range items {
					if item.PlantaoID != tt.plantaoID {
						t.Fatalf("plantão = %s, esperado %s", item.PlantaoID, tt.plantaoID)
					}
				}
			}
			if tt.primeiraData != "" && items[0].Data.Format("2006-01-02") != tt.primeiraData {
				t.Fatalf("primeira data = %s", items[0].Data.Format("2006-01-02"))
			}
			if tt.ultimaData != "" && items[len(items)-1].Data.Format("2006-01-02") != tt.ultimaData {
				t.Fatalf("última data = %s", items[len(items)-1].Data.Format("2006-01-02"))
			}
		})
	}

	items, err := repository.FindRelatorio(ctx, &plantao.RelatorioFiltro{DataInicio: inicio, DataFim: fim})
	if err != nil {
		t.Fatal(err)
	}
	var valorConfiguradoDomingo float64
	if err := pool.QueryRow(ctx, `
		SELECT valor
		FROM config_valores_dia
		WHERE tipo_dia = CASE
			WHEN EXISTS (SELECT 1 FROM feriados WHERE data = '2026-08-02') THEN 'FERIADO'
			ELSE 'DOMINGO'
		END
	`).Scan(&valorConfiguradoDomingo); err != nil {
		t.Fatal(err)
	}
	valores := make(map[string]float64)
	totais := make(map[string]float64)
	for _, item := range items {
		valores[item.PlantaoID+"/"+item.Data.Format("2006-01-02")] = item.Valor
		totais[item.PlantaoID] = item.ValorTotal
	}
	if valores[plantaoA+"/2026-08-01"] != 100 ||
		valores[plantaoA+"/2026-08-02"] != valorConfiguradoDomingo ||
		valores[semDetalhes+"/2026-08-10"] <= 0 {
		t.Fatalf("valores históricos/calculados inesperados: %+v", valores)
	}
	if totais[plantaoA] != 300+valorConfiguradoDomingo ||
		totais[plantaoB] != 999 ||
		totais[semDetalhes] != valores[semDetalhes+"/2026-08-10"] {
		t.Fatalf("totais persistidos/calculados inesperados: %+v", totais)
	}

	itensRecortados, err := repository.FindRelatorio(ctx, &plantao.RelatorioFiltro{
		DataInicio: time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC),
		DataFim:    time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(itensRecortados) != 1 || itensRecortados[0].PlantaoID != plantaoA ||
		itensRecortados[0].ValorTotal != 300+valorConfiguradoDomingo {
		t.Fatalf("recorte não preservou total completo: %+v", itensRecortados)
	}
	var detalhesSemPersistencia int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM plantoes_detalhes WHERE id_plantao = $1`, semDetalhes).Scan(&detalhesSemPersistencia); err != nil {
		t.Fatal(err)
	}
	if detalhesSemPersistencia != 0 {
		t.Fatalf("consulta persistiu %d detalhes calculados", detalhesSemPersistencia)
	}
	var totalA, totalB float64
	var totalSemDetalhes *float64
	if err := pool.QueryRow(ctx, `
		SELECT
			MAX(valor_total) FILTER (WHERE id = $1),
			MAX(valor_total) FILTER (WHERE id = $2),
			MAX(valor_total) FILTER (WHERE id = $3)
		FROM plantoes
		WHERE id IN ($1, $2, $3)
	`, plantaoA, plantaoB, semDetalhes).Scan(&totalA, &totalB, &totalSemDetalhes); err != nil {
		t.Fatal(err)
	}
	if totalA != 0 || totalB != 999 || totalSemDetalhes != nil {
		t.Fatalf("relatório alterou totais persistidos: %.2f/%.2f/%v", totalA, totalB, totalSemDetalhes)
	}
}

func TestPlantaoRepositoryCriacaoComSobreposicaoPorColaborador(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL não configurada")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	colaboradores := []string{
		"10000000-0000-0000-0000-000000000030",
		"10000000-0000-0000-0000-000000000031",
		"10000000-0000-0000-0000-000000000032",
		"10000000-0000-0000-0000-000000000033",
	}
	cleanup := func(cleanupCtx context.Context) {
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM plantoes WHERE id_colaborador IN ($1, $2, $3, $4)`, colaboradores[0], colaboradores[1], colaboradores[2], colaboradores[3])
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM colaboradores WHERE id IN ($1, $2, $3, $4)`, colaboradores[0], colaboradores[1], colaboradores[2], colaboradores[3])
	}
	cleanup(ctx)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		cleanup(cleanupCtx)
	})

	for i, id := range colaboradores {
		if _, err := pool.Exec(ctx, `
			INSERT INTO colaboradores (id, nome, email, telefone, cargo, departamento)
			VALUES ($1, $2, $3, $4, 'Teste', 'Teste')
		`, id, fmt.Sprintf("Colaborador Sobreposição %d", i), fmt.Sprintf("sobreposicao.%d@example.test", i), fmt.Sprintf("+55000000003%d", i)); err != nil {
			t.Fatal(err)
		}
	}

	repository := NewPlantaoRepository(pool)
	novoPlantao := func(id, colaboradorID string, inicio, fim time.Time, status plantao.StatusPlantao) *plantao.Plantao {
		return &plantao.Plantao{
			Id:            id,
			ColaboradorId: colaboradorID,
			Periodo:       &shared.Periodo{Inicio: inicio, Fim: fim},
			Status:        status,
			Auditoria: shared.Auditoria{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
	}
	data := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	hora := func(h int) time.Time { return data.Add(time.Duration(h) * time.Hour) }

	if err := repository.Store(ctx, novoPlantao("20000000-0000-0000-0000-000000000030", colaboradores[0], hora(8), hora(12), plantao.StatusPlantaoAgendado)); err != nil {
		t.Fatal(err)
	}
	if err := repository.Store(ctx, novoPlantao("20000000-0000-0000-0000-000000000031", colaboradores[1], hora(9), hora(11), plantao.StatusPlantaoAgendado)); err != nil {
		t.Fatalf("sobreposição entre colaboradores diferentes foi rejeitada: %v", err)
	}

	conflitos := []*plantao.Plantao{
		novoPlantao("20000000-0000-0000-0000-000000000032", colaboradores[0], hora(9), hora(10), plantao.StatusPlantaoAgendado),
		novoPlantao("20000000-0000-0000-0000-000000000033", colaboradores[0], hora(7), hora(13), plantao.StatusPlantaoAgendado),
		novoPlantao("20000000-0000-0000-0000-000000000034", colaboradores[0], hora(11), hora(14), plantao.StatusPlantaoAgendado),
	}
	for _, p := range conflitos {
		if err := repository.Store(ctx, p); !errors.Is(err, plantao.ErrorExistingPlantao) {
			t.Fatalf("sobreposição do mesmo colaborador = %v, esperado conflito", err)
		}
	}

	if err := repository.Store(ctx, novoPlantao("20000000-0000-0000-0000-000000000035", colaboradores[0], hora(12), hora(14), plantao.StatusPlantaoAgendado)); err != nil {
		t.Fatalf("período adjacente foi rejeitado: %v", err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO plantoes (id, id_colaborador, data_inicio, data_fim, status)
		VALUES ('20000000-0000-0000-0000-000000000036', $1, $2, $3, '3')
	`, colaboradores[0], hora(15), hora(18)); err != nil {
		t.Fatal(err)
	}
	if err := repository.Store(ctx, novoPlantao("20000000-0000-0000-0000-000000000037", colaboradores[0], hora(16), hora(17), plantao.StatusPlantaoAgendado)); err != nil {
		t.Fatalf("plantão cancelado bloqueou o período: %v", err)
	}

	ponto := hora(20)
	if err := repository.Store(ctx, novoPlantao("20000000-0000-0000-0000-000000000038", colaboradores[0], ponto, ponto, plantao.StatusPlantaoAgendado)); err != nil {
		t.Fatal(err)
	}
	if err := repository.Store(ctx, novoPlantao("20000000-0000-0000-0000-000000000039", colaboradores[0], ponto, ponto, plantao.StatusPlantaoAgendado)); !errors.Is(err, plantao.ErrorExistingPlantao) {
		t.Fatalf("período pontual idêntico = %v, esperado conflito", err)
	}

	inexistente := novoPlantao("20000000-0000-0000-0000-000000000040", "10000000-0000-0000-0000-000000000099", hora(8), hora(9), plantao.StatusPlantaoAgendado)
	if err := repository.Store(ctx, inexistente); !errors.Is(err, plantao.ErrorColaboradorNotFound) {
		t.Fatalf("colaborador ausente = %v, esperado não encontrado", err)
	}

	iniciarJuntos := make(chan struct{})
	errosDiferentes := make(chan error, 2)
	for i, colaboradorID := range colaboradores[2:] {
		go func(i int, colaboradorID string) {
			<-iniciarJuntos
			errosDiferentes <- repository.Store(ctx, novoPlantao(
				fmt.Sprintf("20000000-0000-0000-0000-00000000004%d", i+1),
				colaboradorID, hora(8), hora(10), plantao.StatusPlantaoAgendado,
			))
		}(i, colaboradorID)
	}
	close(iniciarJuntos)
	for range 2 {
		if err := <-errosDiferentes; err != nil {
			t.Fatalf("criação concorrente de colaboradores diferentes: %v", err)
		}
	}

	iniciarMesmo := make(chan struct{})
	errosMesmo := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func(i int) {
			<-iniciarMesmo
			errosMesmo <- repository.Store(ctx, novoPlantao(
				fmt.Sprintf("20000000-0000-0000-0000-00000000005%d", i),
				colaboradores[2], hora(12), hora(14), plantao.StatusPlantaoAgendado,
			))
		}(i)
	}
	close(iniciarMesmo)
	sucessos, conflitosConcorrentes := 0, 0
	for range 2 {
		err := <-errosMesmo
		switch {
		case err == nil:
			sucessos++
		case errors.Is(err, plantao.ErrorExistingPlantao):
			conflitosConcorrentes++
		default:
			t.Fatalf("erro concorrente inesperado: %v", err)
		}
	}
	if sucessos != 1 || conflitosConcorrentes != 1 {
		t.Fatalf("concorrência do mesmo colaborador = sucessos %d/conflitos %d", sucessos, conflitosConcorrentes)
	}
	var inseridosConcorrentes int
	if err := pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM plantoes
		WHERE id_colaborador = $1 AND data_inicio = $2 AND data_fim = $3
	`, colaboradores[2], hora(12), hora(14)).Scan(&inseridosConcorrentes); err != nil {
		t.Fatal(err)
	}
	if inseridosConcorrentes != 1 {
		t.Fatalf("criação concorrente inseriu %d plantões, esperado 1", inseridosConcorrentes)
	}
}
