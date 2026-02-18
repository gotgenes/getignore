# Changelog


## 6.0.0 - 2026-02-18

### Added

- *(deps)* Upgrade urfave/cli from v2 to v3

### Changed

- *(agents)* Add project overview, structure, conventions, and tooling guidance
- *(deps)* Bump Ginkgo to v2.28.1 and Gomega to v1.39.1
- *(deps)* Bump google/go-github to v83.0.0
- *(agents)* Recommend ! notation for breaking changes
- *(mise)* Pin Go toolchain version to 1.26

### Fixed

- *(release)* Correctly detect breaking changes for major version bumps


## 5.0.4 - 2026-02-18

### Changed

- *(AI)* Add AGENTS.md.
- Update to latest version of Go.
- Update Ginkgo executable.
- Upgrade Go to 1.25 and update dependencies (ginkgo v2.25.3, gomega v1.38.2, urfave/cli v2.27.7)
- *(deps)* Bump google/go-github to v74.0.0
- *(agents)* Separate the testing instructions.
- Upgrade Go to 1.26
- *(release)* Automate releases with git-cliff
- *(release)* Automate Homebrew formula updates on release

### Fixed

- *(getter)* Use the value of max redirects option.


## 5.0.3 - 2024-01-25

### Added

- Add a release workflow.

### Changed

- Update the minimum Go version to build to 1.21.
- Upgrade urfave/cli.
- Upgrade gomega and ginkgo.
- Upgrade go-github.
- Update CI workflow.
- Name release workflow.


## 5.0.2 - 2022-11-27

### Added

- Add acceptance tests.

### Changed

- Check out submodules.
- Ensure binary built before running acceptance tests.
- Mark flags as part of command.
- Fetch all tags.
- Upgrade test dependencies.
- Upgrade go-github.
- Upgrade cli dependency.
- Upgrade Ginkgo and Gomega.
- Update to newer version of Go.
- Update google/go-github.


## 5.0.1 - 2022-02-15

### Changed

- Update to ginkgo v2.


## 5.0.0 - 2021-12-12

### Changed

- No more masters
- Run CI on pull requests.
- Format with golines.
- Update google/go-github.
- Update dependencies.


## 4.0.0 - 2021-11-16

### Added

- Add test for getting a single file.
- Add tests for ordered results.
- Add new, friendlier errors.
- Add MaxRequests.
- Add concurrent downloads.

### Changed

- DRY up the tests.
- Relocate lister to github package.
- Handle branch and tree errors.
- Returning contents in same order: responses ordered.
- Reorganize tests to make blob error cases easier.
- Check for errors from blob responses.
- Return an error if files don't exist in the tree.
- Rename names to paths.
- Create a cmd package.
- Implement the get subcommand.
- Rename options and change lister to getter.
- Rename NamedIgnoreContents to NamedContents.
- Convert tests to ginkgo.
- Rename contentstructs to contents.
- Accept names that don't include an extension.
- Move packages to pkg subdirectory.
- Move command files under cmd.
- Update README.
- Change target name.
- Describe dev-install command.
- Update dependencies.
- Correct the flag name.
- Cut the release.

### Removed

- Remove HTTPGetter.
- Remove unused tests.


## 3.0.1 - 2021-10-17

### Removed

- Remove unused code.


## 3.0.0 - 2021-10-16

### Added

- Add GitHubLister and ensure it sets correct headers.
- Add test for empty tree.

### Changed

- Start GitHubLister and get a passing test.
- Check headers on both requests.
- Extract file paths from tree.
- Filter files by suffix.
- Filter for files, only.
- Handle errors from branches endpoint.
- Handle errors for tree endpoint.
- Handle empty response for branches endpoint.
- Rename organization to owner per API nomenclature.
- Use app.Run instead of deprecated method.
- Fix setting build version.
- Use GitHubLister for list command.
- Upgrade dependencies.
- Install ginkgo for CI.


## 2.1.1 - 2021-09-05

### Changed

- Include the LDFLAGS for go install.


## 2.1.0 - 2021-09-05

### Added

- Add GitHub actions for CI.
- Add dependabot support for dependency updates.
- Add tag target.

