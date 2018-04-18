/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016
  All Rights Reserved

  For Licensing and Usage information, please see LICENSE.md
*/

package util /* import "github.com/djthorpe/gopi/util" */

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

/*
	This file implements a "Dict" structure which can be read and written
	to XML and JSON files in a similar way to Apple's Property Lists
	(see https://en.wikipedia.org/wiki/Property_list) implements simple
	storage types and compound types like dict. It adds the
	duration scalar type but currently does not implement the
	compound array type.
*/

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Dict defines a dictionary of values, somewhat similar to
// Apple property lists
type Dict struct {
	// Require ability to marshall and unmarshall XML
	xml.Marshaler
	xml.Unmarshaler

	values map[string]*v
}

type xmlState uint
type xmlType uint
type goType uint
type goTypeCastFunction func(*v) interface{}

type v struct {
	k string      // key
	t xmlType     // xml type
	v interface{} // go value
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	xmlStateKeyStart    xmlState = iota // we expect a key start element
	xmlStateKeyString                   // we expect key data
	xmlStateKeyEnd                      // we expect key end element
	xmlStateValueStart                  // we expect value start element
	xmlStateValueString                 // we expect value string
	xmlStateValueEnd                    // we expect value end element
)

const (
	xmlTypeString xmlType = iota
	xmlTypeInteger
	xmlTypeReal
	xmlTypeBool
	xmlTypeData
	xmlTypeDate
	xmlTypeDuration
	xmlTypeDict
)

