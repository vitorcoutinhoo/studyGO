package worker

import (
	"context"
	"sync"
	"time"

	"plantao/internal/domain/log"
	"plantao/internal/domain/plantao"

	"go.uber.org/fx"
)

const plantaoStatusInterval = 5 * time.Minute

type plantaoStarter interface {
	IniciarPlantoesAgendados(ctx context.Context) (int, error)
}

type PlantaoStatusWorker struct {
	service  plantaoStarter
	log      log.Logger
	interval time.Duration

	cancel context.CancelFunc
	done   chan struct{}
	once   sync.Once
}

func NewPlantaoStatusWorker(service *plantao.PlantaoService, logger log.Logger) *PlantaoStatusWorker {
	return newPlantaoStatusWorker(service, logger, plantaoStatusInterval)
}

func newPlantaoStatusWorker(service plantaoStarter, logger log.Logger, interval time.Duration) *PlantaoStatusWorker {
	return &PlantaoStatusWorker{
		service:  service,
		log:      logger,
		interval: interval,
	}
}

func RegisterPlantaoStatusWorker(lifecycle fx.Lifecycle, worker *PlantaoStatusWorker) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			worker.start()
			return nil
		},
		OnStop: worker.stop,
	})
}

func (w *PlantaoStatusWorker) start() {
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

func (w *PlantaoStatusWorker) execute(ctx context.Context) {
	quantidade, err := w.service.IniciarPlantoesAgendados(ctx)
	if err != nil {
		if ctx.Err() == nil {
			w.log.Error("erro no worker de início automático de plantões", "error", err)
		}
		return
	}
	if quantidade > 0 {
		w.log.Info("worker atualizou plantões para em andamento", "quantidade", quantidade)
	} else {
		w.log.Debug("worker de plantões executado sem registros elegíveis")
	}
}

func (w *PlantaoStatusWorker) stop(ctx context.Context) error {
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
