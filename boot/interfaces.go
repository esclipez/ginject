package boot

import "context"

// Initializable Lifecycle interfaces for components
type Initializable interface {
	Init(ctx context.Context) error
}

type Startable interface {
	Start(ctx context.Context) error
}

type Stoppable interface {
	Stop(ctx context.Context) error
}

// Component represents a managed component with priority
type Component interface {
	Priority() int
	Name() string
}
