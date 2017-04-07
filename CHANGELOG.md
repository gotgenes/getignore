# Change Log

## 0.3.0 TBD

### Added

* `list` command now available.
* `go vet` now applied during CI.

### Changed

* Changed the API for `HTTPGetter.GetIgnoreFiles` to take an array of names and return an array of `contentstructs.NamedIgnoreContents` and an `error`. Previously this took a channel as input, on which it sent the results. Channels have been moved internal to the `GetIgnoreFiles` method.
* Updated dependency on [urfave/cli](https://github.com/urfave/cli).
* Updated the minimum Go version to build to 1.8.

### Removed

* Removed `contentstructs.RetrievedContents` as it was no longer necessary.

### Fixed

* Fixed an issue where `getignore get` would deadlock if the number of files requested exceeded the maximum number of connections. (See issue #3.)


## 0.2.0 2016-12-28

### Added

* Additional unit tests.

### Changed

* `get` command now writes output to `STDOUT` by default (previously defaulted to `.gitignore` in the current working directory).
* Reorganized code into separate packages.


## 0.1.0 2016-11-17

Initial release.
