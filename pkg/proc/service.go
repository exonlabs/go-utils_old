package proc

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/exonlabs/go-logging/pkg/xlog"
)

type SimpleService struct {
	*Daemon
	Routines        map[string]*Routine // buffers holding routines handlers
	MonitorInterval time.Duration       // interval to check for routines
	ExitDelay       time.Duration       // delay to wait for routines exit
}

func NewSimpleService(name string, log *xlog.Logger, debug int) *SimpleService {
	s := &SimpleService{
		Routines:        make(map[string]*Routine),
		MonitorInterval: time.Second * 5,
		ExitDelay:       time.Second * 3,
	}
	s.Daemon = NewDaemon(s, name, log, debug)
	return s
}

func (srv *SimpleService) Setup(daemon *Daemon) error { return nil }

func (srv *SimpleService) Initialize() error {
	if len(srv.Routines) > 0 {
		var routines []string
		for k := range srv.Routines {
			routines = append(routines, k)
		}
		srv.Log.Debug("loaded routines: %s",
			strings.Join(routines, ", "))
	} else {
		srv.Log.Warn("no routines loaded")
	}

	return nil
}

func (srv *SimpleService) Execute() error {
	if err := srv.CheckRoutines(); err != nil {
		return err
	}

	// wait monitoring interval
	srv.Sleep(srv.MonitorInterval)

	return nil
}

func (srv *SimpleService) Terminate() error {
	srv.Log.Info("stopping all routines")
	for _, routine := range srv.Routines {
		routine.Stop(true)
	}
	srv.Sleep(time.Millisecond * 50)

	getAlive := func() []string {
		var liveRoutins []string
		for name, routine := range srv.Routines {
			if routine.IsAlive.Load() {
				liveRoutins = append(liveRoutins, name)
			}
		}
		return liveRoutins
	}

	// check running status and wait termination
	ts := time.Now().Add(srv.ExitDelay)
	for time.Now().Before(ts) {
		if len(getAlive()) == 0 {
			return nil
		}
		srv.Sleep(time.Millisecond * 100)
	}

	if len(getAlive()) > 0 {
		srv.Log.Error("failed stopping routines: %s",
			strings.Join(getAlive(), ", "))
	}

	return nil
}

func (srv *SimpleService) AddRoutine(routine *Routine) error {
	if _, ok := srv.Routines[routine.Name]; ok {
		return fmt.Errorf("duplicate routine name: %s", routine.Name)
	}

	srv.Routines[routine.Name] = routine
	srv.Routines[routine.Name].Parent = srv
	return nil
}

func (srv *SimpleService) DelRoutine(name string) error {
	if _, ok := srv.Routines[name]; ok {
		srv.Routines[name].Stop(false)
		srv.Sleep(time.Millisecond * 100)
		if !srv.Routines[name].IsAlive.Load() {
			delete(srv.Routines, name)
		}
	}
	return nil
}

