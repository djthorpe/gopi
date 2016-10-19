/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

package rpi
/*
  #cgo CFLAGS:   -I/opt/vc/include
  #cgo LDFLAGS:  -L/opt/vc/lib -lEGL -lGLESv2
  #include <EGL/egl.h>
  #include <VG/openvg.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"strings"
	"unsafe"
)

import (
	khronos "../../khronos"
	util "../../util"
)

////////////////////////////////////////////////////////////////////////////////

type EGLDisplay uintptr
type EGLContext uintptr
type EGLSurface uintptr
type EGLConfig uintptr

type EGL struct {
	VideoCore *VideoCore
	Logger    util.Logger
}

type EGLState struct {
	major, minor int
	vc           *VideoCore
	logger       util.Logger
	display      EGLDisplay
}

////////////////////////////////////////////////////////////////////////////////

const (
	// General constants
	EGL_DEFAULT_DISPLAY uintptr    = 0
	EGL_NO_DISPLAY      EGLDisplay = 0
	EGL_NO_CONTEXT      EGLContext = 0
	EGL_NO_SURFACE      EGLSurface = 0
	EGL_NO_CONFIG       EGLConfig  = 0
	EGL_FALSE           uint       = 0
	EGL_TRUE            uint       = 1
)

const (
	// BindAPI/QueryAPI targets
	EGL_OPENGL_ES_API uint16 = 0x30A0
	EGL_OPENVG_API    uint16 = 0x30A1
	EGL_OPENGL_API    uint16 = 0x30A2
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
	// EGL_RENDERABLE_TYPE mask bits
	EGL_OPENGL_ES_BIT  uint16 = 0x0001
	EGL_OPENVG_BIT     uint16 = 0x0002
	EGL_OPENGL_ES2_BIT uint16 = 0x0004
	EGL_OPENGL_BIT     uint16 = 0x0008
)

const (
	// Config attribute mask bits
	EGL_WINDOW_BIT uint16 = 0x0004
)

const (
	// QuerySurface / SurfaceAttrib / CreatePbufferSurface targets
	EGL_SWAP_BEHAVIOR uint16 = 0x3093
)

const (
	// QueryString targets
	EGL_VENDOR      uint16 = 0x3053
	EGL_VERSION     uint16 = 0x3054
	EGL_EXTENSIONS  uint16 = 0x3055
	EGL_CLIENT_APIS uint16 = 0x308D
)

const (
	// Back buffer swap behaviors
	EGL_BUFFER_PRESERVED uint16 = 0x3094 /* EGL_SWAP_BEHAVIOR value */
	EGL_BUFFER_DESTROYED uint16 = 0x3095 /* EGL_SWAP_BEHAVIOR value */
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

////////////////////////////////////////////////////////////////////////////////

var (
	ErrorNoDisplay                = errors.New("No display")
	ErrorInitialiseFailed         = errors.New("eglInitialise failed")
	ErrorTerminateFailed          = errors.New("eglTerminate failed")
	ErrorInvalidAPIBind           = errors.New("Invalid API Bind")
	ErrorInvalidFrameBufferConfig = errors.New("Invalid Frame Buffer Configuration")
	ErrorContextError             = errors.New("EGLContext error")
	ErrorSurfaceError             = errors.New("EGLSurface error")
	ErrorUnknown                  = errors.New("EGL Undefined error")
)

