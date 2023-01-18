#!/bin/bash

cd "$(dirname "$0")"

SCRIPT_DIR="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
OUT_DIR="$SCRIPT_DIR/_out_bin"
OUT_FILE="$OUT_DIR/ivpn"

set -e

# make output dir if not exists
mkdir -p $OUT_DIR

# version info variables
VERSION=""
DATE="$(date "+%Y-%m-%d")"
COMMIT=""

# reading version info from arguments
while getopts ":v:c:" opt; do
  case $opt in
    v) VERSION="$OPTARG"
    ;;
    c) COMMIT="$OPTARG"
    ;;
#    \?) echo "Invalid option -$OPTARG" >&2
#   ;;
  esac
done

if [ -z "$COMMIT" ]; then
  COMMIT="$(git rev-list -1 HEAD)"
fi

echo "======================================================"
echo "============== Compiling IVPN CLI ===================="
echo "======================================================"
echo "Version: $VERSION"
echo "Date   : $DATE"
echo "Commit : $COMMIT"

# Check required GLIBC version. 
# Compiling with the new GLIBC version will not allow the program to start on systems with the old GLIBC (error example: "version 'GLIBC_2.34' not found"). 
# Useful links:
#   https://utcc.utoronto.ca/~cks/space/blog/programming/GoAndGlibcVersioning
# Useful commands:
#   ldd -r -v <binary_file> # check shared libraries dependencies
#
# Info: CLI does not use CGO directly (and can be easily disabled manually), but we use the same build environment as a for daemon (to be able to detect same errors as with daemon binary) 
GLIBC_VER_MAX_REQUIRED="2.31"
GLIBC_VER=$(ldd --version | grep "ldd (" | awk '{print $(NF)}')
if [[ "${GLIBC_VER}" > "${GLIBC_VER_MAX_REQUIRED}" ]]; 
then
    echo "[!] GLIBC version '${GLIBC_VER}' is greater than reqired '${GLIBC_VER_MAX_REQUIRED}'"
    echo "[!] Compiling with the new GLIBC version will not allow the program to start on systems with the old GLIBC."
    read -p "[?] Do you want to continue? [y\n] (N - default): " yn
    case $yn in
      [Yy]* ) ;;
      * ) 
      echo "[!] Build interrupted by user"
      exit 1
      ;;
    esac
fi

cd $SCRIPT_DIR/../../

echo "* updating dependencies..."
go get -v

if [[ "$@" == *"-debug"* ]]
then
    echo "Compiling in DEBUG mode"
    go build -tags debug -o "$OUT_FILE" -trimpath -ldflags "-X github.com/ivpn/desktop-app/daemon/version._version=$VERSION -X github.com/ivpn/desktop-app/daemon/version._commit=$COMMIT -X github.com/ivpn/desktop-app/daemon/version._time=$DATE"
else
    go build -o "$OUT_FILE" -trimpath -ldflags "-s -w -X github.com/ivpn/desktop-app/daemon/version._version=$VERSION -X github.com/ivpn/desktop-app/daemon/version._commit=$COMMIT -X github.com/ivpn/desktop-app/daemon/version._time=$DATE"
fi

echo "Compiled CLI binary: '$OUT_FILE'"

set +e
