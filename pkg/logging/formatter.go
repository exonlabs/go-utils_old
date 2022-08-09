package logging

import (
	"fmt"
	"time"
)

func formater(fLog string, name string, level string, format string, args ...any) string {
	// default log format without name
	logFormat := "%[1]s %[2]s %[4]s"
	// default log format with name
	if len(name) > 0 {
		logFormat = "%[1]s %[2]s  [%[3]s] %[4]s"
	}
	// custom format
	if len(fLog) > 0 {
		logFormat = fLog
	}

	return fmt.Sprintf(logFormat+"\n",
		time.Now().Local().Format("2006-01-02 15:04:05.000000"),
		level,
		name,
		fmt.Sprintf(format, args...))
}
