/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	_ "github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/third_party/grpc-go"
	_ "github.com/djthorpe/gopi/third_party/grpc-go/reflection"
)

type Server struct{}
