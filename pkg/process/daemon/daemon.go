package daemon

import (
    "fmt"
    "time"
)


type Daemon interface {
    Initialize() error
    Execute() error
    Terminate() error
    Run() error
    Start() error
    Stop() error
    Sleep(float32)
}

type BaseDaemon struct {
    Daemon

    Name string
    Proctitle string
}


func NewBaseDaemon() *BaseDaemon {
    d := BaseDaemon{}
    d.Daemon = &d
    return &d
}

func (d *BaseDaemon) Initialize() error {
    // nothing todo by default
    return nil
}

func (d *BaseDaemon) Execute() error {
    panic("NOT_IMPLEMENTED")
    return nil
}

func (d *BaseDaemon) Terminate() error {
    // nothing todo by default
    return nil
}

func (d *BaseDaemon) Run() error {
    if err := d.Daemon.Initialize(); err != nil {
        fmt.Println("FATAL", err)
        return err
    }

    for {
        if err := d.Daemon.Execute(); err != nil {
            fmt.Println("FATAL", err)
            break
        }
    }

    if err := d.Daemon.Terminate(); err != nil {
        fmt.Println("ERROR", err)
    }

    return nil
}

func (d *BaseDaemon) Start() error {
    d.Daemon.Run()
    return nil
}

func (d *BaseDaemon) Stop() error {
    return nil
}

func (d *BaseDaemon) Sleep(timeout float32) {
    time.Sleep(time.Millisecond * time.Duration(int64(1000.0 * timeout)))
}
