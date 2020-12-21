// +build egl

package egl_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	egl "github.com/djthorpe/gopi/v3/pkg/sys/egl"
	gbm "github.com/djthorpe/gopi/v3/pkg/sys/gbm"
)

const (
	GPU_NODE = "card1"
)

////////////////////////////////////////////////////////////////////////////////
// TEST EGL INIT

func Test_EGL_000(t *testing.T) {
	fh, err := gbm.OpenDevice(GPU_NODE)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	display := gbm.GBMCreateDevice(fh.Fd())
	if display == nil {
		t.Fatal(err)
	}
	defer display.Free()
	if major, minor, err := egl.EGLInitialize(egl.EGLGetDisplay(display)); err != nil {
		t.Error(err)
	} else if err := egl.EGLTerminate(egl.EGLGetDisplay(display)); err != nil {
		t.Error(err)
	} else {
		t.Log("egl_version=", major, ".", minor)
	}
}

func Test_EGL_002(t *testing.T) {
	fh, err := gbm.OpenDevice(GPU_NODE)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	device := gbm.GBMCreateDevice(fh.Fd())
	if device == nil {
		t.Fatal(err)
	}
	defer device.Free()
	display := egl.EGLGetDisplay(device)
	major, minor, err := egl.EGLInitialize(display)
	if err != nil {
		t.Fatal(err)
	}

	if vendor := egl.EGLQueryString(display, egl.EGL_QUERY_VENDOR); vendor == "" {
		t.Error("Empty value returned for EGL_QUERY_VENDOR")
	} else if version := egl.EGLQueryString(display, egl.EGL_QUERY_VERSION); version == "" {
		t.Error("Empty value returned for EGL_QUERY_VERSION")
	} else if extensions := egl.EGLQueryString(display, egl.EGL_QUERY_EXTENSIONS); extensions == "" {
		t.Error("Empty value returned for EGL_QUERY_EXTENSIONS")
	} else if apis := egl.EGLQueryString(display, egl.EGL_QUERY_CLIENT_APIS); extensions == "" {
		t.Error("Empty value returned for EGL_QUERY_CLIENT_APIS")
	} else {
		t.Log("display=", display)
		t.Logf("egl_version= %v,%v", major, minor)
		t.Logf("EGL_QUERY_VENDOR= %v", vendor)
		t.Logf("EGL_QUERY_VERSION= %v", version)
		t.Logf("EGL_QUERY_EXTENSIONS= %v", extensions)
		t.Logf("EGL_QUERY_CLIENT_APIS= %v", apis)
	}

	if err := egl.EGLTerminate(display); err != nil {
		t.Error(err)
	}

}

func Test_EGL_003(t *testing.T) {
	fh, err := gbm.OpenDevice(GPU_NODE)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	device := gbm.GBMCreateDevice(fh.Fd())
	if device == nil {
		t.Fatal(err)
	}
	defer device.Free()
	display := egl.EGLGetDisplay(device)

	if _, _, err := egl.EGLInitialize(display); err != nil {
		t.Fatal(err)
	}

	if configs, err := egl.EGLGetConfigs(display); err != nil {
		t.Error(err)
	} else if len(configs) == 0 {
		t.Error("Expected at least one config")
	} else {
		t.Log("display=", display)
		t.Logf("configs= %v", configs)
	}

	if err := egl.EGLTerminate(display); err != nil {
		t.Error(err)
	}
}

func Test_EGL_004(t *testing.T) {
	fh, err := gbm.OpenDevice(GPU_NODE)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	device := gbm.GBMCreateDevice(fh.Fd())
	if device == nil {
		t.Fatal(err)
	}
	defer device.Free()
	display := egl.EGLGetDisplay(device)

	if _, _, err := egl.EGLInitialize(display); err != nil {
		t.Fatal(err)
	}

	if configs, err := egl.EGLGetConfigs(display); err != nil {
		t.Error(err)
	} else if len(configs) == 0 {
		t.Error("Expected at least one config")
	} else if attributes, err := egl.EGLGetConfigAttribs(display, configs[0]); err != nil {
		t.Error(err)
	} else {
		t.Logf("attributes[0]= %v", attributes)
	}

	if err := egl.EGLTerminate(display); err != nil {
		t.Error(err)
	}
}

