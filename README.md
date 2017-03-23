# getignore

getignore bootstraps `.gitignore` files from [GitHub .gitignore patterns](https://github.com/github/gitignore).

[![Travis CI Build Status](https://travis-ci.org/gotgenes/getignore.svg?branch=master)](https://travis-ci.org/gotgenes/getignore)


## Installation

### Homebrew users (OS X)

[Homebrew](http://brew.sh) users can install getignore using the following commands:

```shell
brew tap gotgenes/homebrew-gotgenes
brew update
brew install getignore
```

### Chocolatey users (Windows)

[Chocolatey](https://chocolatey.org/) users can install getignore using the following command (as admin):

```shell
choco install getignore
```

### Windows, Linux, and OS X

Download and unpack a pre-compiled executable from [the releases page](https://github.com/gotgenes/getignore/releases). Make sure to place the executable in your shell's `PATH`.


### Other platforms

See [the Building section, below](#building).


## Usage

getignore supports the following commands:

* [`get`](#get)
* [`list`](#list)
* [`help`](#help)


### help

You can get help on getignore itself with

```shell
getignore help
```

You can use `help` to get the full usage information of any other command. For example, to get the full usage of the `get` command, run

```shell
getignore help get
```

### get

Use the `get` command to obtain gitignore patterns from remote repositories. By default, it will obtain patterns from the [GitHub gitignore repository](https://github.com/github/gitignore), and write these patterns to the `STDOUT`. Simply pass in the names of the ignore files you wish to retrieve.

For example,

```shell
getignore get Go.gitignore Global/Vim.gitignore
```

downloads and concatenates the Go and Vim ignore patterns and writes them to `STDOUT`.


Note the `.gitignore` extension on the names optional. Feel free to omit the extension; the previous example could be issued more simply as

```shell
getignore get Go Global/Vim
```

so long as they share a common extension. (See also the `--default-extension` option.)

If you'd like to write the contents directly to a file, use the `-o` option. For example,

```shell
getignore get -o .gitignore Go Global/Vim
```

Would write the contents of the Go and Vim ignore patterns into the `.gitignore` file in the current working directory (`./.gitignore`).

When retrieving many ignore patterns, it can be helpful instead to list names in a file, instead. Given the following file, `names.txt`

```txt
Go
Node
Yeoman
Global/Eclipse
Global/Emacs
Global/JetBrains
Global/Linux
Global/NotepadPP
Global/SublimeText
Global/Tags
Global/TextMate
Global/Vim
Global/Windows
Global/Xcode
Global/macOS
```

we can get all the patterns in this file by passing it via the `--names-file` option

``` shell
getignore get --names-file names.txt
```

Please see the `get` usage via `getignore help get` for explanations of other options available.


### list

Use this command to get a listing of available gitignore patterns files from a remote repository and print the listing to STDOUT. This allows users to standard \*nix tools to manipulate the command's output. For example, the following command line command can be used to download all the "global" gitignore patterns files:

```
getignore list | grep Global/ | xargs getignore get
```

By default, `list` queries the [GitHub .gitignore patterns repository](https://github.com/github/gitignore) using the [GitHub API v3 Trees endpoint](https://developer.github.com/v3/git/trees/). It is possible to pass in a different API URL via the `--api-url` flag, however.

By default, it filters for files that end with the `.gitignore` suffix, however, you can provide an alternative suffix via the `--suffix` flag. To list all files, provide an empty suffix, e.g.

```
getignore list --suffix ''
```


## Building

getignore's dependencies are managed by the [Glide package manager](https://glide.sh/). First install Glide, then build getignore.

``` shell
curl https://glide.sh/get | sh
glide install
go build
```
