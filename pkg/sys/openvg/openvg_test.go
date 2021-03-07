//+build openvg,rpi,egl

package openvg_test

import (
	"testing"

	egl "github.com/djthorpe/gopi/v3/pkg/sys/egl"
	openvg "github.com/djthorpe/gopi/v3/pkg/sys/openvg"
	rpi "github.com/djthorpe/gopi/v3/pkg/sys/rpi"
)

type Context struct {
	display egl.EGLDisplay
	config  egl.EGLConfig
	context egl.EGLContext
	surface egl.EGLSurface
}

func NewContext(t *testing.T, w, h uint) *Context {
	this := new(Context)
	if err := rpi.BCMHostInit(); err != nil {
		t.Error("BCMHostInit Error:", err)
		return nil
	}
	if display := egl.EGLGetDisplay(0); display == 0 {
		t.Error("EGLGetDisplay Error")
		return nil
	} else {
		this.display = display
	}
	if _, _, err := egl.EGLInitialize(this.display); err != nil {
		t.Error("EGLInitialize Error:", err)
		return nil
	}
	if config, err := egl.EGLChooseConfig(this.display, 8, 8, 8, 0, egl.EGL_SURFACETYPE_FLAG_WINDOW|egl.EGL_SURFACETYPE_FLAG_PBUFFER, egl.EGL_RENDERABLE_FLAG_OPENVG); err != nil {
		t.Error("EGLChooseConfig Error:", err)
		if err := egl.EGLTerminate(this.display); err != nil {
			t.Error("EGLTerminate error", err)
		}
		return nil
	} else {
		this.config = config
	}

	if err := egl.EGLBindAPI(egl.EGL_API_OPENVG); err != nil {
		t.Error("EGLBindAPI Error:", err)
		if err := egl.EGLTerminate(this.display); err != nil {
			t.Error("EGLTerminate error", err)
		}
		return nil
	}

	if context, err := egl.EGLCreateContext(this.display, this.config, nil, nil); err != nil {
		t.Error("EGLCreateContext Error:", err)
		if err := egl.EGLTerminate(this.display); err != nil {
			t.Error("EGLTerminate error", err)
		}
		return nil
	} else {
		this.context = context
	}
	if surface, err := egl.EGLCreatePbufferSurface(this.display, this.config, 256, 256); err != nil {
		t.Error("EGLCreatePbufferSurface Error:", err)
		if err := egl.EGLTerminate(this.display); err != nil {
			t.Error("EGLTerminate error", err)
		}
		return nil
	} else {
		this.surface = surface
	}

	if err := egl.EGLMakeCurrent(this.display, this.surface, this.surface, this.context); err != nil {
		t.Error("EGLMakeCurrent Error:", err)
		if err := egl.EGLTerminate(this.display); err != nil {
			t.Error("EGLTerminate error", err)
		}
		return nil
	}

	// Return success
	return this
}

func (this *Context) Dispose(t *testing.T) {
	if err := egl.EGLDestroySurface(this.display, this.surface); err != nil {
		t.Error("EGLDestroySurface error", err)
	}
	if err := egl.EGLDestroyContext(this.display, this.context); err != nil {
		t.Error("EGLDestroyContext error", err)
	}
	if err := egl.EGLTerminate(this.display); err != nil {
		t.Error("EGLTerminate error", err)
	}
}

func Test_OpenVG_001(t *testing.T) {
	ctx := NewContext(t, 100, 100)
	if ctx == nil {
		t.Fatal("NewContext failed")
	}
	defer ctx.Dispose(t)

	for _, query := range []openvg.QueryString{openvg.VG_VENDOR, openvg.VG_RENDERER, openvg.VG_VERSION, openvg.VG_EXTENSIONS} {
		if result, err := openvg.GetString(query); err != nil {
			t.Error("Unexpected nil return: ", query, " Error: ", err)
		} else {
			t.Log(query, "=", result)
		}
	}
}
