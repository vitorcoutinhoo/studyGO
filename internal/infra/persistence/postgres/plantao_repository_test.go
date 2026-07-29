package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"plantao/internal/domain/plantao"

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
