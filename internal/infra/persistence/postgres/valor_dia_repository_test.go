package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"plantao/internal/domain/financeiro"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestTranslateValorDiaTransactionErrorMapeiaConcorrencia(t *testing.T) {
	for _, code := range []string{"55P03", "40001", "40P01"} {
		err := translateValorDiaTransactionError(&pgconn.PgError{Code: code})
		if !errors.Is(err, financeiro.ErrorConflitoValorDia) {
			t.Fatalf("SQLSTATE %s não foi mapeado: %v", code, err)
		}
	}
}

func TestTranslateValorDiaStoreErrorMapeiaDuplicidade(t *testing.T) {
	err := translateValorDiaStoreError(&pgconn.PgError{Code: "23505"})
	if !errors.Is(err, financeiro.ErrorValorDiaAlreadyExists) {
		t.Fatalf("duplicidade não foi mapeada: %v", err)
	}
}

func TestValorDiaTransactionIntegracao(t *testing.T) {
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
		tipoDia  = financeiro.TipoDia("TESTE_PATCH_CODEX")
		tipoNovo = financeiro.TipoDia("TESTE_GLOBAL_CODEX")
	)
	_, _ = pool.Exec(ctx, `DELETE FROM config_valores_dia WHERE tipo_dia IN ($1, $2)`, tipoDia, tipoNovo)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM config_valores_dia WHERE tipo_dia IN ($1, $2)`, tipoDia, tipoNovo)
	})

	repository := NewValorDiaRepository(pool)
	descricaoGlobal := "configuração criada com descrição"
	if err := repository.Store(ctx, &financeiro.ValorDia{
		Id:             uuid.New(),
		TipoDia:        tipoNovo,
		Valor:          125.50,
		Descricao:      &descricaoGlobal,
		VigenciaInicio: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatal(err)
	}
	var inicioGlobal time.Time
	var fimGlobal *time.Time
	var descricaoPersistida *string
	if err := pool.QueryRow(ctx, `SELECT descricao, vigencia_inicio, vigencia_fim FROM config_valores_dia WHERE tipo_dia = $1`, tipoNovo).Scan(&descricaoPersistida, &inicioGlobal, &fimGlobal); err != nil {
		t.Fatal(err)
	}
	if descricaoPersistida == nil || *descricaoPersistida != descricaoGlobal || inicioGlobal.Format("2006-01-02") != "1900-01-01" || fimGlobal != nil {
		t.Fatalf("inserção global = descrição %v, datas %s/%v", descricaoPersistida, inicioGlobal, fimGlobal)
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO config_valores_dia
			(tipo_dia, valor, descricao, vigencia_inicio, vigencia_fim, updated_at)
		VALUES ($1, 100.00, 'original', '2026-01-01', '2026-06-30', '2000-01-01T00:00:00Z')
	`, tipoDia); err != nil {
		t.Fatal(err)
	}

	valores, err := repository.FindAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	encontrouEncerrada := false
	for _, valor := range valores {
		if valor.TipoDia == tipoDia {
			encontrouEncerrada = valor.VigenciaFim != nil
		}
	}
	if !encontrouEncerrada {
		t.Fatal("listagem global omitiu configuração com vigencia_fim preenchida")
	}

	err = repository.Store(ctx, &financeiro.ValorDia{
		Id:             uuid.New(),
		TipoDia:        tipoDia,
		Valor:          200,
		VigenciaInicio: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, financeiro.ErrorValorDiaAlreadyExists) {
		t.Fatalf("store duplicado = %v, esperado conflito", err)
	}

	var atualizado *financeiro.ValorDia
	err = repository.WithTransaction(ctx, func(tx financeiro.ValorDiaTransaction) error {
		valor, err := tx.LockByTipoDia(ctx, tipoDia)
		if err != nil {
			return err
		}
		valor.Valor = 175.50
		atualizado, err = tx.Update(ctx, valor)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if atualizado.Valor != 175.50 || atualizado.UpdatedAt.Equal(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("registro não foi atualizado corretamente: %+v", atualizado)
	}

	sentinel := errors.New("forçar rollback")
	err = repository.WithTransaction(ctx, func(tx financeiro.ValorDiaTransaction) error {
		valor, err := tx.LockByTipoDia(ctx, tipoDia)
		if err != nil {
			return err
		}
		valor.Valor = 200
		if _, err := tx.Update(ctx, valor); err != nil {
			return err
		}
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("erro = %v, esperado sentinel", err)
	}
	var valorPersistido float64
	var inicioPersistido time.Time
	var fimPersistido *time.Time
	if err := pool.QueryRow(ctx, `SELECT valor, vigencia_inicio, vigencia_fim FROM config_valores_dia WHERE tipo_dia = $1`, tipoDia).Scan(&valorPersistido, &inicioPersistido, &fimPersistido); err != nil {
		t.Fatal(err)
	}
	if valorPersistido != 175.50 || inicioPersistido.Format("2006-01-02") != "2026-01-01" || fimPersistido == nil || fimPersistido.Format("2006-01-02") != "2026-06-30" {
		t.Fatalf("rollback/datas internas incorretos: %.2f/%s/%v", valorPersistido, inicioPersistido, fimPersistido)
	}

	lockObtido := make(chan struct{})
	liberarLock := make(chan struct{})
	primeiraFinalizada := make(chan error, 1)
	go func() {
		primeiraFinalizada <- repository.WithTransaction(ctx, func(tx financeiro.ValorDiaTransaction) error {
			if _, err := tx.LockByTipoDia(ctx, tipoDia); err != nil {
				return err
			}
			close(lockObtido)
			<-liberarLock
			return nil
		})
	}()
	<-lockObtido

	err = repository.WithTransaction(ctx, func(tx financeiro.ValorDiaTransaction) error {
		_, err := tx.LockByTipoDia(ctx, tipoDia)
		return err
	})
	if !errors.Is(err, financeiro.ErrorConflitoValorDia) {
		t.Fatalf("segundo lock = %v, esperado conflito", err)
	}
	close(liberarLock)
	if err := <-primeiraFinalizada; err != nil {
		t.Fatalf("primeira transação: %v", err)
	}
}
