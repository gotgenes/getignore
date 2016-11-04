# getignore

A small executable to download and concatenate .gitignore templates from [github](https://github.com/github/gitignore) into a local .gitignore file.


## Usage

Concatenate the Perl and Python templates into a .gitignore file in the current working directory

``` shell
getignore get Perl Python
```

Read the list of names for getignore files from a file, `names.txt`

``` txt
Ruby
Qt
Go
```

``` shell
getignore get --names-file names.txt
```


## Building

getignore's dependencies are managed by the [Glide package manager](https://glide.sh/). First install Glide, then build getignore.

``` shell
curl https://glide.sh/get | sh
glide install
go build
```
