package boot

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ComponentConfig holds configuration for component registration
type ComponentConfig struct {
	Name     string
	Priority int
}

// ParseComponentTag extracts component configuration from struct tags
func ParseComponentTag(tag string) ComponentConfig {
	config := ComponentConfig{
		Priority: 100, // Default priority
	}

	if tag == "" {
		return config
	}

	parts := strings.Split(tag, ",")
	for _, part := range parts {
		kv := strings.Split(strings.TrimSpace(part), "=")
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "name":
			config.Name = value
		case "priority":
			if p, err := strconv.Atoi(value); err == nil {
				config.Priority = p
			}
		}
	}

	return config
}

// AutoRegister scans a struct for component annotations and registers it
func (c *Container) AutoRegister(instance interface{}) error {
	v := reflect.ValueOf(instance)
	t := reflect.TypeOf(instance)

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("component must be a pointer")
	}

	// Check if it implements Component interface
	component, ok := instance.(Component)
	if !ok {
		return fmt.Errorf("instance must implement Component interface")
	}

	// Look for component tag on struct fields or use reflection to find component metadata
	structType := t.Elem()
	if structType.Kind() == reflect.Struct {
		// Option 1: Look for a special field with component tag
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			if tag := field.Tag.Get("component"); tag != "" {
				config := ParseComponentTag(tag)
				if config.Name != "" {
					// Override component name if specified in tag
					if named, ok := instance.(interface{ SetName(string) }); ok {
						named.SetName(config.Name)
					}
				}
				break
			}
		}

		// Option 2: Or look for a metadata field
		if metaField, found := structType.FieldByName("ComponentMeta"); found {
			if tag := metaField.Tag.Get("component"); tag != "" {
				config := ParseComponentTag(tag)
				if config.Name != "" {
					if named, ok := instance.(interface{ SetName(string) }); ok {
						named.SetName(config.Name)
					}
				}
			}
		}
	}

	// Create ComponentInfo for registration
	info := &ComponentInfo{
		Instance: component,
		Name:     component.Name(),
		Priority: component.Priority(),
	}

	return c.registerComponent(info)
}
