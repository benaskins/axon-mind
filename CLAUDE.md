@AGENTS.md

## Conventions
- `engine.Register(name, arity, func)` bridges Go functions into Prolog as built-in predicates
- Keep facts and rules in separate `.pl` files when composing knowledge bases
- Assert auto-declares dynamic — no manual `:- dynamic` directives needed
- Standard Prolog `Exec` replaces predicates on re-consult; load order matters
- CLI at `cmd/mind/` outputs JSON — `mind query -f file.pl "goal."`

## Constraints
- No axon-* dependencies — this is a standalone reasoning library
- Only external dependency is `github.com/ichiban/prolog`
- Do not add HTTP handlers or server code — this is a library + CLI only
- Do not embed domain-specific knowledge — keep it domain-agnostic

## Testing
- `go test ./...` — tests run with no external services
- `go vet ./...` — must be clean
- Test queries against known `.pl` fixtures; assert on structured Go results, not string output
