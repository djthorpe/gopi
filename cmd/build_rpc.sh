#!/bin/bash
##############################################################
# Build RPC binaries
##############################################################

CURRENT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GO=`which go`
PROTOC=`which protoc`
LDFLAGS="-w -s"
TAGS="rpi"
cd "${CURRENT_PATH}/.."

##############################################################
# Sanity checks

if [ ! -d ${CURRENT_PATH} ] ; then
  echo "Not found: ${CURRENT_PATH}" >&2
  exit -1
fi
if [ "${GO}" == "" ] || [ ! -x ${GO} ] ; then
  echo "go not installed or executable" >&2
  exit -1
fi

##############################################################
# Build protobuf go extension

RPC_EXAMPLES=""
if [ "${PROTOC}" != "" ] && [ -x ${PROTOC} ] ; then
  RPC_EXAMPLES="1"
fi

##############################################################
# Install RPC binaries

COMMANDS=(
  rpc/rpc_server.go
  rpc/rpc_client.go
)

if [ "${RPC_EXAMPLES}X" != "X" ] ; then
  echo "go generate github.com/djthorpe/gopi/protobuf"
  go generate -x github.com/djthorpe/gopi/protobuf || exit 1
  echo "go get -u github.com/golang/protobuf/protoc-gen-go" >&2
  go get -u github.com/golang/protobuf/protoc-gen-go || exit 1
  for COMMAND in ${COMMANDS[@]}; do
    echo "go install cmd/${COMMAND}"
    go install -ldflags "${LDFLAGS}" -tags "${TAGS}" "cmd/${COMMAND}" || exit -1
  done
fi


