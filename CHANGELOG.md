# Changelog

## 5.0.0 - 2021-12-11

### Changed

- Changed default branch from `master` to `main`.
- Updated dependencies.

## 4.0.0 - 2021-11-16

### Added

- Added the following options to `get` command:

  - `--base-url`: the base URL to the GitHub REST API v3-compatible server
  - `--owner`: the owner or organization of the repository of gitignore files
  - `--repository`: the name of the repository of gitignore files
  - `--branch`: the branch of the repository from which to list gitignore files

### Removed

- Removed the `--api-url` option from the `get` command.

### Changed

- Refactored `get` command to use [GitHub REST API v3](https://docs.github.com/en/rest), making it similar to the `list` command.
- Renamed `--max-connections` to `--max-requests` for the `get` command.
- Changed the short alias flag for `--owner` from `-o` to `-w` of the `list` command, to avoid conflict with the `--output` alias of the `get` command.
- Moved all non-command code to `pkg` subdirectory, following the recommended layout in [golang-standards/project-layout](https://github.com/golang-standards/project-layout).
  The core code can be found under `pkg/getignore/`.
  The code for interacting with GitHub can be found under `pkg/github/`
- Switched testing framework from [stretch/testify](https://github.com/stretchr/testify) to [Ginkgo](https://onsi.github.io/ginkgo/) and [Gomega](https://onsi.github.io/ginkgo/).
- Renamed `install-deps` rule to `dev-install` in `Makefile`.

## 3.0.1 - 2021-10-16

### Removed

- Removed unused code for `list` command.

## 3.0.0 - 2021-10-16

### Added

- Added the following options to `list` command:

  - `--base-url`: the base URL to the GitHub REST API v3-compatible server
  - `--owner`: the owner or organization of the repository of gitignore files
  - `--repository`: the name of the repository of gitignore files
  - `--branch`: the branch of the repository from which to list gitignore files

### Removed

- Removed the `--api-url` option from the `list` command.

### Changed

- Refactored `list` command to use [GitHub REST API v3](https://docs.github.com/en/rest).
- Added dependency on [go-github](https://github.com/google/go-github) to interact with the GitHub REST API.
- Upgraded dependencies.

## 2.1.1 - 2021-09-04

### Fixed

- Fixed `make install` so that it includes necessary `LDFLAGS`, too.

## 2.1.0 - 2021-09-04

### Added

- Added builds for macOS ARM64 (Apple Silicon)
- Added GitHub Actions for CI.

### Removed

- Removed Travis CI integration.

### Changed

- Upgraded dependencies.

## 2.0.0 - 2020-08-28

### Added

- Added completions for Bash and zsh.

### Changed

- Upgraded dependencies, including to v2 of [urfave/cli](https://github.com/urfave/cli). This is a backwards-incompatible change because [the commands will no longer accept flags after arguments](https://github.com/urfave/cli/blob/master/docs/migrate-v1-to-v2.md#flags-before-args).
- Switched to using `go mod` to manage dependencies, dropping the dependency on [glide](https://glide.sh/).

## 1.0.0 - 2017-07-16

### Changed

- Tests now depend on [testify](https://github.com/stretchr/testify).

### Removed

- `testutils` module replaced by [testify](https://github.com/stretchr/testify).

### Fixed

- Empty lines before and after patterns in the retrieved file contents are now stripped, making line spacing consistent.

## 0.3.0 - 2017-04-07

### Added

- `list` command now available.
- `go vet` now applied during CI.

### Changed

- Changed the API for `HTTPGetter.GetIgnoreFiles` to take an array of names and return an array of `contentstructs.NamedIgnoreContents` and an `error`. Previously this took a channel as input, on which it sent the results. Channels have been moved internal to the `GetIgnoreFiles` method.
- Updated dependency on [cli](https://github.com/urfave/cli).
- Updated the minimum Go version to build to 1.8.

### Removed

- Removed `contentstructs.RetrievedContents` as it was no longer necessary.

### Fixed

- Fixed an issue where `getignore get` would deadlock if the number of files requested exceeded the maximum number of connections. (See issue #3.)

## 0.2.0 - 2016-12-28

### Added

- Additional unit tests.

### Changed

- `get` command now writes output to `STDOUT` by default (previously defaulted to `.gitignore` in the current working directory).
- Reorganized code into separate packages.

## 0.1.0 - 2016-11-17

Initial release.
