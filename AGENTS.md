# AGENTS Guide: getignore

`getignore` is a CLI tool that downloads `.gitignore` templates from GitHub's [gitignore](https://github.com/github/gitignore) repository.
It supports listing available templates and downloading one or more by name, with concurrent fetching and combined output.

## Project Structure

```
getignore/
├── cmd/getignore/          # CLI entrypoint and command definitions
│   ├── main.go             # App setup, registers List and Get commands
│   ├── common.go           # Shared CLI flags and option builders
│   ├── get.go              # Get command (download templates)
│   └── list.go             # List command (list available templates)
├── pkg/
│   ├── getignore/          # Core library (parsing, writing, types)
│   │   ├── identifiers.go  # Version variable, user agent string
│   │   ├── named_contents.go
│   │   ├── parsers.go
│   │   ├── writers.go
│   │   └── failed_file.go
│   └── github/             # GitHub API integration
│       ├── constants.go    # Default owner, repo, branch, suffix
│       └── getter.go       # Getter with List() and Get() methods
├── test/
│   └── acceptance_tests.bats  # Bats acceptance tests (end-to-end)
├── scripts/
│   └── release.sh          # Automated release script (git-cliff)
├── cliff.toml              # git-cliff changelog configuration
├── mise.toml               # Tool version manager (Go 1.26)
├── Makefile                # Build, test, release targets
├── go.mod
└── go.sum
```

## Building

1. Build: `make build` (binary at `./getignore`); multi-platform: `make dist`; install locally: `make install`.
2. Version injection: pass `LDFLAGS` as in Makefile (`-X github.com/gotgenes/getignore/pkg/getignore.Version`); tests/run already include it.
3. Dev deps: `make dev-install` (installs pinned Ginkgo v2.28.1); ensure Go >= 1.26.

## Code Style

1. Linting: rely on `go vet`; keep code `gofmt`/`goimports` clean; group imports std, third-party, internal, each separated by a blank line.
2. Naming: exported identifiers use PascalCase with doc comment sentence; unexported camelCase; keep receiver names short (e.g., `g`, `cwo`). Avoid stutter (`getignore.NamedContents` not `GetignoreNamedContents`).
3. Errors: wrap with `fmt.Errorf("context: %w", err)`; user-facing aggregation uses helper constructors (see `newGetError`, `newListError`); do not leak low-level messages directly.
4. Concurrency: guard ordering determinism (see `contentsWithOrdering`); when adding goroutines, ensure bounded by configuration (pattern: `maxRequests := min(numFiles, g.MaxRequests)`).
5. I/O: trim and validate inputs (e.g., `ParseNamesFile`); always flush buffered writers and propagate error.
6. Option pattern: create `type ...Option func(*params)` and apply in constructor before building clients; follow existing `getterParams` model.
7. CLI flags: define in `cmd/getignore/common.go`; keep flag names kebab-case and map to option builders; update `stringFlagsToOptions` if adding string flag.
8. Formatting: run `go fmt ./...`; avoid inline comments explaining obvious code; keep lines <120 chars.
9. Dependencies: add to `go.mod`; prefer stdlib before adding libs; pin explicit versions.
10. User agent: maintain template `getignore/<version>`; update only via Version variable (never hardcode).
11. Error messages for missing resources: prefer phrase `not present in file tree` or `failed to download` for consistency.
12. Do not introduce magic numbers: add consts or derive (e.g., `DefaultMaxRequests = runtime.NumCPU() - 1`).
13. Security/perf: avoid unbounded parallelism; validate external URLs when introducing new network calls.

## Testing

### Writing tests

1. For new functionality, including bug fixes, add tests before implementation.
2. Add one new test for each behavior (expecting the test to first fail), then implement the behavior, then confirm the test passes, along with existing tests (Test-Driven Development). If multiple behavior changes are expected, test and implement each, sequentially.
3. Testing HTTP: use `ghttp` with `VerifyRequest`, `VerifyHeader`, `RespondWith`; replicate existing patterns for new endpoints.
4. Add new test: prefer Ginkgo/Gomega BDD style; suites named `*suite_test.go` calling `RunSpecs`.
5. Prefer using literal or hard-coded values in tests, rather than variables or constants, to ensure test clarity. Avoid emphasizing DRY in tests. Prefer explicit values in test cases to make it clear to the reader what actual values are being provided as inputs and what values are expected in outputs.

### Running tests

1. Tests: unit/integration via Ginkgo: `make test` (runs `go vet ./...` then `ginkgo -r`). Single Go test file: `ginkgo -r --focus FileName` or run a package: `ginkgo pkg/getignore -focus ParseNamesFile`; single spec: `ginkgo -r -focus 'ParseNamesFile parses a standard file'`.
2. Acceptance tests (CLI behaviors): `make acceptance-test` (requires built binary and `bats` supplied via CI; locally install bats if absent).
3. All tests: `make test-all` (adds Bats acceptance).

## Tool Invocation

This project uses [mise](https://mise.jdx.dev/) to manage the Go toolchain version.
The pinned version is specified in `mise.toml` (currently Go 1.26).

When running commands in a non-interactive shell (e.g., as a coding agent), use `mise exec --` to ensure the correct Go version is available:

```bash
mise exec -- go build ./...
mise exec -- make test
```

Interactive shells with mise activated (via `.zshrc`/`.bashrc`) do not need the prefix.

## Releasing

Releases are automated via `git-cliff` and GitHub Actions.

- `make release-dry-run` -- preview the next version and changelog without making changes.
- `make release` -- regenerates `CHANGELOG.md`, commits, creates an annotated tag, and pushes. The tag push triggers CI which creates a GitHub Release, builds cross-platform binaries, and updates the Homebrew formula.
- Changelog generation is configured in `cliff.toml` and parses Conventional Commit messages.
- Do not manually edit `CHANGELOG.md`; it is generated from commit history.

## Other Conventions

### Git commit messages

- Use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) for commit messages.
- Use the optional scope field to indicate the area of change, e.g., `getter`, `release`, `deps`, `agents`. If the change is broad, omit the scope.
- For breaking changes, use the `!` after the type/scope (e.g., `feat(deps)!:`) and include a `BREAKING CHANGE:` footer in the commit body explaining the user-facing impact.
- Keep the subject line as a summary of what the changes do. If the reasoning behind the changes needs explanation, provide it in the body.
- These conventions matter because `git-cliff` parses commit messages to generate the changelog and determine version bumps.

Examples from the project history:

```
feat(deps)!: upgrade urfave/cli from v2 to v3
chore(deps): bump google/go-github to v74.0.0
build(release): automate releases with git-cliff
chore: upgrade Go to 1.26
```

### Before committing

- Run `make test` (runs `go vet` and all Ginkgo tests).
- Ensure code is formatted (`go fmt ./...`).
- Review changes with `git diff`.
