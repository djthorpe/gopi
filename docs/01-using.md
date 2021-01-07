# Building and Using

This section describes how you might download and integrate units
into your own code. 

## Dependencies

Some units are either platform-dependent or 
dependent on libraries and tools being available. You can satisfy
these dependencies by running the commands as indicated below.

### Debian

These are the commands you should run on Debian to install libraries needed:

```bash
apt install libavfilter-dev libavcodec-dev libavformat-dev libavutil-dev libswscale-dev libswresample-dev
apt install libdrm-dev libegl-dev libgbm-dev libgl-dev libgles-dev
apt install libpulse-dev
apt install libchromaprint1
apt install protobuf-compiler
```

### Macintosh

It is assumed on Macintosh you are using [Homebrew](https://brew.sh/) in order to do package management.
There is no currently supported graphics on Macintosh:

```bash
brew install ffmpeg
brew install pulseaudio
brew install chromaprint
brew install protobuf
```

## Building Examples

There are some examples in the `cmd` folder which demonstrate features of the units, some of which can be built for Raspberry Pi and some which can also be build for other Linux operating systems and Macintosh. Building currently uses the `Makefile`
with the following targets:

  * `make all` will make all the examples;
  * `make clean` will clean the build caches;
  * `make debian` will generate `.deb` packages for the examples;
  * `make test` will run tests in the `pkg` folder.

Build outpput is placed in a temporary `build` folder in the repository. Cross-compiling to other platforms is generally not
supported due to the complexities of binding with libraries,
but where no binding is necessary, you can build for ARM and x64 processors as follows:

```
GOBIN="" GOOS=linux GOARCH=arm make debian 
GOBIN="" GOOS=linux GOARCH=amd64 make debian 
```



