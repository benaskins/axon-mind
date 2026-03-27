---
module: github.com/benaskins/axon-mind
kind: library
---

# axon-mind

Domain-agnostic logical reasoning library embedding Prolog (ichiban/prolog) in Go.

## What it does

- Load facts and rules from `.pl` files or programmatically
- Run queries with full unification and backtracking
- Get results as structured Go types or JSON
- Register Go functions as Prolog built-in predicates

## Architecture

- `mind.go` — Engine type, options, core API
- `cmd/mind/` — CLI binary (`mind query -f file.pl "goal"`)

## Running tests

```bash
go test ./...
go vet ./...
```

## Dependencies

- github.com/ichiban/prolog — embedded Prolog interpreter
- No other external dependencies
