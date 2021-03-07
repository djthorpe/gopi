// +build openvg

package openvg

////////////////////////////////////////////////////////////////////////////////

/*
  #cgo CFLAGS:   -I/opt/vc/include
  #cgo LDFLAGS:  -L/opt/vc/lib -lOpenVG
  #include <VG/openvg.h>
  #include <VG/vgu.h>
*/
import "C"
