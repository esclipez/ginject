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
- **Isolated Containers**: Use the default container or create independent containers with `NewContainer`

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

type App struct {
    UserService UserService `autowire:""`
}

func (a *App) Start(ctx context.Context) error {
    user := a.UserService.GetUser("123")
    fmt.Println("Result:", user)
    return nil
}

func main() {
    boot.Object(&ConsoleLogger{}).
        Export((*Logger)(nil)).
        Name("console-logger")

    boot.Object(&UserServiceImpl{}).
        Export((*UserService)(nil)).
        Name("user-service")

    boot.Object(&App{})

    boot.RunApplication()
}
```

`RunApplication` starts the default container and then waits for `SIGINT`, `SIGTERM`, or a call to `boot.Shutdown()`.

### Advanced Features

#### Default Container vs Custom Containers

The package-level API uses a default container:

```go
boot.Object(&ConsoleLogger{}).Export((*Logger)(nil))
boot.RunApplication()
```

For tests, libraries, or multiple isolated applications in the same process, create your own container:

```go
container := boot.NewContainer()
container.Object(&ConsoleLogger{}).Export((*Logger)(nil))

if err := container.Run(context.Background()); err != nil {
    panic(err)
}
defer container.Stop(context.Background())
```

Register all objects before `Run` or `Start`. Once a container starts, registration is sealed and later calls to `Object` panic.

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
    AuditLog  Logger    `autowire:"audit,optional"`
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

Lifecycle order:

1. Register pending objects
2. Validate exported types and primary selections
3. Inject dependencies
4. Run `Init` methods from high priority to low priority
5. Run `Start` methods from high priority to low priority
6. Run `Stop` methods from low priority to high priority

#### Runtime Logs

The default `RunApplication` lifecycle logs use a compact `ginject:` prefix:

```text
ginject: starting application
ginject: application started
ginject: shutdown requested
ginject: stopping application
ginject: application stopped
```

When shutdown comes from an OS signal, the shutdown request message is `ginject: shutdown requested by OS signal`.

Use `boot.SetLogger` to replace the default logger.

## Documentation

- [Autowiring Guide](./docs/autowiring_guide.md)
- [Container Lifecycle](./docs/container_lifecycle.md)

## Acknowledgments

This project is inspired by [go-spring](https://github.com/go-spring/go-spring) and serves as a simplified version focusing on core dependency injection features.

## License

MIT License
