package boot

import (
	"fmt"
	"log"
	"os"
)

// Logger defines the logging interface
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

// DefaultLogger implements Logger using standard log package
type DefaultLogger struct {
	logger *log.Logger
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *DefaultLogger) Debug(args ...interface{}) {
	l.logger.Print("[DEBUG] ", fmt.Sprint(args...))
}

func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	l.logger.Printf("[DEBUG] "+format, args...)
}

func (l *DefaultLogger) Info(args ...interface{}) {
	l.logger.Print("[INFO] ", fmt.Sprint(args...))
}

func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	l.logger.Printf("[INFO] "+format, args...)
}

func (l *DefaultLogger) Warn(args ...interface{}) {
	l.logger.Print("[WARN] ", fmt.Sprint(args...))
}

func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
	l.logger.Printf("[WARN] "+format, args...)
}

func (l *DefaultLogger) Error(args ...interface{}) {
	l.logger.Print("[ERROR] ", fmt.Sprint(args...))
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	l.logger.Printf("[ERROR] "+format, args...)
}

func (l *DefaultLogger) Fatal(args ...interface{}) {
	l.logger.Fatal("[FATAL] ", fmt.Sprint(args...))
}

func (l *DefaultLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf("[FATAL] "+format, args...)
}

var (
	defaultLogger Logger = NewDefaultLogger()
)

// SetLogger allows replacing the default logger
func SetLogger(logger Logger) {
	defaultLogger = logger
}

// GetLogger returns the current logger instance
func GetLogger() Logger {
	return defaultLogger
}

// Debug Global logging functions
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}
