package main

import (
	"fmt"

	"github.com/exonlabs/go-utils/pkg/logging"
)

func main() {
	fmt.Println("example logger without name:")
	logger := logging.GetLogger("")
	logger.SetFormatter("%(asctime)s %(levelname)s [%(name)s] %(message)s")
	logger.Info("Log Msg for %s", "info")
	logger.Warn("Log Msg for %s", "warn")
	logger.Debug("Log Msg for %s", "debug")
	logger.Fatal("Log Msg for %s", "Critical")

	fmt.Println()
	fmt.Println("example logger name and change format style:")
	logger1 := logging.GetLogger("first main")
	logger1.SetFormatter("(%(name)s): %(asctime)s %(levelname)s %(message)s")
	logger1.Info("Log Msg for %s", "info")
	logger1.Warn("Log Msg for %s", "warn")
	logger1.Debug("Log Msg for %s", "debug")
	logger1.Fatal("Log Msg for %s", "Critical")

	fmt.Println()
	fmt.Println("example logger name and change format style print level and message:")
	logger2 := logging.GetLogger("second main")
	logger2.SetFormatter("%(levelname)s: %(message)s")
	logger2.Info("Log Msg for %s", "info")
	logger2.Warn("Log Msg for %s", "warn")
	logger2.Debug("Log Msg for %s", "debug")
	logger2.Fatal("Log Msg for %s", "Critical")

	fmt.Println()
	fmt.Println("example logger name with default format style:")
	logger3 := logging.GetLogger("third main")
	logger3.Info("Log Msg for %s", "info")
	logger3.Warn("Log Msg for %s", "warn")
	logger3.Debug("Log Msg for %s", "debug")
	logger3.Fatal("Log Msg for %s", "Critical")

	fmt.Println()
	fmt.Println("example logger file with default format style:")
	loggerf1 := logging.GetLogger("file")
	handlerf1 := logging.NewFileHandler("log.txt")
	loggerf1.AddHandler(handlerf1)
	for i := 1; i <= 100; i++ {
		loggerf1.Info("Log Msg for %d", i)
		loggerf1.Warn("Log Msg for %d", i)
		loggerf1.Debug("Log Msg for %d", i)
		loggerf1.Fatal("Log Msg for %d", i)
	}

	fmt.Println()
	fmt.Println("example logger file and change format style print level and message:")
	loggerf2 := logging.GetLogger("file-2")
	loggerf2.SetFormatter("%(levelname)s: %(message)s")
	handlerf2 := logging.NewFileHandler("log.txt")
	loggerf2.AddHandler(handlerf2)
	for i := 1; i <= 100; i++ {
		loggerf2.Info("Log Msg for %d", i)
		loggerf2.Warn("Log Msg for %d", i)
		loggerf2.Debug("Log Msg for %d", i)
		loggerf2.Fatal("Log Msg for %d", i)
	}

}
