package util_test

import (
	"bytes"
	"encoding/xml"
	"math"
	"testing"
	"time"

	"github.com/djthorpe/gopi/util"
)

func TestDict_000(t *testing.T) {
	// Create a dict object
	dict := util.NewDict(0)
	if dict == nil {
		t.Error("NewDict returns <nil>")
	}
}

func TestDict_001(t *testing.T) {
	// Create a dict object with 0 capacity
	dict := util.NewDict(0)
	if dict == nil {
		t.Fatal("NewDict returns <nil>")
	}
	// Create buffer, write out empty XML
	str := xmlstring(t, dict)
	if str != "<dict></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}
}

func TestDict_002(t *testing.T) {
	// Create a dict object with 1 capacity
	dict := util.NewDict(1)
	if dict == nil {
		t.Fatal("NewDict returns <nil>")
	}

	// STRING TEST

	dict.SetString("test_string", "")
	if v, ok := dict.GetString("test_string"); !ok {
		t.Errorf("Error: GetString did not return true")
	} else if v != "" {
		t.Errorf("Error: GetString unexpected value, returned %v", v)
	}

	dict.SetString("test_string", "test")
	if v, ok := dict.GetString("test_string"); !ok {
		t.Errorf("Error: GetString did not return true")
	} else if v != "test" {
		t.Errorf("Error: GetString unexpected value, returned %v", v)
	}

	str := xmlstring(t, dict)
	if str != "<dict><key>test_string</key><string>test</string></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}
}

func TestDict_003(t *testing.T) {
	// Create a dict object with 1 capacity
	dict := util.NewDict(1)
	if dict == nil {
		t.Fatal("NewDict returns <nil>")
	}

	// UINT TEST 1
	dict.SetUint("test_uint", 56)

	if v, ok := dict.GetString("test_uint"); !ok {
		t.Errorf("Error: GetString did not return true")
	} else if v != "56" {
		t.Errorf("Error: GetString unexpected value, returned %v", v)
	}
	if v, ok := dict.GetUint("test_uint"); !ok {
		t.Errorf("Error: GetUint did not return true")
	} else if v != 56 {
		t.Errorf("Error: GetUint unexpected value, returned %v", v)
	}
	if v, ok := dict.GetInt("test_uint"); !ok {
		t.Errorf("Error: GetInt did not return true")
	} else if v != 56 {
		t.Errorf("Error: GetInt unexpected value, returned %v", v)
	}
	if str := xmlstring(t, dict); str != "<dict><key>test_uint</key><integer>56</integer></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// UINT TEST 2
	dict.SetUint("test_uint", 99)
	if v, ok := dict.GetString("test_uint"); !ok {
		t.Errorf("Error: GetString did not return true")
	} else if v != "99" {
		t.Errorf("Error: GetString unexpected value, returned %v", v)
	}
	if v, ok := dict.GetUint("test_uint"); !ok {
		t.Errorf("Error: GetUint did not return true")
	} else if v != 99 {
		t.Errorf("Error: GetUint unexpected value, returned %v", v)
	}
	if v, ok := dict.GetInt("test_uint"); !ok {
		t.Errorf("Error: GetInt did not return true")
	} else if v != 99 {
		t.Errorf("Error: GetInt unexpected value, returned %v", v)
	}
	if v, ok := dict.GetBool("test_uint"); !ok {
		t.Errorf("Error: GetBool did not return true")
	} else if v == false {
		t.Errorf("Error: GetBool unexpected value, returned %v", v)
	}
	if str := xmlstring(t, dict); str != "<dict><key>test_uint</key><integer>99</integer></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// UINT TEST 3
	dict.SetUint("test_uint", 0)
	if v, ok := dict.GetString("test_uint"); !ok {
		t.Errorf("Error: GetString did not return true")
	} else if v != "0" {
		t.Errorf("Error: GetString unexpected value, returned %v", v)
	}
	if v, ok := dict.GetUint("test_uint"); !ok {
		t.Errorf("Error: GetUint did not return true")
	} else if v != 0 {
		t.Errorf("Error: GetUint unexpected value, returned %v", v)
	}
	if v, ok := dict.GetInt("test_uint"); !ok {
		t.Errorf("Error: GetInt did not return true")
	} else if v != 0 {
		t.Errorf("Error: GetInt unexpected value, returned %v", v)
	}
	if v, ok := dict.GetBool("test_uint"); !ok {
		t.Errorf("Error: GetBool did not return true")
	} else if v == true {
		t.Errorf("Error: GetBool unexpected value, returned %v", v)
	}
	if str := xmlstring(t, dict); str != "<dict><key>test_uint</key><integer>0</integer></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}
}

