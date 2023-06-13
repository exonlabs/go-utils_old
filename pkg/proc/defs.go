package proc

import "os"

type SignalHandlers = map[os.Signal]func(os.Signal) error

type IPorcess interface {
	Setup(*Daemon) error
	Execute() error
}

type IPorcessInitialize interface {
	Initialize() error
}

type IPorcessTerminate interface {
	Terminate() error
}

type IRoutine interface {
	Setup(*Routine) error
	Execute() error
}

type IRoutineInitial interface {
	Initialize() error
}

type IRoutineTerminate interface {
	Terminate() error
}
