#!/bin/bash
# Build script to build dynamic FFMPEG libraries for Raspberry Pi

# Edit this to the location you want the build to occur in
export FFMPEG_ROOT="/opt/ffmpeg"

# set up structure
install -d "${FFMPEG_ROOT}/src" || exit -1
cd "${FFMPEG_ROOT}/src"

# download sources
curl ftp://ftp.videolan.org/pub/videolan/x264/snapshots/last_stable_x264.tar.bz2 | tar xj
curl https://ffmpeg.org/releases/ffmpeg-2.8.3.tar.gz | tar xz
export X264_FILENAME=`ls -r | grep x264`
export FFMPEG_FILENAME=`ls -r | grep ffmpeg`

if [ -z "${X264_FILENAME}" ] ; do
  echo "Invalid libx264 path"
  exit -1
fi

if [ -z "${FFMPEG_FILENAME}" ] ; do
  echo "Invalid ffmpeg path"
  exit -1
fi


# build libx264
cd ${FFMPEG_ROOT}/src/${X264_SRC}
./configure --host=arm-unknown-linux-gnueabi --enable-static --disable-opencl --extra-cflags="-fPIC" --prefix=${FFMPEG_ROOT}
make -j4 && make install

# build ffmpeg
cd ${FFMPEG_ROOT}/src/${FFMPEG_SRC}
./configure --prefix=${FFMPEG_ROOT} --enable-nonfree --enable-gpl --enable-libx264 --enable-static --extra-cflags="-I${FFMPEG_ROOT}/include -fPIC" --extra-ldflags="-L${FFMPEG_ROOT}/lib"
make -j4 && make install

