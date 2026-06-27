package boot

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

type runApplicationShutdownComponent struct{}

func (c *runApplicationShutdownComponent) Start(context.Context) error {
	Shutdown()
	return nil
}

type capturingLogger struct {
	info  []string
	error []string
	fatal []string
}

func (l *capturingLogger) Debug(args ...interface{}) {}

func (l *capturingLogger) Debugf(format string, args ...interface{}) {}

func (l *capturingLogger) Info(args ...interface{}) {
	l.info = append(l.info, fmt.Sprint(args...))
}

func (l *capturingLogger) Infof(format string, args ...interface{}) {
	l.info = append(l.info, fmt.Sprintf(format, args...))
}

func (l *capturingLogger) Warn(args ...interface{}) {}

func (l *capturingLogger) Warnf(format string, args ...interface{}) {}

func (l *capturingLogger) Error(args ...interface{}) {
	l.error = append(l.error, fmt.Sprint(args...))
}

func (l *capturingLogger) Errorf(format string, args ...interface{}) {
	l.error = append(l.error, fmt.Sprintf(format, args...))
}

func (l *capturingLogger) Fatal(args ...interface{}) {
	l.fatal = append(l.fatal, fmt.Sprint(args...))
}

func (l *capturingLogger) Fatalf(format string, args ...interface{}) {
	l.fatal = append(l.fatal, fmt.Sprintf(format, args...))
}

func TestRunApplicationLogsCleanLifecycleMessages(t *testing.T) {
	oldContainer := defaultContainer
	oldShutdownChan := shutdownChan
	oldLogger := defaultLogger
	defer func() {
		defaultContainer = oldContainer
		shutdownChan = oldShutdownChan
		defaultLogger = oldLogger
	}()

	logger := &capturingLogger{}
	defaultContainer = NewContainer()
	shutdownChan = make(chan struct{}, 1)
	defaultLogger = logger

	Object(&runApplicationShutdownComponent{})

	RunApplication()

	expectedInfo := []string{
		"ginject: starting application",
		"ginject: application started",
		"ginject: shutdown requested",
		"ginject: stopping application",
		"ginject: application stopped",
	}
	if !reflect.DeepEqual(logger.info, expectedInfo) {
		t.Fatalf("expected info logs %v, got %v", expectedInfo, logger.info)
	}
	if len(logger.error) != 0 {
		t.Fatalf("expected no error logs, got %v", logger.error)
	}
	if len(logger.fatal) != 0 {
		t.Fatalf("expected no fatal logs, got %v", logger.fatal)
	}
}
