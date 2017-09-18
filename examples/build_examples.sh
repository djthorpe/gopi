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
# install

echo "go install app/helloworld_example"
go install -ldflags "${LDFLAGS}" examples/app/helloworld_example.go || exit -1

echo "go install app/gopi_example"
go install -ldflags "${LDFLAGS}" examples/app/gopi_example.go || exit -1

echo "go install app/log_example"
go install -ldflags "${LDFLAGS}" examples/app/log_example.go || exit -1

echo "go install app/task_example"
go install -ldflags "${LDFLAGS}" examples/app/task_example.go || exit -1

