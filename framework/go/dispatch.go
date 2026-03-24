package dispatch

import (
	"context"
	"fmt"
	"sync"
)

// RuneFunc is the signature for any registered callable.
// Input and output are JSON-encoded bytes for isolation.
type RuneFunc func(ctx context.Context, input []byte) ([]byte, error)

// Middleware wraps a callable invocation.
type Middleware func(name string, next RuneFunc) RuneFunc

// Dispatcher is the callable registry and call router.
type Dispatcher struct {
	mu         sync.RWMutex
	runes      map[string]RuneFunc
	middleware []Middleware
}

// New creates an empty dispatcher.
func New() *Dispatcher {
	return &Dispatcher{runes: make(map[string]RuneFunc)}
}

// Register adds a callable to the registry.
func (d *Dispatcher) Register(name string, fn RuneFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.runes[name] = fn
}

// Use appends a middleware to the chain.
// Middleware is applied in the order added, outermost first.
func (d *Dispatcher) Use(mw Middleware) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.middleware = append(d.middleware, mw)
}

// Call invokes a callable by name through the middleware chain.
func (d *Dispatcher) Call(ctx context.Context, name string, input []byte) ([]byte, error) {
	d.mu.RLock()
	fn, ok := d.runes[name]
	mw := make([]Middleware, len(d.middleware))
	copy(mw, d.middleware)
	d.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("callable %q not registered", name)
	}

	for i := len(mw) - 1; i >= 0; i-- {
		fn = mw[i](name, fn)
	}

	return fn(ctx, input)
}
