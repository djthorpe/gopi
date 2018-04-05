/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package protobuf

//go:generate protoc helloworld/helloworld.proto --go_out=plugins=grpc:.

/*
	This folder contains all the protocol buffer definitions including
	the RPC Service definitions. You generate golang code by running:

	go generate -x github.com/djthorpe/gopi/protobuf

	where you have installed the protoc compiler and the GRPC plugin for
	golang. In order to do that on a Mac:

	mac# brew install protobuf
	mac# go get -u github.com/golang/protobuf/protoc-gen-go
	mac# cd ${GOPATH}/src/github.com/djthorpe/gopi
	mac# cmd/build_rpc.sh

	On Debian Linux (including Raspian Linux) use the following commands
	instead:

	rpi# sudo apt install protobuf-compiler
	rpi# go get -u github.com/golang/protobuf/protoc-gen-go
	rpi# cd ${GOPATH}/src/github.com/djthorpe/gopi
	rpi# cmd/build_rpc.sh
*/
