package main

import (
	"flag"
	"os"
	"syscall"
	"time"

	"github.com/exonlabs/go-logging/pkg/xlog"
	"github.com/exonlabs/go-utils/pkg/proc"
)

type sampleProcess struct {
	*proc.Daemon
	counter int
}

func (p *sampleProcess) Setup(daemon *proc.Daemon) error {
	p.Daemon = daemon
	// daemon.Signals[syscall.SIGINT] = p.handleSig
	daemon.Signals[syscall.SIGUSR2] = p.handleSigusr2
	return nil
}

func (p *sampleProcess) Initialize() error {
	p.Log.Info("initializing")

	return nil
}

func (p *sampleProcess) Execute() error {
	p.counter += 1
	p.Log.Debug("running: %d ...", p.counter)
	if p.counter >= 60 {
		p.Log.Info("exit process count = %d", p.counter)
		p.Stop()
	}

	p.Sleep(time.Second)

	return nil
}

func (p *sampleProcess) Terminate() error {
	exitCounts := 2
	p.Log.Info("exit after %d counts", exitCounts)
	for i := 0; i < exitCounts; i++ {
		p.Log.Info("count %d", (i + 1))
		p.Sleep(time.Second * 10)
	}

	p.Log.Info("exit")

	return nil
}

func (p *sampleProcess) handleSigusr2(sig os.Signal) error {
	p.counter = 0
	p.Log.Info("counter reset")

	return nil
}

func (p *sampleProcess) handleSig(os.Signal) error {
	p.Log.Info("Over Write")
	os.Exit(0)
	return nil
}

func main() {
	logger := xlog.NewLogger("main")
	logger.Formatter =
		"%(asctime)s - %(levelname)s [%(name)s] %(message)s"

	debug := flag.Int("x", 0, "set debug modes, (default: 0)")
	flag.Parse()

	if *debug >= 3 {
		logger.Level = xlog.LevelDebug
	}

	logger.Info("**** starting ****")

	daemon := proc.NewDaemon(&sampleProcess{}, "SampleDaemon",
		logger, int(*debug))
	logger.Info("Created service: %s", daemon)

	daemon.Start()
}
