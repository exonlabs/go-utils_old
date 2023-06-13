package main

// import (
// 	"flag"
// 	"sync/atomic"
// 	"time"

// 	"github.com/exonlabs/go-utils/logging"
// 	"github.com/exonlabs/go-utils/logging/handlers"

// 	"github.com/exonlabs/go-utils/process"
// )

// var (
// 	counter atomic.Int32
// )

// type routine1 struct {
// 	*process.Routine
// 	exitCount int32
// }

// func (rt *routine1) Setup(r *process.Routine) error {
// 	rt.Routine = r
// 	return nil
// }

// func (rt *routine1) Execute() error {
// 	counter.Add(1)

// 	rt.Log.Info("new count = %d", counter.Load())

// 	if counter.Load() == 5 {
// 		rt.Log.Info("suspend <%s> at count: %d", "rt2", counter.Load())
// 		rt.Parent.StopRoutine("rt2", true, true)
// 	} else if counter.Load() == 10 {
// 		rt.Log.Info("resume <%s> at count: %d", "rt2", counter.Load())
// 		rt.Parent.StartRoutine("rt2")
// 	}
// 	rt.Sleep(time.Second)
// 	return nil
// }

// func (rt *routine1) Terminate() error {
// 	rt.exitCount = counter.Load() + 2
// 	rt.Log.Info("wait exit count = %d", rt.exitCount)
// 	for counter.Load() < rt.exitCount {
// 		counter.Add(1)
// 		rt.Log.Info("new count = %d", counter.Load())
// 		rt.Sleep(time.Second / 2)
// 	}

// 	return nil
// }

// type routine2 struct {
// 	*process.Routine
// }

// func (rt *routine2) Setup(r *process.Routine) error {
// 	rt.Routine = r
// 	return nil
// }

// func (rt *routine2) Execute() error {
// 	rt.Log.Info("monitoring count = %d", counter.Load())

// 	if counter.Load() == 15 {
// 		rt.Log.Info("stopping myself at count = %d", counter.Load())
// 		rt.Sleep(time.Second)
// 		rt.Stop(false)
// 	}

// 	if counter.Load() >= 20 {
// 		rt.Log.Info("stopping service at count = %d", counter.Load())
// 		rt.Sleep(time.Second)
// 		rt.Parent.Stop()
// 	}
// 	rt.Sleep(time.Second / 2)
// 	return nil
// }

// func main() {
// 	logger := handlers.NewStdoutLogger("main")
// 	logger.Formatter =
// 		"%(asctime)s - %(levelname)s [%(name)s] %(message)s"

// 	debug := flag.Int("x", 0, "set debug modes, (default: 0)")
// 	flag.Parse()

// 	if *debug >= 3 {
// 		logger.Level = logging.LEVEL_DEBUG
// 	}

// 	logger.Info("**** starting ****")

// 	srv := process.NewSimpleService("SampleService", logger, process.LogLevel(*debug))
// 	srv.MonitorInterval = time.Second
// 	srv.ExitDelay = time.Second * 5
// 	srv.Daemon.ExitDelay = time.Second * 10
// 	rt1 := process.NewRoutine(&routine1{}, "rt1", nil, process.LogLevel(*debug))
// 	if err := srv.AddRoutine(rt1); err != nil {
// 		logger.Error(err.Error())
// 	}

// 	rt2 := process.NewRoutine(&routine2{}, "rt2", nil, process.LogLevel(*debug))
// 	if err := srv.AddRoutine(rt2); err != nil {
// 		logger.Error(err.Error())
// 	}

// 	srv.Start()

// 	logger.Info("exit")
// }
