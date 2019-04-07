[![Report Card](https://goreportcard.com/badge/github.com/durandj/git-ignore)](https://goreportcard.com/report/github.com/durandj/git-ignore)
[![Travis CI](https://travis-ci.org/durandj/git-ignore.svg?branch=master)](https://travis-ci.org/durandj/git-ignore)

# git-ignore

A CLI tool for Git to generate `.gitignore` files for your repository.

One common issue when setting up a new Git repository is that at some
point you need to create a `.gitignore` file to make sure that certain
files aren't ever commited. For example, some generated files (`.o`)
or editor specific files (`.swp`/`.swo`). Writing these files isn't
hard but tends to be error prone (did you forget to include a specific
file? Do you know all the files you need to exclude?). Now there are
some solutions to help with this such as
[gitignore.io](https://gitignore.io) or the
[gitignore project](https://github.com/github/gitignore) but both
options have their own problems. gitignore.io requires you to
copy/paste around the file contents that you need and doesn't support
a ton of different services/languages. The gitignore repository does
support a lot of things you might want to ignore but its annoying to
use since you have to dig through all the files to find what you need.
This project works to solve all of those problems by providing a
simple CLI command that can generate a file with a single line.

`git ignore generate C C++ > .gitignore`

You can see all available options for the `generate` command with the
`list` command.

`git ignore list`

## Install

Installation should be pretty straight forward. Just head on over to
the [releases page](https://github.com/durandj/git-ignore/releases)
and download the latest version for your desired platform, mark it
as executable (if on a Unix type system) and put it somewhere on your
`PATH`.

For example, if I wanted to install version `latest` (this isn't a
real version) for 64 bit Linux I would do:

```
wget https://github.com/durandj/git-ignore/releases/download/vLatest/git-ignore_vLatest_linux_amd64
chmod +x git-ignore_vLatest_linux_amd64
sudo mv git-ignore_vLatest_linux_amd64 /usr/local/bin
```

## Developing

Make sure you first install the following dependencies:

### Dependencies

 * Golang 1.11+
 * [Ginkgo](http://onsi.github.io/ginkgo/)
 * [Taskfile](https://taskfile.org)

### Tasks

Taskfile provides a way of running different scripts easily (similar
to how Makefiles work but better).

To build the code you can just do `task build`.
To run tests run `task test`.
