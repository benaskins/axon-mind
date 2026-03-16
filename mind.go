package mind

import (
	"fmt"
	"os"
	"strings"

	"github.com/ichiban/prolog"
)

// Engine wraps a Prolog interpreter for logical reasoning.
type Engine struct {
	interp   *prolog.Interpreter
	declared map[string]bool // tracks declared dynamic predicates (functor/arity)
	err      error           // deferred error from options
}

// Solution holds variable bindings from a successful query.
type Solution struct {
	Bindings map[string]any
}

// NewEngine creates a reasoning engine.
func NewEngine(opts ...Option) *Engine {
	interp := prolog.New(nil, nil)
	e := &Engine{
		interp:   interp,
		declared: make(map[string]bool),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Option configures an Engine.
type Option func(*Engine)

// WithFile loads a .pl file when the engine is created.
func WithFile(path string) Option {
	return func(e *Engine) {
		if e.err != nil {
			return
		}
		e.err = e.Load(path)
	}
}

// WithPrelude loads inline Prolog source when the engine is created.
func WithPrelude(source string) Option {
	return func(e *Engine) {
		if e.err != nil {
			return
		}
		e.err = e.interp.Exec(source)
	}
}

// Load reads and consults a .pl file.
func (e *Engine) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("mind: load %s: %w", path, err)
	}
	if err := e.interp.Exec(string(data)); err != nil {
		return fmt.Errorf("mind: consult %s: %w", path, err)
	}
	return nil
}

// Assert adds a fact programmatically using assertz (appends to database).
// Automatically declares the predicate as dynamic on first use.
func (e *Engine) Assert(functor string, args ...string) error {
	arity := len(args)
	key := fmt.Sprintf("%s/%d", functor, arity)

	// Declare dynamic on first assertion for this functor/arity
	if !e.declared[key] {
		decl := fmt.Sprintf(":- dynamic(%s/%d).", functor, arity)
		if err := e.interp.Exec(decl); err != nil {
			return fmt.Errorf("mind: declare dynamic %s: %w", key, err)
		}
		e.declared[key] = true
	}

	var term string
	if arity == 0 {
		term = functor
	} else {
		quoted := make([]string, arity)
		for i, a := range args {
			quoted[i] = quoteAtom(a)
		}
		term = fmt.Sprintf("%s(%s)", functor, strings.Join(quoted, ", "))
	}

	sol := e.interp.QuerySolution(fmt.Sprintf("assertz(%s).", term))
	if err := sol.Err(); err != nil {
		return fmt.Errorf("mind: assert %s: %w", term, err)
	}
	return nil
}

// Query runs a goal and returns all solutions.
func (e *Engine) Query(goal string) ([]Solution, error) {
	sols, err := e.interp.Query(goal)
	if err != nil {
		return nil, fmt.Errorf("mind: query: %w", err)
	}
	defer sols.Close()

	var results []Solution
	for sols.Next() {
		bindings, err := scanBindings(sols, goal)
		if err != nil {
			return nil, err
		}
		results = append(results, Solution{Bindings: bindings})
	}
	if err := sols.Err(); err != nil {
		if isExistenceError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("mind: query: %w", err)
	}
	return results, nil
}

// QueryOne runs a goal and returns the first solution.
func (e *Engine) QueryOne(goal string) (*Solution, bool, error) {
	sols, err := e.interp.Query(goal)
	if err != nil {
		return nil, false, fmt.Errorf("mind: query: %w", err)
	}
	defer sols.Close()

	if !sols.Next() {
		if err := sols.Err(); err != nil {
			if isExistenceError(err) {
				return nil, false, nil
			}
			return nil, false, fmt.Errorf("mind: query: %w", err)
		}
		return nil, false, nil
	}

	bindings, err := scanBindings(sols, goal)
	if err != nil {
		return nil, false, err
	}
	return &Solution{Bindings: bindings}, true, nil
}

// isExistenceError returns true if the error is a Prolog existence_error
// for an undefined predicate (which we treat as "no solutions").
func isExistenceError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "existence_error(procedure,")
}

// scanBindings extracts variable bindings from the current solution.
func scanBindings(sols *prolog.Solutions, goal string) (map[string]any, error) {
	vars := extractVars(goal)
	if len(vars) == 0 {
		return map[string]any{}, nil
	}

	m := make(map[string]any)
	if err := sols.Scan(m); err != nil {
		return nil, fmt.Errorf("mind: scan: %w", err)
	}

	return m, nil
}

// extractVars finds uppercase variable names in a Prolog goal string.
func extractVars(goal string) []string {
	seen := map[string]bool{}
	var vars []string
	i := 0
	for i < len(goal) {
		c := goal[i]
		// Skip quoted atoms
		if c == '\'' {
			i++
			for i < len(goal) && goal[i] != '\'' {
				if goal[i] == '\\' {
					i++
				}
				i++
			}
			i++ // closing quote
			continue
		}
		// Variable starts with uppercase letter or underscore
		if c >= 'A' && c <= 'Z' {
			start := i
			for i < len(goal) && isVarChar(goal[i]) {
				i++
			}
			name := goal[start:i]
			if name != "_" && !seen[name] {
				seen[name] = true
				vars = append(vars, name)
			}
			continue
		}
		i++
	}
	return vars
}

func isVarChar(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_'
}

// quoteAtom wraps an atom in single quotes if it contains special characters.
func quoteAtom(s string) string {
	if len(s) > 0 && s[0] >= 'a' && s[0] <= 'z' {
		simple := true
		for _, c := range s {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
				simple = false
				break
			}
		}
		if simple {
			return s
		}
	}
	escaped := strings.ReplaceAll(s, "'", "\\'")
	return "'" + escaped + "'"
}
