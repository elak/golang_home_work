package logger

import (
	"fmt"
	"io"
	"os"
)

const (
	LogLevelNone = iota
	LogLevelError
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
)

type Logger struct {
	level       int8
	fileHandler io.WriteCloser
}

var defaultLogger *Logger

func logLevel(level string) int8 {
	switch level {
	case "ERROR":
		return LogLevelError
	case "WARNING":
		return LogLevelWarning
	case "INFO":
		return LogLevelInfo
	case "DEBUG":
		return LogLevelDebug
	default:
		return LogLevelNone
	}
}

func New(level string, logPath string) (*Logger, error) {
	var err error

	log := Logger{}

	log.level = logLevel(level)

	if logPath != "" {
		log.fileHandler, err = os.Create(logPath)
		if err != nil {
			return nil, err
		}
	}

	return &log, nil
}

func (l Logger) Close() error {
	if l.fileHandler == nil {
		return nil
	}

	return l.fileHandler.Close()
}

func (l Logger) LogMessage(msg string, msgLevel int8) {
	if l.level < msgLevel {
		return
	}
	fmt.Println(msg)

	if l.fileHandler != nil {
		l.fileHandler.Write([]byte(msg + "\n"))
	}
}

func (l Logger) Info(msg string) {
	l.LogMessage(msg, LogLevelInfo)
}

func (l Logger) Error(msg string) {
	l.LogMessage(msg, LogLevelError)
}

func (l Logger) Debug(msg string) {
	l.LogMessage(msg, LogLevelDebug)
}

func (l Logger) Warning(msg string) {
	l.LogMessage(msg, LogLevelWarning)
}

func LogMessage(msg string, msgLevel int8) {
	if defaultLogger != nil {
		defaultLogger.LogMessage(msg, msgLevel)
	}
}

func Info(msg string) {
	LogMessage(msg, LogLevelInfo)
}

func Error(msg string) {
	LogMessage(msg, LogLevelError)
}

func Debug(msg string) {
	LogMessage(msg, LogLevelDebug)
}

func Warning(msg string) {
	LogMessage(msg, LogLevelWarning)
}

func Start(level string, logPath string) error {
	var err error

	defaultLogger, err = New(level, logPath)

	return err
}

func Stop() error {
	var err error

	if defaultLogger != nil {
		err = defaultLogger.Close()
	}

	return err
}

func GetDefaultLogger() *Logger {
	return defaultLogger
}
