package main

import (
    "fmt"
    "github.com/exonlabs/go-utils/pkg/utils/console"
)

func main() {
    defer panicHandler()

    ch := console.NewConsole()

    opts := map[string]any{"required": true}
    res, err := ch.GetValue("Enter your first name", opts)
    if err != nil {panic(err)}
    fmt.Printf("  * First name: %s\n", res)

    // ch.Passwd()
}

func panicHandler() {
    err := recover();
    if err != nil && fmt.Sprint(err) == "EOF" {
        fmt.Printf("\n-- terminated --\n")
        return
    }
    panic(err)
}
