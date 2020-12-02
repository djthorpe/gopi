package keycode

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Cache struct {
	gopi.Unit
	sync.RWMutex

	lookup  map[string][]gopi.KeyCode
	device  map[string]gopi.InputDevice
	keycode map[string]gopi.KeyCode
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *Cache) New(gopi.Config) error {
	this.lookup = make(map[string][]gopi.KeyCode, gopi.KEYCODE_MAX*2)
	this.device = make(map[string]gopi.InputDevice)
	this.keycode = make(map[string]gopi.KeyCode, gopi.KEYCODE_MAX)

	// Index
	this.IndexKeycode()
	this.IndexDevice()

	// Return success
	return nil
}

func (this *Cache) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Release resources
	this.lookup = nil
	this.device = nil
	this.keycode = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Cache) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<keycode.cache"
	str += " device=" + fmt.Sprint(this.device)
	str += " keycode=" + fmt.Sprint(this.keycode)
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Cache) IndexKeycode() {
	// Index from names to keycodes
	for k := gopi.KEYCODE_NONE; k <= gopi.KEYCODE_MAX; k++ {
		name := strings.ToUpper(this.KeycodeName(k))
		terms := []string{name}
		if words := strings.Split(name, "_"); len(words) > 1 {
			terms = append(terms, words...)
		}

		// Term lookup
		for _, term := range terms {
			term = strings.ToUpper(term)
			if _, exists := this.lookup[term]; exists {
				this.lookup[term] = append(this.lookup[term], k)
			} else {
				this.lookup[term] = []gopi.KeyCode{k}
			}
		}

		// Keycode lookup
		this.keycode[name] = k
	}
}

func (this *Cache) IndexDevice() {
	// Index from names to input devices
	for d := gopi.INPUT_DEVICE_MIN; d <= gopi.INPUT_DEVICE_MAX; d <<= 1 {
		name := this.DeviceName(d)
		this.device[name] = d
	}
}

func (this *Cache) KeycodeName(k gopi.KeyCode) string {
	return fmt.Sprint(k)
}

func (this *Cache) DeviceName(d gopi.InputDevice) string {
	return strings.TrimPrefix(fmt.Sprint(d), "INPUT_DEVICE_")
}

func (this *Cache) SearchKeycode(name string) []gopi.KeyCode {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Result is an array of possible keycodes
	var result []gopi.KeyCode

	// Split up terms by spaces and commas
	terms := strings.FieldsFunc(name, func(r rune) bool {
		if unicode.IsSpace(r) {
			return true
		}
		if r == ',' {
			return true
		}
		if r == '_' {
			return true
		}
		return false
	})
	for _, term := range terms {
		term = strings.ToUpper(term)
		if keycodes, exists := this.lookup[term]; exists {
			result = append(result, keycodes...)
		}
	}

	// Return keycodes
	return result
}

func (this *Cache) LookupDevice(name string) gopi.InputDevice {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	name = strings.ToUpper(strings.TrimSpace(name))
	if d, exists := this.device[name]; exists {
		return d
	} else {
		return gopi.INPUT_DEVICE_NONE
	}
}

func (this *Cache) LookupKeycode(name string) gopi.KeyCode {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	name = strings.ToUpper(strings.TrimSpace(name))
	if k, exists := this.keycode[name]; exists {
		return k
	} else {
		return gopi.KEYCODE_NONE
	}
}

// DecodeLine returns a *mapentry or nil if the line could not be decoded
// The line is <keycode> <scancode> <device> # comment
// separated by whitepace, where the comment is optional
func (this *Cache) DecodeLine(line string) (*mapentry, error) {
	data, comment := split(line)
	if data == "" {
		return &mapentry{Comment: comment}, nil
	}
	if fields := strings.Fields(data); len(fields) < 3 {
		return nil, gopi.ErrBadParameter
	} else if key := this.LookupKeycode(fields[0]); key == gopi.KEYCODE_NONE {
		return nil, gopi.ErrBadParameter.WithPrefix(fields[0])
	} else if device := this.LookupDevice(fields[1]); device == gopi.INPUT_DEVICE_NONE {
		return nil, gopi.ErrBadParameter.WithPrefix(fields[1])
	} else if code, err := strconv.ParseUint(fields[2], 0, 32); err != nil {
		return nil, gopi.ErrBadParameter.WithPrefix(fields[2], " ", err)
	} else {
		return &mapentry{key, device, uint32(code), comment, 0}, nil
	}
}

func (this *Cache) EncodeLine(value *mapentry) string {
	if value == nil {
		return ""
	} else if value.Device == gopi.INPUT_DEVICE_NONE {
		return value.Comment
	}
	parts := []string{
		this.KeycodeName(value.Key),
		this.DeviceName(value.Device),
		scancodeString(value.Code),
	}
	if value.Comment != "" {
		parts = append(parts, value.Comment)
	}
	return strings.Join(parts, " ")
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Cache) DecodeFile(path string) ([]*mapentry, error) {
	var result []*mapentry

	// Open file
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	// Read lines
	r, lineno, exit := bufio.NewReader(fh), 0, false
	for exit == false {
		lineno++
		line, err := r.ReadString('\n')
		if err == io.EOF {
			exit = true
		} else if err != nil {
			return nil, err
		}
		if entry, err := this.DecodeLine(line); err != nil {
			return nil, fmt.Errorf("Line %v: %w", lineno, err)
		} else {
			result = append(result, entry)
		}
	}
	return result, nil
}

func (this *Cache) EncodeFile(path string, entries []*mapentry) error {
	// Create file
	fh, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fh.Close()

	for _, entry := range entries {
		if _, err := fh.WriteString(this.EncodeLine(entry)); err != nil {
			return err
		} else if _, err := fh.Write([]byte("\n")); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// split returns data and a comment
func split(line string) (string, string) {
	parts := strings.SplitN(line, "#", 2)
	if len(parts) == 1 {
		return strings.TrimSpace(strings.TrimSpace(parts[0])), ""
	} else {
		return strings.TrimSpace(strings.TrimSpace(parts[0])), "# " + strings.TrimSpace(parts[1])
	}
}

func scancodeString(code uint32) string {
	if code <= 0xFF {
		return fmt.Sprintf("0x%02X", code)
	} else if code <= 0xFFFF {
		return fmt.Sprintf("0x%04X", code)
	} else {
		return fmt.Sprintf("0x%08X", code)
	}
}
