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
- Envtest suites live under `test/envtest/<area>/` and are for AM only: controllers and admission against a real kube-apiserver + etcd with the in-process AM mock, run with `make envtest`. There is no APIM mock
- APIM controllers and admission that need the manager cache or a live APIM are tested in `test/integration/<area>/` against kind, run with `make it`
- Do not write e2e tests: platform behaviour end to end is covered by the epic's critical user journey in the platform e2e repo, written separately from the story (see AGENTS.md, Testing)

## Forbidden Patterns
- Do not use `github.com/golang/protobuf` (use `google.golang.org/protobuf`)
- Do not use `github.com/satori/go.uuid` (use `github.com/google/uuid`)