### Changed

- Fix run.
- Update badge to GitHub Actions.
- Make minimum Go version 1.16.
- Compile for macOS ARM64.
- Upgrade dependencies.
- Run go vet before go test.
- Supply version at build time through ldflags.
- Update dependencies and to Go 1.17.
- Apply go mod tidy.

### Removed

- Remove Travis CI integration.


## 2.0.0 - 2020-08-28

### Added

- Add shell completions.

### Changed

- Re-enable x86 32-bit builds
- Fix expected error string for newer version of Go.
- Switch to using go mod for dependency management.
- Upgrade urfave/cli to v1.22.4.
- Upgrade stretchr/testify to v1.6.1.
- Use Makefile with CI.
- Upgrade to urfave/cli v2.2.0.
- Include completions in distribution files.
- Update OS and architecture targets for distribution.
- Handle tags with "v"-prefix.
- Release 2.0.0.

### Removed

- Remove redundant struct name.
- Remove using glide for package discovery.


## 1.0.0 - 2017-07-17

### Changed

- Ensure patterns file contents end in newline.
- Replace testutil with testify.
- Normalize whitespace in patterns file contents.
- Minor formatting changes
- Release 1.0.0


## 0.3.0 - 2017-04-07

### Added

- Add list command

### Changed

- Update urfave/cli dependency
- Use strings.NewReader for in-memory files
- Update the test command
- Fix tests broken after rename
- Increase the minimum Go version
- Internalize channels into the HTTPGetter.
- Split error handling into separate goroutine and channel
- Ensure contents get written in the order requested
- Fix go vet warnings.
- Rearrange order of functions
- Relocates test for NamedIgnoreContents
- Use go vet
- Fix more go vet warnings
- Fix script section
- Fix dist target
- Only build for x86-64
- Mention macOS
- Change mentioned order of commands
- Release 0.3.0

### Removed

- Removes RetrievedContents, no longer used


## 0.2.0 - 2016-12-28

### Added

- Add Installation section
- Add tests for GetIgnoreFiles
- Add a change log

### Changed

- Fix typo
- Updates readme with choco install steps
- Use STDOUT as default output for get
- Rename functions from "fetch" to "download"
- Update .gitignore
- Refactor unit tests
- Correctly restrict the maximum number of simultaneous HTTP requests
- Refactor downloadIgnoreFile to be a method
- Reorganize code into separate packages
- Release 0.2.0

### Removed

- Remove superfluous log statement
- Remove NamesToURLs


## 0.1.0 - 2016-11-18

### Added

- Adds ignoreFetcher
- Adds methods for http requests and output
- Adds flags for file or command arguments
- Adds basic readme
- Adds test for writeIgnoreFile
- Adds option --max-connections
- Adds custom error type FailedURLs
- Adds decorated section headers
- Adds MIT license
- Add more usage information
- Add DisplayName to handle formatting in output
- Add godoc comments
- Add Travis CI support
- Add a CONTRIBUTORS file
- Add Travis CI build status
- Add forgotten go install
- Adds Makefile to help with distribution

### Changed

- Let's get this puppy started
- Implements NamesToUrls via channels
- Renames project to getignore
- Passes file handle and extracts helper method
- Read names from a file
- Extracts method to get names from arguments
- Renames README
- Renames and adds tests for parseNamesFile
- Handle blank lines and spaces in names file
- Resolves golint errors
- Use the cli package for the CLI
- Use arrays instead of channels for NamesToURLs
- Big refactor, program in non-working state
- Program working again
- Have the program exit for failed HTTP responses
- Push error down to individual FailedURL instances
- Ensure patterns output in same order as names
- Permit extensions on input names
- Rename ignoreFetcher to HTTPIgnoreGetter
- Push GetIgnoreFiles to method of HTTPIgnoreGetter
- Rename FetchedContents to RetrievedContents
- Use more abstract wording of "source" over "URL"
- Explain the omission of file extension more thoroughly
- Return back to original directory
- Fix glide configuration and updates dependencies
- Specify command to run tests
- Ignore builds

### Removed

- Removes unused helper functions

