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
	"strconv"
	"unsafe"
)

import (
	openvg "../../openvg"
)

////////////////////////////////////////////////////////////////////////////////

type EGLDisplay uintptr
type EGLContext uintptr
type EGLSurface uintptr
type VGFloat C.VGfloat

type OpenVG struct {
	VideoCore *VideoCore
}

type OpenVGDriver struct {
	major, minor int
	vc           *VideoCore
	element      *Element
	display      EGLDisplay
	context      EGLContext
	surface		 EGLSurface
}

type nativeWindow struct {
	element ElementHandle
	width int
	height int
}

////////////////////////////////////////////////////////////////////////////////

const (
	EGL_DEFAULT_DISPLAY uintptr    = 0
	EGL_NO_DISPLAY      EGLDisplay = 0
	EGL_NO_CONTEXT      EGLContext = 0
	EGL_NO_SURFACE      EGLSurface = 0
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
	EGL_ALPHA_SIZE   uint16 = 0x3021
	EGL_BLUE_SIZE    uint16 = 0x3022
	EGL_GREEN_SIZE   uint16 = 0x3023
	EGL_RED_SIZE     uint16 = 0x3024
	EGL_SURFACE_TYPE uint16 = 0x3033
	EGL_RENDERABLE_TYPE uint16 = 0x3040
	EGL_NONE         uint16 = 0x3038
)

const (
	// EGL_RENDERABLE_TYPE mask bits
	EGL_OPENGL_ES_BIT uint16 = 0x0001
	EGL_OPENVG_BIT uint16 = 0x0002
	EGL_OPENGL_ES2_BIT uint16 = 0x0004
	EGL_OPENGL_BIT uint16 = 0x0008
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
	// Back buffer swap behaviors
	EGL_BUFFER_PRESERVED uint16 = 0x3094	/* EGL_SWAP_BEHAVIOR value */
	EGL_BUFFER_DESTROYED uint16 = 0x3095	/* EGL_SWAP_BEHAVIOR value */
)

const (
	// Color for vgClear
	VG_CLEAR_COLOR uint16 = 0x1121
)

////////////////////////////////////////////////////////////////////////////////
// openvg.Opener interface

