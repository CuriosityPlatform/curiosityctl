#!/bin/bash

# shellcheck disable=SC2046
# shellcheck disable=SC2196
# shellcheck disable=SC2086

set -o errexit

if [ ! -d "$(dirname "$1")" ] ; then
    echo "usage: $(basename "$0") <path-to-cmd-package> <path-to-output-file> <app-name>" 1>&2
    exit 1
fi

CMD_PACKAGE_DIR=$1
EXECUTABLE_PATH=$2
APP_NAME=$3
APP_VERSION=$(git rev-parse HEAD)

GO_SRC_FILES=$(find "$CMD_PACKAGE_DIR" -name "*.go" | tr "\n" " ")

echo_call() {
    echo "$@"
    "$@"
}


echo_call go build -v \
    -o "$EXECUTABLE_PATH" \
    -ldflags="-X main.appID=$APP_NAME -X main.version=$APP_VERSION" \
    $GO_SRC_FILES
