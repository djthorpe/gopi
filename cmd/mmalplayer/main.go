package main

import (
	"fmt"
	"io"
	"os"
	"unsafe"

	"github.com/djthorpe/gopi/v3/pkg/sys/mmal"
)

// https://github.com/t-moe/rpi_mmal_examples/blob/master/example_basic_2.c
// /opt/vc/src/hello_pi/hello_video/test.h264

const (
	headerSizeBytes = 128
)

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
	fmt.Println("Playing", fh.Name())
	if err := Run(fh); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func Run(r io.ReadSeeker) error {
	render, err :=  mmal.MMALComponentCreate(mmal.MMAL_COMPONENT_DEFAULT_VIDEO_RENDERER)
	if err != nil {
		return err
	}
	defer render.Free()

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

	// Connect output to renderer
	connector,err := mmal.MMALConnectionCreate(render.InputPorts()[0],decoder.OutputPorts()[0],mmal.MMAL_CONNECTION_FLAG_TUNNELLING)
	if err != nil {
		return err
	}
	defer connector.Free()

	// Create buffer headers and their associated memory buffers
	pool_in := in.CreatePool(in.BufferMin())
	if pool_in == nil {
		return fmt.Errorf("Failed to create in_pool")
	}
	defer in.FreePool(pool_in)

	// Enable all the input port and the output port.
	// The callback specified here is the function which will be called when the buffer header
	// we sent to the component has been processed.
	if err := in.EnableWithCallback(input_callback); err != nil {
		return err
	}
	defer in.Disable()

	// Enable ports on connector
	if err := connector.Enable(); err != nil {
		return err
	}

	// Enable the decoder and renderer
	if err := decoder.Enable(); err != nil {
		return err
	}
	if err := render.Enable(); err != nil {
		return err
	}

	// Data processing loop, eof = end of input file, eoe = end of encoding
	eof := false 
	i := 0
	for eof == false && i < 500 {
		fmt.Println("LOOP eof=", eof, " i=", i)
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
