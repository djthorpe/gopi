/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This example is the tiger face example
package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"regexp"
	"errors"
	"strconv"
)

import (
	app "github.com/djthorpe/gopi/app"
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////

type Operation struct {
	fill khronos.VGPaint
	stroke khronos.VGPaint
	path khronos.VGPath
}

////////////////////////////////////////////////////////////////////////////////

var (
	opcode_r = regexp.MustCompile("'(\\w)'")
	value_r = regexp.MustCompile("([0-9\\.]*[0-9]+)f?")
)

////////////////////////////////////////////////////////////////////////////////

// Return the opcodes, values and error
func ReadData(filename string) ([]string,[]float32,error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil,nil,err
	}
	// Read opcodes and values
	opcodes := opcode_r.FindAllSubmatch(data,-1)
	if opcodes == nil {
		return nil,nil,errors.New("Invalid data file, no opcodes")
	}
	values := value_r.FindAllSubmatch(data,-1)
	if values == nil {
		return nil,nil,errors.New("Invalid data file, no values")
	}

	opcodes2 := make([]string,len(opcodes))
	values2 := make([]float32,len(values))

	// Convert opcodes to string
	for i,opcode := range opcodes {
		opcodes2[i] = string(opcode[1])
	}

	// Convert values to float32
	for i,value := range values {
		value64, err := strconv.ParseFloat(string(value[1]),32)
		if err != nil {
			return nil,nil,err
		}
		values2[i] = float32(value64)
	}

	// Success
	return opcodes2,values2,nil
}

////////////////////////////////////////////////////////////////////////////////

func (this *Operation) ParseFillOpcode(code string) error {
	switch(code) {
	case "N":
		fmt.Println("<no fill>")
	case "F":
		fmt.Println("<fill non zero>")
	case "E":
		fmt.Println("<fill even odd>")
	default:
		return errors.New("Invalid ParseFillOpcode value")
	}
	return nil
}

func (this *Operation) ParseStrokeOpcode(code string) error {
	switch(code) {
	case "N":
		fmt.Println("<no stroke>")
	case "S":
		fmt.Println("<stroke>")
	default:
		return errors.New("Invalid ParseStrokeOpcode value")
	}
	return nil
}

func (this *Operation) ParseLineCapOpcode(code string) error {
	switch(code) {
	case "B":
		fmt.Println("<butt>")
	case "R":
		fmt.Println("<round>")
	case "S":
		fmt.Println("<square>")
	default:
		return errors.New("Invalid ParseLineCapOpcode value")
	}
	return nil
}

func (this *Operation) ParseLineJoinOpcode(code string) error {
	switch(code) {
	case "M":
		fmt.Println("<miter>")
	case "R":
		fmt.Println("<round>")
	case "B":
		fmt.Println("<bevel>")
	default:
		return errors.New("Invalid ParseLineJoinOpcode value")
	}
	return nil
}

func (this *Operation) ParseMiterLimit(limit float32) error {
	fmt.Println("<miter limit ",limit,">")
	return nil
}

func (this *Operation) ParseStrokeWidth(width float32) error {
	fmt.Println("<stroke width ",width,">")
	return nil
}

func (this *Operation) ParseStrokeColor(r,g,b float32) error {
	fmt.Println("<stroke color ",r,g,b,">")
	return nil
}

func (this *Operation) ParseFillColor(r,g,b float32) error {
	fmt.Println("<fill color ",r,g,b,">")
	return nil
}

func (this *Operation) ParsePathPoint(opcode string,points []float32,i int) (int,error) {
	switch(opcode) {
	case "M":
		fmt.Println("<moveto",points[i],points[i+1],">")
		return 2, nil
	case "L":
		fmt.Println("<lineto",points[i],points[i+1],">")
		return 2, nil
	case "C":
		fmt.Println("<cubicto",points[i],points[i+1],points[i+2],points[i+3],points[i+4],points[i+5],">")
		return 6, nil
	case "E":
		fmt.Println("<closepath>")
		return 0, nil
	default:
		return 0,errors.New("Invalid ParsePathPoint opcode value")
	}
}

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(app *app.App) error {
	args := app.FlagSet.Args()
	if len(args) != 1 {
		return app.Logger.Error("Missing data filename")
	}
	opcodes, values, err := ReadData(args[0])
	if err != nil {
		return err
	}

	c := 0
	v := 0
	for c < len(opcodes) && v < len(values) {
		op := new(Operation)

		// Fill opcode
		if err := op.ParseFillOpcode(opcodes[c]); err != nil {
			return err
		}
		c += 1

		// Stroke opcode
		if err := op.ParseStrokeOpcode(opcodes[c]); err != nil {
			return err
		}
		c += 1

		// Line Cap
		if err := op.ParseLineCapOpcode(opcodes[c]); err != nil {
			return err
		}
		c += 1

		// Line Join
		if err := op.ParseLineJoinOpcode(opcodes[c]); err != nil {
			return err
		}
		c += 1

		// Miter Limit
		if err := op.ParseMiterLimit(values[v]); err != nil {
			return err
		}
		v += 1

		// Stroke Width
		if err := op.ParseStrokeWidth(values[v]); err != nil {
			return err
		}
		v += 1

		// Colors
		if err := op.ParseStrokeColor(values[v],values[v+1],values[v+2]); err != nil {
			return err
		}
		v += 3
		if err := op.ParseFillColor(values[v],values[v+1],values[v+2]); err != nil {
			return err
		}
		v += 3

		// Path elements
		elements := int(values[v])
		v += 1
		for i := 0; i < elements; i++ {
			vinc, err := op.ParsePathPoint(opcodes[c],values,v)
			if err != nil {
				return err
			}
			c += 1
			v += vinc
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_EGL | app.APP_OPENVG)
	config.FlagSet.FlagFloat64("opacity", 1.0, "Image opacity, 0.0 -> 1.0")

	// Create the application
	myapp, err := app.NewApp(config)
	if err == app.ErrHelp {
		return
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer myapp.Close()

	// Run the application
	if err := myapp.Run(MyRunLoop); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
