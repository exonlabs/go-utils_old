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
	logger.Formatter =
		"%(asctime)s [%(name)s] %(levelname)s -- %(message)s"

	hnd2 := handlers.NewFileHandler("/tmp/foobar.log")
	logger.AddHandler(hnd2)

	fmt.Println("\n* logging stdout and file:", hnd2.Filepath)
	log_messages(logger)

	log1 := logging.NewLogger("root.child1")
	log1.Parent = logger
	log1.Level = logging.LEVEL_ERROR
	fmt.Println("\n* logging child logger:", log1.Name)
	log_messages(log1)

	fmt.Println()
}
