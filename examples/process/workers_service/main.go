package main

// import (
// 	"flag"
// 	"fmt"
// 	"math/rand"
// 	"time"

// 	"github.com/exonlabs/go-utils/logging"
// 	"github.com/exonlabs/go-utils/logging/handlers"

// 	"github.com/exonlabs/go-utils/process"
// )

// type worker struct {
// 	*process.Routine
// 	close bool
// }

// func (wr *worker) Setup(r *process.Routine) error {
// 	wr.Routine = r
// 	return nil
// }

// func (wr *worker) Initialize() error {
// 	wr.Log.Info("initializing")
// 	return nil
// }

// func (wr *worker) Execute() error {
// 	if rand.Intn(10) >= 8 {
// 		wr.Log.Info("closing myself")
// 		wr.close = true
// 		wr.Stop(false)
// 		return nil
// 	}

// 	wr.Log.Info("running %s", wr.Name)
// 	wr.Sleep(time.Second * 2)

// 	return nil
// }

// func (wr *worker) Terminate() error {
// 	if !wr.close && rand.Intn(10) >= 5 {
// 		wr.Log.Info("i will not exit")
// 		for {
// 			wr.Sleep(time.Second * 10)
// 		}

// 	}
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

// 	srv := process.NewSimpleService("WorkersService", logger, process.LogLevel(*debug))
// 	srv.MonitorInterval = time.Second
// 	srv.ExitDelay = time.Second

// 	for i := 0; i < 3; i++ {
// 		name := fmt.Sprintf("worker_%d", (i + 1))
// 		r := process.NewRoutine(&worker{}, name, nil, process.LogLevel(*debug))
// 		if err := srv.AddRoutine(r); err != nil {
// 			logger.Error(err.Error())
// 		}
// 	}
// 	srv.Start()

// 	logger.Info("exit")
// }
