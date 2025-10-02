package main

import (
	"context"
	"fmt"

	"github.com/esclipez/ginject/boot"
	"github.com/esclipez/ginject/examples/interfaces"
	"github.com/esclipez/ginject/examples/services"
)

// DatabaseService example component
type DatabaseService struct {
	connected bool
}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

func (d *DatabaseService) Init(ctx context.Context) error {
	fmt.Println("[Database] Initializing database connection...")
	return nil
}

func (d *DatabaseService) Start(ctx context.Context) error {
	fmt.Println("[Database] Starting database service...")
	d.connected = true
	return nil
}

func (d *DatabaseService) Stop(ctx context.Context) error {
	fmt.Println("[Database] Stopping database service...")
	d.connected = false
	return nil
}

// RedisService as alternative cache implementation
type RedisService struct {
	connected bool
}

func NewRedisService() *RedisService {
	return &RedisService{}
}

func (r *RedisService) Init(ctx context.Context) error {
	fmt.Println("[Redis] Initializing Redis connection...")
	return nil
}

func (r *RedisService) Start(ctx context.Context) error {
	fmt.Println("[Redis] Starting Redis service...")
	r.connected = true
	return nil
}

func (r *RedisService) Stop(ctx context.Context) error {
	fmt.Println("[Redis] Stopping Redis service...")
	r.connected = false
	return nil
}

// WebService with different autowiring scenarios
type WebService struct {
	// @Autowired - required by default
	UserService interfaces.UserService `autowire:""`

	// @Autowired @Required - explicitly required
	Database *DatabaseService `autowire:"required"`

	// @Autowired @Qualifier("RedisCache") - specific component by name
	Cache *RedisService `autowire:"RedisCache"`

	// @Autowired(required=false) - optional dependency
	OptionalService *OptionalService `autowire:"optional"`

	// Alternative syntax for optional
	AnotherOptional *AnotherService `autowire:"?"`

	running bool
}

func NewWebService() *WebService {
	return &WebService{}
}

func (w *WebService) Init(ctx context.Context) error {
	fmt.Println("[WebService] Initializing web service...")
	return nil
}

func (w *WebService) Start(ctx context.Context) error {
	fmt.Printf("[WebService] Starting web server...\n")
	fmt.Printf("  - UserService: %v\n", w.UserService != nil)
	fmt.Printf("  - Database connected: %v\n", w.Database.connected)
	fmt.Printf("  - Cache connected: %v\n", w.Cache.connected)
	fmt.Printf("  - OptionalService: %v\n", w.OptionalService != nil)
	fmt.Printf("  - AnotherOptional: %v\n", w.AnotherOptional != nil)

	w.running = true

	// Test the injected UserService
	user, err := w.UserService.GetUser("123")
	if err != nil {
		return err
	}
	fmt.Printf("[WebService] Retrieved user: %+v\n", user)

	return nil
}

func (w *WebService) Stop(ctx context.Context) error {
	fmt.Println("[WebService] Stopping web server...")
	w.running = false
	return nil
}

// OptionalService that may or may not be registered
type OptionalService struct{}

func NewOptionalService() *OptionalService {
	return &OptionalService{}
}

func (o *OptionalService) DoSomething() {
	fmt.Println("[OptionalService] Doing something...")
}

// AnotherService for testing optional injection
type AnotherService struct{}

func NewAnotherService() *AnotherService {
	return &AnotherService{}
}

// CacheService with higher priority
type CacheService struct{}

func NewCacheService() *CacheService {
	return &CacheService{}
}

func (c *CacheService) Init(ctx context.Context) error {
	fmt.Println("[Cache] Initializing cache...")
	return nil
}

func (c *CacheService) Start(ctx context.Context) error {
	fmt.Println("[Cache] Starting cache service...")
	return nil
}

func (c *CacheService) Stop(ctx context.Context) error {
	fmt.Println("[Cache] Stopping cache service...")
	return nil
}

// Register components using fluent API
func init() {
	// Register UserService with interface export
	boot.Object(services.NewUserService()).
		Export((*interfaces.UserService)(nil)).
		Priority(100).
		Name("UserSrv")

	// Register Database service
	boot.Object(NewDatabaseService()).Priority(50)

	// Register Redis with specific name for qualifier injection
	boot.Object(NewRedisService()).
		Priority(80).
		Name("RedisCache")

	// Register cache service with high priority
	boot.Object(NewCacheService()).Priority(200)

	// Register web service
	boot.Object(NewWebService())

	// Register optional service (comment out to test optional injection)
	boot.Object(NewOptionalService()).Name("OptionalSrv")

	// Don't register AnotherService to test optional injection with "?"
	// boot.Object(NewAnotherService()).Name("AnotherSrv")
}
