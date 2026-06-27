# Container Lifecycle

Ginject can be used through the package-level default container or through explicitly created containers.

## Default Container

The package-level helpers register components on a shared default container:

```go
boot.Object(&ConsoleLogger{}).Export((*Logger)(nil))
boot.Object(&App{})

boot.RunApplication()
```

`RunApplication` runs the default container, then waits for `SIGINT`, `SIGTERM`, or `boot.Shutdown()`.

## Custom Containers

Use `NewContainer` when you need isolation, especially in tests or libraries:

```go
container := boot.NewContainer()
container.Object(&ConsoleLogger{}).Export((*Logger)(nil))
container.Object(&App{})

if err := container.Run(context.Background()); err != nil {
    panic(err)
}
defer container.Stop(context.Background())
```

A custom container does not share registrations with the default container.

## Registration Window

Register all objects before `Run` or `Start`:

```go
container.Object(&Database{})
container.Object(&App{})
err := container.Run(ctx)
```

Once a container starts, its registration set is sealed. Calling `Object` after `Run` or `Start` panics because late components would not have participated in validation, dependency injection, initialization, or startup.

## Run Order

`Run` executes these steps:

1. Register pending objects
2. Validate exported types and primary selections
3. Inject fields tagged with `autowire`
4. Call `Init(ctx)` on `Initializable` components from high priority to low priority
5. Call `Start(ctx)` on `Startable` components from high priority to low priority

`Stop` calls `Stop(ctx)` on `Stoppable` components from low priority to high priority.

Higher priority components start earlier and stop later.

## Lifecycle Interfaces

```go
type Initializable interface {
    Init(ctx context.Context) error
}

type Startable interface {
    Start(ctx context.Context) error
}

type Stoppable interface {
    Stop(ctx context.Context) error
}
```

Returning an error from `Init` or `Start` aborts application startup. `Stop` returns the last shutdown error; `RunApplication` logs it.

## Runtime Logs

`RunApplication` logs compact lifecycle messages through the configured logger:

```text
ginject: starting application
ginject: application started
ginject: shutdown requested
ginject: stopping application
ginject: application stopped
```

When shutdown comes from an OS signal, the shutdown request message is `ginject: shutdown requested by OS signal`.

Use `boot.SetLogger` to replace the default logger.
