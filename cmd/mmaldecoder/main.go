package main

import (
	"fmt"
	"io"
	"os"
	"unsafe"

	"github.com/djthorpe/gopi/v3/pkg/sys/mmal"
	"github.com/djthorpe/gopi/v3/pkg/sys/rpi"
)

// https://github.com/t-moe/rpi_mmal_examples/blob/master/example_basic_2.c
// /opt/vc/src/hello_pi/hello_video/test.h264

const (
	headerSizeBytes = 128
)

func main() {
	rpi.BCMHostInit()
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
	fmt.Println("Decoding", fh.Name())
	if err := Run(fh); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func Run(r io.ReadSeeker) error {
	decoder, err := mmal.MMALComponentCreate(mmal.MMAL_COMPONENT_DEFAULT_VIDEO_DECODER)
	if err != nil {
		return err
	}
	defer decoder.Free()

	// Set format on input, assume 25fps and 1:1 pixel aspect ratio (PAR)
	in := decoder.InputPorts()[0]
	in.Format().SetType(mmal.MMAL_ES_TYPE_VIDEO)
	in.Format().SetEncoding(mmal.MMAL_ENCODING_H264)
	in.Format().Video().SetSize(1280, 720)
	in.Format().Video().SetFrameRate(mmal.NewRational(25))
	in.Format().Video().SetPar(mmal.NewRational(1))
	in.Format().SetFlags(mmal.MMAL_ES_FORMAT_FLAG_FRAMED)

	// Read header into extra data
	if err := in.Format().ExtradataRead(r, headerSizeBytes); err != nil {
		return err
	}

	// Commit format
	if err := in.FormatCommit(); err != nil {
		return err
	}

	// Get format on output
	out := decoder.OutputPorts()[0]
	if out.Format().Encoding() == mmal.MMAL_ENCODING_UNKNOWN {
		return fmt.Errorf("Failed to set output format")
	}

	// Control port
	ctrl := decoder.ControlPort()

	fmt.Println("format in=", in.Format())
	fmt.Println("format out=", out.Format())
	fmt.Println("format ctrl=", ctrl.Format())

	// Now we know the format of both ports and the requirements of the decoder, we can create
	// our buffer headers and their associated memory buffers
	pool_in := in.CreatePool(in.BufferMin())
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
	defer queue.Free()

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

	if err := ctrl.EnableWithCallback(control_callback); err != nil {
		return err
	}
	defer ctrl.Disable()

	// Enable the component. Components will only process data when they are enabled.
	if err := decoder.Enable(); err != nil {
		return err
	}

	// Data processing loop, eof = end of input file, eoe = end of encoding
	eof, eoe := false, false
	i := 0
	for eoe == false && i < 1000 {
		fmt.Println("LOOP eof=", eof, " eoe=", eoe, " i=", i)
		i++

		// Send data to decode to the input port of the video decoder
		if eof == false {
			fmt.Println("  SEND DATA INTO DECODER")
			if buffer := pool_in.Get(); buffer != nil {
				if err := buffer.Fill(r); err == io.EOF {
					buffer.SetFlags(mmal.MMAL_BUFFER_HEADER_FLAG_EOS)
					eof = true
				} else if err != nil {
					return err
				} else {
					buffer.ClearPtsDts()
					if err := in.SendBuffer(buffer); err != nil {
						return err
					} else {
						fmt.Println("     Sent to decoder -> ", buffer)
					}
				}
			}
		}

		// Get our decoded frames. We also need to cope with events generated from the
		// component here
		fmt.Println("  GET DECODED FRAMES")
		for {
			buffer := queue.Get()
			// Get a buffer, end loop if no buffers to process
			if buffer == nil {
				fmt.Println("    -> NO FULL OUTPUT BUFFERS")
				break
			}
			if buffer.HasFlags(mmal.MMAL_BUFFER_HEADER_FLAG_EOS) {
				fmt.Println("    END OF DECODING")
				eoe = true
			}
			if evt := buffer.Event(); evt == 0 {
				// We have a frame, do something with it
				fmt.Println("    GOT DECODED DATA ->", buffer)
			} else if buffer.Event() == mmal.MMAL_EVENT_FORMAT_CHANGED {
				event := buffer.AsFormatChangeEvent()
				fmt.Println("    FORMAT CHANGED", event)

				// Assume we can't reuse the buffers, so we resize
				fmt.Println("      DISABLE OUTPUT PORT")
				if err := out.Disable(); err != nil {
					return err
				}

				// Clear queue
				fmt.Println("      RELEASE BUFFERS")
				fmt.Println("      ->",pool_out)
				for {
					if buf := pool_out.Get(); buf == nil {
						break
					} else {
						buf.Release()
					}
				}

				// Resize pool based on new format
				fmt.Println("      RESIZE POOL")
				if err := pool_out.Resize(event.BufferMin()); err != nil {
					return err
				}
				fmt.Println("      ->",pool_out)

				// Copy over the new format and re-enable the port
				fmt.Println("      FORMAT COPY")
				if err := out.FormatFullCopy(event.Format()); err != nil {
					return err
				} else if err := out.FormatCommit(); err != nil {
					return err
				} else if err := out.Enable(); err != nil {
					return err
				}
				fmt.Println("      ->",out.Format())
			} else {
				fmt.Println("    UNHANDLED EVENT", evt)
			}
			// Once we're done with it, we release it. It will magically go back
			// to its original pool so it can be reused for a new video frame.
			buffer.Release()
		}

		// Send empty buffers to the output port of the decoder to allow the decoder to start
		// producing frames as soon as it gets input data
		fmt.Println("  SEND EMPTY BUFFERS")
		for {
			if buffer := pool_out.Get(); buffer == nil {
				break
			} else if err := out.SendBuffer(buffer); err != nil {
				return err
			} else {
				fmt.Println("     ",buffer," -> ",pool_out)
			}
		}
	}

	fmt.Println("FINISHED LOOP")

	// Return success
	return nil
}

// input_callback is called when a buffer should be discarded on an input port
func input_callback(port *mmal.MMALPort, buffer *mmal.MMALBuffer) {
	// The decoder is done with the data, just recycle the buffer header into its pool
	fmt.Println("input_callback done with buffer", buffer)
	buffer.Release()
}
 
// output_callback is called when a buffer is available on an output port
func output_callback(port *mmal.MMALPort, buffer *mmal.MMALBuffer) {
	queue := (*mmal.MMALQueue)(unsafe.Pointer(port.Userdata()))
	// Queue the decoded video frame
	queue.Put(buffer)
	fmt.Println("output_callback, decoded video=", buffer, " => ", queue)
}

func control_callback(port *mmal.MMALPort, buffer *mmal.MMALBuffer) {
	fmt.Println("control_callback, ", buffer)
	if evt := buffer.Event(); evt == mmal.MMAL_EVENT_EOS {
		fmt.Println("  MMAL_EVENT_EOS")
	} else if evt == mmal.MMAL_EVENT_ERROR {
		fmt.Println("  MMAL_EVENT_ERROR: ", buffer.AsError())
	}
	// TODO: Have main loop process the error
	buffer.Release()
}
