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

echo "go install helloworld/helloworld_example"
go install -ldflags "${LDFLAGS}" examples/helloworld/helloworld_example.go || exit -1

echo "go install app/log_example"
go install -ldflags "${LDFLAGS}" examples/app/log_example.go || exit -1

echo "go install display/display_example"
go install -ldflags "${LDFLAGS}" examples/display/display_example.go || exit -1

echo "go install input/input_example"
go install -ldflags "${LDFLAGS}" examples/input/input_example.go || exit -1

echo "go install gpio/gpioctrl"
go install -ldflags "${LDFLAGS}" examples/gpio/gpioctrl.go || exit -1

echo "go install gpio/ledflash"
go install -ldflags "${LDFLAGS}" examples/gpio/ledflash.go || exit -1

echo "go install i2c/i2cdetect"
go install -ldflags "${LDFLAGS}" examples/i2c/i2cdetect.go || exit -1

echo "go install egl/snapshot"
go install -ldflags "${LDFLAGS}" examples/egl/snapshot.go || exit -1

echo "go install egl/cursor_example"
go install -ldflags "${LDFLAGS}" examples/egl/cursor_example.go || exit -1

echo "go install egl/image_example"
go install -ldflags "${LDFLAGS}" examples/egl/image_example.go || exit -1

echo "go install openvg/circle_example"
go install -ldflags "${LDFLAGS}" examples/openvg/circle_example.go || exit -1

echo "go install vgfont/font_example"
go install -ldflags "${LDFLAGS}" examples/vgfont/font_example.go || exit -1

echo "go install vgfont/vgfont_dx_example"
go install -ldflags "${LDFLAGS}" examples/vgfont/vgfont_dx_example.go || exit -1
