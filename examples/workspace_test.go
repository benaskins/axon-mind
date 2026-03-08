package examples

import (
	"testing"

	mind "github.com/benaskins/axon-mind"
)

func loadWorkspace(t *testing.T) *mind.Engine {
	t.Helper()
	e := mind.NewEngine(mind.WithFile("workspace.pl"))
	if e == nil {
		t.Fatal("failed to create engine")
	}
	return e
}

func TestTransitiveDependencies(t *testing.T) {
	e := loadWorkspace(t)

	// axon_chat depends on axon_tool transitively via axon_loop
	results, err := e.Query(`transitive_dep(axon_chat, X).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	deps := map[string]bool{}
	for _, r := range results {
		deps[r.Bindings["X"].(string)] = true
	}

	// Direct deps
	if !deps["axon"] {
		t.Error("expected axon as dependency")
	}
	if !deps["axon_loop"] {
		t.Error("expected axon_loop as dependency")
	}
	if !deps["axon_tool"] {
		t.Error("expected axon_tool as dependency")
	}
}

func TestAffectedBy(t *testing.T) {
	e := loadWorkspace(t)

	// Changing axon_tool should affect axon_loop, axon_talk, axon_chat, axon_lens
	results, err := e.Query(`affected_by(axon_tool, X).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	affected := map[string]bool{}
	for _, r := range results {
		affected[r.Bindings["X"].(string)] = true
	}

	expected := []string{"axon_loop", "axon_talk", "axon_chat", "axon_lens"}
	for _, name := range expected {
		if !affected[name] {
			t.Errorf("expected %s to be affected by axon_tool change", name)
		}
	}
}

func TestStandaloneModules(t *testing.T) {
	e := loadWorkspace(t)

	results, err := e.Query(`standalone(X).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	standalones := map[string]bool{}
	for _, r := range results {
		standalones[r.Bindings["X"].(string)] = true
	}

	// axon has no dependencies
	if !standalones["axon"] {
		t.Error("expected axon to be standalone")
	}
	// axon_eval has no dependencies
	if !standalones["axon_eval"] {
		t.Error("expected axon_eval to be standalone")
	}
	// axon_chat has dependencies, should NOT be standalone
	if standalones["axon_chat"] {
		t.Error("axon_chat should not be standalone")
	}
}

func TestServiceUsingLibrary(t *testing.T) {
	e := loadWorkspace(t)

	// Which services use axon_tool?
	results, err := e.Query(`service_using(axon_tool, X).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	services := map[string]bool{}
	for _, r := range results {
		services[r.Bindings["X"].(string)] = true
	}

	if !services["axon_chat"] {
		t.Error("expected axon_chat to use axon_tool")
	}
}

func TestJSONOutput(t *testing.T) {
	e := loadWorkspace(t)

	results, err := e.Query(`kind(X, library).`)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	data, err := mind.SolutionsJSON(results)
	if err != nil {
		t.Fatalf("SolutionsJSON failed: %v", err)
	}

	// Should be valid JSON
	if len(data) == 0 {
		t.Fatal("expected non-empty JSON")
	}
}
