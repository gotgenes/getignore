# getignore

A small executable to concatenate .gitignore templates from [github](https://github.com/github/gitignore) into a local .gitignore file.

## Usage

Concatenate the Perl and Python templates into a .gitignore file in the current working directory

``` shell
./getignore Perl Python
```

Read the list of templates from a file, `templates.txt`

``` txt
Ruby
Qt
Go
```

``` shell
./getignore -file=templates.txt
```

## Notes

This command will overwrite any existing .gitignore file in the current working directory
