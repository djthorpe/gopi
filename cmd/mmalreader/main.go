package main

import (
	"fmt"
	"io"
	"os"
	"unsafe"

	"github.com/djthorpe/gopi/v3/pkg/sys/mmal"
)

// /opt/vc/src/hello_pi/hello_video/test.h264

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Expected file argument")
		os.Exit(-1)
	}
	fh, err := os.Open(args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	defer fh.Close()
	fmt.Println("Decoding",fh.Name())
	if err := Run(fh); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func Run(r io.Reader) error {
	decoder, err := mmal.MMALComponentCreate(mmal.MMAL_COMPONENT_DEFAULT_VIDEO_DECODER)
	if err != nil {
		return err
	}
	defer decoder.Free()

	// Set format on input
	in := decoder.InputPorts()[0]
	in.Format().SetType(mmal.MMAL_ES_TYPE_VIDEO)
	in.Format().SetEncoding(mmal.MMAL_ENCODING_H264)
	in.Format().Video().SetSize(1280, 720)
	in.Format().Video().SetFrameRate(mmal.NewRational(30))
	in.Format().Video().SetPar(mmal.NewRational(1))
	in.Format().SetFlags(mmal.MMAL_ES_FORMAT_FLAG_FRAMED)

	if err := in.FormatCommit(); err != nil {
		return err
	}

	// Get format on output
	out := decoder.OutputPorts()[0]
	fmt.Println("decoder in=", in.Format())
	fmt.Println("decoder out=", out.Format())
	if out.Format().Encoding() == mmal.MMAL_ENCODING_UNKNOWN {
		return fmt.Errorf("Failed to set output format")
	}

	// Now we know the format of both ports and the requirements of the decoder, we can create
	// our buffer headers and their associated memory buffers
	pool_in := in.CreatePool(in.BufferPreferred())
	if pool_in == nil {
		return fmt.Errorf("Failed to create in_pool")
	}
	defer in.FreePool(pool_in)

	pool_out := out.CreatePool(out.BufferPreferred())
	if pool_out == nil {
		return fmt.Errorf("Failed to create out_pool")
	}
	defer out.FreePool(pool_out)

	// Create a queue to store our decoded video frames. The callback we will get when
	// a frame has been decoded will put the frame into this queue.
	queue := mmal.MMALQueueCreate()
	out.SetUserdata(uintptr(unsafe.Pointer(queue)))

	// Enable all the input port and the output port.
	// The callback specified here is the function which will be called when the buffer header
	// we sent to the component has been processed.
	if err := in.EnableWithCallback(input_callback); err != nil {
		return err
	}
	defer in.Disable()

	if err := out.EnableWithCallback(output_callback); err != nil {
		return err
	}
	defer out.Disable()

	// Enable the component. Components will only process data when they are enabled.
	if err := decoder.Enable(); err != nil {
		return err
	}

	// Data processing loop, eof = end of input file, eoe = end of encoding
	eof, eoe := false, false
	i := 0
	for eoe == false && i < 500 {
		fmt.Println("LOOP eof=", eof, " eoe=", eoe," i=",i)
		i++

		// Send empty buffers to the output port of the decoder to allow the decoder to start
		// producing frames as soon as it gets input data
		fmt.Println("  SEND EMPTY BUFFERS")
		for {
			if buffer := pool_out.Get(); buffer == nil {
				break
			} else if err := out.SendBuffer(buffer); err != nil {
				return err
			} else {
				fmt.Println("     -> ", buffer)
			}
		}

		// Send data to decode to the input port of the video decoder
		fmt.Println("  SEND DATA INTO DECODER")
		if buffer := pool_in.Get(); buffer != nil {
			if err := buffer.Fill(r); err == io.EOF {
				eof = true
			} else if err != nil {
				return err
			} else if err := in.SendBuffer(buffer); err != nil {
				return err
			} else {
				fmt.Println("     Sent to decoder -> ", buffer)
			}
		} else {
			fmt.Println("    -> NO EMPTY INPUT BUFFERS")
		}

		// Get our decoded frames. We also need to cope with events generated from the
		// component here
		fmt.Println("  GET DECODED FRAMES")
		for {
			// Get a buffer, end loop if no buffers to process
			if buffer := queue.Get(); buffer == nil {
				fmt.Println("    -> NO FULL OUTPUT BUFFERS")
				break
			} else {
				if evt := buffer.Event(); evt != 0 {
					// This is an event. Do something with it and release the buffer.
					fmt.Println("    -> EVT=", buffer.Event())
					if evt&mmal.MMAL_EVENT_EOS == mmal.MMAL_EVENT_EOS || evt&mmal.MMAL_EVENT_ERROR == mmal.MMAL_EVENT_ERROR {
						eoe = true
					} else if evt&mmal.MMAL_EVENT_FORMAT_CHANGED == mmal.MMAL_EVENT_FORMAT_CHANGED {
						fmt.Println("     ->",out.Format())
					}
				} else {
					// We have a frame, do something with it (why not display it for instance?).
					fmt.Println("    -> DATA=", buffer)
				}
				// Once we're done with it, we release it. It will magically go back
				// to its original pool so it can be reused for a new video frame.
				buffer.Release()
			}
		}
	}

	// Return success
	return nil
}

func input_callback(port *mmal.MMALPort, buffer *mmal.MMALBuffer) {
	// The decoder is done with the data, just recycle the buffer header into its pool
	fmt.Println("input_callback done with buffer",buffer)
	buffer.Release()
}

func output_callback(port *mmal.MMALPort, buffer *mmal.MMALBuffer) {
	queue := (*mmal.MMALQueue)(unsafe.Pointer(port.Userdata()))
	// Queue the decoded video frame
	fmt.Println("output_callback, decoded video=",buffer," => ",queue)
	queue.Put(buffer)
}
