// Package mind provides a domain-agnostic logical reasoning engine
// embedding Prolog for structured inference over facts and rules.
package mind

import (
	"github.com/ichiban/prolog"
)

// Engine wraps a Prolog interpreter for logical reasoning.
type Engine struct {
	interp *prolog.Interpreter
}

// NewEngine creates a reasoning engine.
func NewEngine(opts ...Option) *Engine {
	interp := new(prolog.Interpreter)
	e := &Engine{interp: interp}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Option configures an Engine.
type Option func(*Engine)
