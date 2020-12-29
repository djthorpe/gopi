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
