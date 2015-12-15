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

# set paths for sources
export X264_FILENAME=`ls -r "${FFMPEG_ROOT}/src" | grep x264`
export FFMPEG_FILENAME=`ls -r "${FFMPEG_ROOT}/src" | grep ffmpeg`

if [ -z "${X264_FILENAME}" ] ; then
  echo "Invalid libx264 path"
  exit -1
fi

if [ -z "${FFMPEG_FILENAME}" ] ; then
  echo "Invalid ffmpeg path"
  exit -1
fi


# build libx264
cd ${FFMPEG_ROOT}/src/${X264_FILENAME}
./configure --host=arm-unknown-linux-gnueabi --enable-shared --enable-pic --disable-opencl --prefix=${FFMPEG_ROOT}
make -j4 && make install

# build ffmpeg
export PKG_CONFIG_PATH="${PKG_CONFIG_PATH}:${FFMPEG_ROOT}/lib/pkgconfig"
cd ${FFMPEG_ROOT}/src/${FFMPEG_FILENAME}
./configure --prefix=${FFMPEG_ROOT} --enable-nonfree --enable-gpl --enable-libx264 --disable-static --enable-shared --extra-cflags="-I${FFMPEG_ROOT}/include" --extra-ldflags="-L${FFMPEG_ROOT}/lib"
make -j4 && make install

