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

	fmt.Println("\n* logging level: DEBUG")
	logger.Level = logging.LEVEL_DEBUG
	log_messages(logger)

	fmt.Println("\n* logging level: INFO")
	logger.Level = logging.LEVEL_INFO
	log_messages(logger)

	fmt.Println("\n* logging level: WARN")
	logger.Level = logging.LEVEL_WARN
	log_messages(logger)

	fmt.Println("\n* logging level: ERROR")
	logger.Level = logging.LEVEL_ERROR
	log_messages(logger)

	fmt.Println("\n* logging level: FATAL")
	logger.Level = logging.LEVEL_FATAL
	log_messages(logger)

	fmt.Println()
}
