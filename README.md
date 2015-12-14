# gopi

## Introduction

This repository contains Raspberry Pi Go Language Experiments. In order to build the examples, use:

```
go install github.com/djthorpe/gopi/examples
```

This will put the binary programs in your ```${GOPATH}/bin``` subdirectory. You can then run the programs from there.

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
  cd ${FFMPEH_ROOT}/src

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



