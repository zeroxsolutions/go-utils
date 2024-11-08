// Package builderutil provides utility functions and interfaces for building and
// configuring instances of generic types using a list of functional options.
package builderutil

import "reflect"

// Lister is a generic interface that requires a method to return a list of functions.
// These functions are used to configure or modify an instance of type T.
// The functions in the list are expected to take a pointer to T and return an error.
type Lister[T any] interface {
	// List returns a slice of functions, each of which modifies the instance of T
	// or returns an error if the modification fails.
	List() []func(*T) error
}

// Build constructs and configures an instance of type T using the provided Lister options.
// Each Lister option can return a list of functions that are called in sequence to modify
// the instance of T. If any function returns an error, the Build function stops and returns
// that error. If an option is nil or its List method returns nil, it is skipped.
// Parameters:
// - opts: Variadic arguments of type Lister[T] that provide configuration functions.
//
// Returns:
// - A pointer to the newly constructed instance of T.
// - An error if any configuration function fails.
func Build[T any](opts ...Lister[T]) (*T, error) {

	t := new(T)

	for _, opt := range opts {
		if opt == nil || reflect.ValueOf(opt).IsNil() {
			continue
		}

		for _, setArgs := range opt.List() {

			if setArgs == nil {
				continue
			}

			if err := setArgs(t); err != nil {
				return nil, err
			}

		}

	}

	return t, nil
}