func TestDict_004(t *testing.T) {
	// Create a dict object with 1 capacity
	dict := util.NewDict(1)
	if dict == nil {
		t.Fatal("NewDict returns <nil>")
	}

	// INT TEST 1
	dict.SetInt("test_int", 56)
	if v, ok := dict.GetString("test_int"); !ok {
		t.Errorf("Error: GetString did not return true")
	} else if v != "56" {
		t.Errorf("Error: GetString unexpected value, returned %v", v)
	}
	if v, ok := dict.GetUint("test_int"); !ok {
		t.Errorf("Error: GetUint did not return true")
	} else if v != 56 {
		t.Errorf("Error: GetUint unexpected value, returned %v", v)
	}
	if v, ok := dict.GetInt("test_int"); !ok {
		t.Errorf("Error: GetInt did not return true")
	} else if v != 56 {
		t.Errorf("Error: GetInt unexpected value, returned %v", v)
	}
	if v, ok := dict.GetBool("test_int"); !ok {
		t.Errorf("Error: GetBool did not return true")
	} else if v != true {
		t.Errorf("Error: GetBool unexpected value, returned %v", v)
	}
	if str := xmlstring(t, dict); str != "<dict><key>test_int</key><integer>56</integer></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// INT TEST 2
	dict.SetInt("test_int", -99)
	if v, ok := dict.GetString("test_int"); !ok {
		t.Errorf("Error: GetString did not return true")
	} else if v != "-99" {
		t.Errorf("Error: GetString unexpected value, returned %v", v)
	}
	if v, ok := dict.GetInt("test_int"); !ok {
		t.Errorf("Error: GetInt did not return true")
	} else if v != -99 {
		t.Errorf("Error: GetInt unexpected value, returned %v", v)
	}
	if v, ok := dict.GetBool("test_int"); !ok {
		t.Errorf("Error: GetBool did not return true")
	} else if v != true {
		t.Errorf("Error: GetBool unexpected value, returned %v", v)
	}
	if str := xmlstring(t, dict); str != "<dict><key>test_int</key><integer>-99</integer></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// INT TEST 3
	dict.SetInt("test_int", 0)
	if v, ok := dict.GetString("test_int"); !ok {
		t.Errorf("Error: GetString did not return true")
	} else if v != "0" {
		t.Errorf("Error: GetString unexpected value, returned %v", v)
	}
	if v, ok := dict.GetInt("test_int"); !ok {
		t.Errorf("Error: GetInt did not return true")
	} else if v != 0 {
		t.Errorf("Error: GetInt unexpected value, returned %v", v)
	}
	if v, ok := dict.GetBool("test_int"); !ok {
		t.Errorf("Error: GetBool did not return true")
	} else if v == true {
		t.Errorf("Error: GetBool unexpected value, returned %v", v)
	}
	if str := xmlstring(t, dict); str != "<dict><key>test_int</key><integer>0</integer></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}
}

func TestDict_005(t *testing.T) {
	// Create a dict object with 1 capacity
	dict := util.NewDict(1)
	if dict == nil {
		t.Fatal("NewDict returns <nil>")
	}

	// BOOL TEST 1
	dict.SetBool("test_bool", true)
	str := xmlstring(t, dict)
	if str != "<dict><key>test_bool</key><true></true></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// BOOL TEST 2
	dict.SetBool("test_bool", false)
	str = xmlstring(t, dict)
	if str != "<dict><key>test_bool</key><false></false></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}
}

