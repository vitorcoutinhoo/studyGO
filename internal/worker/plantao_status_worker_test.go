package worker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"plantao/internal/domain/log"

	"go.uber.org/fx"
)

type starterFake struct {
	mu            sync.Mutex
	chamadas      int
	emExecucao    int
	maxSimultaneo int
	erros         []error
	duracao       time.Duration
	notificacoes  chan struct{}
}

func (s *starterFake) IniciarPlantoesAgendados(ctx context.Context) (int, error) {
	s.mu.Lock()
	s.chamadas++
	s.emExecucao++
	if s.emExecucao > s.maxSimultaneo {
		s.maxSimultaneo = s.emExecucao
	}
	indice := s.chamadas - 1
	s.mu.Unlock()

	select {
	case s.notificacoes <- struct{}{}:
	default:
	}
	if s.duracao > 0 {
		select {
		case <-time.After(s.duracao):
		case <-ctx.Done():
		}
	}

	s.mu.Lock()
	s.emExecucao--
	var err error
	if indice < len(s.erros) {
		err = s.erros[indice]
	}
	s.mu.Unlock()
	return 0, err
}

type workerLoggerFake struct{}

func (workerLoggerFake) Info(string, ...any)      {}
func (workerLoggerFake) Warn(string, ...any)      {}
func (workerLoggerFake) Error(string, ...any)     {}
func (workerLoggerFake) Debug(string, ...any)     {}
func (workerLoggerFake) Fatal(string, ...any)     {}
func (l workerLoggerFake) With(...any) log.Logger { return l }
func (workerLoggerFake) Sync() error              { return nil }

type lifecycleFake struct {
	hooks []fx.Hook
}

func (l *lifecycleFake) Append(hook fx.Hook) {
	l.hooks = append(l.hooks, hook)
}

func TestPlantaoStatusWorkerExecutaImediatamenteRepeteESemSobreposicao(t *testing.T) {
	starter := &starterFake{
		erros:        []error{errors.New("falha inicial")},
		duracao:      15 * time.Millisecond,
		notificacoes: make(chan struct{}, 4),
	}
	worker := newPlantaoStatusWorker(starter, workerLoggerFake{}, 5*time.Millisecond)
	lifecycle := &lifecycleFake{}
	RegisterPlantaoStatusWorker(lifecycle, worker)
	if len(lifecycle.hooks) != 1 {
		t.Fatalf("hooks registrados = %d, esperado 1", len(lifecycle.hooks))
	}
	if err := lifecycle.hooks[0].OnStart(context.Background()); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 2; i++ {
		select {
		case <-starter.notificacoes:
		case <-time.After(300 * time.Millisecond):
			t.Fatal("worker não executou/repetiu no prazo")
		}
	}
	stopCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := lifecycle.hooks[0].OnStop(stopCtx); err != nil {
		t.Fatal(err)
	}

	starter.mu.Lock()
	defer starter.mu.Unlock()
	if starter.chamadas < 2 {
		t.Fatalf("chamadas = %d, esperado ao menos 2", starter.chamadas)
	}
	if starter.maxSimultaneo != 1 {
		t.Fatalf("execuções simultâneas = %d, esperado 1", starter.maxSimultaneo)
	}
}
