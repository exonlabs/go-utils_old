package main

import (
    "fmt"
    "github.com/exonlabs/go-utils/pkg/process/daemon"
)

type SampleDaemon struct {
    daemon.BaseDaemon

    counter uint32
}


func NewSampleDaemon() *SampleDaemon {
    d := SampleDaemon{}
    d.Daemon = &d
    return &d
}

func (d *SampleDaemon) Initialize() error {
    fmt.Println("initializing")
    d.counter = 0
    return nil
}

func (d *SampleDaemon) Execute() error {
    d.counter++
    fmt.Printf("running: %v ...\n", d.counter)
    if d.counter >= 10 {
        fmt.Println("exit process count =", d.counter)
        return fmt.Errorf("error breaking")
    }
    d.Daemon.Sleep(1)
    return nil
}

func (d *SampleDaemon) Terminate() error {
    fmt.Println("shutting down")
    fmt.Println("exit after 5 counts")
    for i := 0; i < 5; i++ {
        fmt.Printf("count %v\n", i + 1)
        d.Daemon.Sleep(1)
    }
    fmt.Println("exit")
    return nil
}


func main() {
    defer panicHandler()

    d := NewSampleDaemon()
    d.Name = "SampleDaemon"
    d.Proctitle = "SampleDaemon"

    d.Start()
}

func panicHandler() {
    err := recover();
    if err != nil {
        if fmt.Sprint(err) == "EOF" {
            fmt.Printf("\n-- terminated --\n")
            return
        }
        panic(err)
    }
}