func TestDict_006(t *testing.T) {
	// Create a dict object with 1 capacity
	dict := util.NewDict(1)
	if dict == nil {
		t.Fatal("NewDict returns <nil>")
	}

	// FLOAT32 TEST 1
	dict.SetFloat32("test_float32", 3.1415927)
	str := xmlstring(t, dict)
	if str != "<dict><key>test_float32</key><real>3.1415927</real></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// FLOAT32 TEST 2
	dict.SetFloat32("test_float32", -3.1415927)
	str = xmlstring(t, dict)
	if str != "<dict><key>test_float32</key><real>-3.1415927</real></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// FLOAT32 TEST 3
	dict.SetFloat32("test_float32", float32(math.NaN()))
	str = xmlstring(t, dict)
	if str != "<dict><key>test_float32</key><real>NaN</real></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// FLOAT32 TEST 4
	dict.SetFloat32("test_float32", float32(math.Inf(1)))
	str = xmlstring(t, dict)
	if str != "<dict><key>test_float32</key><real>+Inf</real></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// FLOAT32 TEST 5
	dict.SetFloat32("test_float32", float32(math.Inf(-1)))
	str = xmlstring(t, dict)
	if str != "<dict><key>test_float32</key><real>-Inf</real></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

}

func TestDict_007(t *testing.T) {
	// Create a dict object with 1 capacity
	dict := util.NewDict(1)
	if dict == nil {
		t.Fatal("NewDict returns <nil>")
	}

	// FLOAT64 TEST 1
	dict.SetFloat64("test_float64", 3.1415927)
	str := xmlstring(t, dict)
	if str != "<dict><key>test_float64</key><real>3.1415927</real></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// FLOAT64 TEST 2
	dict.SetFloat64("test_float64", -3.1415927)
	str = xmlstring(t, dict)
	if str != "<dict><key>test_float64</key><real>-3.1415927</real></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// FLOAT64 TEST 3
	dict.SetFloat64("test_float64", math.NaN())
	str = xmlstring(t, dict)
	if str != "<dict><key>test_float64</key><real>NaN</real></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// FLOAT64 TEST 4
	dict.SetFloat64("test_float64", math.Inf(1))
	str = xmlstring(t, dict)
	if str != "<dict><key>test_float64</key><real>+Inf</real></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// FLOAT64 TEST 5
	dict.SetFloat64("test_float64", math.Inf(-1))
	str = xmlstring(t, dict)
	if str != "<dict><key>test_float64</key><real>-Inf</real></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

}

func TestDict_008(t *testing.T) {
	// Create a dict object with 1 capacity
	dict := util.NewDict(1)
	if dict == nil {
		t.Fatal("NewDict returns <nil>")
	}

	// DATA TEST 1
	dict.SetData("test_data", []byte{})
	if v, ok := dict.GetData("test_data"); !ok {
		t.Errorf("GetData unexpected return value")
	} else if len(v) != 0 {
		t.Errorf("GetData unexpected return value")
	} else if str := xmlstring(t, dict); str != "<dict><key>test_data</key><data></data></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// DATA TEST 2
	data := []byte("test_data")
	dict.SetData("test_data", data)
	if v, ok := dict.GetData("test_data"); !ok {
		t.Errorf("error: unexpected return value\n")
	} else if bytes.Compare(data, v) != 0 {
		t.Errorf("error: unexpected return value\n")
	} else if str := xmlstring(t, dict); str != "<dict><key>test_data</key><data>746573745F64617461</data></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

}

func TestDict_009(t *testing.T) {
	// Create a dict object with 1 capacity
	dict := util.NewDict(1)
	if dict == nil {
		t.Fatal("NewDict returns <nil>")
	}

	// DATE TEST 1
	dict.SetDate("test_date", time.Time{})
	str := xmlstring(t, dict)
	if str != "<dict><key>test_date</key><date>0001-01-01T00:00:00Z</date></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}
}