// Concrete Open method
func (config OpenVG) Open() (openvg.Driver, error) {
	driver := new(OpenVGDriver)
	driver.vc = config.VideoCore

	// Get EGL Display
	driver.display = EGLDisplay(unsafe.Pointer(C.eglGetDisplay(EGL_DEFAULT_DISPLAY)))
	if driver.display == EGL_NO_DISPLAY {
		return nil, errors.New("No Display")
	}

	// Initialise
	var major, minor C.EGLint
	result := C.eglInitialize(driver.display, (*C.EGLint)(unsafe.Pointer(&major)), (*C.EGLint)(unsafe.Pointer(&minor)))
	if result == C.EGLBoolean(EGL_FALSE) {
		return nil, errors.New("eglInitialise failed")
	}
	driver.major = int(major)
	driver.minor = int(minor)

	// Bind OpenVG API
	C.eglBindAPI(C.EGLenum(EGL_OPENVG_API))

	// Get EGL frame buffer configuration
	attribute_list := []C.EGLint{
		C.EGLint(EGL_RED_SIZE), C.EGLint(8),
		C.EGLint(EGL_GREEN_SIZE), C.EGLint(8),
		C.EGLint(EGL_BLUE_SIZE), C.EGLint(8),
		C.EGLint(EGL_ALPHA_SIZE), C.EGLint(8),
		C.EGLint(EGL_SURFACE_TYPE), C.EGLint(EGL_WINDOW_BIT),
		C.EGLint(EGL_RENDERABLE_TYPE), C.EGLint(EGL_OPENVG_BIT),
		C.EGLint(EGL_NONE),
	}
	var num_config C.EGLint
	var bufferconfig C.EGLConfig
	result = C.eglChooseConfig(driver.display, &attribute_list[0], &bufferconfig, C.EGLint(1), &num_config)
	if result == C.EGLBoolean(EGL_FALSE) {
		C.eglTerminate(driver.display)
		return nil, errors.New("eglChooseConfig failed")
	}
	if num_config != C.EGLint(1) {
		C.eglTerminate(driver.display)
		return nil, errors.New("eglChooseConfig failed")
	}

	// Create rendering context
	driver.context = EGLContext(C.eglCreateContext(driver.display,bufferconfig,C.EGLContext(EGL_NO_CONTEXT),(*C.EGLint)(unsafe.Pointer(nil))))
	if driver.context == EGL_NO_CONTEXT {
		C.eglTerminate(driver.display)
		return nil, errors.New("eglCreateContext failed")
	}

	// Start update
	update, err := driver.vc.UpdateBegin()
	if err != nil {
		C.eglDestroyContext(driver.display, driver.context)
		C.eglTerminate(driver.display)
		return nil,err
	}

	// Add element on Layer 0, Stretch to fill screen
	var dst_rect = &Rectangle{}
	var src_rect = &Rectangle{}
	dst_rect.Set(Point{ 0, 0 },driver.vc.GetSize())
	src_rect.Set(Point{ 0, 0 },Size{ dst_rect.Size.Width << 16, dst_rect.Size.Height << 16 })
	driver.element, err = driver.vc.AddElement(update, 0, dst_rect, nil, src_rect)
	if err != nil {
		C.eglDestroyContext(driver.display, driver.context)
		C.eglTerminate(driver.display)
		return nil,err
	}

	// Update
	driver.vc.UpdateSubmit(update)

	// Connect window surface
	nativewindow := nativeWindow{ driver.element.GetHandle(), int(dst_rect.Size.Width), int(dst_rect.Size.Height) }
	driver.surface = EGLSurface(C.eglCreateWindowSurface(driver.display,bufferconfig,(*C.EGLNativeWindowType)(unsafe.Pointer(&nativewindow)),nil))
	if driver.surface == EGL_NO_SURFACE {
		C.eglDestroyContext(driver.display, driver.context)
		C.eglTerminate(driver.display)
		return nil,errors.New("eglCreateWindowSurface failed")
	}

	// Preserve the buffers on swap
	result = C.eglSurfaceAttrib(driver.display,driver.surface, C.EGLint(EGL_SWAP_BEHAVIOR), C.EGLint(EGL_BUFFER_PRESERVED));
	if result == C.EGLBoolean(EGL_FALSE) {
		C.eglDestroyContext(driver.display, driver.context)
		C.eglTerminate(driver.display)
		return nil, errors.New("eglSurfaceAttrib failed")
	}

	// Connect the context to the surface
	result = C.eglMakeCurrent(driver.display,driver.surface,driver.surface,driver.context)
	if result == C.EGLBoolean(EGL_FALSE) {
		C.eglDestroyContext(driver.display, driver.context)
		C.eglTerminate(driver.display)
		return nil, errors.New("eglMakeCurrent failed")
	}

	// clear
    clearColor := []VGFloat{ 0.8, 0.3, 0.4, 1.0 }
    C.vgSetfv(C.VGParamType(VG_CLEAR_COLOR), 4, (*C.VGfloat)(unsafe.Pointer(&clearColor[0])));
    C.vgClear(0, 0, 1920, 1080);
    C.vgFlush();

	return driver, nil
}

////////////////////////////////////////////////////////////////////////////////
// openvg.Driver interface

func (this *OpenVGDriver) Close() error {
	// Close rendering context
	result := C.eglDestroyContext(this.display, this.context)
	if result == C.EGLBoolean(EGL_FALSE) {
		C.eglTerminate(this.display)
		return errors.New("eglDestroyContext failed")
	}

	// Terminate display
	result = C.eglTerminate(this.display)
	if result == C.EGLBoolean(EGL_FALSE) {
		return errors.New("eglTerminate failed")
	}

	return nil
}

func (this *OpenVGDriver) String() string {
	return "<OpenVGDriver>{ version=" + strconv.Itoa(this.major) + "." + strconv.Itoa(this.minor) + " }"
}
