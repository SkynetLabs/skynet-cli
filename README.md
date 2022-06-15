# skynet-cli

[![Go](https://img.shields.io/github/go-mod/go-version/SkynetLabs/skynet-cli)](https://github.com/SkynetLabs/skynet-cli)
[![Build Status](https://img.shields.io/github/workflow/status/SkynetLabs/skynet-cli/Go)](https://github.com/SkynetLabs/skynet-cli/actions)
[![Contributors](https://img.shields.io/github/contributors/SkynetLabs/skynet-cli)](https://github.com/SkynetLabs/skynet-cli/graphs/contributors)
[![License](https://img.shields.io/github/license/SkynetLabs/skynet-cli)](https://github.com/SkynetLabs/skynet-cli)

skynet-cli is a lightweight CLI to interact with [Skynet](https://siasky.net).

Skynet is the decentralized CDN and file sharing platform for devs and the
storage foundation for a Free Internet!

## Installing

The following methods will install a binary called `skynet` to your machine.

### Using Go Install

If you have [Go](https://golang.org/cmd/go/) installed you can run:

```
go install github.com/SkynetLabs/skynet-cli/v2/cmd/skynet@latest
```

### Homebrew (MacOS Users)

If you have [Homebrew](https://brew.sh/) installed you can run:

```shell
brew tap SkynetLabs/skynet-cli https://github.com/SkynetLabs/skynet-cli.git
brew install skynet-cli
```

### Downloading Release Binary

Alternatively, you can pull the appropriate binary from our [Releases](https://github.com/SkynetLabs/skynet-cli/releases) page.

### Building From Source

You can build the source code yourself with:

```
make release
```

## Usage

skynet-cli is designed to be simple and easy to use, just like Skynet.

Uploading a file or directory is a simple as using the following command:

```shell
skynet upload [source path]
```

This will return a `skylink`, which then in turn can be used to download the
file with the following command:

```shell
skynet download [skylink] [destination]
```

## Documentation

For comprehensive documentation complete with examples, please see [the Skynet SDK docs](https://siasky.net/docs/?shell--cli#introduction).

To learn about additional commands you can always use the `-h` flag or check out
the [documentation](./doc).