func TestDict_010(t *testing.T) {
	// Create a dict object with 1 capacity
	dict := util.NewDict(1)
	if dict == nil {
		t.Fatal("NewDict returns <nil>")
	}

	// DURATION TEST
	dict.SetDuration("test_duration", 24*time.Hour)
	str := xmlstring(t, dict)
	if str != "<dict><key>test_duration</key><duration>24h0m0s</duration></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}
}

func TestDict_011(t *testing.T) {
	// Create a dict object with 1 capacity
	dict := util.NewDict(1)
	if dict == nil {
		t.Fatal("NewDict returns <nil>")
	}

	// DICT TEST 1
	dict.SetDict("test_dict", util.NewDict(0))
	str := xmlstring(t, dict)
	if str != "<dict><key>test_dict</key><dict></dict></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// DICT TEST 2
	child := util.NewDict(1)
	child.SetString("test_string", "test")
	dict.SetDict("test_dict", child)
	if str := xmlstring(t, dict); str != "<dict><key>test_dict</key><dict><key>test_string</key><string>test</string></dict></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}

	// DICT TEST 3
	child2 := util.NewDict(1)
	child2.SetDict("test_child", child2)
	dict.SetDict("test_dict", child2)
	if str := xmlstring(t, dict); str != "<dict><key>test_dict</key><dict><key>test_child</key><dict></dict></dict></dict>" {
		t.Errorf("error: unexpected XML: %v\n", str)
	}
}

func TestUnmarshall_001(t *testing.T) {
	// Create an empty dict object
	var dict *util.Dict
	if err := xml.Unmarshal([]byte("<dict></dict>"), &dict); err != nil {
		t.Errorf("Unmarshal error: %v", err)
	} else if dict.IsEmpty() == false {
		t.Errorf("Dictionary not empty")
	}
	// Create a dict object, and ensure it's empty when reading in XML
	dict2 := util.NewDict(1)
	dict2.SetBool("test_value", true)
	dict2.SetString("test_value2", "test string")
	if dict2.IsEmpty() == true {
		t.Errorf("Dictionary Empty")
	}
	if err := xml.Unmarshal([]byte("<dict></dict>"), &dict2); err != nil {
		t.Errorf("Unmarshal error: %v", err)
	} else if dict2.IsEmpty() == false {
		t.Errorf("Dictionary not empty")
	}
}

func TestUnmarshall_002(t *testing.T) {
	// Create an empty dict object
	var dict *util.Dict
	if err := xml.Unmarshal([]byte("<dict><key>test</key><string test=\"1\">test<!-- comment --></string></dict>"), &dict); err != nil {
		t.Errorf("Unmarshal error: %v", err)
	} else if len(dict.Keys()) != 1 {
		t.Errorf("Dictionary should contain 1 key: %v", dict)
	}
	if value, ok := dict.GetString("test"); !ok {
		t.Errorf("Expected 'test' value to be retrieved from dict: %v", dict)
	} else if value != "test" {
		t.Errorf("Expected 'test' value to be retrieved from dict: %v", dict)
	}
}

func TestUnmarshall_003(t *testing.T) {
	// Create an empty dict object
	var dict *util.Dict
	if err := xml.Unmarshal([]byte("<dict><key>test_true</key><true/><key>test_false</key><false/></dict>"), &dict); err != nil {
		t.Errorf("Unmarshal error: %v", err)
	} else if len(dict.Keys()) != 2 {
		t.Errorf("Dictionary should contain 2 keys: %v", dict)
	}
	if value, ok := dict.GetBool("test_true"); !ok {
		t.Errorf("Expected 'test_true' value to be retrieved from dict: %v", dict)
	} else if value != true {
		t.Errorf("Expected 'test_true' value to be retrieved from dict: %v", dict)
	}
	if value, ok := dict.GetBool("test_false"); !ok {
		t.Errorf("Expected 'test_false' value to be retrieved from dict: %v", dict)
	} else if value != false {
		t.Errorf("Expected 'test_false' value to be retrieved from dict: %v", dict)
	}
}

