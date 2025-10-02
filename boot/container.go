package boot

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"sync"
)

// ComponentInfo holds metadata about a registered component
type ComponentInfo struct {
	Instance      interface{}
	InstanceType  reflect.Type
	Name          string
	Priority      int
	ExportedTypes []reflect.Type
	IsPrimary     bool
}

// Container manages the IoC lifecycle
type Container struct {
	componentsByName map[string]*ComponentInfo
	componentsByType map[reflect.Type]*ComponentInfo
	components       []*ComponentInfo
	mu               sync.RWMutex
	started          bool
}

// NewContainer creates a new IoC container
func NewContainer() *Container {
	return &Container{
		componentsByName: make(map[string]*ComponentInfo),
		componentsByType: make(map[reflect.Type]*ComponentInfo),
		components:       make([]*ComponentInfo, 0),
	}
}

// Object starts the fluent API for component registration
func (c *Container) Object(instance interface{}) *ObjectBuilder {
	return newObjectBuilder(c, instance)
}

// registerComponent adds a component to the container
func (c *Container) registerComponent(info *ComponentInfo) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check name uniqueness
	if existing, exists := c.componentsByName[info.Name]; exists {
		return fmt.Errorf("component with name '%s' already registered: existing type %s, new type %s",
			info.Name, existing.InstanceType, info.InstanceType)
	}

	// Register by name
	c.componentsByName[info.Name] = info
	c.components = append(c.components, info)

	return nil
}

// GetByName retrieves a component by name
func (c *Container) GetByName(name string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	info, exists := c.componentsByName[name]
	if !exists {
		return nil, fmt.Errorf("component '%s' not found", name)
	}
	return info.Instance, nil
}

// GetByType retrieves a component by type
func (c *Container) GetByType(componentType reflect.Type) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	info, exists := c.componentsByType[componentType]
	if !exists {
		return nil, fmt.Errorf("no component of type '%s' found", componentType)
	}
	return info.Instance, nil
}

// InjectDependencies performs dependency injection on all components
func (c *Container) InjectDependencies() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, info := range c.components {
		if err := c.injectComponentUnsafe(info.Instance); err != nil {
			return fmt.Errorf("failed to inject dependencies for '%s': %w", info.Name, err)
		}
	}
	return nil
}

// injectComponentUnsafe performs injection without locking (assumes caller holds lock)
func (c *Container) injectComponentUnsafe(component interface{}) error {
	v := reflect.ValueOf(component)
	if v.Kind() != reflect.Ptr {
		return nil
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Check if autowire tag exists (including empty values)
		if tag, exists := fieldType.Tag.Lookup("autowire"); exists {
			if !field.CanSet() {
				continue
			}

			// Parse for optional syntax
			isOptional := tag == "optional" || tag == "?" ||
				(len(tag) > 9 && tag[len(tag)-9:] == ",optional")

			dependency, err := c.resolveDependencyUnsafe(field.Type(), tag)
			if err != nil {
				if isOptional {
					// Optional dependency - skip if not found, no error
					continue
				}
				// Required dependency - fail if not found
				return fmt.Errorf("failed to autowire required field %s: %w", fieldType.Name, err)
			}

			if dependency != nil {
				field.Set(reflect.ValueOf(dependency))
			}
		}
	}
	return nil
}

// resolveDependencyUnsafe resolves dependency without locking (assumes caller holds lock)
func (c *Container) resolveDependencyUnsafe(fieldType reflect.Type, qualifier string) (interface{}, error) {
	// Parse qualifier for optional syntax: "ComponentName,optional"
	componentName := qualifier
	isOptional := false

	if qualifier == "optional" || qualifier == "?" {
		// Pure optional - resolve by type
		dependency, err := c.getByTypeUnsafe(fieldType)
		if err != nil {
			return nil, nil // Return nil without error for optional
		}
		return dependency, nil
	}

	// Check for "ComponentName,optional" syntax
	if len(qualifier) > 9 && qualifier[len(qualifier)-9:] == ",optional" {
		componentName = qualifier[:len(qualifier)-9]
		isOptional = true
	}

	switch componentName {
	case "required", "":
		// Default is required - resolve by type
		return c.getByTypeUnsafe(fieldType)
	default:
		// Specific component name
		component, err := c.getByNameUnsafe(componentName)
		if err != nil {
			if isOptional {
				return nil, nil // Return nil without error for optional named component
			}
			return nil, err
		}

		// Type check: verify component can be assigned to target type
		componentValue := reflect.ValueOf(component)
		if !componentValue.Type().AssignableTo(fieldType) {
			if isOptional {
				return nil, nil // Return nil without error for optional incompatible type
			}
			return nil, fmt.Errorf("component '%s' (type %s) is not assignable to field type %s",
				componentName, componentValue.Type(), fieldType)
		}

		return component, nil
	}
}

