package boot

import (
	"context"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

var (
	defaultContainer = NewContainer()
	pendingBuilders  []*ObjectBuilder // Pending builders to be registered
	shutdownChan     = make(chan struct{}, 1)
)

// Object provides global access to component registration
func Object(instance interface{}) *ObjectBuilder {
	builder := defaultContainer.Object(instance)
	// Don't register immediately, add to pending list instead
	pendingBuilders = append(pendingBuilders, builder)
	return builder
}

// GetByName retrieves a component by name from the default container
func GetByName(name string) (interface{}, error) {
	return defaultContainer.GetByName(name)
}

// GetByType retrieves a component by type from the default container
func GetByType(componentType interface{}) (interface{}, error) {
	t := reflect.TypeOf(componentType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // Get the interface type
	}
	return defaultContainer.GetByType(t)
}

// GetAllByType retrieves all components by type from the default container
func GetAllByType(componentType interface{}) ([]interface{}, error) {
	t := reflect.TypeOf(componentType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // Get the interface type
	}
	return defaultContainer.GetAllByType(t)
}

// Shutdown triggers graceful shutdown of the application
func Shutdown() {
	select {
	case shutdownChan <- struct{}{}:
		// Signal sent successfully
	default:
		// Channel already has a signal, ignore
	}
}

// RunApplication starts the application with the default container
func RunApplication() {
	ctx := context.Background()

	// Run the complete lifecycle
	Info("=== Starting Application ===")
	if err := defaultContainer.Run(ctx); err != nil {
		Fatalf("Application startup failed: %v", err)
	}

	Info("=== Application Started ===")

	// Wait for shutdown signal (either OS signal or programmatic shutdown)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		Info("=== Received OS Signal ===")
	case <-shutdownChan:
		Info("=== Received Shutdown Signal ===")
	}

	// Graceful shutdown
	Info("=== Shutting Down Application ===")
	if err := defaultContainer.Stop(ctx); err != nil {
		Errorf("Shutdown error: %v", err)
	}
	Info("=== Application Stopped ===")
}
