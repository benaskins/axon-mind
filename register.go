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
// Arity 1-8 is supported, matching the ichiban/prolog Predicate types.
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
		e.interp.Register1(atom, makeHandler1(fnVal, fnType))
	case 2:
		e.interp.Register2(atom, makeHandler2(fnVal, fnType))
	case 3:
		e.interp.Register3(atom, makeHandler3(fnVal, fnType))
	case 4:
		e.interp.Register4(atom, makeHandler4(fnVal, fnType))
	case 5:
		e.interp.Register5(atom, makeHandler5(fnVal, fnType))
	case 6:
		e.interp.Register6(atom, makeHandler6(fnVal, fnType))
	case 7:
		e.interp.Register7(atom, makeHandler7(fnVal, fnType))
	case 8:
		e.interp.Register8(atom, makeHandler8(fnVal, fnType))
	default:
		return fmt.Errorf("mind: Register %s: arity %d not supported (max 8)", name, arity)
	}

	return nil
}

// hasStringOutput returns true if the function returns a string as its first output,
// meaning the last Prolog argument is the output variable.
func hasStringOutput(fnType reflect.Type) bool {
	return fnType.NumOut() > 0 && fnType.Out(0).Kind() == reflect.String
}

// splitArgs separates Prolog terms into Go input args and an optional output variable.
// When the Go function returns a string, the last Prolog argument is the output.
func splitArgs(all []engine.Term, hasOutput bool) (inputs []engine.Term, output *engine.Term) {
	if hasOutput {
		last := all[len(all)-1]
		return all[:len(all)-1], &last
	}
	return all, nil
}

func makeHandler1(fnVal reflect.Value, fnType reflect.Type) engine.Predicate1 {
	return func(vm *engine.VM, a1 engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
		args := []engine.Term{a1}
		inputs, output := splitArgs(args, hasStringOutput(fnType))
		return callGo(vm, fnVal, fnType, inputs, output, k, env)
	}
}

func makeHandler2(fnVal reflect.Value, fnType reflect.Type) engine.Predicate2 {
	return func(vm *engine.VM, a1, a2 engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
		args := []engine.Term{a1, a2}
		inputs, output := splitArgs(args, hasStringOutput(fnType))
		return callGo(vm, fnVal, fnType, inputs, output, k, env)
	}
}

func makeHandler3(fnVal reflect.Value, fnType reflect.Type) engine.Predicate3 {
	return func(vm *engine.VM, a1, a2, a3 engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
		args := []engine.Term{a1, a2, a3}
		inputs, output := splitArgs(args, hasStringOutput(fnType))
		return callGo(vm, fnVal, fnType, inputs, output, k, env)
	}
}

func makeHandler4(fnVal reflect.Value, fnType reflect.Type) engine.Predicate4 {
	return func(vm *engine.VM, a1, a2, a3, a4 engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
		args := []engine.Term{a1, a2, a3, a4}
		inputs, output := splitArgs(args, hasStringOutput(fnType))
		return callGo(vm, fnVal, fnType, inputs, output, k, env)
	}
}

func makeHandler5(fnVal reflect.Value, fnType reflect.Type) engine.Predicate5 {
	return func(vm *engine.VM, a1, a2, a3, a4, a5 engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
		args := []engine.Term{a1, a2, a3, a4, a5}
		inputs, output := splitArgs(args, hasStringOutput(fnType))
		return callGo(vm, fnVal, fnType, inputs, output, k, env)
	}
}

func makeHandler6(fnVal reflect.Value, fnType reflect.Type) engine.Predicate6 {
	return func(vm *engine.VM, a1, a2, a3, a4, a5, a6 engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
		args := []engine.Term{a1, a2, a3, a4, a5, a6}
		inputs, output := splitArgs(args, hasStringOutput(fnType))
		return callGo(vm, fnVal, fnType, inputs, output, k, env)
	}
}

func makeHandler7(fnVal reflect.Value, fnType reflect.Type) engine.Predicate7 {
	return func(vm *engine.VM, a1, a2, a3, a4, a5, a6, a7 engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
		args := []engine.Term{a1, a2, a3, a4, a5, a6, a7}
		inputs, output := splitArgs(args, hasStringOutput(fnType))
		return callGo(vm, fnVal, fnType, inputs, output, k, env)
	}
}

func makeHandler8(fnVal reflect.Value, fnType reflect.Type) engine.Predicate8 {
	return func(vm *engine.VM, a1, a2, a3, a4, a5, a6, a7, a8 engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
		args := []engine.Term{a1, a2, a3, a4, a5, a6, a7, a8}
		inputs, output := splitArgs(args, hasStringOutput(fnType))
		return callGo(vm, fnVal, fnType, inputs, output, k, env)
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
