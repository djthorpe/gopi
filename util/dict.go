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
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Dict defines a dictionary of values, somewhat similar to
// Apple property lists
type Dict struct {
	// Require ability to marshall XML
	xml.Marshaler

	values map[string]*v
}

type v struct {
	v interface{}
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	ErrUnsupportedType = errors.New("Unsupported type")
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

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - SETTERS

// SetString sets string value for key
func (this *Dict) SetString(key string, value string) {
	this.values[key] = &v{v: value}
}

// SetInt sets int value for key
func (this *Dict) SetInt(key string, value int) {
	this.values[key] = &v{v: value}
}

// SetUint sets uint value for key
func (this *Dict) SetUint(key string, value uint) {
	this.values[key] = &v{v: value}
}

// SetFloat64 sets float64 value for key
func (this *Dict) SetFloat64(key string, value float64) {
	this.values[key] = &v{v: value}
}

// SetFloat32 sets float32 value for key
func (this *Dict) SetFloat32(key string, value float32) {
	this.values[key] = &v{v: value}
}

// SetData sets []byte value for key
func (this *Dict) SetData(key string, value []byte) {
	this.values[key] = &v{v: value}
}

// SetDate sets time.Time value for key
func (this *Dict) SetDate(key string, value time.Time) {
	this.values[key] = &v{v: value}
}

// SetDuration sets time.Duration value for key
func (this *Dict) SetDuration(key string, value time.Duration) {
	this.values[key] = &v{v: value}
}

// SetBool sets bool value for key
func (this *Dict) SetBool(key string, value bool) {
	this.values[key] = &v{v: value}
}

// SetDict sets dict value for key
func (this *Dict) SetDict(key string, value *Dict) {
	if value != nil {
		this.values[key] = &v{v: CopyDict(value)}
	} else {
		this.values[key] = &v{v: NewDict(0)}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - STRINGIFY

func (this *Dict) String() string {
	parts := make([]string, 0, len(this.values))
	for k, v := range this.values {
		s, _ := v.String()
		parts = append(parts, fmt.Sprintf("%v=%v", k, s))
	}
	return fmt.Sprintf("<dict>{ %v }", strings.Join(parts, ","))
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - GETTERS

func (this *Dict) GetString(key string) (string, bool) {
	if v, exists := this.values[key]; !exists {
		return "", false
	} else if s, err := v.String(); err != nil {
		return s, false
	} else {
		return s, true
	}
}

func (this *Dict) GetInt(key string) (int, bool) {
	if v, exists := this.values[key]; !exists {
		return 0, false
	} else if i, err := v.Int(); err != nil {
		return i, false
	} else {
		return i, true
	}
}

func (this *Dict) GetUint(key string) (uint, bool) {
	if v, exists := this.values[key]; !exists {
		return 0, false
	} else if i, err := v.Uint(); err != nil {
		return i, false
	} else {
		return i, true
	}
}

func (this *Dict) GetFloat64(key string) (float64, bool) {
	if v, exists := this.values[key]; !exists {
		return 0, false
	} else {
		i, err := v.Float64()
		return i, err == nil
	}
}

func (this *Dict) GetFloat32(key string) (float32, bool) {
	if v, exists := this.values[key]; !exists {
		return 0, false
	} else {
		i, err := v.Float32()
		return i, err == nil
	}
}
func (this *Dict) GetBool(key string) (bool, bool) {
	if v, exists := this.values[key]; !exists {
		return false, false
	} else {
		i, err := v.Bool()
		return i, err == nil
	}
}

func (this *Dict) GetData(key string) ([]byte, bool) {
	if v, exists := this.values[key]; !exists {
		return nil, false
	} else {
		i, err := v.Data()
		return i, err == nil
	}
}

func (this *Dict) GetDate(key string) (time.Time, bool) {
	if v, exists := this.values[key]; !exists {
		return time.Time{}, false
	} else {
		i, err := v.Date()
		return i, err == nil
	}
}

func (this *Dict) GetDuration(key string) (time.Duration, bool) {
	if v, exists := this.values[key]; !exists {
		return time.Duration(0), false
	} else {
		i, err := v.Duration()
		return i, err == nil
	}
}

func (this *Dict) GetDict(key string) (*Dict, bool) {
	if v, exists := this.values[key]; !exists {
		return nil, false
	} else {
		switch v.v.(type) {
		case *Dict:
			return v.v.(*Dict), true
		default:
			return nil, false
		}
	}
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

	value, err := this.String()
	if err != nil {
		return err
	}
	switch name {
	case "string", "integer", "real", "data", "date", "duration":
		return e.EncodeElement(value, xml.StartElement{Name: xml.Name{Local: name}})
	case "bool":
		// TODO: Currently outputs <true></true> but we want <true/> for example
		return e.EncodeElement("", xml.StartElement{Name: xml.Name{Local: value}})
	default:
		return ErrUnsupportedType
	}
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

// Return string value, which doesn't support
// *Dict or array values
func (this *v) String() (string, error) {
	switch this.v.(type) {
	case string:
		return this.v.(string), nil
	case int:
		return strconv.FormatInt(int64(this.v.(int)), 10), nil
	case uint:
		return strconv.FormatUint(uint64(this.v.(uint)), 10), nil
	case float32:
		return strconv.FormatFloat(float64(this.v.(float32)), 'G', -1, 32), nil
	case float64:
		return strconv.FormatFloat(this.v.(float64), 'G', -1, 64), nil
	case bool:
		if this.v.(bool) {
			return "true", nil
		} else {
			return "false", nil
		}
	case []byte:
		return strings.ToUpper(hex.EncodeToString(this.v.([]byte))), nil
	case time.Time:
		return this.v.(time.Time).Format(time.RFC3339), nil
	case time.Duration:
		return this.v.(time.Duration).String(), nil
	default:
		return "", ErrUnsupportedType
	}
}

// Return int value
func (this *v) Int() (int, error) {
	switch this.v.(type) {
	case int:
		return this.v.(int), nil
	case uint:
		return int(this.v.(uint)), nil
	case float32:
		return int(this.v.(float32)), nil
	case float64:
		return int(this.v.(float64)), nil
	default:
		return 0, ErrUnsupportedType
	}
}

// Return uint value
func (this *v) Uint() (uint, error) {
	switch this.v.(type) {
	case int:
		return uint(this.v.(int)), nil
	case uint:
		return this.v.(uint), nil
	case float32:
		return uint(this.v.(float32)), nil
	case float64:
		return uint(this.v.(float64)), nil
	default:
		return 0, ErrUnsupportedType
	}
}

// Return float64 value
func (this *v) Float64() (float64, error) {
	switch this.v.(type) {
	case int:
		return float64(this.v.(int)), nil
	case uint:
		return float64(this.v.(uint)), nil
	case float32:
		return float64(this.v.(float32)), nil
	case float64:
		return float64(this.v.(float64)), nil
	default:
		return 0, ErrUnsupportedType
	}
}

// Return float32 value
func (this *v) Float32() (float32, error) {
	switch this.v.(type) {
	case int:
		return float32(this.v.(int)), nil
	case uint:
		return float32(this.v.(uint)), nil
	case float32:
		return float32(this.v.(float32)), nil
	case float64:
		return float32(this.v.(float64)), nil
	default:
		return 0, ErrUnsupportedType
	}
}

// Return bool value
func (this *v) Bool() (bool, error) {
	switch this.v.(type) {
	case int:
		return this.v.(int) != 0, nil
	case uint:
		return this.v.(uint) != 0, nil
	case float32:
		return this.v.(float32) != 0, nil
	case float64:
		return this.v.(float64) != 0, nil
	case string:
		return this.v.(string) != "", nil
	default:
		return false, ErrUnsupportedType
	}
}

// Return time.Time value
func (this *v) Date() (time.Time, error) {
	switch this.v.(type) {
	case time.Time:
		return this.v.(time.Time), nil
	default:
		return time.Time{}, ErrUnsupportedType
	}
}

// Return time.Duration value
func (this *v) Duration() (time.Duration, error) {
	switch this.v.(type) {
	case time.Duration:
		return this.v.(time.Duration), nil
	default:
		return time.Duration(0), ErrUnsupportedType
	}
}

// Return data value
func (this *v) Data() ([]byte, error) {
	switch this.v.(type) {
	case []byte:
		return this.v.([]byte), nil
	case string:
		return []byte(this.v.(string)), nil
	default:
		return nil, ErrUnsupportedType
	}
}

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
