# Go Coding Standards

## Error Handling
- Always check errors from type assertions: use `v, ok := x.(T)` form
- Wrap errors with context using `fmt.Errorf("...: %w", err)`
- Sentinel errors must be prefixed with `Err` (e.g. `ErrNotFound`)

## Concurrency
- Never copy a sync.Mutex or sync.WaitGroup after first use
- Always pass context.Context as the first parameter

## Style
- Max function length: 100 lines
- Max cyclomatic complexity: 30
- No naked returns
- No dot-imports except in test files for ginkgo/v2 and gomega
- Imports ordered: stdlib, external, internal (enforced by goimports)

## Naming
- `Api/Url/Http` casing is accepted (not forced to `API/URL/HTTP`)
- Avoid shadowing predeclared identifiers (e.g. `error`, `len`, `new`)

## Testing
- Use Ginkgo v2 and Gomega for unit and integration tests
- Dot-imports for `github.com/onsi/ginkgo/v2` and `github.com/onsi/gomega` are allowed in test files

## Forbidden Patterns
- Do not use `github.com/golang/protobuf` (use `google.golang.org/protobuf`)
- Do not use `github.com/satori/go.uuid` (use `github.com/google/uuid`)
