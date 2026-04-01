package mind

import (
	"os"
	"testing"
)

func TestRegisterBoolPredicate(t *testing.T) {
	e := NewEngine(WithPrelude(`
		item(apple).
		item(banana).
		item(cherry).
	`))

	// Register a Go function that filters items
	err := e.Register("starts_with_b", 1, func(s string) bool {
		return len(s) > 0 && s[0] == 'b'
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	results, err := e.Query(`item(X), starts_with_b(X).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Bindings["X"] != "banana" {
		t.Errorf("expected X=banana, got X=%v", results[0].Bindings["X"])
	}
}

func TestRegisterFileExists(t *testing.T) {
	// Create a temp file to check
	f, err := os.CreateTemp("", "mind-test-*")
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	f.Close()
	defer os.Remove(path)

	e := NewEngine()
	err = e.Register("file_exists", 1, func(path string) bool {
		_, err := os.Stat(path)
		return err == nil
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Assert the path as a fact
	if err := e.Assert("check_path", path); err != nil {
		t.Fatalf("Assert failed: %v", err)
	}

	result, ok, err := e.QueryOne(`check_path(P), file_exists(P).`)
	if err != nil {
		t.Fatalf("QueryOne failed: %v", err)
	}
	if !ok {
		t.Fatal("expected a solution for existing file")
	}
	if result.Bindings["P"] != path {
		t.Errorf("expected P=%s, got P=%v", path, result.Bindings["P"])
	}

	// Non-existent file should return no solution
	if err := e.Assert("check_path", "/nonexistent/file/path"); err != nil {
		t.Fatalf("Assert failed: %v", err)
	}
	results, err := e.Query(`check_path(P), file_exists(P).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	// Should only find the real file
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestRegisterArity4BoolPredicate(t *testing.T) {
	e := NewEngine()

	// Register a 4-arity predicate: all_different(A, B, C, D)
	err := e.Register("all_different", 4, func(a, b, c, d string) bool {
		return a != b && a != c && a != d && b != c && b != d && c != d
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	result, ok, err := e.QueryOne(`all_different(a, b, c, d).`)
	if err != nil {
		t.Fatalf("QueryOne failed: %v", err)
	}
	if !ok {
		t.Fatal("expected success for all different atoms")
	}
	_ = result

	_, ok, err = e.QueryOne(`all_different(a, b, a, d).`)
	if err != nil {
		t.Fatalf("QueryOne failed: %v", err)
	}
	if ok {
		t.Fatal("expected failure for duplicate atoms")
	}
}

func TestRegisterArity5StringPredicate(t *testing.T) {
	e := NewEngine()

	// Register a 5-arity predicate: concat4(A, B, C, D, Result)
	// Go function takes 4 string inputs, returns string (arity = 4 inputs + 1 output = 5)
	err := e.Register("concat4", 5, func(a, b, c, d string) string {
		return a + b + c + d
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	result, ok, err := e.QueryOne(`concat4(h, e, l, p, X).`)
	if err != nil {
		t.Fatalf("QueryOne failed: %v", err)
	}
	if !ok {
		t.Fatal("expected a solution")
	}
	if result.Bindings["X"] != "help" {
		t.Errorf("expected X=help, got X=%v", result.Bindings["X"])
	}
}

func TestRegisterUnsupportedArity(t *testing.T) {
	e := NewEngine()

	err := e.Register("too_many", 9, func(a string) bool { return true })
	if err == nil {
		t.Fatal("expected error for arity 9")
	}
	if err.Error() != "mind: Register too_many: arity 9 not supported (max 8)" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRegisterStringPredicate(t *testing.T) {
	e := NewEngine()

	// Register a predicate that transforms: upper(Input, Output)
	err := e.Register("reverse_atom", 2, func(input string) string {
		runes := []rune(input)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	result, ok, err := e.QueryOne(`reverse_atom(hello, X).`)
	if err != nil {
		t.Fatalf("QueryOne failed: %v", err)
	}
	if !ok {
		t.Fatal("expected a solution")
	}
	if result.Bindings["X"] != "olleh" {
		t.Errorf("expected X=olleh, got X=%v", result.Bindings["X"])
	}
}
