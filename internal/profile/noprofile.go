//go:build !profile

// Package profile provides tooling to easily do profiling.
package profile

func Start() func() {
	return func() {}
}
