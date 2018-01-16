#!/bin/bash
##############################################################
# Build Raspberry Pi Flavours
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
# Install

COMMANDS=(
    helloworld/helloworld.go
    timer/timer_tester.go
    hw/hw_list.go
    hw/vcgencmd_list.go
    hw/display_list.go
    gpio/gpio_ctrl.go
    i2c/i2c_detect.go
    spi/spi_ctrl.go
    input/input_tester.go
    lirc/lirc_receive.go        
)

for COMMAND in ${COMMANDS[@]}; do
    echo "go install cmd/${COMMAND}"
    go install -ldflags "${LDFLAGS}" -tags "${TAGS}" "cmd/${COMMAND}" || exit -1
done



##############################################################
# RPC EXAMPLES

COMMANDS=(
  rpc/rpc_server.go
)

if [ "${RPC_EXAMPLES}X" != "X" ] ; then
  echo "go get -u github.com/golang/protobuf/protoc-gen-go" >&2
  go get -u github.com/golang/protobuf/protoc-gen-go || exit 1
  for COMMAND in ${COMMANDS[@]}; do
    echo "go generate cmd/${COMMAND}"
    go generate "cmd/${COMMAND}" || exit -1    
    echo "go install cmd/${COMMAND}"
    go install -ldflags "${LDFLAGS}" -tags "${TAGS}" "cmd/${COMMAND}" || exit -1
  done
fi


