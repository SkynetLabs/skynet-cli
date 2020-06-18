#!/usr/bin/env bash
set -e

# version and keys are supplied as arguments
version="$1"
rc=`echo $version | awk -F - '{print $2}'`
if [[ -z $version ]]; then
	echo "Usage: $0 VERSION"
	exit 1
fi

# setup build-time vars
ldflags="-s -w -X 'github.com/NebulousLabs/skynet-cli/build.GitRevision=`git rev-parse --short HEAD`' -X 'github.com/NebulousLabs/skynet-cli/build.BuildTime=`git show -s --format=%ci HEAD`' -X 'github.com/NebulousLabs/skynet-cli/build.ReleaseTag=${rc}'"

function build {
  os=$1
  arch=$2

	echo Building ${os}...
	# create workspace
  subFolder=skynet-cli-$version-$os-$arch
  folder=release/$subFolder
	rm -rf $folder
	mkdir -p $folder
	# compile and hash binaries
	bin=skynet
	if [ "$os" == "windows" ]; then
		bin=skynet.exe
	fi
	GOOS=${os} GOARCH=${arch} go build -a -tags 'netgo' -trimpath -ldflags="$ldflags" -o $folder/$bin ./cmd/skynet
	(
		cd release/
		sha256sum $subFolder/$bin >> skynet-cli-$version-SHA256SUMS.txt
	)

	cp -r doc LICENSE README.md $folder
}

# Build amd64 binaries.
for os in darwin linux windows; do
  build "$os" "amd64"
done

# Build Raspberry Pi binaries.
build "linux" "arm64"
