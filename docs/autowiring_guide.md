# Autowiring Guide

## Autowire Tag Syntax

Use the `autowire` tag on struct fields to enable dependency injection:

### Basic Patterns

```go
type MyService struct {
    // @Autowired - required by default, inject by type
    Database DatabaseService `autowire:""`
    
    // @Autowired @Required - explicitly required
    Cache CacheService `autowire:"required"`
    
    // @Autowired @Qualifier("SpecificName") - inject specific component by name
    Logger Logger `autowire:"FileLogger"`
    
    // @Autowired(required=false) - optional dependency
    OptionalService OptionalService `autowire:"optional"`
    
    // Alternative optional syntax
    MaybeService MaybeService `autowire:"?"`
}
```

### Autowire Tag Values

| Tag Value | Behavior | Spring Equivalent |
|-----------|----------|-------------------|
| `""` | Required, inject by type | `@Autowired` |
| `"required"` | Explicitly required, inject by type | `@Autowired @Required` |
| `"optional"` | Optional, inject by type if available | `@Autowired(required=false)` |
| `"?"` | Optional (alternative syntax) | `@Autowired(required=false)` |
| `"ComponentName"` | Required, inject specific component by name | `@Autowired @Qualifier("ComponentName")` |
| `"ComponentName,optional"` | Optional, inject specific component by name if available | `@Autowired(required=false) @Qualifier("ComponentName")` |

### Error Handling

- **Required dependencies**: Application fails to start if not found
- **Optional dependencies**: Field remains nil if not found, no error
- **Type mismatch**: Error if qualified component doesn't match field type
- **Invalid exports**: Registration fails if `Export` names a type the component cannot be assigned to
- **Ambiguous exports**: If multiple components export the same type, mark exactly one with `Primary`

### Exported Types

Each component is exported as its concrete type by default. Use `Export` to make it available through an interface:

```go
type Logger interface {
    Log(message string)
}

type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(message string) {}

boot.Object(&ConsoleLogger{}).Export((*Logger)(nil))
```

`Export` validates assignability during registration. This fails during `Run`:

```go
type Metrics interface {
    Count(name string)
}

boot.Object(&ConsoleLogger{}).Export((*Metrics)(nil))
```

### Examples

```go
// Multiple logger implementations
boot.Object(NewFileLogger()).Export((*Logger)(nil)).Name("FileLog")
boot.Object(NewConsoleLogger()).Export((*Logger)(nil)).Name("ConsoleLog").Primary()

type Service struct {
    // Injects primary Logger
    DefaultLogger Logger `autowire:""`
    
    // Injects specific FileLogger
    FileLogger Logger `autowire:"FileLog"`
    
    // Won't fail if "DebugLog" doesn't exist
    DebugLogger Logger `autowire:"DebugLog,optional"`
}
```

### Nested Structs

Ginject also scans exported nested structs and non-nil pointers for `autowire` fields:

```go
type Dependencies struct {
    Logger Logger `autowire:""`
}

type Service struct {
    Dependencies
}
```

Nil pointer fields are skipped unless the pointer field itself has an `autowire` tag.
