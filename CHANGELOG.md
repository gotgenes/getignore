# Changelog

## 2.0.0 - 2020-08-28

### Added

* Added completions for Bash and zsh.

### Changed

* Upgraded dependencies, including to v2 of [urfave/cli](https://github.com/urfave/cli). This is a backwards-incompatible change because [the commands will no longer accept flags after arguments](https://github.com/urfave/cli/blob/master/docs/migrate-v1-to-v2.md#flags-before-args).
* Switched to using `go mod` to manage dependencies, dropping the dependency on [glide](https://glide.sh/).


## 1.0.0 - 2017-07-16

### Changed

* Tests now depend on [testify](https://github.com/stretchr/testify).

### Removed

* `testutils` module replaced by [testify](https://github.com/stretchr/testify).

### Fixed

* Empty lines before and after patterns in the retrieved file contents are now stripped, making line spacing consistent.


## 0.3.0 - 2017-04-07

### Added

* `list` command now available.
* `go vet` now applied during CI.

### Changed

* Changed the API for `HTTPGetter.GetIgnoreFiles` to take an array of names and return an array of `contentstructs.NamedIgnoreContents` and an `error`. Previously this took a channel as input, on which it sent the results. Channels have been moved internal to the `GetIgnoreFiles` method.
* Updated dependency on [cli](https://github.com/urfave/cli).
* Updated the minimum Go version to build to 1.8.

### Removed

* Removed `contentstructs.RetrievedContents` as it was no longer necessary.

### Fixed

* Fixed an issue where `getignore get` would deadlock if the number of files requested exceeded the maximum number of connections. (See issue #3.)


## 0.2.0 - 2016-12-28

### Added

* Additional unit tests.

### Changed

* `get` command now writes output to `STDOUT` by default (previously defaulted to `.gitignore` in the current working directory).
* Reorganized code into separate packages.


## 0.1.0 - 2016-11-17

Initial release.
