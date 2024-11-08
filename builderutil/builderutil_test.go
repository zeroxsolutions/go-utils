package builderutil_test

import (
	"errors"
	"testing"

	"github.com/zeroxsolutions/go-utils/builderutil"
)

// MockLister is a mock implementation of the Lister interface for testing purposes.
type MockLister[T any] struct {
	Funcs []func(*T) error
}

// List returns the list of functions that MockLister holds for testing.
func (m *MockLister[T]) List() []func(*T) error {
	return m.Funcs
}

// TestBuild_Success tests if Build applies configurations successfully and returns the expected result.
func TestBuild_Success(t *testing.T) {
	type Config struct {
		Value int
	}

	setValue := func(value int) func(*Config) error {
		return func(c *Config) error {
			c.Value = value
			return nil
		}
	}

	// Create a mock Lister with a function to set the value to 42.
	mockLister := &MockLister[Config]{Funcs: []func(*Config) error{setValue(42)}}

	// Build the Config instance
	config, err := builderutil.Build[Config](mockLister)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify that the configuration was applied
	if config.Value != 42 {
		t.Errorf("Expected config.Value to be 42, got %d", config.Value)
	}
}

// TestBuild_NilOption tests if Build handles nil options gracefully.
func TestBuild_NilOption(t *testing.T) {
	type Config struct {
		Value int
	}

	// Pass a nil Lister as an option
	config, err := builderutil.Build[Config](nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify that config is returned with default values
	if config.Value != 0 {
		t.Errorf("Expected config.Value to be 0, got %d", config.Value)
	}
}

// TestBuild_NilFunction tests if Build ignores nil functions in the Lister.
func TestBuild_NilFunction(t *testing.T) {
	type Config struct {
		Value int
	}

	setValue := func(value int) func(*Config) error {
		return func(c *Config) error {
			c.Value = value
			return nil
		}
	}

	// Create a mock Lister with a nil function and a valid function
	mockLister := &MockLister[Config]{Funcs: []func(*Config) error{nil, setValue(42)}}

	// Build the Config instance
	config, err := builderutil.Build[Config](mockLister)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify that the configuration was applied, and nil function was ignored
	if config.Value != 42 {
		t.Errorf("Expected config.Value to be 42, got %d", config.Value)
	}
}

// TestBuild_ErrorInFunction tests if Build stops and returns an error when a function fails.
func TestBuild_ErrorInFunction(t *testing.T) {
	type Config struct {
		Value int
	}

	errFunc := func(*Config) error {
		return errors.New("error in function")
	}

	setValue := func(value int) func(*Config) error {
		return func(c *Config) error {
			c.Value = value
			return nil
		}
	}

	// Create a mock Lister with an error function followed by a valid function
	mockLister := &MockLister[Config]{Funcs: []func(*Config) error{errFunc, setValue(42)}}

	// Build the Config instance
	config, err := builderutil.Build[Config](mockLister)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Verify that the configuration was not applied due to the error
	if config != nil {
		t.Errorf("Expected config to be nil, got %v", config)
	}
}

// TestBuild_MultipleOptions tests if Build correctly applies multiple options in sequence.
func TestBuild_MultipleOptions(t *testing.T) {
	type Config struct {
		Value int
	}

	setValue := func(value int) func(*Config) error {
		return func(c *Config) error {
			c.Value += value
			return nil
		}
	}

	// Create two mock Listers with different functions to modify the Config instance
	mockLister1 := &MockLister[Config]{Funcs: []func(*Config) error{setValue(10)}}
	mockLister2 := &MockLister[Config]{Funcs: []func(*Config) error{setValue(15)}}

	// Build the Config instance with multiple options
	config, err := builderutil.Build[Config](mockLister1, mockLister2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify that both configurations were applied in sequence
	expectedValue := 25
	if config.Value != expectedValue {
		t.Errorf("Expected config.Value to be %d, got %d", expectedValue, config.Value)
	}
}

// TestBuild_EmptyOptions tests if Build returns a default-initialized instance of T when no options are provided.
func TestBuild_EmptyOptions(t *testing.T) {
	type Config struct {
		Value int
	}

	// Build the Config instance with no options
	config, err := builderutil.Build[Config]()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify that config is default-initialized
	if config.Value != 0 {
		t.Errorf("Expected config.Value to be 0, got %d", config.Value)
	}
}
