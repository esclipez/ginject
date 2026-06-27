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
	shutdownChan     = make(chan struct{}, 1)
)

// Object provides global access to component registration
func Object(instance interface{}) *ObjectBuilder {
	return defaultContainer.Object(instance)
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
	Info("ginject: starting application")
	if err := defaultContainer.Run(ctx); err != nil {
		Fatalf("ginject: startup failed: %v", err)
	}

	Info("ginject: application started")

	// Wait for shutdown signal (either OS signal or programmatic shutdown)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		Info("ginject: shutdown requested by OS signal")
	case <-shutdownChan:
		Info("ginject: shutdown requested")
	}

	// Graceful shutdown
	Info("ginject: stopping application")
	if err := defaultContainer.Stop(ctx); err != nil {
		Errorf("ginject: shutdown failed: %v", err)
	}
	Info("ginject: application stopped")
}
