package adapter

import (
	"fmt"
	"os"
)

type LogLevel int

func (r LogLevel) String() string {
	switch r {
	case LogLevelTrace:
		return "TRACE"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
	LogLevelTrace
)

type StdLogger struct {
	instance string
	logLevel LogLevel
}

func NewStdLogger(instance string, logLevel LogLevel) *StdLogger {
	return &StdLogger{
		instance,
		logLevel,
	}
}

func (r *StdLogger) Trace(msg string) {
	r.log(LogLevelTrace, msg)
}

func (r *StdLogger) Debug(msg string) {
	r.log(LogLevelDebug, msg)
}

func (r *StdLogger) Info(msg string) {
	r.log(LogLevelInfo, msg)
}

func (r *StdLogger) Warn(msg string) {
	r.log(LogLevelWarn, msg)
}

func (r *StdLogger) Error(msg string) {
	r.log(LogLevelError, msg)
}

func (r *StdLogger) log(level LogLevel, msg string) {
	if level > r.logLevel {
		return
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s: (%s) %s\n", level.String(), r.instance, msg)
}
