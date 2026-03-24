package dispatch

import (
	"context"
	"fmt"
)

// RuneFunc is the signature for any registered callable.
// Input and output are JSON-encoded bytes for isolation.
type RuneFunc func(ctx context.Context, input []byte) ([]byte, error)

// Middleware wraps a callable invocation.
type Middleware func(name string, next RuneFunc) RuneFunc

// Dispatcher is an immutable callable registry and call router.
type Dispatcher struct {
	runes      map[string]RuneFunc
	middleware []Middleware
}

// New creates an immutable Dispatcher with the given runes and middleware.
func New(runes map[string]RuneFunc, middleware []Middleware) *Dispatcher {
	return &Dispatcher{runes: runes, middleware: middleware}
}

// Call invokes a callable by name through the middleware chain.
func (d *Dispatcher) Call(ctx context.Context, name string, input []byte) ([]byte, error) {
	fn, ok := d.runes[name]
	if !ok {
		return nil, fmt.Errorf("callable %q not registered", name)
	}

	for i := len(d.middleware) - 1; i >= 0; i-- {
		fn = d.middleware[i](name, fn)
	}

	return fn(ctx, input)
}