const (
	goTypeString goType = iota
	goTypeInt
	goTypeUint
	goTypeFloat32
	goTypeFloat64
	goTypeBool
	goTypeByteArray
	goTypeTime
	goTypeDuration
	goTypeDict
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	xmlScalarTypes = map[string]xmlType{
		"string":   xmlTypeString,
		"integer":  xmlTypeInteger,
		"real":     xmlTypeReal,
		"bool":     xmlTypeBool,
		"data":     xmlTypeData,
		"date":     xmlTypeDate,
		"duration": xmlTypeDuration,
	}
	xmlCompoundTypes = map[string]xmlType{
		"dict": xmlTypeDict,
	}
	goTypeCast = map[goType]goTypeCastFunction{
		goTypeString:    castString,
		goTypeInt:       castInt,
		goTypeUint:      castUint,
		goTypeFloat32:   castFloat32,
		goTypeFloat64:   castFloat64,
		goTypeBool:      castBool,
		goTypeByteArray: castData,
		goTypeDuration:  castDuration,
		goTypeTime:      castDate,
		goTypeDict:      castDict,
	}
	xmlBoolValue = map[string]bool{
		"true":  true,
		"false": false,
	}
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - NEW

// NewDict creates a new empty dictionary. Use the capacity argument to provide
// a hint on how many items the dictionary will contain
func NewDict(capacity uint) *Dict {
	this := new(Dict)
	this.values = make(map[string]*v, capacity)
	return this
}

// CopyDict copies values from one dictionary into a new dictionary
func CopyDict(src *Dict) *Dict {
	capacity := uint(len(src.values))
	if this := NewDict(capacity); this == nil {
		return nil
	} else {
		for k, v := range src.values {
			this.values[k] = v.Copy()
		}
		return this
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Keys returns the keys associated with the dictionary
func (this *Dict) Keys() []string {
	keys := make([]string, 0, len(this.values))
	for k := range this.values {
		keys = append(keys, k)
	}
	return keys
}

func (this *Dict) IsEmpty() bool {
	return len(this.values) == 0
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - SETTERS

// SetString sets string value for key
func (this *Dict) SetString(key string, value string) {
	this.values[key] = &v{k: key, t: xmlTypeString, v: value}
}

// SetInt sets int value for key
func (this *Dict) SetInt(key string, value int) {
	this.values[key] = &v{k: key, t: xmlTypeInteger, v: value}
}

// SetUint sets uint value for key
func (this *Dict) SetUint(key string, value uint) {
	this.values[key] = &v{k: key, t: xmlTypeInteger, v: value}
}

// SetFloat64 sets float64 value for key
func (this *Dict) SetFloat64(key string, value float64) {
	this.values[key] = &v{k: key, t: xmlTypeReal, v: value}
}

// SetFloat32 sets float32 value for key
func (this *Dict) SetFloat32(key string, value float32) {
	this.values[key] = &v{k: key, t: xmlTypeReal, v: value}
}

// SetData sets []byte value for key
func (this *Dict) SetData(key string, value []byte) {
	this.values[key] = &v{k: key, t: xmlTypeData, v: value}
}

// SetDate sets time.Time value for key
func (this *Dict) SetDate(key string, value time.Time) {
	this.values[key] = &v{k: key, t: xmlTypeDate, v: value}
}

// SetDuration sets time.Duration value for key
func (this *Dict) SetDuration(key string, value time.Duration) {
	this.values[key] = &v{k: key, t: xmlTypeDuration, v: value}
}

// SetBool sets bool value for key
func (this *Dict) SetBool(key string, value bool) {
	this.values[key] = &v{k: key, t: xmlTypeBool, v: value}
}

// SetDict sets dict value for key
func (this *Dict) SetDict(key string, value *Dict) {
	if value != nil {
		this.values[key] = &v{k: key, t: xmlTypeDict, v: CopyDict(value)}
	} else {
		this.values[key] = &v{k: key, t: xmlTypeDict, v: NewDict(0)}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - STRINGIFY

func (this *Dict) String() string {
	if this.values == nil {
		return "<dict>{ nil }"
	}
	parts := make([]string, 0, len(this.values))
	for k, v := range this.values {
		s := v.cast(goTypeString).(string)
		parts = append(parts, fmt.Sprintf("%v=<%v>%v", k, v.t, s))
	}
	return fmt.Sprintf("<dict>{ %v }", strings.Join(parts, ","))
}

func (v xmlState) String() string {
	switch v {
	case xmlStateKeyStart:
		return "xmlStateKeyStart"
	case xmlStateKeyString:
		return "xmlStateKeyString"
	case xmlStateKeyEnd:
		return "xmlStateKeyEnd"
	case xmlStateValueStart:
		return "xmlStateValueStart"
	case xmlStateValueString:
		return "xmlStateValueString"
	case xmlStateValueEnd:
		return "xmlStateValueEnd"
	default:
		return "[?? Invalid xmlState value]"
	}
}

func (v xmlType) String() string {
	switch v {
	case xmlTypeString:
		return "xmlTypeString"
	case xmlTypeInteger:
		return "xmlTypeInteger"
	case xmlTypeReal:
		return "xmlTypeReal"
	case xmlTypeBool:
		return "xmlTypeBool"
	case xmlTypeData:
		return "xmlTypeData"
	case xmlTypeDate:
		return "xmlTypeDate"
	case xmlTypeDuration:
		return "xmlTypeDuration"
	case xmlTypeDict:
		return "xmlTypeDict"
	default:
		return "[?? Invalid xmlType value]"
	}
}

func (v *v) String() string {
	return fmt.Sprintf("v{ key=%v type=%v value=%v }", v.k, v.t, v.v)
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - GETTERS

func (this *Dict) GetString(key string) (string, bool) {
	if v := this.values[key]; v != nil {
		s, ok := v.cast(goTypeString).(string)
		return s, ok
	}
	return "", false
}

func (this *Dict) GetInt(key string) (int, bool) {
	if v := this.values[key]; v != nil {
		i, ok := v.cast(goTypeInt).(int)
		return i, ok
	}
	return 0, false
}

func (this *Dict) GetUint(key string) (uint, bool) {
	if v := this.values[key]; v != nil {
		u, ok := v.cast(goTypeUint).(uint)
		return u, ok
	}
	return 0, false
}

func (this *Dict) GetFloat64(key string) (float64, bool) {
	if v := this.values[key]; v != nil {
		r, ok := v.cast(goTypeFloat64).(float64)
		return r, ok
	}
	return 0, false
}

func (this *Dict) GetFloat32(key string) (float32, bool) {
	if v := this.values[key]; v != nil {
		r, ok := v.cast(goTypeFloat32).(float32)
		return r, ok
	}
	return 0, false
}

func (this *Dict) GetBool(key string) (bool, bool) {
	if v := this.values[key]; v != nil {
		b, ok := v.cast(goTypeBool).(bool)
		return b, ok
	}
	return false, false
}

func (this *Dict) GetData(key string) ([]byte, bool) {
	if v := this.values[key]; v != nil {
		a, ok := v.cast(goTypeByteArray).([]byte)
		return a, ok
	}
	return nil, false
}

func (this *Dict) GetDate(key string) (time.Time, bool) {
	if v := this.values[key]; v != nil {
		t, ok := v.cast(goTypeTime).(time.Time)
		return t, ok
	}
	return time.Time{}, false
}

func (this *Dict) GetDuration(key string) (time.Duration, bool) {
	if v := this.values[key]; v != nil {
		t, ok := v.cast(goTypeDuration).(time.Duration)
		return t, ok
	}
	return 0, false
}

func (this *Dict) GetDict(key string) (*Dict, bool) {
	if v := this.values[key]; v != nil {
		d, ok := v.cast(goTypeDict).(*Dict)
		return d, ok
	}
	return nil, false
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - XML

func (this *Dict) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// start dict element
	start.Name = xml.Name{Local: "dict"}
	e.EncodeToken(start)

	// iterate through values
	for k, v := range this.values {
		e.EncodeElement(k, xml.StartElement{Name: xml.Name{Local: "key"}})
		if err := e.Encode(v); err != nil {
			return err
		}
	}

	// end dict element
	e.EncodeToken(xml.EndElement{Name: start.Name})

	// return success
	return nil
}

func (this *v) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	name, err := this.xmlTypeString()
	if err != nil {
		return err
	}

	if name == "dict" {
		return e.Encode(this.v.(*Dict))
	}

	value, ok := this.cast(goTypeString).(string)
	if !ok {
		return err
	}
	switch name {
	case "string", "integer", "real", "data", "date", "duration":
		return e.EncodeElement(value, xml.StartElement{Name: xml.Name{Local: name}})
	case "bool":
		// TODO: Currently outputs <true></true> but we want <true/> for example
		// encoding/xml may not support this
		return e.EncodeElement("", xml.StartElement{Name: xml.Name{Local: value}})
	default:
		return ErrUnsupportedType
	}
}

func (this *Dict) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	// Create the this.values member variable
	this.values = make(map[string]*v, 0)

	// Check for name being 'dict'
	if start.Name.Local != "dict" {
		return ErrParseError
	}
	// Read key/value pairs in and end when we reach the 'dict' element
	// also ignore comments and attributes
	state := xmlStateKeyStart
	value := &v{}
	for true {
		t, err := d.Token()
		if err != nil {
			// Error occurred
			return err
		}
		if t == nil {
			// At EOF but without closing dict tag
			return ErrParseError
		}

		switch t.(type) {
		case xml.Comment:
			break // Ignore all comments
		case xml.Attr:
			break // Ignore all attributes
		case xml.StartElement:
			name := t.(xml.StartElement).Name.Local
			if state == xmlStateKeyStart && name == "key" {
				state = xmlStateKeyString
			} else if state == xmlStateValueStart && isXMLScalarNameOrTrueFalse(name, &value.t) {
				state = xmlStateValueString
				// Handle true and false values
				if value.t == xmlTypeBool {
					if v, ok := xmlBoolValue[name]; !ok {
						return ErrParseError
					} else {
						value.v = v
					}
				}
			} else {
				return ErrParseError
			}
		case xml.EndElement:
			name := t.(xml.EndElement).Name.Local
			if state == xmlStateKeyStart && name == "dict" {
				// Successfully completed
				return nil
			} else if (state == xmlStateKeyEnd || state == xmlStateKeyString) && name == "key" {
				state = xmlStateValueStart
			} else if (state == xmlStateValueEnd || state == xmlStateValueString) && isXMLScalarNameOrTrueFalse(name, &value.t) {
				// Eject the value, and move back to the start state
				state = xmlStateKeyStart
				this.values[value.k] = value
				value = &v{}
			} else {
				return ErrParseError
			}
		case xml.CharData:
			if state == xmlStateKeyString {
				value.k = string(t.(xml.CharData))
				state = xmlStateKeyEnd
			} else if state == xmlStateValueString {
				if ok := value.setValue(string(t.(xml.CharData))); !ok {
					return ErrParseError
				}
				state = xmlStateValueEnd
			} else {
				return ErrParseError
			}
		default:
			return ErrParseError
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - VALUES

// Return the XML type name for this value
func (this *v) xmlTypeString() (string, error) {
	switch this.v.(type) {
	case string:
		return "string", nil
	case uint, int:
		return "integer", nil
	case float32, float64:
		return "real", nil
	case bool:
		return "bool", nil
	case []byte:
		return "data", nil
	case time.Time:
		return "date", nil
	case time.Duration:
		return "duration", nil
	case *Dict:
		return "dict", nil
	}
	return "", fmt.Errorf("%v: %v", ErrUnsupportedType, this.v)
}

func (this *v) setValue(s string) bool {
	switch this.t {
	case xmlTypeString:
		this.v = s
		return true
	case xmlTypeBool:
		// we expect an empty string and the value has already been set
		if s != "" {
			return false
		}
		switch this.v.(type) {
		case bool:
			return true
		default:
			return false
		}
	default:
		return false
	}
	return true
}

// Return true if it's a scalar
func isXMLScalarName(name string) bool {
	_, exists := xmlScalarTypes[name]
	return exists
}

// Return true if it's not bool or a scalar, true or false
func isXMLScalarNameOrTrueFalse(name string, t *xmlType) bool {
	if name == "bool" {
		return false
	}
	if name == "true" || name == "false" {
		*t = xmlTypeBool
		return true
	}
	exists := false
	*t, exists = xmlScalarTypes[name]
	return exists
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - COPY VALUES

// Copy creates a copy of v if necessary, which is only necessary
// for *Dict since it's the only mutable type
func (this *v) Copy() *v {
	switch this.v.(type) {
	case int, uint, float32, float64, string, bool:
		// We can safely provide existing value
		return this
	case time.Time, time.Duration:
		// We can safely provide existing value
		return this
	case *Dict:
		// We make a copy of the dict and return that
		return &v{v: CopyDict(this.v.(*Dict))}
	default:
		panic(fmt.Sprint("Cannot copy: ", this.v))
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - CASTING FROM ONE GO TYPE TO ANOTHER

// Return value cast to a different value, or nil
func (this *v) cast(t goType) interface{} {
	if f, ok := goTypeCast[t]; ok {
		return f(this)
	} else {
		return nil
	}
}

func castString(this *v) interface{} {
	switch this.v.(type) {
	case string:
		return this.v.(string)
	case int:
		return strconv.FormatInt(int64(this.v.(int)), 10)
	case uint:
		return strconv.FormatUint(uint64(this.v.(uint)), 10)
	case float32:
		return strconv.FormatFloat(float64(this.v.(float32)), 'G', -1, 32)
	case float64:
		return strconv.FormatFloat(this.v.(float64), 'G', -1, 64)
	case bool:
		return fmt.Sprint(this.v.(bool))
	case []byte:
		return strings.ToUpper(hex.EncodeToString(this.v.([]byte)))
	case time.Time:
		return this.v.(time.Time).Format(time.RFC3339)
	case time.Duration:
		return this.v.(time.Duration).String()
	default:
		return nil
	}
}

func castInt(this *v) interface{} {
	switch this.v.(type) {
	case int:
		return this.v.(int)
	case uint:
		return int(this.v.(uint))
	case float32:
		return int(this.v.(float32))
	case float64:
		return int(this.v.(float64))
	default:
		return nil
	}
}

func castUint(this *v) interface{} {
	switch this.v.(type) {
	case int:
		return uint(this.v.(int))
	case uint:
		return uint(this.v.(uint))
	case float32:
		return uint(this.v.(float32))
	case float64:
		return uint(this.v.(float64))
	default:
		return nil
	}
}

func castFloat64(this *v) interface{} {
	switch this.v.(type) {
	case int:
		return float64(this.v.(int))
	case uint:
		return float64(this.v.(uint))
	case float32:
		return float64(this.v.(float32))
	case float64:
		return float64(this.v.(float64))
	default:
		return nil
	}
}

func castFloat32(this *v) interface{} {
	switch this.v.(type) {
	case int:
		return float32(this.v.(int))
	case uint:
		return float32(this.v.(uint))
	case float32:
		return float32(this.v.(float32))
	case float64:
		return float32(this.v.(float64))
	default:
		return nil
	}
}

func castBool(this *v) interface{} {
	switch this.v.(type) {
	case int:
		return this.v.(int) != 0
	case uint:
		return this.v.(uint) != 0
	case float32:
		return this.v.(float32) != 0
	case float64:
		return this.v.(float64) != 0
	case bool:
		return this.v.(bool)
	case string:
		return this.v.(string) != ""
	default:
		return nil
	}
}

func castDate(this *v) interface{} {
	switch this.v.(type) {
	case string:
		if time, err := time.Parse(time.RFC3339, this.v.(string)); err == nil {
			return time
		} else {
			return nil
		}
	case time.Time:
		return this.v.(time.Time)
	default:
		return nil
	}
}

func castDuration(this *v) interface{} {
	switch this.v.(type) {
	case string:
		if duration, err := time.ParseDuration(this.v.(string)); err == nil {
			return duration
		} else {
			return nil
		}
	case time.Duration:
		return this.v.(time.Duration)
	default:
		return nil
	}
}

func castData(this *v) interface{} {
	switch this.v.(type) {
	case []byte:
		return this.v.([]byte)
	case string:
		return []byte(this.v.(string))
	default:
		return nil
	}
}

func castDict(this *v) interface{} {
	switch this.v.(type) {
	case *Dict:
		return this.v.(*Dict)
	default:
		return nil
	}
}