var (
	ClientAPI = map[string]C.EGLenum{
		"OpenGL_ES": C.EGLenum(EGL_OPENGL_ES_API),
		"OpenVG":    C.EGLenum(EGL_OPENVG_API),
		"OpenGL":    C.EGLenum(EGL_OPENGL_API),
	}
	// EGL_RENDERABLE_TYPE
	RenderableType = map[string]uint16{
		"OpenGL":    EGL_OPENGL_BIT,
		"OpenVG":    EGL_OPENVG_BIT,
		"OpenGL_ES": EGL_OPENGL_ES2_BIT,
	}
	// Errors
	EGLError = map[uint16]error{
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
)

////////////////////////////////////////////////////////////////////////////////
// Opener interface

// Initialise the EGL interface
func (config EGL) Open() (khronos.EGLDriver, error) {
	this := &EGLState{vc: config.VideoCore, logger: config.Logger}

	// Get EGL Display
	this.display = EGLDisplay(unsafe.Pointer(C.eglGetDisplay(EGL_DEFAULT_DISPLAY)))
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

	return this, nil
}

// Terminate the EGL interface
func (this *EGLState) Close() error {
	result := C.eglTerminate(this.display)
	if result == C.EGLBoolean(EGL_FALSE) {
		return this.GetError()
	}
	return nil
}

// Return major and minor version of EGL
func (this *EGLState) GetVersion() (int, int) {
	return this.major, this.minor
}

// Return vendor information
func (this *EGLState) GetVendorString() string {
	return C.GoString(C.eglQueryString(this.display, C.EGLint(EGL_VENDOR)))
}

// Return vendor information
func (this *EGLState) GetVersionString() string {
	return C.GoString(C.eglQueryString(this.display, C.EGLint(EGL_VERSION)))
}

// Return extensions information
func (this *EGLState) GetExtensions() []string {
	return strings.Split(C.GoString(C.eglQueryString(this.display, C.EGLint(EGL_EXTENSIONS))), " ")
}

// Return API's information
func (this *EGLState) GetSupportedClientAPIs() []string {
	return strings.Split(C.GoString(C.eglQueryString(this.display, C.EGLint(EGL_CLIENT_APIS))), " ")
}

// Bind API
func (this *EGLState) BindAPI(api string) error {
	value, ok := ClientAPI[api]
	if !ok {
		return ErrorInvalidAPIBind
	}
	result := C.eglBindAPI(value)
	if result == C.EGLBoolean(EGL_FALSE) {
		return this.GetError()
	}
	return nil
}

// Query currently bound API
func (this *EGLState) QueryAPI() (string, error) {
	api := C.eglQueryAPI()
	for k, v := range ClientAPI {
		if api == v {
			return k, nil
		}
	}
	return "", ErrorInvalidAPIBind
}

// Return string version of the EGL interface
func (this *EGLState) String() string {
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

	return "<ConcreteEGLState>{ " + strings.Join(parts, " ") + "}"
}

// Return error
func (this *EGLState) GetError() error {
	code := C.eglGetError()
	err, ok := EGLError[uint16(code)]
	if !ok {
		return ErrorUnknown
	}
	return err
}

// Return framesize of the display
func (this *EGLState) GetFrame() *khronos.EGLFrame {
	size := this.vc.GetSize()
	return &khronos.EGLFrame{ khronos.EGLPoint{ }, khronos.EGLSize{ uint(size.Width), uint(size.Height)  } }
}

////////////////////////////////////////////////////////////////////////////////
// private methods

// Choose EGL frame buffer configuration
func (this *EGLState) getFrameBufferConfiguration() (EGLConfig, error) {
	// Get bound API
	api, err := this.QueryAPI()
	if err != nil {
		return EGL_NO_CONFIG, err
	}
	// Get Renderable type depending on the API
	renderable_type, ok := RenderableType[api]
	if !ok {
		return EGL_NO_CONFIG, ErrorInvalidFrameBufferConfig
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
	return EGLConfig(config), nil
}

// Create EGL Context with API value
func (this *EGLState) createContext(api string) (EGLConfig, EGLContext, error) {
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
	context := EGLContext(C.eglCreateContext(this.display, config, C.EGLContext(EGL_NO_CONTEXT), (*C.EGLint)(unsafe.Pointer(nil))))
	if context == EGL_NO_CONTEXT {
		return EGL_NO_CONFIG, EGL_NO_CONTEXT, this.GetError()
	}

	return config, context, nil
}

// Destroy EGL Context
func (this *EGLState) destroyContext(context EGLContext) error {
	result := C.eglDestroyContext(this.display, C.EGLContext(context))
	if result == C.EGLBoolean(EGL_FALSE) {
		return this.GetError()
	}
	return nil
}

// Create surface
func (this *EGLState) createSurface(config EGLConfig, window *EGLNativeWindow) (EGLSurface, error) {
	// Create EGL window surface given a native window
	surface := EGLSurface(C.eglCreateWindowSurface(this.display, C.EGLConfig(config), (*C.EGLNativeWindowType)(unsafe.Pointer(window)), nil))
	if surface == EGL_NO_SURFACE {
		return EGL_NO_SURFACE, this.GetError()
	}
	return surface, nil
}

// Destroy EGL Surface
func (this *EGLState) destroySurface(surface EGLSurface) error {
	result := C.eglDestroySurface(this.display, C.EGLSurface(surface))
	if result == C.EGLBoolean(EGL_FALSE) {
		return this.GetError()
	}
	return nil
}

// Attach context to surface
func (this *EGLState) attachContextToSurface(context EGLContext, surface EGLSurface) error {
	result := C.eglMakeCurrent(this.display, C.EGLSurface(surface), C.EGLSurface(surface), C.EGLContext(context))
	if result == C.EGLBoolean(EGL_FALSE) {
		return this.GetError()
	}
	return nil
}

// Swap buffer
func (this *EGLState) swapBuffer(surface EGLSurface) error {
	result := C.eglSwapBuffers(this.display, surface);
	if result == C.EGLBoolean(EGL_FALSE) {
		return this.GetError()
	}
	return nil
}

