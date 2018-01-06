#!/bin/bash
##############################################################
# Build Raspberry Pi Flavours
##############################################################

CURRENT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GO=`which go`
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
# Install

COMMANDS=(
    helloworld/helloworld.go
    timer/timer_tester.go
    gpio/gpio_ctrl.go
    i2c/i2c_detect.go
    spi/spi_ctrl.go
    input/input_tester.go    
    hw/hw_list.go
    hw/display_list.go
    rpc/rpc_discovery.go
)

for COMMAND in ${COMMANDS[@]}; do
    echo "go install cmd/${COMMAND}"
    go install -ldflags "${LDFLAGS}" -tags "${TAGS}" "cmd/${COMMAND}" || exit -1
done


