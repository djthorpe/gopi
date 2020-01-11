// +build egl

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package egl_test

import (
	"sync"
	"testing"

	// Frameworks
	egl "github.com/djthorpe/gopi/v2/sys/egl"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

var (
	dxinit sync.Once
)

func DXInit(t *testing.T) {
	dxinit.Do(func() {
		rpi.DXInit()
	})
}

////////////////////////////////////////////////////////////////////////////////
// TEST EGL INIT

func Test_EGL_000(t *testing.T) {
	DXInit(t)
	if display, err := rpi.DXDisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if err := rpi.DXDisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log(display)
	}
}

func Test_EGL_001(t *testing.T) {
	DXInit(t)
	if display, err := rpi.DXDisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if major, minor, err := egl.EGLInitialize(egl.EGLGetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD))); err != nil {
		t.Error(err)
	} else if err := egl.EGLTerminate(egl.EGLGetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD))); err != nil {
		t.Error(err)
	} else if err := rpi.DXDisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log("display=", display)
		t.Log("egl_version=", major, ".", minor)
	}
}

func Test_EGL_002(t *testing.T) {
	DXInit(t)
	if display, err := rpi.DXDisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if handle := egl.EGLGetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD)); handle == 0 {
		t.Error("EGLGetDisplay returned nil")
	} else if major, minor, err := egl.EGLInitialize(handle); err != nil {
		t.Error(err)
	} else if vendor := egl.EGLQueryString(handle, egl.EGL_QUERY_VENDOR); vendor == "" {
		t.Error("Empty value returned for EGL_QUERY_VENDOR")
	} else if version := egl.EGLQueryString(handle, egl.EGL_QUERY_VERSION); version == "" {
		t.Error("Empty value returned for EGL_QUERY_VERSION")
	} else if extensions := egl.EGLQueryString(handle, egl.EGL_QUERY_EXTENSIONS); extensions == "" {
		t.Error("Empty value returned for EGL_QUERY_EXTENSIONS")
	} else if apis := egl.EGLQueryString(handle, egl.EGL_QUERY_CLIENT_APIS); extensions == "" {
		t.Error("Empty value returned for EGL_QUERY_CLIENT_APIS")
	} else if err := egl.EGLTerminate(handle); err != nil {
		t.Error(err)
	} else if err := rpi.DXDisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log("display=", display)
		t.Logf("egl_version= %v,%v", major, minor)
		t.Logf("EGL_QUERY_VENDOR= %v", vendor)
		t.Logf("EGL_QUERY_VERSION= %v", version)
		t.Logf("EGL_QUERY_EXTENSIONS= %v", extensions)
		t.Logf("EGL_QUERY_CLIENT_APIS= %v", apis)
	}
}

func Test_EGL_003(t *testing.T) {
	DXInit(t)
	if display, err := rpi.DXDisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if handle := egl.EGLGetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD)); handle == 0 {
		t.Error("EGL_GetDisplay returned nil")
	} else if major, minor, err := egl.EGLInitialize(handle); err != nil {
		t.Error(err)
	} else if configs, err := egl.EGLGetConfigs(handle); err != nil {
		t.Error(err)
	} else if len(configs) == 0 {
		t.Error("Expected at least one config")
	} else if err := egl.EGLTerminate(handle); err != nil {
		t.Error(err)
	} else if err := rpi.DXDisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log("display=", display)
		t.Logf("egl_version= %v,%v", major, minor)
		t.Logf("configs= %v", configs)
	}
}

func Test_EGL_004(t *testing.T) {
	DXInit(t)
	if display, err := rpi.DXDisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if handle := egl.EGLGetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD)); handle == 0 {
		t.Error("EGL_GetDisplay returned nil")
	} else if major, minor, err := egl.EGLInitialize(handle); err != nil {
		t.Error(err)
	} else if configs, err := egl.EGLGetConfigs(handle); err != nil {
		t.Error(err)
	} else if len(configs) == 0 {
		t.Error("Expected at least one config")
	} else if attributes, err := egl.EGLGetConfigAttribs(handle, configs[0]); err != nil {
		t.Error(err)
	} else if err := egl.EGLTerminate(handle); err != nil {
		t.Error(err)
	} else if err := rpi.DXDisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log("display=", display)
		t.Logf("egl_version= %v,%v", major, minor)
		t.Logf("attributes[0]= %v", attributes)
	}
}

