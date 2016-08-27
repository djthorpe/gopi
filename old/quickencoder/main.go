/*
    quickencoder - stream file directly to YouTube
*/
package main

////////////////////////////////////////////////////////////////////////////////

import (
	"flag"
	"github.com/djthorpe/goav/avformat"
	"github.com/djthorpe/goav/avcodec"
	"fmt"
	"log"
//	"os"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////

var (
    filename = flag.String("input","","Input filename")
)

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Read flags
    flag.Parse()

	// Register all formats and codecs
	avformat.AvRegisterAll()

	var (
		ctxtFormat    *avformat.Context
		url           string
		ctxtSource    *avcodec.Context
//		ctxtDest      *avcodec.Context
//		videoCodec    *avcodec.Codec
//		videoFrame    *avutil.Frame
//		videoFrameRGB *avutil.Frame
//		packet        *avcodec.Packet
//		ctxtSws       *swscale.Context
		videoStream   int
//		frameFinished int
//		numBytes      int
	)

	// Open video file
	if avformat.AvformatOpenInput(&ctxtFormat,*filename,nil,nil) != 0 {
		fmt.Printf("Error: Couldn't open file: %s\n",*filename)
		return
	}

	// Retrieve stream information
	if ctxtFormat.AvformatFindStreamInfo(nil) < 0 {
		log.Println("Error: Couldn't find stream information.")
		return
	}

	// Dump information about file onto standard error
	ctxtFormat.AvDumpFormat(0, url, 0)

	// Find the first video stream
	videoStream = -1
	n := ctxtFormat.NbStreams()
	s := ctxtFormat.Streams()
	codec := s.Codec()
	log.Print("Number of Streams: ", n)

	for i := 0; i < int(n); i++ {
		// ctxtFormat->streams[i]->codec->codec_type
		log.Println("Stream Number:", i)

		//FIX: AvMEDIA_TYPE_VIDEO
		if (*avformat.CodecContext)(codec) != nil {
			videoStream = i
			break
		}
	}

	if videoStream == -1 {
		log.Println("Couldn't find a video stream")
		return
	}

	// Get a pointer to the codec context for the video stream
	ctxtSource = (*avcodec.Context)(unsafe.Pointer(&codec))

	fmt.Printf("ctx = %+v\n",ctxtSource)

	log.Println("Bit Rate:", ctxtSource.BitRate())
	log.Println("Channels:", ctxtSource.Channels())
	log.Println("Coded_height:", ctxtSource.CodedHeight())
	log.Println("Coded_width:", ctxtSource.CodedWidth())
	log.Println("Coder_type:", ctxtSource.CoderType())
	log.Println("Height:", ctxtSource.Height())
	log.Println("Profile:", ctxtSource.Profile())
	log.Println("Width:", ctxtSource.Width())
	log.Printf("Codec ID: %x", ctxtSource.CodecId())


}
/*
	// C.enum_AVCodecID
	codec_id := ctxtSource.CodecId()

	// Find the decoder for the video stream
	videoCodec = avcodec.AvcodecFindDecoder(codec_id)
	if videoCodec == nil {
		log.Printf("Error: Unsupported codec: 0x%x",codec_id)
		return // Codec not found
	}

	// Copy context
	ctxtDest = videoCodec.AvcodecAllocContext3()

	if ctxtDest.AvcodecCopyContext(ctxtSource) != 0 {
		log.Println("Error: Couldn't copy codec context")
		return // Error copying codec context
	}

	// Open codec
	if ctxtDest.AvcodecOpen2(videoCodec, nil) < 0 {
		return // Could not open codec
	}

	// Allocate video frame
	videoFrame = avutil.AvFrameAlloc()

	// Allocate an Frame structure
	if videoFrameRGB = avutil.AvFrameAlloc(); videoFrameRGB == nil {
		return
	}

	//##TODO
	var a swscale.PixelFormat
	var b int
	//avcodec.PixelFormat
	//avcodec.PIX_FMT_RGB24
	//avcodec.SWS_BILINEAR

	w := ctxtDest.Width()
	h := ctxtDest.Height()
	pix_fmt := ctxtDest.PixFmt()

	// Determine required buffer size and allocate buffer
	numBytes = avcodec.AvpictureGetSize((avcodec.PixelFormat)(a), w, h)

	buffer := avutil.AvMalloc(uintptr(numBytes))

	// Assign appropriate parts of buffer to image planes in videoFrameRGB
	// Note that videoFrameRGB is an Frame, but Frame is a superset
	// of Picture
	avp := (*avcodec.Picture)(unsafe.Pointer(videoFrameRGB))
	avp.AvpictureFill((*uint8)(buffer), (avcodec.PixelFormat)(a), w, h)

	// initialize SWS context for software scaling
	ctxtSws = swscale.SwsGetcontext(w,
		h,
		(swscale.PixelFormat)(pix_fmt),
		w,
		h,
		a,
		b,
		nil,
		nil,
		nil,
	)

	// Read frames and save first five frames to disk
	i := 0

	for ctxtFormat.AvReadFrame(packet) >= 0 {
		// Is this a packet from the video stream?
		s := packet.StreamIndex()
		if s == videoStream {
			// Decode video frame
			ctxtDest.AvcodecDecodeVideo2((*avcodec.Frame)(unsafe.Pointer(videoFrame)), &frameFinished, packet)

			// Did we get a video frame?
			if frameFinished > 0 {
				// Convert the image from its native format to RGB
				d := avutil.Data(videoFrame)
				l := avutil.Linesize(videoFrame)
				dr := avutil.Data(videoFrameRGB)
				lr := avutil.Linesize(videoFrameRGB)
				swscale.SwsScale(ctxtSws,
					d,
					l,
					0,
					h,
					dr,
					lr,
				)

				// Save the frame to disk
				if i <= 5 {
					saveFrame(videoFrameRGB, w, h, i)
				}
				i++
			}
		}

		// Free the packet that was allocated by av_read_frame
		packet.AvFreePacket()
	}

	// Free the RGB image
	avutil.AvFree(buffer)
	avutil.AvFrameFree(videoFrameRGB)

	// Free the YUV frame
	avutil.AvFrameFree(videoFrame)

	// Close the codecs
	ctxtDest.AvcodecClose()
	ctxtSource.AvcodecClose()

	// Close the video file
	ctxtFormat.AvformatCloseInput()

}

func saveFrame(videoFrame *avutil.Frame, width int, height int, iFrame int) {

	var szFilename string
	var y int
	var file *os.File
	var err error

	szFilename = ""

	// Open file
	szFilename = fmt.Sprintf("frame%d.ppm", iFrame)

	if file, err = os.Open(szFilename); err != nil {
		log.Println("Error Reading")
	}

	// Write header
	fh := fmt.Sprintf("P6\n%d %d\n255\n", width, height)
	log.Println(fh)

	// Write pixel data
	for y = 0; y < height; y++ {
		// d := avutil.Data(videoFrame)
		// l := avutil.Linesize(videoFrame)
		//##TODO
		f := make([]byte, 100)
		file.Write(f)
	}

	file.Close()

}
*/

