# go versions

govs. A tool for installing and managing multiple go versions. A cross platform
solution written in go.

## Overview

govs can be used to manage multiple go versions on your system. It is
designed to be simple with a very small api. It implements the [recommended
approach](https://go.dev/doc/manage-install) for managing multiple installations
according the go team with some convient extras.

## Extras

Go's recommended approach does come with some minor pain points that we aim to
improve:

1. The first is the requirement for go to already be installed on the
machine. This causes issues when attempting to symlink an older version to the
`go` command. By writing this tool in go, we no longer need it as a dependency.

2. The second is that symlinks need to be handled manually. It is easy to forget
the exact syntax and paths if you don't switch versions very often. Our `govs
set` command simplifies this process.

3. The third is that the recommended approach requires manually deleting old
versions. We now have a `govs remove` command to handle that.

## Installation

1. Download and extract the relevant binary for your system
2. Make it executable `chmod +x govs` ?? (is this needed with a go binary?)
3. Add `$GOPATH/bin` to your path

If you already have multiple go versions installed using the recommended
approach, they will be automatically detected by govs.

## API

### get

Install and set a specific go version as the default with `govs get {version}`.
Eg: `govs get 1.18.5`. This will `install` and `set` in a single command.

### install

Install a specific go version with `govs install {version}. Eg: `govs install
1.18.5`. This will add a version specific binary to your path. Eg: `go1.18.5`.

### set

Set a default go version with `govs set {version}`. Eg: `govs set 1.18.5`. This
will set 1.18.5 as the default version when using the `go` command directly.

### remove

Remove a specific go version with `govs remove {version}. Eg: `govs remove
1.18.5`. This will remove the version from your path.

### list

List the currently installed go versions with `govs list`.

### versions

List all available go versions with `govs versions`.

## Details

### GOPATH

If you have a `$GOPATH` set, it will be used. If you do not have it set,
`$HOME/go` will be the default. This follows the convention set by the go
team.

### Source

Source code is installed to `$HOME/sdk/go{version}`.

### Binaries

Binaries are installed to `$GOPATH/bin/go{version}`.

Each version can be called directly using it's version specific executable name.
Eg `go1.18.5`.

Use the `govs set` command to symlink a specific version to the `$GOPATH/bin/go`
executable.

## Uninstall govs

### Keeping go versions

If you want to remove govs while keeping your go versions and symlinks intact,
simply remove the binary. You will need to manually manage your own versions
going forward.

1. Remove the binary

### Removing go versions

If you want to remove go completely, you can use govs to delete all go versions
before removing govs itself.

1. Remove all go versions: `govs versions | xargs govs remove`
2. Remove the binary
