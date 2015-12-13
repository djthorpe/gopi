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

  # set up structure
  mkdir ${FFMPEG_ROOT}
  cd ${FFMPEG_ROOT}
  mkdir src
  cd src

  # download sources
  curl ftp://ftp.videolan.org/pub/videolan/x264/snapshots/last_stable_x264.tar.bz2 | tar xj
  curl https://ffmpeg.org/releases/ffmpeg-2.8.3.tar.gz | tar xz  
  export X264_FILENAME=`ls -r | grep x264`; echo "X264=${X264_FILENAME}"
  export FFMPEG_FILENAME=`ls -r | grep ffmpeg`; echo "FFMPEG=${FFMPEG_FILENAME}"

  # build libx264
  cd ${FFMPEG_ROOT}/src/${X264_FILENAME}
  ./configure --host=arm-unknown-linux-gnueabi --enable-static --disable-opencl --prefix=${FFMPEG_ROOT}
  make && make install

  # build ffmpeg
  cd ${FFMPEG_ROOT}/src/${FFMPEG_FILENAME}
  ./configure --prefix=${FFMPEG_ROOT} --enable-nonfree --enable-gpl --enable-libx264 --enable-shared --enable-static --extra-cflags="-I${FFMPEG_ROOT}/include" --extra-ldflags="-L${FFMPEG_ROOT}/lib"
  make && make install
  
```

The resulting binaries and libraries will be under `/opt/ffmpeg` or whereever you
indicated `${FFMPEG_ROOT}` should be. You'll then need to add the location to
your `${PATH}` variable:

```
  export FFMPEG_ROOT="/opt/ffmpeg"
  export PATH="${PATH}:${FFMPEG_ROOT}/bin"
```



