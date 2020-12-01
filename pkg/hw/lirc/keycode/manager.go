package keycode

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"unicode"

	gopi "github.com/djthorpe/gopi/v3"
	codec "github.com/djthorpe/gopi/v3/pkg/hw/lirc/codec"

	_ "github.com/djthorpe/gopi/v3/pkg/event"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	gopi.Publisher
	gopi.Logger
	sync.Mutex

	folder, ext *string
	db          []*keycodedb
	keycodes    map[string][]gopi.KeyCode
	codecs      []codec.Codec
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	// Name for keycode files must be alphanumeric with one or more chars
	reKeycodeName = regexp.MustCompile("^[\\w\\-\\_\\.]+$")
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *Manager) Define(cfg gopi.Config) error {
	this.folder = cfg.FlagString("lirc.db", "", "Folder for keycode database")
	this.ext = cfg.FlagString("lirc.dbext", ".keycodes", "Extension for keycode files")
	return nil
}

func (this *Manager) New(cfg gopi.Config) error {
	if this.Publisher == nil {
		return gopi.ErrBadParameter.WithPrefix("Publisher")
	}
	if ext := "." + strings.Trim(strings.TrimSpace(*this.ext), "."); ext == "." {
		return gopi.ErrBadParameter.WithPrefix("-lirc.dbext")
	} else if *this.folder != "" {
		if stat, err := os.Stat(*this.folder); os.IsNotExist(err) {
			if folder, err := filepath.Abs(*this.folder); err != nil {
				return gopi.ErrBadParameter.WithPrefix(*this.folder)
			} else {
				return gopi.ErrBadParameter.WithPrefix(folder)
			}
		} else if err != nil {
			return err
		} else if stat.IsDir() == false {
			return gopi.ErrBadParameter.WithPrefix(*this.folder)
		} else if db, err := readAll(*this.folder, ext); err != nil {
			return err
		} else {
			this.db = db
		}
	}

	// Index keycodes in the background
	this.keycodes = make(map[string][]gopi.KeyCode, gopi.KEYCODE_MAX)
	go func() {
		this.Mutex.Lock()
		defer this.Mutex.Unlock()
		this.indexKeycodes()
	}()

	// Add codecs
	this.codecs = append(this.codecs, codec.NewRC5(gopi.INPUT_DEVICE_RC5_14))

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Release resources
	this.db = nil
	this.keycodes = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *Manager) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	var cancels []context.CancelFunc

	// Run codecs in background
	for _, codec := range this.codecs {
		ctx, cancel := context.WithCancel(context.Background())
		cancels = append(cancels, cancel)
		wg.Add(1)
		go func() {
			codec.Run(ctx, this.Publisher)
			wg.Done()
		}()
	}

	// Subscribe for CodecEvents
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

	// Receive CodecEvents until done
FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		case evt := <-ch:
			if codecevt, ok := evt.(*codec.CodecEvent); ok {
				inputevt := NewInputEvent(codecevt.Name(), gopi.KEYCODE_NONE, codecevt)
				if kc := this.Lookup(codecevt.Device, codecevt.Code); len(kc) > 0 {
					inputevt.KeyCode = kc[0]
				}
				if err := this.Publisher.Emit(inputevt, true); err != nil {
					this.Print(err)
				}
			}
		}
	}

	// cancel codecs and wait until they complete
	for _, cancel := range cancels {
		cancel()
	}
	wg.Wait()

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

// Keycode returns keycodes which match a name
func (this *Manager) Keycode(name string) []gopi.KeyCode {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Index the keycodes
	if len(this.keycodes) == 0 {
		this.indexKeycodes()
	}

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
		if keycodes, exists := this.keycodes[term]; exists {
			result = append(result, keycodes...)
		}
	}

	// Return keycodes
	return result
}

// Lookup one or more keycodes for a device and scancode
func (this *Manager) Lookup(gopi.InputDevice, uint32) []gopi.KeyCode {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return nil
}

// Set keycode for name,device and scancode
func (this *Manager) Set(gopi.InputDevice, uint32, string, gopi.KeyCode) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return gopi.ErrNotImplemented
}

func (this *Manager) DatabaseForName(string) (*keycodedb, error) {
	// Return an existing database or create a new one
	return nil, gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	str := "<keycode.manager"
	for _, db := range this.db {
		str += " " + db.name + "=" + fmt.Sprint(db)
	}
	str += " " + fmt.Sprint(this.keycodes)
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func readAll(path, ext string) ([]*keycodedb, error) {
	dbs := []*keycodedb{}

	// Files are sorted alphabetically, which is how we
	// will determine priority order
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.Mode().IsRegular() == false {
			continue
		}
		if ext_ := filepath.Ext(file.Name()); ext != ext_ {
			continue
		} else if name := strings.TrimSuffix(file.Name(), ext); reKeycodeName.MatchString(name) == false {
			continue
		} else if db, err := NewKeycodeDatabase(filepath.Join(path, file.Name()), name); err != nil {
			return nil, err
		} else {
			dbs = append(dbs, db)
		}
	}

	// Return databases
	return dbs, nil
}

func (this *Manager) indexKeycodes() {
	// Index from names to keycodes
	for k := gopi.KEYCODE_NONE; k <= gopi.KEYCODE_MAX; k++ {
		name := strings.TrimPrefix(fmt.Sprint(k), "KEYCODE_")
		terms := []string{name}
		if words := strings.Split(name, "_"); len(words) > 1 {
			terms = append(terms, words...)
		}
		for _, term := range terms {
			term = strings.ToUpper(term)
			if _, exists := this.keycodes[term]; exists {
				this.keycodes[term] = append(this.keycodes[term], k)
			} else {
				this.keycodes[term] = []gopi.KeyCode{k}
			}
		}
	}
}
