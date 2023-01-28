package logging

type Logger struct {
	Name      string
	Level     LogLevel
	Formatter string
	Handlers  []IHandler
	Parent    *Logger
	Propagate bool
}

func NewLogger(name string) *Logger {
	return &Logger{
		Name:      name,
		Level:     LEVEL_INFO,
		Formatter: "%(asctime)s %(levelname)s %(message)s",
		Handlers:  []IHandler{},
		Propagate: false,
	}
}

func (log *Logger) AddHandler(handler IHandler) {
	log.Handlers = append(log.Handlers, handler)
}

func (log *Logger) ClearHandlers() {
	log.Handlers = []IHandler{}
}

func (log *Logger) Log(record *Record) error {
	// handle record with loaded handlers
	for _, hnd := range log.Handlers {
		if err := hnd.HandleRecord(log, record); err != nil {
			return err
		}
	}

	// propagate to parent logger
	if log.Parent != nil && record.Level >= log.Level &&
		(log.Propagate || len(log.Handlers) == 0) {
		log.Parent.Log(record)
	}

	return nil
}

func (log *Logger) Debug(msg string, args ...any) error {
	return log.Log(
		NewRecord(log.Name, LEVEL_DEBUG, msg, args...))
}

func (log *Logger) Info(msg string, args ...any) error {
	return log.Log(
		NewRecord(log.Name, LEVEL_INFO, msg, args...))
}

func (log *Logger) Warn(msg string, args ...any) error {
	return log.Log(
		NewRecord(log.Name, LEVEL_WARN, msg, args...))
}

func (log *Logger) Error(msg string, args ...any) error {
	return log.Log(
		NewRecord(log.Name, LEVEL_ERROR, msg, args...))
}

func (log *Logger) Fatal(msg string, args ...any) error {
	return log.Log(
		NewRecord(log.Name, LEVEL_FATAL, msg, args...))
}
