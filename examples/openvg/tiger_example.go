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
	value_r = regexp.MustCompile("([0-9\\.]+)f?")
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

func MyRunLoop(app *app.App) error {
	args := app.FlagSet.Args()
	if len(args) != 1 {
		return app.Logger.Error("Missing data filename")
	}
	opcodes, values, err := ReadData(args[0])
	if err != nil {
		return err
	}
	fmt.Println(opcodes)
	fmt.Println(values)
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