func Test_EGL_005(t *testing.T) {
	DXInit(t)

	// RGBA32
	if display, err := rpi.DXDisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if handle := egl.EGLGetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD)); handle == 0 {
		t.Error("EGL_GetDisplay returned nil")
	} else if major, minor, err := egl.EGLInitialize(handle); err != nil {
		t.Error(err)
	} else if config, err := egl.EGLChooseConfig(handle, 8, 8, 8, 8, 0, 0); err != nil {
		t.Error(err)
	} else if attributes, err := egl.EGLGetConfigAttribs(handle, config); err != nil {
		t.Error(err)
	} else if err := egl.EGLTerminate(handle); err != nil {
		t.Error(err)
	} else if err := rpi.DXDisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log("display=", display)
		t.Logf("egl_version= %v,%v", major, minor)
		t.Logf("attributes= %v", attributes)
	}
}

func Test_EGL_006(t *testing.T) {
	DXInit(t)

	// RGB565
	if display, err := rpi.DXDisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if handle := egl.EGLGetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD)); handle == 0 {
		t.Error("EGL_GetDisplay returned nil")
	} else if major, minor, err := egl.EGLInitialize(handle); err != nil {
		t.Error(err)
	} else if config, err := egl.EGLChooseConfig(handle, 5, 6, 5, 0, 0, 0); err != nil {
		t.Error(err)
	} else if attributes, err := egl.EGLGetConfigAttribs(handle, config); err != nil {
		t.Error(err)
	} else if err := egl.EGLTerminate(handle); err != nil {
		t.Error(err)
	} else if err := rpi.DXDisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log("display=", display)
		t.Logf("egl_version= %v,%v", major, minor)
		t.Logf("attributes= %v", attributes)
	}
}

func Test_EGL_007(t *testing.T) {
	DXInit(t)

	// RGB888
	if display, err := rpi.DXDisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if handle := egl.EGLGetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD)); handle == 0 {
		t.Error("EGL_GetDisplay returned nil")
	} else if major, minor, err := egl.EGLInitialize(handle); err != nil {
		t.Error(err)
	} else if config, err := egl.EGLChooseConfig(handle, 8, 8, 8, 0, 0, 0); err != nil {
		t.Error(err)
	} else if attributes, err := egl.EGLGetConfigAttribs(handle, config); err != nil {
		t.Error(err)
	} else if err := egl.EGLTerminate(handle); err != nil {
		t.Error(err)
	} else if err := rpi.DXDisplayClose(display); err != nil {
		t.Error(err)
	} else {
		t.Log("display=", display)
		t.Logf("egl_version= %v,%v", major, minor)
		t.Logf("attributes= %v", attributes)
	}
}

/*
func Test_EGL_008(t *testing.T) {
	DXInit(t)

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
	DXInit(t)

	rpi.DX_Init()
	if display, err := rpi.DX_DisplayOpen(rpi.DX_DISPLAYID_MAIN_LCD); err != nil {
		t.Error(err)
	} else if handle := egl.EGL_GetDisplay(uint(rpi.DX_DISPLAYID_MAIN_LCD)); handle == nil {
		t.Error("EGL_GetDisplay returned nil")
	} else if _, _, err := egl.EGL_Initialize(handle); err != nil {
		t.Error(err)
	} else {
		types := strings.Split(egl.EGL_QueryString(handle, egl.EGL_QUERY_CLIENT_APIS), " ")
		for _, api_string := range types {
			if surface_type, exists := egl.EGL_SurfaceTypeMap[api_string]; exists == false {
				t.Error("Does not exist in EGL_SurfaceTypeMap:", api_string)
			} else if api, exists := egl.EGL_APIMap[surface_type]; exists == false {
				t.Error("Does not exist in EGL_APIMap:", api_string)
			} else if renderable, exists := egl.EGL_RenderableMap[surface_type]; exists == false {
				t.Error("Does not exist in EGL_Renderable_Map:", api_string)
			} else if err := egl.EGL_BindAPI(api); err != nil {
				t.Error("Error in EGL_BindAPI:", err)
			} else if api_, err := egl.EGL_QueryAPI(); err != nil {
				t.Error(err)
			} else if api != api_ {
				t.Error("Unexpected mismatch", api, api_)
			} else {
				t.Logf("%v => %v => %v => %v, %v", api_string, surface_type, api, api_, renderable)
			}
		}
		if err := egl.EGL_Terminate(handle); err != nil {
			t.Error(err)
		} else if err := rpi.DX_DisplayClose(display); err != nil {
			t.Error(err)
		}
	}
}
*/
