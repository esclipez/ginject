# Ginject

A lightweight dependency injection container for Go applications with fluent API and lifecycle management.

Inspired by [go-spring](https://github.com/go-spring/go-spring), Ginject is a simplified version that focuses on core dependency injection features with an intuitive API.

## Features

- **Fluent API**: Chain method calls for intuitive component registration
- **Type-based Resolution**: Automatic dependency injection by type or interface
- **Lifecycle Management**: Built-in initialization, startup, and shutdown phases
- **Priority Control**: Configure component startup/shutdown order
- **Primary Components**: Resolve ambiguity when multiple components implement the same interface
- **Optional Dependencies**: Support for optional autowiring with graceful fallback

## Quick Start

### Installation

```bash
go get github.com/esclipez/ginject
```

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/esclipez/ginject/boot"
)

// Define interfaces
type Logger interface {
    Log(message string)
}

type UserService interface {
    GetUser(id string) string
}

// Implement components
type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(message string) {
    fmt.Println("[LOG]", message)
}

type UserServiceImpl struct {
    Logger Logger `autowire:""`
}

func (s *UserServiceImpl) GetUser(id string) string {
    s.Logger.Log(fmt.Sprintf("Getting user: %s", id))
    return "User-" + id
}

func main() {
    // Register components with fluent API
    boot.Object(&ConsoleLogger{}).
        Export((*Logger)(nil)).
        Name("console-logger")

    boot.Object(&UserServiceImpl{}).
        Export((*UserService)(nil)).
        Name("user-service")
    
    go func() {
        // Use components
        userService, _ := boot.GetByType((*UserService)(nil))
        user := userService.(UserService).GetUser("123")
        fmt.Println("Result:", user)
    }

    // Start the application
    boot.RunApplication()

}
```

### Advanced Features

#### Multiple Implementations with Primary

```go
// Register multiple loggers
boot.Object(&FileLogger{}).
    Export((*Logger)(nil)).
    Name("file-logger")

boot.Object(&ConsoleLogger{}).
    Export((*Logger)(nil)).
    Primary().  // Mark as primary for Logger interface
    Name("console-logger")
```

#### Optional Dependencies

```go
type Service struct {
    Logger    Logger    `autowire:""`           // Required
    Cache     Cache     `autowire:"optional"`   // Optional
    Metrics   Metrics   `autowire:"?"`          // Optional (short form)
}
```

#### Lifecycle Management

```go
type DatabaseService struct{}

func (d *DatabaseService) Init(ctx context.Context) error {
    // Initialization logic
    return nil
}

func (d *DatabaseService) Start(ctx context.Context) error {
    // Startup logic
    return nil
}

func (d *DatabaseService) Stop(ctx context.Context) error {
    // Cleanup logic
    return nil
}

boot.Object(&DatabaseService{}).
    Priority(100).  // Higher priority starts first, stops last
    Name("database")
```

## Documentation

For detailed documentation, examples, and best practices, see the [docs](./docs) directory.

## Acknowledgments

This project is inspired by [go-spring](https://github.com/go-spring/go-spring) and serves as a simplified version focusing on core dependency injection features.

## License

MIT License
