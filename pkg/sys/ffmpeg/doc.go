/*

This package provides ffmpeg bindings, targetting new versions
of the ffmpeg API.  In order to use this package, you will need
to install the ffmpeg development libraries. On Darwin (Mac)
with Homebrew installed:

	brew install ffmpeg

For Debian, the following is sufficient:

	sudo apt install libavcodec-dev libavformat-dev libavdevice-dev

You will also need to use -tags ffmpeg when testing, building or
installing.

References:
https://ffmpeg.org/doxygen/trunk/index.html

*/
package ffmpeg
