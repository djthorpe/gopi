// +build ffmpeg

package ffmpeg

import (
	"bytes"
	"fmt"
	"strconv"
	"syscall"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libavutil
#include <libavutil/error.h>
#include <libavutil/dict.h>
#include <libavutil/mem.h>
#include <libavutil/frame.h>
#include <libavutil/error.h>
#include <stdlib.h>
#define MAX_LOG_BUFFER 1024

extern void av_log_cb_(int level,char* message,void* userInfo);

static void av_log_cb(void* userInfo,int level,const char* fmt,va_list args) {
	static char buf[MAX_LOG_BUFFER];
	vsnprintf(buf,MAX_LOG_BUFFER,fmt,args);
	av_log_cb_(level,buf,userInfo);
}
static void av_log_set_callback_(int def) {
	// true if the default callback should be set
	if (def) {
		av_log_set_callback(av_log_default_callback);
	} else {
		av_log_set_callback(av_log_cb);
	}
}
static int av_error_matches(int av,int en) {
	return av == AVERROR(en);
}
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	AVError           int
	AVDictionaryEntry C.struct_AVDictionaryEntry
	AVDictionaryFlag  int
	AVRational        C.struct_AVRational
)

type AVLogCallback func(level AVLogLevel, message string, userInfo uintptr)

type AVDictionary struct {
	ctx *C.struct_AVDictionary
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	BUF_SIZE = 1024
)

const (
	AV_DICT_NONE            AVDictionaryFlag = 0
	AV_DICT_MATCH_CASE      AVDictionaryFlag = 1
	AV_DICT_IGNORE_SUFFIX   AVDictionaryFlag = 2
	AV_DICT_DONT_STRDUP_KEY AVDictionaryFlag = 4
	AV_DICT_DONT_STRDUP_VAL AVDictionaryFlag = 8
	AV_DICT_DONT_OVERWRITE  AVDictionaryFlag = 16
	AV_DICT_APPEND          AVDictionaryFlag = 32
	AV_DICT_MULTIKEY        AVDictionaryFlag = 64
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

var (
	log_callback AVLogCallback
)

////////////////////////////////////////////////////////////////////////////////
// ERROR HANDLINE

func (this AVError) Error() string {
	cbuffer := make([]byte, BUF_SIZE)
	if err := C.av_strerror(C.int(this), (*C.char)(unsafe.Pointer(&cbuffer[0])), BUF_SIZE); err == 0 {
		if n := bytes.IndexByte(cbuffer, 0); n >= 0 {
			return string(cbuffer[:n])
		} else {
			return string(cbuffer)
		}
	} else {
		return fmt.Sprintf("Error code: %v", int(this))
	}
}

func (this AVError) IsErrno(err syscall.Errno) bool {
	c := int(C.av_error_matches(C.int(this), C.int(err)))
	return c == 1
}

////////////////////////////////////////////////////////////////////////////////
// DICTIONARY

func NewAVDictionary() *AVDictionary {
	return new(AVDictionary)
}

func (this *AVDictionary) Close() {
	if this.ctx != nil {
		C.av_dict_free(&this.ctx)
	}
}

func (this *AVDictionary) Count() int {
	if this.ctx == nil {
		return 0
	} else {
		return int(C.av_dict_count(this.ctx))
	}
}

func (this *AVDictionary) Get(key string, prev *AVDictionaryEntry, flags AVDictionaryFlag) *AVDictionaryEntry {
	if this.ctx == nil {
		return nil
	} else {
		key_ := C.CString(key)
		defer C.free(unsafe.Pointer(key_))
		return (*AVDictionaryEntry)(C.av_dict_get(this.ctx, key_, (*C.struct_AVDictionaryEntry)(prev), C.int(flags)))
	}
}

func (this *AVDictionary) Set(key, value string, flags AVDictionaryFlag) error {
	key_ := C.CString(key)
	value_ := C.CString(value)
	defer C.free(unsafe.Pointer(key_))
	defer C.free(unsafe.Pointer(value_))
	if err := AVError(C.av_dict_set(&this.ctx, key_, value_, C.int(flags))); err != 0 {
		return err
	} else {
		return nil
	}
}

func (this *AVDictionary) Keys() []string {
	keys := make([]string, 0, this.Count())
	entry := this.Get("", nil, AV_DICT_IGNORE_SUFFIX)
	for entry != nil {
		keys = append(keys, entry.Key())
		entry = this.Get("", entry, AV_DICT_IGNORE_SUFFIX)
	}
	return keys
}

func (this *AVDictionary) Entries() []*AVDictionaryEntry {
	keys := make([]*AVDictionaryEntry, 0, this.Count())
	entry := this.Get("", nil, AV_DICT_IGNORE_SUFFIX)
	for entry != nil {
		keys = append(keys, entry)
		entry = this.Get("", entry, AV_DICT_IGNORE_SUFFIX)
	}
	return keys
}

func (this *AVDictionary) String() string {
	return fmt.Sprintf("<AVDictionary>{ count=%v entries=%v }", this.Count(), this.Entries())
}

func (this *AVDictionary) context() *C.struct_AVDictionary {
	return this.ctx
}

////////////////////////////////////////////////////////////////////////////////
// DICTIONARY ENTRY

func (this *AVDictionaryEntry) Key() string {
	return C.GoString(this.key)
}

func (this *AVDictionaryEntry) Value() string {
	return C.GoString(this.value)
}

func (this *AVDictionaryEntry) String() string {
	return fmt.Sprintf("%v=%v", this.Key(), strconv.Quote(this.Value()))
}

////////////////////////////////////////////////////////////////////////////////
// RATIONAL NUMBER

func (this AVRational) Num() int {
	return int(this.num)
}

func (this AVRational) Den() int {
	return int(this.den)
}

func (this AVRational) String() string {
	if this.Num() == 0 {
		return "0"
	} else {
		return fmt.Sprintf("<AVRational>{ num=%v den=%v }", this.Num(), this.Den())
	}
}

// Float is used to convert an int64 value multipled by the rational to a float64
func (this AVRational) Float(multiplier int64) float64 {
	return float64(int64(this.num)*multiplier) / float64(this.den)
}

////////////////////////////////////////////////////////////////////////////////
// LOGGING

// AVLogSetCallback sets both the callback function and the level of output
// for logging. Where the callback is nil, the default ffmpeg logging is used.
func AVLogSetCallback(level AVLogLevel, cb AVLogCallback) {
	log_callback = cb
	if cb == nil {
		C.av_log_set_callback_(1)
	} else {
		C.av_log_set_callback_(0)
	}
	C.av_log_set_level(C.int(level))
}

func AVGetLogLevel() AVLogLevel {
	return AVLogLevel(C.av_log_get_level())
}

//export av_log_cb_
func av_log_cb_(level C.int, message *C.char, userInfo unsafe.Pointer) {
	if log_callback != nil && message != nil {
		level_ := AVLogLevel(level)
		if level_ <= AVGetLogLevel() {
			log_callback(level_, C.GoString(message), uintptr(userInfo))
		}
	}
}
