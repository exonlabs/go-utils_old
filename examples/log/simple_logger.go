package main

import (
	"fmt"
	"github.com/exonlabs/go-utils/pkg/logging"
)

func main() {
	var logger *logging.Logger

	fmt.Println("\nExample default logger:")
	logger = logging.NewLogger()
	logger.Debug("Log Msg for %s", "debug")
	logger.Info("Log Msg for %s", "info")
	logger.Warn("Log Msg for %s", "warn")
	logger.Error("Log Msg for %s", "error")
	logger.Fatal("Log Msg for %s", "fatal")

	fmt.Println("\nExample logger with debug: ")
	logger = logging.NewLogger()
	logger.Level = logging.DEBUG
	logger.Debug("Log Msg for %s", "debug")
	logger.Info("Log Msg for %s", "info")
	logger.Warn("Log Msg for %s", "warn")
	logger.Error("Log Msg for %s", "error")
	logger.Fatal("Log Msg for %s", "fatal")

	fmt.Println("\nExample logger with error: ")
	logger = logging.NewLogger()
	logger.Level = logging.ERROR
	logger.Debug("Log Msg for %s", "debug")
	logger.Info("Log Msg for %s", "info")
	logger.Warn("Log Msg for %s", "warn")
	logger.Error("Log Msg for %s", "error")
	logger.Fatal("Log Msg for %s", "fatal")

	fmt.Println("\nExample default logger and change formatting:")
	logger = logging.NewLogger()
	logger.Name = "main"
	logger.SetFormatter("%(asctime)s - %(levelname)s [%(name)s] -- %(message)s")
	logger.Debug("Log Msg for %s", "debug")
	logger.Info("Log Msg for %s", "info")
	logger.Warn("Log Msg for %s", "warn")
	logger.Error("Log Msg for %s", "error")
	logger.Fatal("Log Msg for %s", "fatal")

	fmt.Println("\nExample logger file with default format style:")
	fmt.Println("--> see log1.txt")
	logger = logging.NewLogger()
	handlerf1 := logging.NewFileHandler("log1.txt")
	logger.ClearHandlers()
	logger.AddHandler(handlerf1)
	for i := 1; i <= 3; i++ {
		logger.Debug("Log Msg for %d", i)
		logger.Info("Log Msg for %d", i)
		logger.Warn("Log Msg for %d", i)
		logger.Fatal("Log Msg for %d", i)
	}

	fmt.Println("\nExample logger stdout and file with change format style:")
	fmt.Println("--> see log2.txt")
	logger = logging.NewLogger()
	handlerf2 := logging.NewFileHandler("log2.txt")
	handlerf2.SetFormatter("%(levelname)s: %(message)s")
	logger.AddHandler(handlerf2)
	for i := 1; i <= 3; i++ {
		logger.Debug("Log Msg for %d", i)
		logger.Info("Log Msg for %d", i)
		logger.Warn("Log Msg for %d", i)
		logger.Fatal("Log Msg for %d", i)
	}
}
