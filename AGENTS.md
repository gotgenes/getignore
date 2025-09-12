# AGENTS

1. Build: `make build` (binary at `cmd/getignore`); multi-platform: `make dist`; install locally: `make install`.
2. Version injection: pass `LDFLAGS` as in Makefile (`-X github.com/gotgenes/getignore/pkg/getignore.Version`); tests/run already include it.
3. Dev deps: `make dev-install` (installs pinned Ginkgo v2.25.3); ensure Go ≥ 1.25.
4. Linting: rely on `go vet`; keep code `gofmt`/`goimports` clean; group imports std, third-party, internal, each separated by a blank line.
5. Naming: exported identifiers use PascalCase with doc comment sentence; unexported camelCase; keep receiver names short (e.g., `g`, `cwo`). Avoid stutter (`getignore.NamedContents` not `GetignoreNamedContents`).
6. Errors: wrap with `fmt.Errorf("context: %w", err)`; user-facing aggregation uses helper constructors (see `newGetError`, `newListError`); do not leak low-level messages directly.
7. Concurrency: guard ordering determinism (see `contentsWithOrdering`); when adding goroutines, ensure bounded by configuration (pattern: `maxRequests := min(numFiles, g.MaxRequests)`).
8. I/O: trim and validate inputs (e.g., `ParseNamesFile`); always flush buffered writers and propagate error.
9. Option pattern: create `type ...Option func(*params)` and apply in constructor before building clients; follow existing `getterParams` model.
10. CLI flags: define in `cmd/getignore/common.go`; keep flag names kebab-case and map to option builders; update `stringFlagsToOptions` if adding string flag.
11. Formatting: run `go fmt ./...`; avoid inline comments explaining obvious code; keep lines <120 chars.
12. Dependencies: add to `go.mod`; prefer stdlib before adding libs; pin explicit versions.
13. User agent: maintain template `getignore/<version>`; update only via Version variable (never hardcode).
14. Error messages for missing resources: prefer phrase `not present in file tree` or `failed to download` for consistency.
15. Do not introduce magic numbers: add consts or derive (e.g., `DefaultMaxRequests = runtime.NumCPU() - 1`).
16. Security/perf: avoid unbounded parallelism; validate external URLs when introducing new network calls.

## Testing

### Writing tests

1. For new functionality, including bug fixes, add tests before implementation.
2. Add one new test for each behavior (expecting the test to first fail), then implement the behavior, then confirm the test passes, along with existing tests (Test-Driven Development). If multiple behavior changes are expected, test and implement each, sequentially.
3. Testing HTTP: use `ghttp` with `VerifyRequest`, `VerifyHeader`, `RespondWith`; replicate existing patterns for new endpoints.
4. Add new test: prefer Ginkgo/Gomega BDD style; suites named `*suite_test.go` calling `RunSpecs`.
5. Prefer using literal or hard-coded values in tests, rather than variables or constants, to ensure test clarity. Avoid emphasizing DRY in tests. Prefer explicit values in test cases to make it clear to the reader what actual values are being provided as inputs and what values are expected in outputs.

### Running tests

3. Tests: unit/integration via Ginkgo: `make test` (runs `go vet ./...` then `ginkgo -r`). Single Go test file: `ginkgo -r --focus FileName` or run a package: `ginkgo pkg/getignore -focus ParseNamesFile`; single spec: `ginkgo -r -focus 'ParseNamesFile parses a standard file'`.
4. Acceptance tests (CLI behaviors): `make acceptance-test` (requires built binary and `bats` supplied via CI; locally install bats if absent).
5. All tests: `make test-all` (adds Bats acceptance).
