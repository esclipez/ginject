package main

import (
	"context"
	"fmt"

	"github.com/esclipez/ginject/boot"
)

// Logger Interface for demonstration
type Logger interface {
	Log(message string)
}

// FileLogger implementation
type FileLogger struct {
	filename string
}

func NewFileLogger() *FileLogger {
	return &FileLogger{filename: "app.log"}
}

func (f *FileLogger) Log(message string) {
	fmt.Printf("[FileLogger] %s\n", message)
}

// ConsoleLogger implementation
type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (c *ConsoleLogger) Log(message string) {
	fmt.Printf("[ConsoleLogger] %s\n", message)
}

// DebuggerLogger implementation
type DebuggerLogger struct{}

func NewDebuggerLogger() *DebuggerLogger {
	return &DebuggerLogger{}
}

func (c *DebuggerLogger) Log(message string) {
	fmt.Printf("[DebuggerLogger] %s\n", message)
}

type StructComponent struct {
	// Another specific qualifier
	ConsoleLogger Logger `autowire:"ConsoleLog"`

	// Required dependency (explicit)
	RequiredLogger Logger `autowire:"required"`
}

type PointerComponent struct {
	// Optional dependency - won't fail if not found
	OptionalLogger Logger `autowire:"optional"`

	// Alternative optional syntax
	MaybeLogger *FileLogger `autowire:"?"`
}

// AdvancedService Service demonstrating different autowiring patterns
type AdvancedService struct {
	// Default autowiring - injects by type
	DefaultLogger Logger `autowire:""`

	// Specific qualifier - injects component by name
	FileLogger Logger `autowire:"FileLog"`

	StructComponent

	*PointerComponent

	// Won't fail if "NonExistentLogger" doesn't exist
	SpecificOptional Logger `autowire:"NonExistentLogger,optional"`
}

func NewAdvancedService() *AdvancedService {
	return &AdvancedService{
		PointerComponent: &PointerComponent{},
	}
}

func (a *AdvancedService) Start(_ context.Context) error {
	fmt.Println("[AdvancedService] Testing different injection patterns:")

	if a.DefaultLogger != nil {
		a.DefaultLogger.Log("Default logger works")
	}

	if a.FileLogger != nil {
		a.FileLogger.Log("File logger works")
	}

	if a.ConsoleLogger != nil {
		a.ConsoleLogger.Log("Console logger works")
	}

	if a.RequiredLogger != nil {
		a.RequiredLogger.Log("Required logger works")
	}

	if a.OptionalLogger != nil {
		a.OptionalLogger.Log("Optional logger works")
	} else {
		fmt.Println("[AdvancedService] Optional logger not injected")
	}

	if a.MaybeLogger != nil {
		a.MaybeLogger.Log("Maybe logger works")
	} else {
		fmt.Println("[AdvancedService] Maybe logger not injected")
	}

	loggers, err := boot.GetAllByType((*Logger)(nil))
	if err != nil {
		return err
	}
	fmt.Printf("[AdvancedService] Number of work loggers:%d\n", len(loggers))
	return nil
}

// Register advanced components
func init() {
	// Register FileLogger with specific name
	boot.Object(NewFileLogger()).
		Export((*Logger)(nil)).
		Name("FileLog")

	// Register ConsoleLogger with specific name
	boot.Object(NewConsoleLogger()).
		Export((*Logger)(nil)).
		Primary().
		Name("ConsoleLog")

	boot.Object(NewDebuggerLogger()).
		Export((*Logger)(nil)).
		Name("DebuggerLog")

	// Register the advanced service
	boot.Object(NewAdvancedService()).Priority(5)
}
