# axon-mind

> Primitives · Part of the [lamina](https://github.com/benaskins/lamina-mono) workspace

Embedded Prolog engine for structured inference over facts and rules. axon-mind wraps [ichiban/prolog](https://github.com/ichiban/prolog) in a small Go API, letting you load `.pl` files, assert facts programmatically, run queries with full unification and backtracking, and bridge Go functions into Prolog as built-in predicates.

## Getting started

```bash
go get github.com/benaskins/axon-mind@latest
```

Requires Go 1.26+.

```go
package main

import (
    "fmt"
    mind "github.com/benaskins/axon-mind"
)

func main() {
    e := mind.NewEngine(mind.WithFile("facts.pl"))

    e.Assert("likes", "alice", "bob")

    results, _ := e.Query(`likes(alice, X).`)
    for _, r := range results {
        fmt.Println(r.Bindings["X"])
    }
}
```

See [`examples/`](examples/) for a complete workspace dependency graph modelled in Prolog.

## CLI

The `mind` CLI at `cmd/mind/` queries Prolog files from the command line and outputs JSON:

```bash
mind query -f examples/workspace.pl "depends_on(axon_chat, X)."
```

Install with `go install github.com/benaskins/axon-mind/cmd/mind@latest`.

## Key types

- **`Engine`** -- wraps the Prolog interpreter. Created with `NewEngine(opts...)`.
- **`Solution`** -- holds variable bindings from a successful query. Serialises to JSON via `Solution.JSON()`.
- **`Option`** -- functional options: `WithFile(path)`, `WithPrelude(source)`.
- **`Engine.Register(name, arity, fn)`** -- bridges a Go function into Prolog as a built-in predicate.
- **`Engine.Assert(functor, args...)`** -- adds facts programmatically with automatic dynamic declaration.

## License

MIT