func (srv *SimpleService) StartRoutine(name string) error {
	if _, ok := srv.Routines[name]; ok {
		if !srv.Routines[name].IsAlive.Load() {
			srv.Log.Info("starting routine: %s", name)
			if err := srv.Routines[name].Start(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (srv *SimpleService) StopRoutine(name string, suspend bool, waitExit bool) error {
	if _, ok := srv.Routines[name]; ok {
		if srv.Routines[name].IsAlive.Load() {
			srv.Log.Info("stopping routine: %s", name)
			srv.Routines[name].Stop(suspend)
		}

		srv.Sleep(time.Millisecond * 50)

		if !waitExit {
			return nil
		}

		ts := time.Now().Add(srv.ExitDelay)
		for time.Now().Before(ts) {
			if !srv.Routines[name].IsAlive.Load() {
				return nil
			}
			srv.Sleep(time.Millisecond * 100)
		}

		srv.Log.Error("failed stopping routine: %s", name)
	}

	return nil
}

func (srv *SimpleService) CheckRoutines() error {
	for name, routine := range srv.Routines {
		// check if suspended routine
		if routine.IsSuspended.Load() {
			if routine.IsAlive.Load() {
				routine.Stop(true)
			}
			continue
		}

		// check routine status
		if !routine.IsAlive.Load() {
			if routine.InitialRun {
				srv.Log.Warn("found dead routine: %s", name)
			}
			srv.Log.Info("starting routine: %s", name)
			if err := routine.Start(); err != nil {
				return err
			}
		}
	}
	return nil
}

type ManagedService struct {
	*SimpleService
	ManagePipe  string                                                    // management pipe path
	CmdHandler  func(srv *ManagedService, command string) (string, error) // command handler callback function
	LastMonitor time.Duration
	RecvTimeout time.Duration
	SendTimeout time.Duration
	inPipeName  string
	outPipeName string
}

func NewManagedService(name string, log *xlog.Logger, debug int) *ManagedService {
	srv := &ManagedService{
		SimpleService: NewSimpleService(name, log, debug),
	}
	srv.Daemon = NewDaemon(srv, name, log, debug)

	return srv
}

func (srv *ManagedService) Setup(daemon *Daemon) error { return nil }

func (srv *ManagedService) Initialize() error {
	if err := srv.SimpleService.Initialize(); err != nil {
		return err
	}

	if len(srv.ManagePipe) == 0 && srv.CmdHandler == nil {
		return errors.New("ManagePipe or CmdHandler not set")
	}

	srv.LastMonitor = 0

	// init input pipes
	srv.inPipeName = srv.ManagePipe + ".in"
	os.Remove(srv.inPipeName)
	if err := syscall.Mkfifo(srv.inPipeName, 0666); err != nil {
		return err
	}
	recvTimeout := time.Second
	if (srv.MonitorInterval / 2) > recvTimeout {
		recvTimeout = (srv.MonitorInterval / 2)
	}
	srv.RecvTimeout = recvTimeout

	// init output pipes
	srv.outPipeName = srv.ManagePipe + ".out"
	os.Remove(srv.outPipeName)
	if err := syscall.Mkfifo(srv.outPipeName, 0666); err != nil {
		return err
	}
	srv.SendTimeout = time.Second * 3

	return nil
}

func (srv *ManagedService) Execute() error {
	if len(srv.ManagePipe) == 0 && srv.CmdHandler == nil {
		return srv.SimpleService.Execute()
	}

	// check routines
	if srv.LastMonitor == 0 ||
		(time.Duration(time.Now().UnixNano())-srv.LastMonitor) >= srv.MonitorInterval {
		srv.CheckRoutines()
		srv.LastMonitor = time.Duration(time.Now().UnixNano())
	}

	inPipeCh := make(chan *os.File)
	errCh := make(chan error)
	go func(inPipeCh chan *os.File, errCh chan error) {
		// Open the named pipe for reading
		pipe, err := os.OpenFile(srv.inPipeName, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			pipe.Close()
			errCh <- err
			return
		}
		inPipeCh <- pipe
	}(inPipeCh, errCh)

	// wait and handle commands
	select {
	case err := <-errCh:
		return err
	case <-time.After(srv.RecvTimeout):
	case inPipe := <-inPipeCh:
		command := bufio.NewScanner(inPipe)
		for command.Scan() {
			srv.handleCommand(command.Text())
		}
		if err := inPipe.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (srv *ManagedService) Terminate() error {
	if len(srv.ManagePipe) != 0 && srv.CmdHandler != nil {
		// remove inPipe file
		if err := os.Remove(srv.inPipeName); err != nil {
			return err
		}

		// remove outPipe file
		if err := os.Remove(srv.outPipeName); err != nil {
			return err
		}
	}

	if err := srv.SimpleService.Terminate(); err != nil {
		return err
	}

	return nil
}

func (srv *ManagedService) handleCommand(command string) {
	reply, err := srv.CmdHandler(srv, command)
	if err != nil {
		reply = "INTERNAL_ERROR"
		srv.Log.Error("COMMAND_ERROR: %s", err.Error())
	}
	if len(reply) == 0 {
		reply = "DONE"
	}

	outPipeCh := make(chan *os.File)

	go func(outPipeCh chan *os.File) {
		// Open the named pipe for writing
		pipe, err := os.OpenFile(srv.outPipeName, os.O_WRONLY, os.ModeNamedPipe)
		if err != nil {
			pipe.Close()
			srv.Log.Error(err.Error())
			return
		}
		outPipeCh <- pipe
	}(outPipeCh)

	select {
	case <-time.After(srv.SendTimeout):
	case outPipe := <-outPipeCh:
		// Write reply to the named pipe
		_, err = outPipe.WriteString(reply + "\n")
		if err != nil {
			srv.Log.Error(err.Error())
		}

		// Close the named pipe
		if err := outPipe.Close(); err != nil {
			srv.Log.Error(err.Error())
		}
	}
}