// getByNameUnsafe retrieves a component by name without locking
func (c *Container) getByNameUnsafe(name string) (interface{}, error) {
	info, exists := c.componentsByName[name]
	if !exists {
		return nil, fmt.Errorf("component '%s' not found", name)
	}
	return info.Instance, nil
}

// getByTypeUnsafe retrieves a component by type without locking
func (c *Container) getByTypeUnsafe(componentType reflect.Type) (interface{}, error) {
	info, exists := c.componentsByType[componentType]
	if !exists {
		return nil, fmt.Errorf("no component of type '%s' found", componentType)
	}
	return info.Instance, nil
}

// Initialize runs init phase in descending priority order (higher priority first)
func (c *Container) Initialize(ctx context.Context) error {
	components := c.getSortedComponents(false) // descending order

	for _, info := range components {
		if initializable, ok := info.Instance.(Initializable); ok {
			if err := initializable.Init(ctx); err != nil {
				return fmt.Errorf("initialization failed for '%s': %w", info.Name, err)
			}
		}
	}
	return nil
}

// Start runs startup phase in descending priority order (higher priority first)
func (c *Container) Start(ctx context.Context) error {
	if c.started {
		return fmt.Errorf("container already started")
	}

	components := c.getSortedComponents(false) // descending order

	for _, info := range components {
		if startable, ok := info.Instance.(Startable); ok {
			if err := startable.Start(ctx); err != nil {
				return fmt.Errorf("startup failed for '%s': %w", info.Name, err)
			}
		}
	}

	c.started = true
	return nil
}

// Stop runs shutdown phase in ascending priority order (lower priority first)
func (c *Container) Stop(ctx context.Context) error {
	if !c.started {
		return nil
	}

	components := c.getSortedComponents(true) // ascending order

	var lastErr error
	for _, info := range components {
		if stoppable, ok := info.Instance.(Stoppable); ok {
			if err := stoppable.Stop(ctx); err != nil {
				lastErr = fmt.Errorf("shutdown failed for '%s': %w", info.Name, err)
			}
		}
	}

	c.started = false
	return lastErr
}

// getSortedComponents returns components sorted by priority
func (c *Container) getSortedComponents(ascending bool) []*ComponentInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	components := make([]*ComponentInfo, len(c.components))
	copy(components, c.components)

	sort.Slice(components, func(i, j int) bool {
		if ascending {
			return components[i].Priority < components[j].Priority
		}
		return components[i].Priority > components[j].Priority
	})

	return components
}

// validateTypeRegistrations validates type mappings and resolves conflicts
func (c *Container) validateTypeRegistrations() error {
	// Clear existing type mappings
	c.componentsByType = make(map[reflect.Type]*ComponentInfo)

	// Group components by exported type
	typeGroups := make(map[reflect.Type][]*ComponentInfo)

	for _, info := range c.components {
		for _, exportedType := range info.ExportedTypes {
			typeGroups[exportedType] = append(typeGroups[exportedType], info)
		}
	}

	// Validate each type group
	for exportedType, components := range typeGroups {
		if len(components) == 1 {
			// Single component - always use it
			c.componentsByType[exportedType] = components[0]
			continue
		}

		// Multiple components - find primary
		var primaryComponents []*ComponentInfo
		for _, comp := range components {
			if comp.IsPrimary {
				primaryComponents = append(primaryComponents, comp)
			}
		}

		if len(primaryComponents) == 0 {
			// No primary - ambiguous
			names := make([]string, len(components))
			for i, comp := range components {
				names[i] = comp.Name
			}
			return fmt.Errorf("ambiguous components for type '%s': %v (mark one as Primary())",
				exportedType, names)
		}

		if len(primaryComponents) > 1 {
			// Multiple primaries - conflict
			names := make([]string, len(primaryComponents))
			for i, comp := range primaryComponents {
				names[i] = comp.Name
			}
			return fmt.Errorf("multiple primary components for type '%s': %v",
				exportedType, names)
		}

		// Exactly one primary - use it
		c.componentsByType[exportedType] = primaryComponents[0]
	}

	return nil
}

// Run executes the complete lifecycle: register pending → validate → inject → init → start
func (c *Container) Run(ctx context.Context) error {
	// First register all pending builders
	if err := c.registerPendingBuilders(); err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	// Then validate all type registrations
	if err := c.validateTypeRegistrations(); err != nil {
		return fmt.Errorf("type validation failed: %w", err)
	}

	if err := c.InjectDependencies(); err != nil {
		return fmt.Errorf("dependency injection failed: %w", err)
	}

	if err := c.Initialize(ctx); err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	if err := c.Start(ctx); err != nil {
		return fmt.Errorf("startup failed: %w", err)
	}

	return nil
}

// registerPendingBuilders registers all pending ObjectBuilders
func (c *Container) registerPendingBuilders() error {
	for _, builder := range pendingBuilders {
		if err := builder.register(); err != nil {
			return err
		}
	}
	// Clear pending builders after registration
	pendingBuilders = nil
	return nil
}
