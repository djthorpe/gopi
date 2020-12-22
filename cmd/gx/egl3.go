package main

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: glesv2
#include <GLES3/gl3.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type egl3 struct {
	vertexShader, fragmentShader C.GLuint
	program                      C.GLuint
}

type Drawable interface {
	Draw()
	Dispose()
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	vShaderStr = `#version 300 es
#extension GL_EXT_separate_shader_objects : require
layout(location = 0) in vec4 vPosition;
void main() 
{
	gl_Position = vPosition;
}
`

	fShaderStr = `#version 300 es
precision mediump float;
out vec4 fragColor;
void main() 
{
	fragColor = vec4 ( 1.0, 0.0, 0.0, 1.0 );
}
`
)

////////////////////////////////////////////////////////////////////////////////
// METHODS

func NewEGL3Drawable() (Drawable, error) {
	this := new(egl3)
	if shader, err := this.LoadShader(C.GL_VERTEX_SHADER, vShaderStr); err != nil {
		return nil, err
	} else {
		this.vertexShader = shader
	}
	if shader, err := this.LoadShader(C.GL_FRAGMENT_SHADER, fShaderStr); err != nil {
		return nil, err
	} else {
		this.fragmentShader = shader
	}

	this.program = C.glCreateProgram()
	if this.program == 0 {
		return nil, fmt.Errorf("Unable to create program")
	}

	C.glAttachShader(this.program, this.vertexShader)
	C.glAttachShader(this.program, this.fragmentShader)
	C.glLinkProgram(this.program)

	linked := C.GLint(0)
	C.glGetProgramiv(this.program, C.GL_LINK_STATUS, &linked)
	if linked == C.GLint(0) {
		C.glDeleteProgram(this.program)
		return nil, fmt.Errorf("Unable to link program")
	}

	C.glClearColor(C.float(0), C.float(0), C.float(0), C.float(0))

	return this, nil
}

func (this *egl3) Dispose() {
	if this.program != 0 {
		C.glDeleteProgram(this.program)
	}
	if this.vertexShader != 0 {
		C.glDeleteShader(this.vertexShader)
	}
	if this.fragmentShader != 0 {
		C.glDeleteShader(this.fragmentShader)
	}
}

func (this *egl3) Draw() {
	vVertices := []C.GLfloat{
		C.float(0), C.float(0.5), C.float(0),
		C.float(-0.5), C.float(-0.5), C.float(0),
		C.float(0.5), C.float(-0.5), C.float(0),
	}

	C.glViewport(0, 0, 1920, 1080)
	C.glClear(C.GL_COLOR_BUFFER_BIT)

	// Use the program object
	C.glUseProgram(this.program)

	// Load the vertex data
	data := unsafe.Pointer(&vVertices[0])
	C.glVertexAttribPointer(0, 3, C.GL_FLOAT, C.GL_FALSE, 0, data)
	C.glEnableVertexAttribArray(0)
	C.glDrawArrays(C.GL_TRIANGLES, 0, 3)
}

func (this *egl3) LoadShader(shaderType C.GLenum, source string) (C.GLuint, error) {
	shader := C.glCreateShader(shaderType)
	if shader == 0 {
		return 0, fmt.Errorf("Unable to create shader")
	}

	cSource := C.CString(source)
	defer C.free(unsafe.Pointer(cSource))

	// Load the shader source
	compiled := C.GLint(0)
	C.glShaderSource(shader, 1, &cSource, nil)
	C.glCompileShader(shader)
	C.glGetShaderiv(shader, C.GL_COMPILE_STATUS, &compiled)
	if compiled == 0 {
		var err error
		infoLen := C.GLint(0)
		C.glGetShaderiv(shader, C.GL_INFO_LOG_LENGTH, &infoLen)
		if infoLen > 0 {
			data := make([]byte, int(infoLen))
			ptr := unsafe.Pointer(&data[0])
			C.glGetShaderInfoLog(shader, infoLen, nil, (*C.char)(ptr))
			err = fmt.Errorf("LoadShader: %v", string(data))
		} else {
			err = fmt.Errorf("Unable to create shader")
		}
		C.glDeleteShader(shader)
		return 0, err
	} else {
		return shader, nil
	}
}
