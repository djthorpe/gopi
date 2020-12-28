# Building and Using

## Dependencies: Debian

These are the commands you should run on Debian to install libraries needed:

```bash
apt install libavfilter-dev libavcodec-dev libavformat-dev libavutil-dev libswscale-dev libswresample-dev
apt install libdrm-dev libegl-dev libgbm-dev libgl-dev libgles-dev
apt install libpulse-dev
apt install libchromaprint1
apt install protobuf-compiler
```

## Dependencies: Macintosh

It is assumed on Macintosh you are using [Homebrew](https://brew.sh/) in order to do package management.
There is no currently supported graphics on Macintosh:

```bash
brew install ffmpeg
brew install pulseaudio
brew install chromaprint
brew install protobuf
```

