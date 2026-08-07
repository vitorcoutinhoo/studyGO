package worker

import (
	"context"
	"sync"
	"time"

	"plantao/internal/domain/financeiro"
	"plantao/internal/domain/log"
	"plantao/internal/infra/feriadoapi"

	"go.uber.org/fx"
)

const feriadoSyncCheckInterval = 24 * time.Hour

type feriadoFetcher interface {
	FetchFeriados(ano int) ([]financeiro.Feriado, error)
}

type FeriadoSyncWorker struct {
	fetcher    feriadoFetcher
	repository financeiro.FeriadoRepository
	log        log.Logger
	interval   time.Duration
	now        func() time.Time

	lastSyncedYear int
	mu             sync.Mutex

	cancel context.CancelFunc
	done   chan struct{}
	once   sync.Once
}

func NewFeriadoSyncWorker(fetcher *feriadoapi.InvertextoClient, repository financeiro.FeriadoRepository, logger log.Logger) *FeriadoSyncWorker {
	return newFeriadoSyncWorker(fetcher, repository, logger, feriadoSyncCheckInterval, time.Now)
}

func newFeriadoSyncWorker(fetcher feriadoFetcher, repository financeiro.FeriadoRepository, logger log.Logger, interval time.Duration, now func() time.Time) *FeriadoSyncWorker {
	return &FeriadoSyncWorker{
		fetcher:    fetcher,
		repository: repository,
		log:        logger,
		interval:   interval,
		now:        now,
	}
}

func RegisterFeriadoSyncWorker(lifecycle fx.Lifecycle, worker *FeriadoSyncWorker) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			worker.start()
			return nil
		},
		OnStop: worker.stop,
	})
}

func (w *FeriadoSyncWorker) start() {
	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel
	w.done = make(chan struct{})

	go func() {
		defer close(w.done)
		w.execute(ctx)

		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.execute(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (w *FeriadoSyncWorker) execute(ctx context.Context) {
	ano := w.now().Year()

	w.mu.Lock()
	if w.lastSyncedYear == ano {
		w.mu.Unlock()
		return
	}
	w.mu.Unlock()

	existentes, err := w.repository.FindByAno(ctx, ano)
	if err != nil {
		w.log.Error("erro ao verificar feriados existentes", "ano", ano, "error", err)
		return
	}
	if len(existentes) > 0 {
		w.log.Debug("feriados do ano já sincronizados, ignorando chamada à API externa", "ano", ano, "quantidade", len(existentes))
		w.mu.Lock()
		w.lastSyncedYear = ano
		w.mu.Unlock()
		return
	}

	feriados, err := w.fetcher.FetchFeriados(ano)
	if err != nil {
		w.log.Error("erro ao buscar feriados na API externa", "ano", ano, "error", err)
		return
	}

	if err := w.repository.Upsert(ctx, feriados); err != nil {
		w.log.Error("erro ao salvar feriados sincronizados", "ano", ano, "error", err)
		return
	}

	w.mu.Lock()
	w.lastSyncedYear = ano
	w.mu.Unlock()

	w.log.Info("feriados sincronizados com sucesso", "ano", ano, "quantidade", len(feriados))
}

func (w *FeriadoSyncWorker) stop(ctx context.Context) error {
	w.once.Do(func() {
		if w.cancel != nil {
			w.cancel()
		}
	})
	if w.done == nil {
		return nil
	}
	select {
	case <-w.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