func TestUnmarshall_004(t *testing.T) {
	// Create an empty dict object
	var dict *util.Dict
	if err := xml.Unmarshal([]byte("<dict><key>test_int1</key><integer>100</integer><key>test_int2</key><integer>-100</integer></dict>"), &dict); err != nil {
		t.Errorf("Unmarshal error: %v", err)
	} else if len(dict.Keys()) != 2 {
		t.Errorf("Dictionary should contain 2 keys: %v", dict)
	}
	if value, ok := dict.GetUint("test_int1"); !ok {
		t.Errorf("Expected 'test_int1' value to be retrieved from dict: %v", dict)
	} else if value != 100 {
		t.Errorf("Expected 'test_int1' value to be retrieved from dict: %v", dict)
	}
	if value, ok := dict.GetInt("test_int1"); !ok {
		t.Errorf("Expected 'test_int1' value to be retrieved from dict: %v", dict)
	} else if value != 100 {
		t.Errorf("Expected 'test_int1' value to be retrieved from dict: %v", dict)
	}
	if value, ok := dict.GetInt("test_int2"); !ok {
		t.Errorf("Expected 'test_int2' value to be retrieved from dict: %v", dict)
	} else if value != -100 {
		t.Errorf("Expected 'test_int2' value to be retrieved from dict: %v", dict)
	}
}

func TestUnmarshall_005(t *testing.T) {
	// Create an empty dict object
	var dict *util.Dict
	if err := xml.Unmarshal([]byte("<dict><key>test_real1</key><real>3.14</real><key>test_real2</key><real>1E10</real></dict>"), &dict); err != nil {
		t.Errorf("Unmarshal error: %v", err)
	} else if len(dict.Keys()) != 2 {
		t.Errorf("Dictionary should contain 2 keys: %v", dict)
	}
	if value, ok := dict.GetFloat32("test_real1"); !ok {
		t.Errorf("Expected 'test_real1' value to be retrieved from dict: %v", dict)
	} else if value != 3.14 {
		t.Errorf("Expected 'test_real1' value to be retrieved from dict: %v", dict)
	}
	if value, ok := dict.GetFloat64("test_real1"); !ok {
		t.Errorf("Expected 'test_real1' value to be retrieved from dict: %v", dict)
	} else if value != 3.14 {
		t.Errorf("Expected 'test_real1' value to be retrieved from dict: %v", dict)
	}
	if value, ok := dict.GetInt("test_real1"); !ok {
		t.Errorf("Expected 'test_real1' value to be retrieved from dict: %v", dict)
	} else if value != 3 {
		t.Errorf("Expected 'test_real1' value to be retrieved from dict: %v", dict)
	}

	if value, ok := dict.GetFloat32("test_real2"); !ok {
		t.Errorf("Expected 'test_real2' value to be retrieved from dict: %v", dict)
	} else if value != 1E10 {
		t.Errorf("Expected 'test_real2' value to be retrieved from dict: %v", dict)
	}
	if value, ok := dict.GetFloat64("test_real2"); !ok {
		t.Errorf("Expected 'test_real2' value to be retrieved from dict: %v", dict)
	} else if value != 1E10 {
		t.Errorf("Expected 'test_real2' value to be retrieved from dict: %v", dict)
	}
	if value, ok := dict.GetInt("test_real2"); !ok {
		t.Errorf("Expected 'test_real2' value to be retrieved from dict: %v", dict)
	} else if value != 1E10 {
		t.Errorf("Expected 'test_real2' value to be retrieved from dict: %v", dict)
	}
}

////////////////////////////////////////////////////////////////////////////////

func xmlstring(t *testing.T, dict *util.Dict) string {
	// Create buffer, write out empty XML
	buf := bytes.NewBuffer([]byte{})
	encoder := xml.NewEncoder(buf)
	if err := encoder.Encode(dict); err != nil {
		t.Errorf("error: %v\n", err)
	}
	if err := encoder.Flush(); err != nil {
		t.Errorf("error: %v\n", err)
	}
	return buf.String()
}
