package mind

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewEngine(t *testing.T) {
	e := NewEngine()
	if e == nil {
		t.Fatal("NewEngine returned nil")
	}
}

func TestWithPrelude(t *testing.T) {
	e := NewEngine(WithPrelude(`human(socrates). human(plato).`))
	if e == nil {
		t.Fatal("NewEngine returned nil")
	}

	results, err := e.Query(`human(X).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestAssert(t *testing.T) {
	e := NewEngine()

	if err := e.Assert("color", "red"); err != nil {
		t.Fatalf("Assert failed: %v", err)
	}
	if err := e.Assert("color", "blue"); err != nil {
		t.Fatalf("Assert failed: %v", err)
	}

	results, err := e.Query(`color(X).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Bindings["X"] != "red" {
		t.Errorf("expected X=red, got X=%v", results[0].Bindings["X"])
	}
	if results[1].Bindings["X"] != "blue" {
		t.Errorf("expected X=blue, got X=%v", results[1].Bindings["X"])
	}
}

func TestAssertBinary(t *testing.T) {
	e := NewEngine()

	if err := e.Assert("parent", "tom", "bob"); err != nil {
		t.Fatalf("Assert failed: %v", err)
	}

	result, ok, err := e.QueryOne(`parent(tom, X).`)
	if err != nil {
		t.Fatalf("QueryOne failed: %v", err)
	}
	if !ok {
		t.Fatal("expected a solution")
	}
	if result.Bindings["X"] != "bob" {
		t.Errorf("expected X=bob, got X=%v", result.Bindings["X"])
	}
}

func TestQueryNoSolutions(t *testing.T) {
	e := NewEngine()

	results, err := e.Query(`nonexistent(X).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestQueryOneNoSolution(t *testing.T) {
	e := NewEngine()

	_, ok, err := e.QueryOne(`nonexistent(X).`)
	if err != nil {
		t.Fatalf("QueryOne failed: %v", err)
	}
	if ok {
		t.Fatal("expected no solution")
	}
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "facts.pl")
	if err := os.WriteFile(path, []byte(`fruit(apple). fruit(banana).`), 0644); err != nil {
		t.Fatal(err)
	}

	e := NewEngine()
	if err := e.Load(path); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	results, err := e.Query(`fruit(X).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestWithFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "facts.pl")
	if err := os.WriteFile(path, []byte(`animal(cat). animal(dog).`), 0644); err != nil {
		t.Fatal(err)
	}

	e := NewEngine(WithFile(path))

	results, err := e.Query(`animal(X).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestQueryWithRules(t *testing.T) {
	e := NewEngine(WithPrelude(`
		parent(tom, bob).
		parent(bob, ann).
		ancestor(X, Y) :- parent(X, Y).
		ancestor(X, Y) :- parent(X, Z), ancestor(Z, Y).
	`))

	results, err := e.Query(`ancestor(tom, X).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results (bob, ann), got %d", len(results))
	}

	names := map[string]bool{}
	for _, r := range results {
		names[r.Bindings["X"].(string)] = true
	}
	if !names["bob"] || !names["ann"] {
		t.Errorf("expected bob and ann, got %v", names)
	}
}
