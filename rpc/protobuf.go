package rpc

//go:generate protoc --go_out=../pkg/rpc --go_opt=paths=source_relative --go-grpc_out=../pkg/rpc --go-grpc_opt=paths=source_relative input/input.proto
//go:generate protoc --go_out=../pkg/rpc --go_opt=paths=source_relative --go-grpc_out=../pkg/rpc --go-grpc_opt=paths=source_relative ping/ping.proto

/*
	This folder contains all the protocol buffer definitions. You
	can generate golang code for these definitions by running:

	  go generate -x github.com/djthorpe/gopi/v3/rpc

	where you have installed the protoc compiler and the GRPC plugin for
	golang. In order to do that on a Mac:

  	  mac# brew install protobuf
	  mac# go get -u github.com/golang/protobuf/protoc-gen-go

	On Debian Linux (including Raspian Linux) use the following commands
	instead:

	  linux# sudo apt install protobuf-compiler
	  linux# sudo apt install libprotobuf-dev
	  linux# go get -u github.com/golang/protobuf/protoc-gen-go
*/
