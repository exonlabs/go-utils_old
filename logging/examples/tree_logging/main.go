package main

import (
	"fmt"

	"github.com/exonlabs/go-utils/logging"
	"github.com/exonlabs/go-utils/logging/handlers"
)

func log_messages(logger *logging.Logger) {
	logger.Debug("logging message type: %s", "debug")
	logger.Info("logging message type: %s", "info")
	logger.Warn("logging message type: %s", "warn")
	logger.Error("logging message type: %s", "error")
	logger.Fatal("logging message type: %s", "fatal")
}

func main() {
	logger := handlers.NewStdoutLogger("root")
	logger.Formatter =
		"%(asctime)s %(levelname)s [%(name)s] -- %(message)s"

	fmt.Println("\n* logging parent logger:", logger.Name)
	log_messages(logger)

	log1 := logging.NewLogger("root.child1")
	log1.Parent = logger
	fmt.Println("\n* logging child logger:", log1.Name)
	log_messages(log1)

	log2 := handlers.NewStdoutLogger("root.child2")
	log2.Parent = logger
	log2.Formatter =
		"%(asctime)s %(levelname)s +++ (%(name)s) +++ %(message)s"
	fmt.Println("\n* logging child logger (+handlers):", log2.Name)
	log_messages(log2)

	log3 := handlers.NewStdoutLogger("child3")
	log3.Parent = logger
	log3.Propagate = true
	log3.Formatter =
		"%(asctime)s %(levelname)s +++ (%(name)s) +++ %(message)s"
	fmt.Println("\n* logging child logger (+handlers) (+propagate):", log3.Name)
	log_messages(log3)

	fmt.Println()
}
