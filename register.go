package mind

import (
	"fmt"
	"reflect"

	"github.com/ichiban/prolog/engine"
)

// Register adds a Go function as a Prolog built-in predicate.
//
// Supported function signatures:
//   - func(args ...string) bool          → succeeds/fails based on return value
//   - func(args ...string) string        → unifies result with the last argument
//   - func(args ...string) (string, bool) → returns result + success flag
//
// The arity parameter must match: for bool-returning functions, arity == number
// of Go parameters. For string-returning functions, arity == parameters + 1
// (the extra argument is the output variable).
func (e *Engine) Register(name string, arity int, fn any) error {
	fnVal := reflect.ValueOf(fn)
	fnType := fnVal.Type()

	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("mind: Register %s: expected function, got %T", name, fn)
	}

	atom := engine.NewAtom(name)

	switch arity {
	case 1:
		e.interp.Register1(atom, makeHandler(fnVal, fnType, arity))
	case 2:
		e.interp.Register2(atom, makeHandler2(fnVal, fnType, arity))
	case 3:
		e.interp.Register3(atom, makeHandler3(fnVal, fnType, arity))
	default:
		return fmt.Errorf("mind: Register %s: arity %d not supported (max 3)", name, arity)
	}

	return nil
}

// makeHandler creates a Predicate1 from a Go function.
func makeHandler(fnVal reflect.Value, fnType reflect.Type, arity int) engine.Predicate1 {
	return func(vm *engine.VM, arg1 engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
		args := []engine.Term{arg1}
		return callGo(vm, fnVal, fnType, args, nil, k, env)
	}
}

// makeHandler2 creates a Predicate2 from a Go function.
func makeHandler2(fnVal reflect.Value, fnType reflect.Type, arity int) engine.Predicate2 {
	hasOutput := fnType.NumOut() > 0 && fnType.Out(0).Kind() == reflect.String
	return func(vm *engine.VM, arg1, arg2 engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
		if hasOutput {
			// Last arg is output, only pass first to Go function
			args := []engine.Term{arg1}
			return callGo(vm, fnVal, fnType, args, &arg2, k, env)
		}
		args := []engine.Term{arg1, arg2}
		return callGo(vm, fnVal, fnType, args, nil, k, env)
	}
}

// makeHandler3 creates a Predicate3 from a Go function.
func makeHandler3(fnVal reflect.Value, fnType reflect.Type, arity int) engine.Predicate3 {
	hasOutput := fnType.NumOut() > 0 && fnType.Out(0).Kind() == reflect.String
	return func(vm *engine.VM, arg1, arg2, arg3 engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
		if hasOutput {
			args := []engine.Term{arg1, arg2}
			return callGo(vm, fnVal, fnType, args, &arg3, k, env)
		}
		args := []engine.Term{arg1, arg2, arg3}
		return callGo(vm, fnVal, fnType, args, nil, k, env)
	}
}

// callGo invokes the Go function with resolved Prolog arguments.
// If output is non-nil, the function's string return value is unified with it.
func callGo(vm *engine.VM, fnVal reflect.Value, fnType reflect.Type, args []engine.Term, output *engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
	goArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		resolved := env.Resolve(arg)
		goArg, ok := termToGo(resolved, fnType.In(i))
		if !ok {
			return engine.Bool(false)
		}
		goArgs[i] = goArg
	}

	results := fnVal.Call(goArgs)

	switch fnType.NumOut() {
	case 0:
		// No return value — always succeeds
		return k(env)

	case 1:
		ret := results[0]
		switch ret.Kind() {
		case reflect.Bool:
			if ret.Bool() {
				return k(env)
			}
			return engine.Bool(false)
		case reflect.String:
			if output == nil {
				return engine.Bool(false)
			}
			atom := engine.NewAtom(ret.String())
			return engine.Unify(vm, *output, atom, k, env)
		default:
			return engine.Bool(false)
		}

	case 2:
		// (string, bool) — return value + success flag
		val := results[0]
		ok := results[1]
		if !ok.Bool() {
			return engine.Bool(false)
		}
		if output != nil && val.Kind() == reflect.String {
			atom := engine.NewAtom(val.String())
			return engine.Unify(vm, *output, atom, k, env)
		}
		return k(env)

	default:
		return engine.Bool(false)
	}
}

// termToGo converts a Prolog term to a Go reflect.Value of the target type.
func termToGo(term engine.Term, targetType reflect.Type) (reflect.Value, bool) {
	switch targetType.Kind() {
	case reflect.String:
		switch v := term.(type) {
		case engine.Atom:
			return reflect.ValueOf(v.String()), true
		default:
			return reflect.Value{}, false
		}
	case reflect.Int, reflect.Int64:
		switch v := term.(type) {
		case engine.Integer:
			return reflect.ValueOf(int64(v)).Convert(targetType), true
		default:
			return reflect.Value{}, false
		}
	case reflect.Float64:
		switch v := term.(type) {
		case engine.Float:
			return reflect.ValueOf(float64(v)), true
		default:
			return reflect.Value{}, false
		}
	case reflect.Bool:
		switch v := term.(type) {
		case engine.Atom:
			return reflect.ValueOf(v.String() == "true"), true
		default:
			return reflect.Value{}, false
		}
	default:
		return reflect.Value{}, false
	}
}
