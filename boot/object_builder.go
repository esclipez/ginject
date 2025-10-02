package boot

import (
	"reflect"
)

// ObjectBuilder provides fluent API for component configuration
type ObjectBuilder struct {
	container     *Container
	instance      interface{}
	name          string
	priority      int
	exportedTypes []reflect.Type
	nameSet       bool
	prioritySet   bool
	isPrimary     bool // 新增：标记为主要实现
}

// newObjectBuilder creates a new ObjectBuilder
func newObjectBuilder(container *Container, instance interface{}) *ObjectBuilder {
	instanceType := reflect.TypeOf(instance)

	builder := &ObjectBuilder{
		container:     container,
		instance:      instance,
		exportedTypes: []reflect.Type{instanceType}, // 默认导出自身类型
		priority:      0,
		isPrimary:     false,
	}

	// Check if instance implements Named interface
	if named, ok := instance.(Named); ok {
		builder.name = named.Name()
	}

	return builder
}

// Name sets the component name (must be unique)
func (b *ObjectBuilder) Name(name string) *ObjectBuilder {
	b.name = name
	b.nameSet = true
	return b
}

// Priority sets the execution priority (higher values = higher priority)
func (b *ObjectBuilder) Priority(priority int) *ObjectBuilder {
	b.priority = priority
	b.prioritySet = true
	return b
}

// Export adds a type that this component should be registered for
func (b *ObjectBuilder) Export(typePtr interface{}) *ObjectBuilder {
	t := reflect.TypeOf(typePtr)
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // Get the interface type
	}
	b.exportedTypes = append(b.exportedTypes, t)
	return b
}

// Primary marks this component as the primary implementation for its exported types
func (b *ObjectBuilder) Primary() *ObjectBuilder {
	b.isPrimary = true
	return b
}

// Register completes the component registration
func (b *ObjectBuilder) register() error {
	// Use type name as default if no name is set
	if b.name == "" {
		instanceType := reflect.TypeOf(b.instance)
		if instanceType.Kind() == reflect.Ptr {
			instanceType = instanceType.Elem()
		}
		b.name = instanceType.String()
	}

	info := &ComponentInfo{
		Instance:      b.instance,
		InstanceType:  reflect.TypeOf(b.instance),
		Name:          b.name,
		Priority:      b.priority,
		ExportedTypes: b.exportedTypes,
		IsPrimary:     b.isPrimary,
	}

	return b.container.registerComponent(info)
}
