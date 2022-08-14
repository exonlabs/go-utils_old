package main

import (
	"fmt"

	"github.com/exonlabs/go-utils/pkg/logging"
)

func main() {
	var parent, child *logging.Logger

	fmt.Println("\nExample parent logger:")
	parent = logging.NewLogger()
	parent.Name = "parent"
	parent.SetFormatter("%(asctime)s - %(levelname)s [%(name)s] -- %(message)s")
	parent.Debug("Parent Log Msg for %s", "debug")
	parent.Info("Parent Log Msg %s", "info")
	parent.Warn("Parent Log Msg %s", "warn")
	parent.Error("Parent Log Msg %s", "error")
	parent.Fatal("Parent Log Msg %s", "fatal")

	fmt.Println("\nExample child logger with debug:")
	child = logging.NewLogger()
	child.Name = "child"
	child.Parent = parent
	child.Level = logging.DEBUG
	child.SetFormatter("%(asctime)s : %(levelname)s [%(name)s] : %(message)s")
	child.Debug("Child Log Msg for %s", "debug")
	child.Info("Child Log Msg %s", "info")
	child.Warn("Child Log Msg %s", "warn")
	child.Error("Child Log Msg %s", "error")
	child.Fatal("Child Log Msg %s", "fatal")

}
