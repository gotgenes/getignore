# getignore

A small executable to download and concatenate [.gitignore templates from GitHub](https://github.com/github/gitignore) into a local .gitignore file.

[![Travis CI Build Status](https://travis-ci.org/gotgenes/getignore.svg?branch=master)](https://travis-ci.org/gotgenes/getignore)

## Usage

getignore supports the following commands:

* [`get`](#get)
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

Use the `get` command to obtain gitignore patterns from remote repositories. By default, it will obtain patterns from the [GitHub gitignore repository](https://github.com/github/gitignore), and write these patterns to the `.gitignore` file in the current working directory. Simply pass in the names of the ignore files you wish to retrieve.

For example,

```shell
getignore get Go.gitignore Global/Vim.gitignore
```

downloads and concatenates the Go and Vim ignore patterns and writes them into the `.gitignore` file in the current working directory (`./.gitignore`).

Note the `.gitignore` extension on the names optional. Feel free to omit the extension; the previous example could be issued more simply as

```shell
getignore get Go Global/Vim
```

so long as they share a common extension. (See also the `--default-extension` option.)


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


## Building

getignore's dependencies are managed by the [Glide package manager](https://glide.sh/). First install Glide, then build getignore.

``` shell
curl https://glide.sh/get | sh
glide install
go build
```
