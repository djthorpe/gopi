// DVB (Digital Video Broadcasting) bindings for Go
//
// When you depend on this code, use -tags dvb when
// running go build test and install. In order to install
// scan tables for multiplexers, use:
//
// apt install dtv-scan-tables
//
package dvb

// Ref: https://www.kernel.org/doc/html/v4.12/media/uapi/dvb/dvbapi.html

/* DVB-T Parameters:
DTV_FREQUENCY
DTV_MODULATION
DTV_BANDWIDTH_HZ
DTV_INVERSION
DTV_CODE_RATE_HP
DTV_CODE_RATE_LP
DTV_GUARD_INTERVAL
DTV_TRANSMISSION_MODE
DTV_HIERARCHY
DTV_BANDWIDTH_HZ
DTV_STREAM_ID
DTV_LNA
*/

/*
DVB-C Parameters
DTV_FREQUENCY
DTV_MODULATION
DTV_INVERSION
DTV_SYMBOL_RATE
DTV_INNER_FEC
DTV_LNA
*/

/*

How to determine the API version?

  Check in your configure script for #include <linux/dvb/version.h>,
  include it and check the DVB_API_VERSION #define.

  Currently we use version 3, it will be incremented whenever an API change
  meets the CVS main branch.

-------------------------------------------------------------------------------

What is a demultiplexer?

  The demultiplexer in your DVB system cares about the routing of an MPEG2
  stream you feed into the DVB adapter either by read/write system calls,
  by using stream inputs of the demultiplexer or the frontend(s).

  Using the demux API you can set up the stream routing and set up filters to
  filter the interesting parts of the input streams only.

-------------------------------------------------------------------------------

I have set up the frontend and now I want to see some video!
What do I have to do?

  When you have an MPEG video decoder on board you can set up the demultiplexer
  to feed the related PES packets into the MPEG decoder:

	#include <linux/dvb/dmx.h>

	struct dmx_pes_filter_params pesfilter;

	if ((fd = open("/dev/dvb/adapter0/demux0", O_RDWR) < 0) {
		perror ("open failed");
		return -1;
	}

	pesfilter.pid = pid;
	pesfilter.input = DMX_IN_FRONTEND;
	pesfilter.output = DMX_OUT_DECODER;
	pesfilter.pes_type = DMX_PES_VIDEO;
	pesfilter.flags = DMX_IMMEDIATE_START;

	if (ioctl(fd, DMX_SET_PES_FILTER, &pesfilter) < 0) {
		perror ("ioctl DMX_SET_PES_FILTER failed");
		return -1;
	}

  This will unpack the payload from all transport stream packets with
  packet ID <pid> and feed it into the MPEG decoder. When pes_type is set
  to DMX_PES_VIDEO it will be handled as video data. Other types are
  DMX_PES_AUDIO, DMX_PES_TELETEXT, DMX_PES_SUBTITLE which will be fed into
  the corresponding decoders (if these deocders exist in hardware or firmware).
  DMX_PES_TELETEXT usually means VBI insertion by the PAL/NTSC encoder for display
  on a connected TV set. If you want to avoid sending the data to one of the
  decoders, use pes_type = DMX_PES_OTHER.

  You must open the demux device once for each PES filter you want to set.
  E.g. if you want audio and video, you must have two distinct file descriptors
  for the two filters.

  DMX_PES_PCR is used by the decoder to achieve a correct timing syncronisation
  between the audio/video/... substreams.

  Note that you have to keep the frontend and demux filedescriptor open until
  you are not interested in the stream anymore. Old API versions did not shut
  down the demodulator and decoder, new driver versions can be configured to
  power down after the device is closed.

-------------------------------------------------------------------------------

I want to record a stream to disk! How?

  Set up a filter for each substream you want to record as above but set
  the pesfilter.output field to DMX_OUT_TAP. Then you can use read() calls
  to receive PES data.

  When you want to receive transport stream packets use DMX_OUT_TS_TAP and
  read the stream from /dev/dvb/adapterX/dvrY. The dvr device gives you
  a multiplex of all filtered PES data with DMX_OUT_TS_TAP. E.g. if you
  want to record video and audio, open demuxX twice and set two PEs filters
  with DMX_OUT_TS_TAP, and open dvrX once to read the TS.

  [ The current API does not allow you to specify you an input/output
    routing for section filters. So you can't get multiplexed section
    data from the dvr device. ]

  Don't forget to keep all device filedescriptors you use open.

-------------------------------------------------------------------------------

I want to play back a recorded stream from disk! How?

  Just do the opposite as above. pesfilter.input is now DMX_IN_DVR.
  Write your transport stream into the dvr device.

-------------------------------------------------------------------------------

What the heck are section filters?

  On some pid's in an MPEG2 stream data is transmitted in form of sections.
  These sections describe the stream content, provider data, service
  information and other things.

  Here a short list of some pid's where section data is transmitted on DVB
  streams:

      0x00  PAT (Program Association Table - refers to PMT's)
      0x10  NIT (Network Information Table - frequency lists etc)
      0x11  SDT (Service Description Table - service names etc)
      0x12  EIT (Event Information Table - event descriptions etc)
      0x14  TDT/TOT (Time and Date Table, Time Offset Table - time and timezone)

  For a complete list look into the ITU H222.0 (MPEG2) and ETSI EN300468
  standards, there you also find some informations how to parse these sections.

  When you want to receive this data the simple way you can set up a section
  filter.

-------------------------------------------------------------------------------

How to set up a section filter?

	#include <linux/dvb/dmx.h>

	struct dmx_sct_filter_params sctfilter;

	if ((fd = open("/dev/dvb/adapter0/demux0", O_RDWR) < 0) {
		perror ("open failed");
		return -1;
	}

	memset(&sctfilter, 0, sizeof(struct dmx_sct_filter_params));

	sctfilter.pid = pid;
	sctfilter.flags = DMX_IMMEDIATE_START;

	if (ioctl(fd, DMX_SET_FILTER, &sctfilter) < 0) {
		perror ("ioctl DMX_SET_FILTER failed");
		return -1;
	}

  Now start subsequent read() calls to receive your sections and parse them.
  Your read-buffer should be at least 4096 bytes if you want to receive complete
  sections, otherwise you have to put the parts together manually.

  If your read() call returns -EOVERFLOW you were not fast enough to read/
  process the section data, the internal driver ringbuffer was overflown.

  This error is usually not critical since section data is transmitted
  periodically. Anyway, you can adjust the ringbuffer size with the
  DMX_SET_BUFFER_SIZE ioctl.

-------------------------------------------------------------------------------

How to do table id filtering?

  The struct dmx_sct_filter_params contains two fields filter.filter and
  filter.mask. set those mask bits to '1' which should be equal to the filter
  bits you set:

	// set up a TDT filter, table id 0x70
	sctfilter.pid = pid;
	sctfilter.filter.filter[0] = 0x70;
	sctfilter.filter.mask[0]   = 0xff;
	sctfilter.flags = DMX_IMMEDIATE_START;

  Then submit the DMX_SET_FILTER ioctl.

  The filter comprises 16 bytes covering byte 0 and byte 3..17 in a section,
  thus excluding bytes 1 and 2 (the length field of a section).

-------------------------------------------------------------------------------

What are not-equal filters?

  When you want to get notified about a new version of a section you can set
  up a not-equal filter. Set those filter.mode bits to '1' for which the filter
  should pass a section when the corresponding section bit is not equal to the
  corresponding filter bit.

-------------------------------------------------------------------------------

*/
