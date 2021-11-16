# getignore

getignore bootstraps `.gitignore` files from [GitHub gitignore patterns](https://github.com/github/gitignore).

[![CI workflow status](https://github.com/gotgenes/getignore/actions/workflows/ci.yaml/badge.svg)](https://github.com/gotgenes/getignore/actions/workflows/ci.yaml)


## Installation

### Homebrew users (macOS / OS X)

[Homebrew](http://brew.sh) users can install getignore using the following commands:

```shell
brew tap gotgenes/homebrew-gotgenes
brew update
brew install getignore
```

### Windows, Linux, and macOS / OS X

Download and unpack a pre-compiled executable from [the releases page](https://github.com/gotgenes/getignore/releases).
Make sure to place the executable in your shell's `PATH`.


### Other platforms

See [the Building section, below](#building).


## Usage

getignore supports the following commands:

* [`help`](#help)
* [`get`](#get)
* [`list`](#list)


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

Use the `get` command to obtain gitignore patterns from remote repositories, concatenate their contents, and output them.
Simply pass in the names of the ignore files you wish to retrieve.

For example,

```shell
getignore get Go.gitignore Global/Vim.gitignore
```

downloads and concatenates the Go and Vim ignore patterns and writes them to `STDOUT`.

Note the `.gitignore` extension on the names optional.
Feel free to omit the extension; the previous example could be issued more simply as

```shell
getignore get Go Global/Vim
```

`get` will add the `.gitignore` extension for you when retrieving the files.
You can use the `--suffix` flag to choose a different default suffix.
If you want no suffix added, pass the empty string (`--suffix ''`).

By default, `get` downloads the files from the [GitHub gitignore patterns repository](https://github.com/github/gitignore) using the [GitHub API v3 Trees endpoint](https://developer.github.com/v3/git/trees/).
You can use a different owner, repository name, branch, or combination of all of them via the respective `--owner`, `--repository`, and `--branch` flags.
It is also possible to pass in a different API URL via the `--base-url` flag.

By default, `get` writes the contents to `STDOUT`.
If you'd like to write the contents directly to a file, you can use the `-o` option.
For example,

```shell
getignore get -o .gitignore Go Global/Vim
```

Would write the contents of the Go and Vim ignore patterns into the `.gitignore` file in the current working directory (`./.gitignore`).

When retrieving many ignore patterns, it can be helpful instead to list names in a file, instead.
Suppose we create a file `names.txt` with the following contents:

```txt
Go
Node
Global/Vim
Global/macOS
```

We can get all the patterns in this file by passing it in to `get` via the `--names-file` option

```shell
getignore get --names-file names.txt
```

Please see the `get` usage via `getignore help get` for explanations of other options available.


### list

Use this command to get a listing of available gitignore patterns files from a remote repository and print the listing to `STDOUT`.
This allows users to use standard command line tools to manipulate the command's output.
For example, the following command line command can be used to download all the "global" gitignore patterns files:

```
getignore list | grep Global/ | xargs getignore get
```

By default, `list` queries the [GitHub gitignore patterns repository](https://github.com/github/gitignore) using the [GitHub API v3 Trees endpoint](https://developer.github.com/v3/git/trees/).
You can use a different owner, repository name, branch, or combination of all of them via the respective `--owner`, `--repository`, and `--branch` flags.
It is possible to pass in a different API URL via the `--base-url` flag.

By default, `list` filters for files that end with the `.gitignore` suffix, however, you can provide an alternative suffix via the `--suffix` flag.
Alternatively, to list all files in the repository, regardless of suffix, provide an empty string as the value, e.g.

```
getignore list --suffix ''
```


## Completion

getignore supports completion of the command line for [Bash](completions/bash/getignore-completion.bash) and [zsh](completions/zsh/_getignore). If completions were not installed by default, please place the respective completion file in the appropriate location for completion scripts on your system.


## Building

```shell
make
```


## Testing

Ensure you have testing dependencies with

```shell
make dev-install
```

Then run the tests with

```shell
make test
```
