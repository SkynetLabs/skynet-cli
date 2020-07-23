# skynet-cli

![](https://github.com/NebulousLabs/skynet-cli/workflows/Go/badge.svg)

skynet-cli is a lightweight CLI to interact with [Skynet](https://siasky.net).

Skynet is the decentralized CDN and file sharing platform for devs and the
storage foundation for a Free Internet!

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

For comprehensive documentation complete with examples, please see [the Skynet SDK docs](https://nebulouslabs.github.io/skynet-docs/?shell--cli#introduction).

To learn about additional commands you can always use the `-h` flag or check out
the [documentation](./doc).
