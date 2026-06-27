package boot

import (
	"context"
	"strings"
	"testing"
	"time"
)

type containerRunLogger interface {
	Log(string)
}

type containerRunConsoleLogger struct{}

func (l *containerRunConsoleLogger) Log(string) {}

type containerRunApp struct {
	Logger  containerRunLogger `autowire:""`
	Started bool
}

func (a *containerRunApp) Start(context.Context) error {
	if a.Logger != nil {
		a.Started = true
	}
	return nil
}

func TestContainerRunRegistersPendingBuildersFromContainerObject(t *testing.T) {
	container := NewContainer()
	app := &containerRunApp{}

	container.Object(&containerRunConsoleLogger{}).Export((*containerRunLogger)(nil))
	container.Object(app)

	if err := container.Run(context.Background()); err != nil {
		t.Fatalf("expected custom container to run registered components, got %v", err)
	}
	if !app.Started {
		t.Fatal("expected custom container object to be registered, injected, and started")
	}
	if app.Logger == nil {
		t.Fatal("expected dependency to be injected")
	}
}

type blockingStartComponent struct {
	started chan struct{}
	release chan struct{}
}

func (c *blockingStartComponent) Start(context.Context) error {
	c.started <- struct{}{}
	<-c.release
	return nil
}

func TestContainerStartSerializesConcurrentCalls(t *testing.T) {
	container := NewContainer()
	component := &blockingStartComponent{
		started: make(chan struct{}, 2),
		release: make(chan struct{}),
	}

	container.Object(component)
	if err := container.registerPendingBuilders(); err != nil {
		t.Fatalf("expected pending builders to register, got %v", err)
	}

	firstErr := make(chan error, 1)
	go func() {
		firstErr <- container.Start(context.Background())
	}()

	<-component.started

	secondErr := make(chan error, 1)
	go func() {
		secondErr <- container.Start(context.Background())
	}()

	select {
	case <-component.started:
		t.Fatal("expected concurrent Start call not to start components twice")
	case <-time.After(50 * time.Millisecond):
	}

	close(component.release)

	if err := <-firstErr; err != nil {
		t.Fatalf("expected first Start to succeed, got %v", err)
	}
	err := <-secondErr
	if err == nil {
		t.Fatal("expected second Start to fail")
	}
	if !strings.Contains(err.Error(), "container already started") {
		t.Fatalf("expected already started error, got %v", err)
	}
}

func TestContainerObjectPanicsAfterRun(t *testing.T) {
	container := NewContainer()
	container.Object(&containerRunConsoleLogger{})

	if err := container.Run(context.Background()); err != nil {
		t.Fatalf("expected container to run, got %v", err)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected Object to panic after container has run")
		}
	}()

	container.Object(&containerRunConsoleLogger{})
}
