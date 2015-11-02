package csvdata

/* This package implements a CSV reader which outputs records, rather than
   arrays of strings.
*/

import (
    "io"
    "errors"
    "reflect"
    "strconv"
    "strings"
    "log"
)

////////////////////////////////////////////////////////////////////////////////

// The data source is any object that has a Read method which can
// return a row as a slice of strings. This matches csv.Reader in particular.
type Reader interface {
    Read() ([]string,error)
}

// Custom data types can be implemented by implementing Value; these
// methods must be defined on a pointer receiver.
// The interface is also used by flag package for a similar purpose.
type Value interface {
    String() string
    Set(string) bool
}

// ReadIter encapsulates an iterator over a Reader source that fills a
// pointer to a user struct with data.
type ReadIter struct {
    Reader Reader
    Headers []string
    Error error
    Line, Column int
    fields []reflect.Value
    kinds []int
    tags []int
}

const (
    none_k = iota
    string_k
    int_k
    float_k
    uint_k
    bool_k
    value_k
)

// Creates a new iterator from a Reader source and a user-defined struct.
func NewReadIter(rdr Reader,ps interface{}) (*ReadIter,error) {
    // create a "this" object" and read in the header
    var this = new(ReadIter)
    var err error
    this.Line = 1
    this.Reader = rdr
    this.Headers, err = rdr.Read()
    if err != nil {
        return nil,err
    }

    // Reflect data structure names and values
    st := reflect.TypeOf(ps).Elem()
    sv := reflect.ValueOf(ps).Elem()
    nf := st.NumField()
    this.kinds = make([]int, nf)
    this.tags = make([]int, nf)
    this.fields = make([]reflect.Value, nf)
    for i := 0; i < nf; i++ {
        // Set up the field
        f := st.Field(i)
        val := sv.Field(i)
        tag := f.Tag.Get("field")
        required,_ := strconv.ParseBool(f.Tag.Get("required"))

        // get the corresponding field name as lowercase
        if len(tag) == 0 {
            tag = f.Name
        }
        tag = strings.ToLower(tag)
        tag = strings.Replace(tag," ","_",-1)

        // determine index of tag in header
        itag := -1
        for k, h := range this.Headers {
            h = strings.ToLower(h)
            h = strings.Replace(h," ","-",-1)
            if h == tag {
                itag = k
                break
            }
        }
        if itag == -1 && required {
            return nil,errors.New("Missing required field: " + f.Name)
        }

        // Work out type of the field
        kind := none_k
        Kind := f.Type.Kind()
        // this is necessary because Kind can't tell distinguish between a primitive type
        // and a type derived from it. We're looking for a Value interface defined on
        // the pointer to this value
        _, ok := val.Addr().Interface().(Value)
        if ok {
            val = val.Addr()
            kind = value_k
        } else {
            switch Kind {
                case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64:
                    kind = int_k
                case reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32, reflect.Uint64:
                    kind = uint_k
                case reflect.Float32, reflect.Float64:
                    kind = float_k
                case reflect.String:
                    kind = string_k
                case reflect.Bool:
                    kind = bool_k
                default:
                    kind = value_k
                    _, ok := val.Interface().(Value)
                    if !ok {
                        return nil,errors.New("Cannot determine type for field: " + f.Name)
                    }
            }
        }
        this.kinds[i] = kind
        this.tags[i] = itag
        this.fields[i] = val
    }
    return this,nil
}

// The Get method reads the next row. If there was an error or EOF, it
// will return false.  Client code must then check that ReadIter.Error is
// not nil to distinguish between normal EOF and specific errors.
func (this *ReadIter) Get() bool {
    row, err := this.Reader.Read()
    this.Line = this.Line + 1
    if err != nil {
        if err != io.EOF {
            this.Error = err
        }
        return false
    }
    var ival int64
    var fval float64
    var uval uint64
    var bval bool
    var v Value
    var ok bool

    log.Printf("tags = %v",this.tags)

    for fi, ci := range this.tags {
        if ci == -1 {
            continue
        }
        vals := row[ci] // string at column ci of current row
        f := this.fields[fi]
        switch this.kinds[fi] {
            case string_k:
                f.SetString(vals)
            case int_k:
                ival, err = strconv.ParseInt(vals,10,64)
                f.SetInt(ival)
            case uint_k:
                uval, err = strconv.ParseUint(vals,10,64)
                f.SetUint(uval)
            case float_k:
                fval, err = strconv.ParseFloat(vals,64)
                f.SetFloat(fval)
            case bool_k:
                bval, err = parseBool(vals)
                f.SetBool(bval)
            case value_k:
                v, ok = f.Interface().(Value)
                if !ok {
                    err = errors.New("Not a Value object")
                    break
                }
                v.Set(vals)
        }
        if err != nil {
            this.Column = ci + 1
            this.Error = err
            return false
        }
    }
    return true
}

// The parseBool method will do wider boolean parsing than strconv method
// It accepts y,n,YES,NO,true,false,0,1
func parseBool(s string) (bool,error) {
    switch strings.ToLower(s) {
        case "t","true","y","yes","1":
            return true,nil
        case "f","false","n","no","0":
            return false,nil
    }
    return false,errors.New("Invalid value: " + s)
}

