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
- Unit tests live under `test/unit/<area>/` (Ginkgo suite + specs), colocated with the other unit suites — never next to controller/`internal` packages
- Use Ginkgo v2 and Gomega (not raw `testing.T` with `t.Errorf`/`t.Fatal`)
- Dot-imports for `github.com/onsi/ginkgo/v2` and `github.com/onsi/gomega` are allowed in test files
- Prefer testing through importable packages (`internal/...`, `api/...`); do not place `_test.go` under nested `controllers/**/internal`
- Unit tests are the only Go tests to add. Anything needing a cluster or a live APIM belongs in the [platform e2e repo](https://github.com/gravitee-io/gravitee-platform-e2e); do not extend `test/integration/`

## Forbidden Patterns
- Do not use `github.com/golang/protobuf` (use `google.golang.org/protobuf`)
- Do not use `github.com/satori/go.uuid` (use `github.com/google/uuid`)