func Test_EGL_005(t *testing.T) {
	// RGBA32

	fh, err := gbm.OpenDevice(GPU_NODE)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	device := gbm.GBMCreateDevice(fh.Fd())
	if device == nil {
		t.Fatal(err)
	}
	defer device.Free()
	display := egl.EGLGetDisplay(device)

	if _, _, err := egl.EGLInitialize(display); err != nil {
		t.Fatal(err)
	}

	if config, err := egl.EGLChooseConfig(display, 8, 8, 8, 8, 0, 0); err != nil {
		t.Error(err)
	} else if attributes, err := egl.EGLGetConfigAttribs(display, config); err != nil {
		t.Error(err)
	} else {
		t.Logf("RGBA32 attributes= %v", attributes)
	}

	if err := egl.EGLTerminate(display); err != nil {
		t.Error(err)
	}
}

func Test_EGL_006(t *testing.T) {
	// RGB565

	fh, err := gbm.OpenDevice(GPU_NODE)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	device := gbm.GBMCreateDevice(fh.Fd())
	if device == nil {
		t.Fatal(err)
	}
	defer device.Free()
	display := egl.EGLGetDisplay(device)

	if _, _, err := egl.EGLInitialize(display); err != nil {
		t.Fatal(err)
	}

	if config, err := egl.EGLChooseConfig(display, 5, 6, 5, 0, 0, 0); err != nil {
		t.Error(err)
	} else if attributes, err := egl.EGLGetConfigAttribs(display, config); err != nil {
		t.Error(err)
	} else {
		t.Logf("RGB565 attributes= %v", attributes)
	}

	if err := egl.EGLTerminate(display); err != nil {
		t.Error(err)
	}
}

func Test_EGL_007(t *testing.T) {
	// RGB888

	fh, err := gbm.OpenDevice(GPU_NODE)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	device := gbm.GBMCreateDevice(fh.Fd())
	if device == nil {
		t.Fatal(err)
	}
	defer device.Free()
	display := egl.EGLGetDisplay(device)

	if _, _, err := egl.EGLInitialize(display); err != nil {
		t.Fatal(err)
	}

	if config, err := egl.EGLChooseConfig(display, 8, 8, 8, 0, 0, 0); err != nil {
		t.Error(err)
	} else if attributes, err := egl.EGLGetConfigAttribs(display, config); err != nil {
		t.Error(err)
	} else {
		t.Logf("RGB888 attributes= %v", attributes)
	}

	if err := egl.EGLTerminate(display); err != nil {
		t.Error(err)
	}
}

func Test_EGL_008(t *testing.T) {
	for api := egl.EGL_API_MIN; api <= egl.EGL_API_MAX; api++ {
		api_string := fmt.Sprint(api)
		if strings.HasPrefix(api_string, "EGL_API_") == true {
			t.Logf("%v => %v", api, api_string)
		} else {
			t.Errorf("Error for %v => %v", api, api_string)
		}
	}
}

func Test_EGL_009(t *testing.T) {

	fh, err := gbm.OpenDevice(GPU_NODE)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	device := gbm.GBMCreateDevice(fh.Fd())
	if device == nil {
		t.Fatal(err)
	}
	defer device.Free()
	display := egl.EGLGetDisplay(device)

	if _, _, err := egl.EGLInitialize(display); err != nil {
		t.Fatal(err)
	}

	types := strings.Split(egl.EGLQueryString(display, egl.EGL_QUERY_CLIENT_APIS), " ")
	for _, api_string := range types {
		if api_string == "" {
			continue
		}
		if surface_type, exists := egl.EGLSurfaceTypeMap[api_string]; exists == false {
			t.Error("Does not exist in EGLSurfaceTypeMap:", strconv.Quote(api_string))
		} else if api, exists := egl.EGLAPIMap[surface_type]; exists == false {
			t.Error("Does not exist in EGL_APIMap:", api_string)
		} else if renderable, exists := egl.EGLRenderableMap[surface_type]; exists == false {
			t.Error("Does not exist in EGLRenderable_Map:", api_string)
		} else if err := egl.EGLBindAPI(api); err != nil {
			t.Error("Error in EGLBindAPI:", err)
		} else if api_, err := egl.EGLQueryAPI(); err != nil {
			t.Error(err)
		} else if api != api_ {
			t.Error("Unexpected mismatch", api, api_)
		} else {
			t.Logf("%v => %v => %v => %v, %v", api_string, surface_type, api, api_, renderable)
		}
	}

	if err := egl.EGLTerminate(display); err != nil {
		t.Error(err)
	}

}
