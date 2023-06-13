package proc

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/exonlabs/go-logging/pkg/xlog"
)

type Daemon struct {
	Process   IPorcess
	Ctx       context.Context
	cancelCtx context.CancelFunc
	Name      string       // daemon name
	Log       *xlog.Logger // daemon logger
	Debug     int          // debug mode
	exitCh    chan bool
	Signals   SignalHandlers
	ExitDelay time.Duration
}

func NewDaemon(process IPorcess, name string, log *xlog.Logger, debug int) *Daemon {
	ctx, cancel := context.WithCancel(context.Background())

	daemon := &Daemon{
		Process:   process,
		Name:      name,
		Log:       log,
		Debug:     debug,
		Ctx:       ctx,
		cancelCtx: cancel,
		exitCh:    make(chan bool, 1),
		ExitDelay: time.Second * 5,
		Signals: SignalHandlers{
			syscall.SIGINT:  nil,
			syscall.SIGTERM: nil,
			syscall.SIGQUIT: nil,
			syscall.SIGHUP:  nil,

			syscall.SIGUSR1: nil,
			syscall.SIGUSR2: nil,
		},
	}

	process.Setup(daemon)

	return daemon
}

func (d *Daemon) String() string {
	return fmt.Sprintf("<%T: %s>", d.Process, d.Name)
}

func (d *Daemon) Start() {
	// set signal handlers
	var sigs []os.Signal
	for sig := range d.Signals {
		sigs = append(sigs, sig)
	}

	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl, sigs...)
	go func() {
		for {
			d.handelSignal(<-sigchnl)
		}
	}()

	// initialize
	if v, ok := d.Process.(IPorcessInitialize); ok {
		if err := v.Initialize(); err != nil {
			d.Log.Error(err.Error())
		}
	}

	// execute
	for {
		select {
		case <-d.Ctx.Done():
			d.Log.Info("-- terminate event --")

			// graceful terminate
			if v, ok := d.Process.(IPorcessTerminate); ok {
				go func() {
					defer func() {
						d.exitCh <- true
					}()

					if err := v.Terminate(); err != nil {
						d.Log.Error(err.Error())
					}
				}()
				select {
				case <-time.After(d.ExitDelay):
				case <-d.exitCh:
				}
			}

			return
		default:
			// execute
			if err := d.Process.Execute(); err != nil {
				d.Log.Error(err.Error())
				d.Sleep(time.Second)
			}
		}
	}
}

func (d *Daemon) Stop() {
	d.cancelCtx()
	d.exitCh <- true
}

func (d *Daemon) Sleep(timeout time.Duration) {
	select {
	case <-time.After(timeout):
	case <-d.exitCh:
	}
}

func (d *Daemon) handelSignal(s os.Signal) {
	fun, _ := d.Signals[s]
	if fun != nil {
		d.Log.Debug("<received: %s>", s)
		d.Signals[s](s)
		return
	}

	switch s {
	case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP:
		d.Log.Info("<received: %s>", s)
		d.Stop()
	default:
		d.Log.Debug("<received: %s> - ignoring", s)
	}
}
