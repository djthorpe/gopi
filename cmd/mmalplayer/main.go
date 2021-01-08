package main

import (
	"fmt"
	"io"
	"os"
	"unsafe"

	"github.com/djthorpe/gopi/v3/pkg/sys/mmal"
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
	fmt.Println("Decoding", fh.Name())
	if err := Run(fh); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func Run(r io.ReadSeeker) error {
	decoder, err := mmal.MMALComponentCreate(mmal.MMAL_COMPONENT_DEFAULT_IMAGE_DECODER)
	if err != nil {
		return err
	}
	defer decoder.Free()

	// Input
	in := decoder.InputPorts()[0]
	in.Format().Video().SetSize(1280, 720)

	// Commit format
	if err := in.FormatCommit(); err != nil {
		return err
	}

	pool_in := in.CreatePool(in.BufferPreferred())
	if pool_in == nil {
		return fmt.Errorf("Failed to create pool_in")
	}
	defer in.FreePool(pool_in)

	// Output
	out := decoder.OutputPorts()[0]
	if err := out.FormatFullCopy(in.Format()); err != nil {
		return err
	}

	// Commit format
	if err := in.FormatCommit(); err != nil {
		return err
	}

	pool_out := in.CreatePool(out.BufferMin())
	if pool_out == nil {
		return fmt.Errorf("Failed to create pool_out")
	}
	defer out.FreePool(pool_out)
	defer out.Flush()

	// Enable all the input port and the output port.
	if err := in.EnableWithCallback(input_callback); err != nil {
		return err
	}
	defer in.Disable()

	if err := out.EnableWithCallback(output_callback); err != nil {
		return err
	}
	defer out.Disable()

	// Enable the decoder and renderer
	if err := decoder.Enable(); err != nil {
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
