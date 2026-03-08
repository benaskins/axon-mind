package mind

import (
	"encoding/json"
	"testing"
)

func TestSolutionJSON(t *testing.T) {
	e := NewEngine(WithPrelude(`
		person(alice, 30).
		person(bob, 25).
	`))

	results, err := e.Query(`person(Name, Age).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	data, err := results[0].JSON()
	if err != nil {
		t.Fatalf("JSON failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if m["Name"] != "alice" {
		t.Errorf("expected Name=alice, got %v", m["Name"])
	}
	// Numbers come back from Prolog as integers
	if m["Age"] != float64(30) {
		t.Errorf("expected Age=30, got %v (type %T)", m["Age"], m["Age"])
	}
}

func TestSolutionsJSON(t *testing.T) {
	e := NewEngine(WithPrelude(`
		color(red).
		color(green).
		color(blue).
	`))

	results, err := e.Query(`color(X).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	data, err := SolutionsJSON(results)
	if err != nil {
		t.Fatalf("SolutionsJSON failed: %v", err)
	}

	var arr []map[string]any
	if err := json.Unmarshal(data, &arr); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(arr) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(arr))
	}

	colors := []string{}
	for _, m := range arr {
		colors = append(colors, m["X"].(string))
	}
	if colors[0] != "red" || colors[1] != "green" || colors[2] != "blue" {
		t.Errorf("expected [red green blue], got %v", colors)
	}
}

func TestSolutionJSONEmptyBindings(t *testing.T) {
	s := Solution{Bindings: map[string]any{}}
	data, err := s.JSON()
	if err != nil {
		t.Fatalf("JSON failed: %v", err)
	}
	if string(data) != "{}" {
		t.Errorf("expected {}, got %s", data)
	}
}

func TestSolutionsJSONEmpty(t *testing.T) {
	data, err := SolutionsJSON(nil)
	if err != nil {
		t.Fatalf("SolutionsJSON failed: %v", err)
	}
	if string(data) != "[]" {
		t.Errorf("expected [], got %s", data)
	}
}
