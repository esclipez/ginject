package boot

import (
	"reflect"
	"strings"
	"testing"
)

type exportValidationFoo struct{}

type exportValidationBar interface {
	Bar()
}

type exportValidationLogger interface {
	Log(string)
}

type exportValidationConsoleLogger struct{}

func (l *exportValidationConsoleLogger) Log(string) {}

func TestObjectBuilderRegisterRejectsUnassignableExportType(t *testing.T) {
	container := NewContainer()

	err := container.Object(&exportValidationFoo{}).Export((*exportValidationBar)(nil)).register()

	if err == nil {
		t.Fatal("expected invalid export to fail registration")
	}
	if !strings.Contains(err.Error(), "cannot be exported as") {
		t.Fatalf("expected assignability error, got %q", err.Error())
	}
}

func TestObjectBuilderRegisterAcceptsAssignableInterfaceExportType(t *testing.T) {
	container := NewContainer()

	if err := container.Object(&exportValidationConsoleLogger{}).Export((*exportValidationLogger)(nil)).register(); err != nil {
		t.Fatalf("expected valid export to register, got %v", err)
	}
	if err := container.validateTypeRegistrations(); err != nil {
		t.Fatalf("expected valid export type to validate, got %v", err)
	}

	component, err := container.GetByType(reflect.TypeOf((*exportValidationLogger)(nil)).Elem())
	if err != nil {
		t.Fatalf("expected exported component to resolve by interface, got %v", err)
	}
	if _, ok := component.(exportValidationLogger); !ok {
		t.Fatalf("expected resolved component to implement exportValidationLogger, got %T", component)
	}
}
