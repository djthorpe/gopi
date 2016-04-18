# Introduction

This repository contains Raspberry Pi Go Language Experiments. Firstly, you'll need to build Go on
your Raspberry Pi, please see an Appendix on order to see the instructions for that. In order to 
retrieve the source code, use:

```
go get github.com/djthorpe/gopi
```

Once you have your Go environment setup...TODO

There is a single `rpi` module and several submodules:

  * `rpi` - Contains code for interfacing with the Raspberry Pi hardware
  * `dispmanx` - Low-level VideoCore hardware interface
  * `egl` - Native interface to link OpenGL and OpenVG to the hardware
  * `gles` - OpenGL ES for rendering 3D
  * `vg` - OpenVG for rendering 2D vector graphics 
  * `omx` - OpenMAX Media Library
  * `gpio` - Interface to the General Purpose IO connector

Most of these are still to be written or completed. There are a set of examples
of using these in the `examples` folder.

# Running the example programs

# Appendices

## Appendix: Links

Please see the following locations for more information:

  * [Building Go on your Raspberry Pi](http://dave.cheney.net/2015/09/04/building-go-1-5-on-the-raspberry-pi) to
    get your Go environment set-up
  * [How to write Go Code](http://golang.org/doc/code.html) in order to work out how to structure your Go folder
 
## Appendix: Building golang

In order to get Go working on your Raspberry Pi, you'll need to "bootstrap" it from an existing Go binary, since
the package is compiled using Go itself. You can use the following command line sequence:

```
  export GO_ROOT="/opt/go"
  
  # set up the structure
  install -d ${GO_ROOT}/build
  cd ${GO_ROOT}/build
  curl https://storage.googleapis.com/golang/go1.5.src.tar.gz | tar xz
  ulimit -s 1024
  export GO_TEST_TIMEOUT_SCALE=10
  export GOROOT_BOOTSTRAP=${GO_ROOT}/build/go
  cd ${GOROOT_BOOTSTRAP}/src
  . ./all.bash
```

TODO

## Appendix: Building ffmpeg

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



