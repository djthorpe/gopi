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
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	codec "github.com/djthorpe/gopi/v3/pkg/hw/lirc/codec"
	"github.com/hashicorp/go-multierror"

	_ "github.com/djthorpe/gopi/v3/pkg/event"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	gopi.Publisher
	gopi.Logger
	sync.Mutex
	*Cache

	folder, ext *string
	db          []*keycodedb
	codecs      []codec.Codec
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	// Databases get written once every thirty seconds
	writeDelta = 30 * time.Second
)

var (
	// Name for keycode files must be alphanumeric with one or more chars
	reKeycodeName = regexp.MustCompile("^[\\w\\-\\_\\.]+$")
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *Manager) Define(cfg gopi.Config) error {
	this.folder = cfg.FlagString("lirc.db", "", "Folder for keycode database")
	this.ext = cfg.FlagString("lirc.ext", ".keycode", "Extension for keycode files")
	return nil
}

func (this *Manager) New(cfg gopi.Config) error {
	if this.Publisher == nil {
		return gopi.ErrBadParameter.WithPrefix("Publisher")
	}
	if ext := "." + strings.Trim(strings.TrimSpace(*this.ext), "."); ext == "." {
		return gopi.ErrBadParameter.WithPrefix("-lirc.ext")
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
		} else if db, err := this.readAll(*this.folder, ext); err != nil {
			return err
		} else {
			this.db = db
			*this.ext = ext
		}
	}

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

	// Timer to write databases occasionally on modified
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

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
		case <-timer.C:
			if err := this.writeDirty(); err != nil {
				this.Print("LIRCKeycodeManager: ", err)
			}
			timer.Reset(writeDelta)
		}
	}

	// cancel codecs and wait until they complete
	for _, cancel := range cancels {
		cancel()
	}
	wg.Wait()

	// Write databases if necessary
	if err := this.writeDirty(); err != nil {
		this.Print("LIRCKeycodeManager: ", err)
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

// Keycode returns keycodes which match a name
func (this *Manager) Keycode(name string) []gopi.KeyCode {
	if k := this.Cache.LookupKeycode(name); k != gopi.KEYCODE_NONE {
		return []gopi.KeyCode{k}
	} else {
		return this.Cache.SearchKeycode(name)
	}
}

// Lookup one or more keycodes for a device and scancode
func (this *Manager) Lookup(device gopi.InputDevice, code uint32) []gopi.KeyCode {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	keys := []gopi.KeyCode{}
	for _, db := range this.db {
		if k := db.Lookup(device, code); k != gopi.KEYCODE_NONE {
			keys = append(keys, k)
		}
	}
	return keys
}

// Set keycode for name,device and scancode
func (this *Manager) Set(device gopi.InputDevice, code uint32, key gopi.KeyCode, name string) error {
	// Set mapping in existing or new database
	if db, err := this.DatabaseForName(name); err != nil {
		return err
	} else if err := db.Set(device, code, key); err != nil {
		return err
	}

	// Return sucess
	return nil
}

func (this *Manager) DatabaseForName(name string) (*keycodedb, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Return existing database
	for _, db := range this.db {
		if db.Name() == name {
			return db, nil
		}
	}

	// Check name parameter
	if reKeycodeName.MatchString(name) == false {
		return nil, gopi.ErrBadParameter.WithPrefix("DatabaseForName", name)
	}

	// Create a new database
	path := filepath.Join(*this.folder, name+*this.ext)
	if db, err := NewDatabase(path, name, this.Cache); err != nil {
		return nil, err
	} else {
		this.db = append(this.db, db)
		return db, nil
	}
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
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Manager) writeDirty() error {
	var result error

	// Write all modified databases
	for _, db := range this.db {
		if db.Modified() {
			this.Debug("LIRCKeycodeManager: Writing:", db.Name())
			if err := db.Write(this.Cache); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	// Return any errors
	return result
}

func (this *Manager) readAll(path, ext string) ([]*keycodedb, error) {
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
		} else if db, err := NewDatabase(filepath.Join(path, file.Name()), name, this.Cache); err != nil {
			return nil, err
		} else {
			this.Debug("LIRCKeycodeManager: Read:", db.Name())
			dbs = append(dbs, db)
		}
	}

	// Return databases
	return dbs, nil
}
