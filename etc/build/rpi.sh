#!/bin/bash
##############################################################
# RPI BUILD SCRIPT
##############################################################

CURRENT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GO=`which go`
LDFLAGS="-w -s"
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
# go get dependencies

##############################################################
# install

echo "go install helloworld"
go install -ldflags "${LDFLAGS}" cmd/helloworld.go || exit -1

echo "go install vcgencmd"
go install -ldflags "${LDFLAGS}" cmd/vcgencmd.go || exit -1

#echo "go install gpio"
#go install -ldflags "${LDFLAGS}" cmd/gpio.go || exit -1

#echo "go install i2c"
#go install -ldflags "${LDFLAGS}" cmd/i2c.go || exit -1

echo "go install dispmanx"
go install -ldflags "${LDFLAGS}" cmd/dispmanx.go || exit -1

echo "go install touchscreen"
go install -ldflags "${LDFLAGS}" cmd/touchscreen.go || exit -1


