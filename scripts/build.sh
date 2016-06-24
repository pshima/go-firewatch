#!/bin/bash
#
# This script builds the application from source for multiple platforms.
set -e

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Get the name of the binary
NAME="$(basename "$DIR")"

# Change into that directory
cd "$DIR"

# Get the git commit
if [ "${GITHUB_COMMIT}x" != "x" ]; then
  GIT_COMMIT="${GITHUB_COMMIT}"
else
  GIT_COMMIT=$(git rev-parse HEAD)
fi

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"amd64"}
XC_OS=${XC_OS:-"darwin linux windows"}

# Delete the old dir
echo "==> Removing old directory..."
rm -f bin/*
rm -rf pkg/*
mkdir -p bin/

# If its dev mode, only build for ourself
if [ "${DEV}x" != "x" ]; then
  XC_OS=$(go env GOOS)
  XC_ARCH=$(go env GOARCH)
fi

# Build!
echo "==> Building ${NAME}..."
gox \
  -os="${XC_OS}" \
  -arch="${XC_ARCH}" \
  -ldflags "-X main.GitCommit ${GIT_COMMIT:0:7}" \
  -output "pkg/{{.OS}}_{{.Arch}}/${NAME}" \
  .

# Move all the compiled things to the $GOPATH/bin
GOPATH=${GOPATH:-$(go env GOPATH)}
case $(uname) in
    CYGWIN*)
        GOPATH="$(cygpath $GOPATH)"
        ;;
esac
OLDIFS=$IFS
IFS=: MAIN_GOPATH=($GOPATH)
IFS=$OLDIFS

# Copy our OS/Arch to the bin/ directory
DEV_PLATFORM="./pkg/$(go env GOOS)_$(go env GOARCH)"
for F in $(find ${DEV_PLATFORM} -mindepth 1 -maxdepth 1 -type f); do
    cp ${F} bin/
    cp ${F} ${MAIN_GOPATH}/bin/
done

if [ "${DEV}x" = "x" ]; then
  echo "==> Packaging..."
  for PLATFORM in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
    OSARCH=$(basename ${PLATFORM})
    echo "--> ${OSARCH}"

    pushd $PLATFORM >/dev/null 2>&1
    tar -cvzf "../${OSARCH}.tar.gz" ./*
    popd >/dev/null 2>&1
  done
fi

# Done!
echo
echo "==> Results:"
ls -hl bin/
