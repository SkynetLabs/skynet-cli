# Release Scripts

This directory contains all the scripts used to build releases of skynet-cli.

`./release.sh` compiles `skynet` binaries for each supported system. It places
those binaries along with documentation into a separate directory for each
system type. It also creates a file of SHA256 hashes containing the hash of each
binary produced.

`./package.sh` zips the directories created by `release.sh` and adds their
SHA256 hashes to the checksums file.

These procedures are split out so that compilation can be done in a more
reproducible environment and reduce the risk of adding non-determinism into the
compilation process. Deterministic builds are desirable so that the binaries
produced and signed as part of the release process can be verified by any
contributor as having been produced using the expected/correct source code.

`./build-in-docker.sh`  uses the reproducible build environment defined in the
`Dockerfile`. The `Dockerfile` uses a base image with just the prerequisites for
Go compilation pre-installed. This script creates a docker container which
checks out a specified branch of this repository and then runs `./release.sh`.
The outputs are copied out of the container before it is deleted.

Currently this process can only be reproduced and verified on `linux/amd64`
systems. That is, anyone with a `linux/amd64` system can reproduce the binaries
created for **all** systems using this script. In the future this can be
improved to allow reproducibility on more systems.

## Example
The following list of commands shows how one could use the scripts in this
directory to build skynet-cli.

1. `./build-in-docker.sh master v1.1.0`
3. `cd ../release && sha256sum --check skynet-cli-v1.1.0-SHA256SUMS.txt`
3. `gpg --clearsign --armor skynet-cli-v1.1.0-SHA256SUMS.txt`
