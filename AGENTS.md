# AGENTS

1. Build: `make build` (binary at `cmd/getignore`); multi-platform: `make dist`; install locally: `make install`.
2. Version injection: pass `LDFLAGS` as in Makefile (`-X github.com/gotgenes/getignore/pkg/getignore.Version`); tests/run already include it.
3. Tests: unit/integration via Ginkgo: `make test` (runs `go vet ./...` then `ginkgo -r`). All tests: `make test-all` (adds Bats acceptance). Single Go test file: `ginkgo -r --focus FileName` or run a package: `ginkgo pkg/getignore -focus ParseNamesFile`; single spec: `ginkgo -r -focus 'ParseNamesFile parses a standard file'`.
4. Add new test: prefer Ginkgo/Gomega BDD style; suites named `*suite_test.go` calling `RunSpecs`.
5. Acceptance tests (CLI behaviors): `make acceptance-test` (requires built binary and `bats` supplied via CI; locally install bats if absent).
6. Dev deps: `make dev-install` (installs pinned Ginkgo); ensure Go ≥ 1.25.
7. Linting: rely on `go vet`; keep code `gofmt`/`goimports` clean; group imports std, third-party, internal, each separated by a blank line.
8. Naming: exported identifiers use PascalCase with doc comment sentence; unexported camelCase; keep receiver names short (e.g., `g`, `cwo`). Avoid stutter (`getignore.NamedContents` not `GetignoreNamedContents`).
9. Errors: wrap with `fmt.Errorf("context: %w", err)`; user-facing aggregation uses helper constructors (see `newGetError`, `newListError`); do not leak low-level messages directly.
10. Concurrency: guard ordering determinism (see `contentsWithOrdering`); when adding goroutines, ensure bounded by configuration (pattern: `maxRequests := min(numFiles, g.MaxRequests)`).
11. I/O: trim and validate inputs (e.g., `ParseNamesFile`); always flush buffered writers and propagate error.
12. Option pattern: create `type ...Option func(*params)` and apply in constructor before building clients; follow existing `getterParams` model.
13. CLI flags: define in `cmd/getignore/common.go`; keep flag names kebab-case and map to option builders; update `stringFlagsToOptions` if adding string flag.
14. Formatting: run `go fmt ./...`; avoid inline comments explaining obvious code; keep lines <120 chars.
15. Dependencies: add to `go.mod`; prefer stdlib before adding libs; pin explicit versions.
16. Testing HTTP: use `ghttp` with `VerifyRequest`, `VerifyHeader`, `RespondWith`; replicate existing patterns for new endpoints.
17. User agent: maintain template `getignore/<version>`; update only via Version variable (never hardcode).
18. Error messages for missing resources: prefer phrase `not present in file tree` or `failed to download` for consistency.
19. Do not introduce magic numbers: add consts or derive (e.g., `DefaultMaxRequests = runtime.NumCPU() - 1`).
20. Security/perf: avoid unbounded parallelism; validate external URLs when introducing new network calls.
