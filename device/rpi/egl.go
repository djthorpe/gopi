/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

import (
	gopi "../.."            /* import "github.com/djthorpe/gopi" */
	khronos "../../khronos" /* import "github.com/djthorpe/gopi/khronos" */
	util "../../util"       /* import "github.com/djthorpe/gopi/util" */
)

////////////////////////////////////////////////////////////////////////////////

/*
  #cgo CFLAGS:   -I/opt/vc/include
  #cgo LDFLAGS:  -L/opt/vc/lib -lEGL -lGLESv2
  #include <EGL/egl.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Configuration when creating the EGL driver
type EGL struct {
	Display gopi.DisplayDriver
}

// Display handle
type eglDisplay uintptr

// Context handle
type eglContext uintptr

// Surface handle
type eglSurface uintptr

// Configuration handle
type eglConfig uintptr

// EGL driver
type eglDriver struct {
	major, minor int
	dx           *DXDisplay
	display      eglDisplay
	log          *util.LoggerDevice
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// General constants
	EGL_DEFAULT_DISPLAY uintptr    = 0
	EGL_NO_DISPLAY      eglDisplay = 0
	EGL_NO_CONTEXT      eglContext = 0
	EGL_NO_SURFACE      eglSurface = 0
	EGL_NO_CONFIG       eglConfig  = 0
	EGL_FALSE           uint       = 0
	EGL_TRUE            uint       = 1
)

const (
	// QueryString targets
	EGL_VENDOR      uint16 = 0x3053
	EGL_VERSION     uint16 = 0x3054
	EGL_EXTENSIONS  uint16 = 0x3055
	EGL_CLIENT_APIS uint16 = 0x308D
)

const (
	// EGL_RENDERABLE_TYPE mask bits
	EGL_OPENGL_ES_BIT  uint16 = 0x0001
	EGL_OPENVG_BIT     uint16 = 0x0002
	EGL_OPENGL_ES2_BIT uint16 = 0x0004
	EGL_OPENGL_BIT     uint16 = 0x0008
)

const (
	// BindAPI/QueryAPI targets
	EGL_OPENGL_ES_API uint16 = 0x30A0
	EGL_OPENVG_API    uint16 = 0x30A1
	EGL_OPENGL_API    uint16 = 0x30A2
)

const (
	// Errors / GetError return values
	EGL_SUCCESS             uint16 = 0x3000
	EGL_NOT_INITIALIZED     uint16 = 0x3001
	EGL_BAD_ACCESS          uint16 = 0x3002
	EGL_BAD_ALLOC           uint16 = 0x3003
	EGL_BAD_ATTRIBUTE       uint16 = 0x3004
	EGL_BAD_CONFIG          uint16 = 0x3005
	EGL_BAD_CONTEXT         uint16 = 0x3006
	EGL_BAD_CURRENT_SURFACE uint16 = 0x3007
	EGL_BAD_DISPLAY         uint16 = 0x3008
	EGL_BAD_MATCH           uint16 = 0x3009
	EGL_BAD_NATIVE_PIXMAP   uint16 = 0x300A
	EGL_BAD_NATIVE_WINDOW   uint16 = 0x300B
	EGL_BAD_PARAMETER       uint16 = 0x300C
	EGL_BAD_SURFACE         uint16 = 0x300D
	EGL_CONTEXT_LOST        uint16 = 0x300E
)

const (
	// Configuration attributes
	EGL_ALPHA_SIZE      uint16 = 0x3021
	EGL_BLUE_SIZE       uint16 = 0x3022
	EGL_GREEN_SIZE      uint16 = 0x3023
	EGL_RED_SIZE        uint16 = 0x3024
	EGL_SURFACE_TYPE    uint16 = 0x3033
	EGL_RENDERABLE_TYPE uint16 = 0x3040
	EGL_NONE            uint16 = 0x3038
)

const (
	// Config attribute mask bits
	EGL_WINDOW_BIT uint16 = 0x0004
)

const (
	// QuerySurface / SurfaceAttrib / CreatePbufferSurface targets
	EGL_SWAP_BEHAVIOR uint16 = 0x3093
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	// Errors
	eglError = map[uint16]error{
		EGL_SUCCESS:             nil,
		EGL_NOT_INITIALIZED:     errors.New("EGL_NOT_INITIALIZED"),
		EGL_BAD_ACCESS:          errors.New("EGL_BAD_ACCESS"),
		EGL_BAD_ALLOC:           errors.New("EGL_BAD_ALLOC"),
		EGL_BAD_ATTRIBUTE:       errors.New("EGL_BAD_ATTRIBUTE"),
		EGL_BAD_CONFIG:          errors.New("EGL_BAD_CONFIG"),
		EGL_BAD_CONTEXT:         errors.New("EGL_BAD_CONTEXT"),
		EGL_BAD_CURRENT_SURFACE: errors.New("EGL_BAD_CURRENT_SURFACE"),
		EGL_BAD_DISPLAY:         errors.New("EGL_BAD_DISPLAY"),
		EGL_BAD_MATCH:           errors.New("EGL_BAD_MATCH"),
		EGL_BAD_NATIVE_PIXMAP:   errors.New("EGL_BAD_NATIVE_PIXMAP"),
		EGL_BAD_NATIVE_WINDOW:   errors.New("EGL_BAD_NATIVE_WINDOW"),
		EGL_BAD_PARAMETER:       errors.New("EGL_BAD_PARAMETER"),
		EGL_BAD_SURFACE:         errors.New("EGL_BAD_SURFACE"),
		EGL_CONTEXT_LOST:        errors.New("EGL_CONTEXT_LOST"),
	}
	EGLErrorUnknown                  = errors.New("Unknown EGL error")
	EGLErrorInvalidDisplayDriver     = errors.New("Invalid display driver parameter")
	EGLErrorInvalidAPIBind           = errors.New("Invalid EGL API binding parameter")
	EGLErrorInvalidFrameBufferConfig = errors.New("Invalid EGL framebuffer parameter")
	EGLErrorInvalidParameter         = errors.New("Invalid parameter")
)

var (
	// Map API names to API values
	eglClientAPI = map[string]C.EGLenum{
		"OpenGL_ES": C.EGLenum(EGL_OPENGL_ES_API),
		"OpenVG":    C.EGLenum(EGL_OPENVG_API),
		"OpenGL":    C.EGLenum(EGL_OPENGL_API),
	}

	// Map API names to EGL_RENDERABLE_TYPE
	eglRenderableType = map[string]uint16{
		"OpenGL":    EGL_OPENGL_BIT,
		"OpenVG":    EGL_OPENVG_BIT,
		"OpenGL_ES": EGL_OPENGL_ES2_BIT,
	}
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

// Open
func (config EGL) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	this := new(eglDriver)
	this.log = log

	// Get EGL Display
	this.display = eglDisplay(unsafe.Pointer(C.eglGetDisplay(EGL_DEFAULT_DISPLAY)))
	if this.display == EGL_NO_DISPLAY {
		return nil, this.GetError()
	}

	// Initialise
	var major, minor C.EGLint
	result := C.eglInitialize(this.display, (*C.EGLint)(unsafe.Pointer(&major)), (*C.EGLint)(unsafe.Pointer(&minor)))
	if result == C.EGLBoolean(EGL_FALSE) {
		return nil, this.GetError()
	}
	this.major = int(major)
	this.minor = int(minor)

	// DX
	var ok bool
	if config.Display == nil {
		return nil, EGLErrorInvalidDisplayDriver
	}
	this.dx, ok = config.Display.(*DXDisplay)
	if ok != true {
		return nil, EGLErrorInvalidParameter
	}

	log.Debug2("<rpi.EGL>OpenEGL version=%v.%v", major, minor)

	// Success
	return this, nil
}

// Close the driver
func (this *eglDriver) Close() error {
	this.log.Debug2("<rpi.EGL>Close")

	result := C.eglTerminate(this.display)
	if result == C.EGLBoolean(EGL_FALSE) {
		return this.GetError()
	}
	return nil
}

// Return string version of the EGL interface
func (this *eglDriver) String() string {
	var parts = make([]string, 0)
	parts = append(parts,
		"version="+this.GetVersionString(),
		"vendor="+this.GetVendorString(),
		"apis="+strings.Join(this.GetSupportedClientAPIs(), ","),
		"extensions="+strings.Join(this.GetExtensions(), ","),
		"display_frame="+this.GetFrame().String(),
	)
	api, err := this.QueryAPI()
	if err == nil {
		parts = append(parts, "bound_api="+api)
	}

	return "<rpi.EGL>{ " + strings.Join(parts, " ") + "}"
}

// Return the logging object
func (this *eglDriver) Log() *util.LoggerDevice {
	return this.log
}

// Return major and minor version of EGL
func (this *eglDriver) GetVersion() (int, int) {
	return this.major, this.minor
}

// Return vendor information
func (this *eglDriver) GetVendorString() string {
	return C.GoString(C.eglQueryString(this.display, C.EGLint(EGL_VENDOR)))
}

// Return version information
func (this *eglDriver) GetVersionString() string {
	return C.GoString(C.eglQueryString(this.display, C.EGLint(EGL_VERSION)))
}

// Return extensions information
func (this *eglDriver) GetExtensions() []string {
	return strings.Split(C.GoString(C.eglQueryString(this.display, C.EGLint(EGL_EXTENSIONS))), " ")
}

// Return API's information
func (this *eglDriver) GetSupportedClientAPIs() []string {
	return strings.Split(C.GoString(C.eglQueryString(this.display, C.EGLint(EGL_CLIENT_APIS))), " ")
}

// Bind API
func (this *eglDriver) BindAPI(api string) error {
	value, ok := eglClientAPI[api]
	if !ok {
		return EGLErrorInvalidAPIBind
	}
	result := C.eglBindAPI(value)
	if result == C.EGLBoolean(EGL_FALSE) {
		return this.GetError()
	}
	return nil
}

// Query currently bound API
func (this *eglDriver) QueryAPI() (string, error) {
	api := C.eglQueryAPI()
	for k, v := range eglClientAPI {
		if api == v {
			return k, nil
		}
	}
	return "", EGLErrorInvalidAPIBind
}

// Return error
func (this *eglDriver) GetError() error {
	code := C.eglGetError()
	err, ok := eglError[uint16(code)]
	if !ok {
		return EGLErrorUnknown
	}
	return err
}

// Return framesize of the display
func (this *eglDriver) GetFrame() khronos.EGLFrame {
	size := this.dx.GetSize()
	return khronos.EGLFrame{khronos.EGLPoint{}, khronos.EGLSize{uint(size.Width), uint(size.Height)}}
}

// Human-readable version of the eglDisplay
func (h eglDisplay) String() string {
	return fmt.Sprintf("<rpi.eglDisplay>{%08X}", uint32(h))
}

// Human-readable version of the eglContext
func (h eglContext) String() string {
	return fmt.Sprintf("<rpi.eglContext>{%08X}", uint32(h))
}

// Human-readable version of the eglSurface
func (h eglSurface) String() string {
	return fmt.Sprintf("<rpi.eglSurface>{%08X}", uint32(h))
}

// Human-readable version of the eglSurface
func (h eglConfig) String() string {
	return fmt.Sprintf("<rpi.eglConfig>{%08X}", uint32(h))
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Choose EGL frame buffer configuration
func (this *eglDriver) getFrameBufferConfiguration() (eglConfig, error) {
	// Get bound API
	api, err := this.QueryAPI()
	if err != nil {
		return EGL_NO_CONFIG, err
	}
	// Get Renderable type depending on the API
	renderable_type, ok := eglRenderableType[api]
	if !ok {
		return EGL_NO_CONFIG, EGLErrorInvalidFrameBufferConfig
	}
	attribute_list := []C.EGLint{
		C.EGLint(EGL_RED_SIZE), C.EGLint(8),
		C.EGLint(EGL_GREEN_SIZE), C.EGLint(8),
		C.EGLint(EGL_BLUE_SIZE), C.EGLint(8),
		C.EGLint(EGL_ALPHA_SIZE), C.EGLint(8),
		C.EGLint(EGL_SURFACE_TYPE), C.EGLint(EGL_WINDOW_BIT),
		C.EGLint(EGL_RENDERABLE_TYPE), C.EGLint(renderable_type),
		C.EGLint(EGL_NONE),
	}
	var num_config C.EGLint
	var config C.EGLConfig
	result := C.eglChooseConfig(this.display, &attribute_list[0], &config, C.EGLint(1), &num_config)
	if result == C.EGLBoolean(EGL_FALSE) {
		return EGL_NO_CONFIG, this.GetError()
	}
	if num_config != C.EGLint(1) {
		return EGL_NO_CONFIG, this.GetError()
	}

	// success
	return eglConfig(config), nil
}

// Create EGL Context with API value
func (this *eglDriver) createContext(api string) (eglConfig, eglContext, error) {
	// Bind API
	err := this.BindAPI(api)
	if err != nil {
		return EGL_NO_CONFIG, EGL_NO_CONTEXT, err
	}

	// Get configuration
	config, err := this.getFrameBufferConfiguration()
	if err != nil {
		return EGL_NO_CONFIG, EGL_NO_CONTEXT, err
	}

	// Create rendering context
	context := eglContext(C.eglCreateContext(this.display, config, C.EGLContext(EGL_NO_CONTEXT), (*C.EGLint)(unsafe.Pointer(nil))))
	if context == EGL_NO_CONTEXT {
		return EGL_NO_CONFIG, EGL_NO_CONTEXT, this.GetError()
	}

	return config, context, nil
}

// Destroy EGL Context
func (this *eglDriver) destroyContext(context eglContext) error {
	result := C.eglDestroyContext(this.display, C.EGLContext(context))
	if result == C.EGLBoolean(EGL_FALSE) {
		return this.GetError()
	}
	return nil
}

// Set current EGL Context
func (this *eglDriver) makeCurrent(surface eglSurface,context eglContext) error {
	result := C.eglMakeCurrent(this.display, C.EGLSurface(surface), C.EGLSurface(surface), C.EGLContext(context))
	if result == C.EGLBoolean(EGL_FALSE) {
		return this.GetError()
	}
	return nil
}

// Create surface
func (this *eglDriver) createSurface(config eglConfig, window *eglNativeWindow) (eglSurface, error) {
	// Create EGL window surface given a native window
	surface := eglSurface(C.eglCreateWindowSurface(this.display, C.EGLConfig(config), (*C.EGLNativeWindowType)(unsafe.Pointer(window)), nil))
	if surface == EGL_NO_SURFACE {
		return EGL_NO_SURFACE, this.GetError()
	}
	return surface, nil
}

// Destroy EGL Surface
func (this *eglDriver) destroySurface(surface eglSurface) error {
	result := C.eglDestroySurface(this.display, C.EGLSurface(surface))
	if result == C.EGLBoolean(EGL_FALSE) {
		return this.GetError()
	}
	return nil
}

// Swap buffer
func (this *eglDriver) swapBuffer(surface eglSurface) error {
	result := C.eglSwapBuffers(this.display, surface)
	if result == C.EGLBoolean(EGL_FALSE) {
		return this.GetError()
	}
	return nil
}
