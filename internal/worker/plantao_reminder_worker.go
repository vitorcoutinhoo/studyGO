package worker

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"plantao/internal/domain/colaborador"
	"plantao/internal/domain/comunicacao"
	"plantao/internal/domain/log"
	"plantao/internal/domain/plantao"
)

const plantaoReminderInterval = 30 * time.Minute

type PlantaoReminderWorker struct {
	plantoes      plantao.PlantaoRepository
	colaboradores colaborador.ColaboradorRepository
	envio         *comunicacao.EnvioService
	log           log.Logger
	interval      time.Duration
	now           func() time.Time
	cancel        context.CancelFunc
	done          chan struct{}
	once          sync.Once
}

func NewPlantaoReminderWorker(plantoes plantao.PlantaoRepository, colaboradores colaborador.ColaboradorRepository, envio *comunicacao.EnvioService, logger log.Logger) *PlantaoReminderWorker {
	return &PlantaoReminderWorker{plantoes: plantoes, colaboradores: colaboradores, envio: envio, log: logger, interval: plantaoReminderInterval, now: time.Now}
}

func RegisterPlantaoReminderWorker(lifecycle fx.Lifecycle, worker *PlantaoReminderWorker) {
	lifecycle.Append(fx.Hook{OnStart: func(context.Context) error { worker.start(); return nil }, OnStop: worker.stop})
}

func (w *PlantaoReminderWorker) start() {
	ctx, cancel := context.WithCancel(context.Background())
	w.cancel, w.done = cancel, make(chan struct{})
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

func (w *PlantaoReminderWorker) execute(ctx context.Context) {
	status := plantao.StatusPlantaoEmAndamento
	plantoes, err := w.plantoes.Find(ctx, &plantao.Filtro{Status: &status})
	if err != nil {
		w.log.Error("erro ao buscar plantões em andamento para lembrete", "error", err)
		return
	}
	agora := w.now()
	for _, p := range plantoes {
		if !p.Periodo.Fim.Before(agora) {
			continue
		}
		id, err := uuid.Parse(p.ColaboradorId)
		if err != nil {
			w.log.Warn("plantão com colaborador inválido para lembrete", "id_plantao", p.Id)
			continue
		}
		col, err := w.colaboradores.FindById(ctx, id)
		if err != nil || col == nil {
			w.log.Warn("colaborador não encontrado para lembrete de plantão", "id_plantao", p.Id, "error", err)
			continue
		}
		data := map[string]any{string(comunicacao.Nome): col.Nome, string(comunicacao.DataInicio): p.Periodo.Inicio, string(comunicacao.DataFim): p.Periodo.Fim}
		if err := w.envio.SendEmailComunicacao(ctx, comunicacao.PlantaoAindaAberto, col.Email, col.Id, data); err != nil {
			w.log.Warn("erro ao enviar lembrete de plantão aberto", "id_plantao", p.Id, "error", err)
		}
	}
}

func (w *PlantaoReminderWorker) stop(ctx context.Context) error {
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
