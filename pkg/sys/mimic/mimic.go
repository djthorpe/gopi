// +build mimic

package mimic

import (
	"errors"
	"fmt"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mimic
#include <ttsmimic/mimic.h>
#include <stdlib.h>

void mimic_set_lang_list(void);
cst_val *mimic_set_voice_list(const char *voxdir);
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	Value C.cst_val
)

type (
	ValueType C.int
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	CST_VAL_TYPE_CONS   ValueType = 0
	CST_VAL_TYPE_INT    ValueType = 1
	CST_VAL_TYPE_FLOAT  ValueType = 3
	CST_VAL_TYPE_STRING ValueType = 5
)

var (
	errError = errors.New("General Error")
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func Init() {
	C.mimic_init()
	C.mimic_set_lang_list()
}

func Exit() {
	C.mimic_exit()
}

func SetVoiceList(path string) (*Value, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	if value := C.mimic_set_voice_list(cPath); value == nil {
		return nil, errError
	} else {
		fmt.Println(value)
		return (*Value)(value), nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// VALUE

func (this *Value) Type() ValueType {
	ctx := (*C.int)(unsafe.Pointer(this))
	return ValueType(*ctx)
}

func (this *Value) Ref() int {
	ctx := (*C.int)(unsafe.Pointer(this))
	return int(*ctx)
}

func (this *Value) String() string {
	str := "<value"
	if t := this.Type(); t != 0 {
		str += fmt.Sprint(" type=", t)
		str += fmt.Sprint(" ref=", this.Ref())
	}
	return str + ">"
}
