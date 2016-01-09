# gopi

## Introduction

This repository contains Raspberry Pi Go Language Experiments. In order to retrieve the source code, use:

```
go get github.com/djthorpe/gopi
```

There is a single `rpi` module and several submodules:

  * `rpi` - Contains code for interfacing with the Raspberry Pi hardware
  * `dispmanx` - Low-level VideoCore interface
  * `egl` - Native interface to link OpenGL ES and OpenVG to the GPU
  * `gles` - OpenGL ES for rendering 3D
  * `vg` - OpenVG for rendering 2D vector graphics 
  * `openmax` - OpenMAX Media Library
  * `gpio` - Interface to the General Purpose IO connector

Most of these are still to be written or completed. There are a set of examples
of using these in the `examples` folder.

## Links

Please see the following locations for more information:

  * [How to write Go Code](http://golang.org/doc/code.html) in order to work out how to structure your Go folder
  

## Building ffmpeg

In order to build ffmpeg with libx264 for your Raspberry Pi, you can use the 
following command line sequence:

```  
  export FFMPEG_ROOT="/opt/ffmpeg"
  export PKG_CONFIG_PATH="${FFMPEG_ROOT}/lib/pkgconfig"
  
  # set up structure
  install -d ${FFMPEG_ROOT}/src
  cd ${FFMPEG_ROOT}/src

  # download sources
  curl ftp://ftp.videolan.org/pub/videolan/x264/snapshots/last_stable_x264.tar.bz2 | tar xj
  curl https://ffmpeg.org/releases/ffmpeg-2.8.3.tar.gz | tar xz  
  export X264_SRC=`ls -r ${FFMPEG_ROOT}/src | grep x264`
  export FFMPEG_SRC=`ls -r ${FFMPEG_ROOT}/src | grep ffmpeg`

  # build libx264
  cd ${FFMPEG_ROOT}/src/${X264_SRC}
  ./configure --host=arm-unknown-linux-gnueabi --enable-static --disable-opencl --extra-cflags="-fPIC" --prefix=${FFMPEG_ROOT}
  make -j4 && make install

  # build ffmpeg
  cd ${FFMPEG_ROOT}/src/${FFMPEG_SRC}
  ./configure --prefix=${FFMPEG_ROOT} --enable-nonfree --enable-gpl --enable-libx264 --enable-static --extra-cflags="-I${FFMPEG_ROOT}/include -fPIC" --extra-ldflags="-L${FFMPEG_ROOT}/lib"
  make -j4 && make install
  
```

The resulting binaries and libraries will be under `/opt/ffmpeg` or whereever you
indicated `${FFMPEG_ROOT}` should be. You'll then need to add the location to
your `${PATH}` variable:

```
  export FFMPEG_ROOT="/opt/ffmpeg"
  export PATH="${PATH}:${FFMPEG_ROOT}/bin"
```

## Commands to set up a Raspberry Pi

Start with the following commands to update your OS, then run through 
expanding SD card space, changing the password for the 'pi' user, etc:

```
sudo apt-get update
sudo apt-get upgrade
sudo raspi-config
```

This will then require a reboot. Then add a user for yourself, and add to the relevant groups:

```
sudo useradd -m -g pi -G adm,sudo,audio,video,gpio,input,i2c,spi -s /bin/bash <USERNAME>
sudo passwd <USERNAME>
```

You should then logout and login again as your own user.

## Setting up your golang environment

Please see here http://dave.cheney.net/2015/09/04/building-go-1-5-on-the-raspberry-pi
for more information. Here are the commands you can run to have Go installed in a
`/opt/go` folder:

```
cd $HOME
install -d go/bin
install -d go/src
install -d go/pkg
curl http://dave.cheney.net/paste/go-linux-arm-bootstrap-c788a8e.tbz | tar xj
sudo install -o $USER -d /opt/go
curl https://storage.googleapis.com/golang/go1.5.src.tar.gz | tar xz -C /opt
ulimit -s 1024
cd /opt/go/src
env GO_TEST_TIMEOUT_SCALE=10 GOROOT_BOOTSTRAP=$HOME/go-linux-arm-bootstrap ./all.bash

```

Once this is completed (it can take a few hours) you can add the following lines
to your `~/.bash_profile` file:

```
# Raspberry Pi
export PIROOT="/opt/vc"
export PATH="${PATH}:${PIROOT}/bin"

# Go Language
export GOROOT="/opt/go"
export GOPATH="${HOME}/go"
export GOBIN="${GOPATH}/bin"
export PATH="${GOROOT}/bin:${GOPATH}/bin:${PATH}"
```




 

