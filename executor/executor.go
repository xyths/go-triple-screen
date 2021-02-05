package executor

import (
	"context"
	"github.com/xyths/go-triple-screen/exchange"
	"go.uber.org/zap"
)

type Executor struct {
	command chan int
	ex      exchange.Exchange
	Sugar   *zap.SugaredLogger

	stopCh  chan struct{}
	running bool
}

func NewExecutor(command chan int, ex exchange.Exchange, Sugar *zap.SugaredLogger) *Executor {
	return &Executor{
		command: command,
		ex:      ex,
		Sugar:   Sugar,
		stopCh:  make(chan struct{}),
	}
}

func (e *Executor) Start(ctx context.Context) {
	go e.serve(ctx)
	e.running = true
	e.Sugar.Info("executor started")
}

func (e *Executor) Stop(ctx context.Context) {
	if e.running {
		e.stopCh <- struct{}{}
	}
}

func (e *Executor) serve(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			e.running = false
			e.Sugar.Infof("executor stopped: %s", ctx.Err())
			return
		case <-e.stopCh:
			e.running = false
			e.Sugar.Infof("executor stopped successfully")
			return
		case signal := <-e.command:
			e.Sugar.Infof("executor got command %d", signal)
		}
	}
}
