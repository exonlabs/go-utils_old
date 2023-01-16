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
	logger.Level = logging.LEVEL_DEBUG

	fmt.Println("\n* logging default formatter:", logger.Formatter)
	log_messages(logger)

	logger.Formatter =
		"%(asctime)s [%(name)s] %(levelname)s -- %(message)s"
	fmt.Println("\n* logging formatter:", logger.Formatter)
	log_messages(logger)

	hnd2 := handlers.NewStdoutHandler()
	hnd2.Formatter =
		"%(asctime)s [%(name)s] %(levelname)s -- %(message)s -- (hnd2)"
	logger.AddHandler(hnd2)

	fmt.Println("\n* logging with 2 handlers")
	log_messages(logger)

	fmt.Println()
}
