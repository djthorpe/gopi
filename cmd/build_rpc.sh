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

PROTOC_GEN_GO=`which protoc-gen-go`
echo 
if [ ! -x "${PROTOC_GEN_GO}" ] ; then
  echo "go get -u github.com/golang/protobuf/protoc-gen-go" >&2
  go get -u github.com/golang/protobuf/protoc-gen-go || exit 1
fi

##############################################################
# Install RPC binaries

COMMANDS=(
  rpc/helloworld_server.go
  rpc/helloworld_client.go
  rpc/rpc_discovery.go
)

echo "go generate github.com/djthorpe/gopi/rpc/protobuf"
go generate -x github.com/djthorpe/gopi/rpc/protobuf || exit 1
for COMMAND in ${COMMANDS[@]}; do
  echo "go install cmd/${COMMAND}"
  go install -ldflags "${LDFLAGS}" -tags "${TAGS}" "cmd/${COMMAND}" || exit -1
done



