package main

// import (
// 	"errors"
// 	"flag"
// 	"fmt"
// 	"os"
// 	"path"
// 	"strings"
// 	"sync/atomic"
// 	"time"

// 	"github.com/exonlabs/go-utils/logging"
// 	"github.com/exonlabs/go-utils/logging/handlers"

// 	"github.com/exonlabs/go-utils/process"
// )

// var (
// 	WorkersIndex atomic.Int32
// )

// type worker struct {
// 	*process.Routine
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
// 	wr.Log.Info("running %s", wr.Name)
// 	wr.Sleep(time.Second * 2)
// 	return nil
// }

// func commandHandler(srv *process.ManagedService, command string) (string, error) {
// 	if command == "EXIT" {
// 		srv.Stop()
// 		return "", nil
// 	}

// 	if strings.Contains(command, "ADD_WORKER") {
// 		WorkersIndex.Add(1)

// 		rt := process.NewRoutine(&worker{},
// 			fmt.Sprintf("worker_%d", WorkersIndex.Load()), nil, 0)
// 		if err := srv.AddRoutine(rt); err != nil {
// 			return "", err
// 		}

// 		return "DONE", nil
// 	}

// 	if strings.Contains(command, "DEL_WORKER") {
// 		index := strings.Split(command, ":")[1]
// 		if err := srv.DelRoutine(fmt.Sprintf("worker_%s", index)); err != nil {
// 			return "", err
// 		}
// 		return "DONE", nil
// 	}

// 	if strings.Contains(command, "START_WORKER") {
// 		index := strings.Split(command, ":")[1]
// 		if err := srv.StartRoutine(fmt.Sprintf("worker_%s", index)); err != nil {
// 			return "", err
// 		}
// 		return "DONE", nil
// 	}

// 	if strings.Contains(command, "STOP_WORKER") {
// 		index := strings.Split(command, ":")[1]
// 		if err := srv.StopRoutine(fmt.Sprintf("worker_%s", index),
// 			false, false); err != nil {
// 			return "", err
// 		}
// 		return "DONE", nil
// 	}

// 	if strings.Contains(command, "SUSPEND_WORKER") {
// 		index := strings.Split(command, ":")[1]
// 		if err := srv.StopRoutine(fmt.Sprintf("worker_%s", index),
// 			true, false); err != nil {
// 			return "", err
// 		}
// 		return "DONE", nil
// 	}

// 	return "", errors.New("INVALID_COMMAND")
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

// 	srv := process.NewManagedService("ManagedService",
// 		logger, process.LogLevel(*debug))
// 	srv.MonitorInterval = time.Second * 3
// 	srv.ManagePipe = path.Join(os.TempDir(),
// 		fmt.Sprintf("%s.pipe", strings.ToLower(srv.Name)))
// 	srv.CmdHandler = commandHandler

// 	for i := 0; i < 3; i++ {
// 		WorkersIndex.Add(1)
// 		name := fmt.Sprintf("worker_%d", WorkersIndex.Load())
// 		r := process.NewRoutine(&worker{}, name, nil, process.LogLevel(*debug))
// 		if err := srv.AddRoutine(r); err != nil {
// 			logger.Error(err.Error())
// 		}
// 	}

// 	srv.Start()

// 	logger.Info("exit")
// }
