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

### Error Handling

- **Required dependencies**: Application fails to start if not found
- **Optional dependencies**: Field remains nil if not found, no error
- **Type mismatch**: Error if qualified component doesn't match field type

### Examples

```go
// Multiple logger implementations
boot.Object(NewFileLogger()).Export((*Logger)(nil)).Name("FileLog").Primary()
boot.Object(NewConsoleLogger()).Export((*Logger)(nil)).Name("ConsoleLog")

type Service struct {
    // Injects primary Logger
    DefaultLogger Logger `autowire:""`
    
    // Injects specific FileLogger
    FileLogger Logger `autowire:"FileLog"`
    
    // Won't fail if "DebugLog" doesn't exist
    DebugLogger Logger `autowire:"DebugLog,optional"`
}
```