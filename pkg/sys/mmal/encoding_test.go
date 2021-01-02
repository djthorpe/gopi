//+build mmal

package mmal_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/sys/mmal"
)

func Test_Encoding_001(t *testing.T) {
	encs := []mmal.MMALEncodingType{
		mmal.MMAL_ENCODING_H264,
		mmal.MMAL_ENCODING_MVC,
		mmal.MMAL_ENCODING_H263,
		mmal.MMAL_ENCODING_MP4V,
		mmal.MMAL_ENCODING_MP2V,
		mmal.MMAL_ENCODING_MP1V,
		mmal.MMAL_ENCODING_WMV3,
		mmal.MMAL_ENCODING_WMV2,
		mmal.MMAL_ENCODING_WMV1,
		mmal.MMAL_ENCODING_WVC1,
		mmal.MMAL_ENCODING_VP8,
		mmal.MMAL_ENCODING_VP7,
		mmal.MMAL_ENCODING_VP6,
		mmal.MMAL_ENCODING_THEORA,
		mmal.MMAL_ENCODING_SPARK,
		mmal.MMAL_ENCODING_MJPEG,
	}

	for i, enc := range encs {
		t.Log(i, enc)
	}
}
