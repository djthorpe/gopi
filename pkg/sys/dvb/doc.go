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
