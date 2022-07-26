#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/merlion/${BINARY:-merliond}
ID=${ID:-0}
LOG=${LOG:-merliond.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'merliond' E.g.: -e BINARY=merliond_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export MERLION_HOME="/merlion/node${ID}/merliond"

if [ -d "$(dirname "${MERLION_HOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${MERLION_HOME}" --trace "$@" | tee "${MERLION_HOME}/${LOG}"
else
  "${BINARY}" --home "${MERLION_HOME}" --trace "$@"
fi
