package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	// Build the binary once for all tests
	build := exec.Command("go", "build", "-o", "mind_test_bin", ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		panic("build failed: " + string(out))
	}
	code := m.Run()
	os.Remove("mind_test_bin")
	os.Exit(code)
}

func mindCmd(args ...string) *exec.Cmd {
	abs, _ := filepath.Abs("mind_test_bin")
	return exec.Command(abs, args...)
}

func TestQueryWithFile(t *testing.T) {
	dir := t.TempDir()
	plFile := filepath.Join(dir, "test.pl")
	os.WriteFile(plFile, []byte(`fruit(apple). fruit(banana). fruit(cherry).`), 0644)

	cmd := mindCmd("query", "-f", plFile, "fruit(X).")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	var results []map[string]any
	if err := json.Unmarshal(out, &results); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0]["X"] != "apple" {
		t.Errorf("expected apple, got %v", results[0]["X"])
	}
}

func TestQueryMultipleFiles(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "facts.pl")
	os.WriteFile(f1, []byte(`
		parent(tom, bob).
		parent(bob, ann).
	`), 0644)

	f2 := filepath.Join(dir, "rules.pl")
	os.WriteFile(f2, []byte(`
		ancestor(X, Y) :- parent(X, Y).
		ancestor(X, Y) :- parent(X, Z), ancestor(Z, Y).
	`), 0644)

	cmd := mindCmd("query", "-f", f1, "-f", f2, "ancestor(tom, X).")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	var results []map[string]any
	json.Unmarshal(out, &results)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d: %s", len(results), out)
	}
}

func TestQueryNoSolutionsExitCode(t *testing.T) {
	dir := t.TempDir()
	plFile := filepath.Join(dir, "empty.pl")
	os.WriteFile(plFile, []byte(`fruit(apple).`), 0644)

	cmd := mindCmd("query", "-f", plFile, "fruit(banana).")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit code")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
	}
	// Should still output empty JSON array
	if string(out) != "[]\n" {
		t.Errorf("expected []\n, got %q", string(out))
	}
}

func TestQueryNoArgs(t *testing.T) {
	cmd := mindCmd()
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit code")
	}
	exitErr := err.(*exec.ExitError)
	if exitErr.ExitCode() != 2 {
		t.Errorf("expected exit code 2, got %d", exitErr.ExitCode())
	}
}
