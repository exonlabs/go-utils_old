package proc

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/exonlabs/go-logging/pkg/xlog"
)

type Routine struct {
	routine     IRoutine
	ctx         context.Context
	cancelCtx   context.CancelFunc
	Name        string
	Log         *xlog.Logger
	Debug       int
	Parent      *SimpleService
	InitialRun  bool
	IsSuspended atomic.Bool
	IsAlive     atomic.Bool
	exitCh      chan bool
}

func NewRoutine(routine IRoutine, name string, log *xlog.Logger, debug int) *Routine {
	rt := &Routine{
		routine: routine,
		Name:    name,
		Log:     log,
		Debug:   debug,
		exitCh:  make(chan bool, 1),
	}

	routine.Setup(rt)

	return rt
}

func (r *Routine) String() string {
	return fmt.Sprintf("<%T: %s>", r.routine, r.Name)
}

func (r *Routine) run() {
	r.InitialRun = true
	r.IsAlive.Store(true)
	defer func() {
		if recover := recover(); recover != nil {
			r.Log.Error("panic: %v\n", recover)
		}
	}()

	// initialize
	if v, ok := r.routine.(IRoutineInitial); ok {
		if err := v.Initialize(); err != nil {
			r.Log.Error("terminated reason: %s", err.Error())
			return
		}
	}

	for {
		select {
		case <-r.ctx.Done():
			r.Log.Info("-- terminate event --")

			// graceful terminate
			if v, ok := r.routine.(IRoutineTerminate); ok {
				if err := v.Terminate(); err != nil {
					r.Log.Error(err.Error())
				}
			}

			r.Log.Info("terminated")

			r.IsAlive.Store(false)
			return
		default:
			// execute
			if !r.IsSuspended.Load() {
				if err := r.routine.Execute(); err != nil {
					r.Log.Error(err.Error())
					r.Sleep(time.Second)
				}
			}
		}
	}
}

func (r *Routine) Start() error {
	if r.Parent == nil {
		return fmt.Errorf("no parent service handler")
	}

	r.ctx, r.cancelCtx = context.WithCancel(r.Parent.Ctx)

	if r.Log == nil {
		r.Log = xlog.NewLogger(r.Name)
		r.Log.Parent = r.Parent.Log
		r.Log.Level = r.Parent.Log.Level
	}

	r.Debug = r.Parent.Debug

	if r.Debug == 0 && r.Log.Level == xlog.LevelDebug {
		r.Debug = 1
	}

	r.IsSuspended.Store(false)

	go r.run()

	return nil
}

func (r *Routine) Stop(suspend bool) {
	if suspend {
		r.IsSuspended.Store(true)
	}

	r.cancelCtx()
	r.exitCh <- true
}

func (r *Routine) Sleep(timeout time.Duration) {
	select {
	case <-time.After(timeout):
	case <-r.exitCh:
	}
}
