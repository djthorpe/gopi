package layout /* import "github.com/djthorpe/gopi/sys/default/layout" */

import (
	"errors"
	"io"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Config defines a default configuration for a layout
type Config struct {
	Direction gopi.LayoutDirection
}

type driver struct {
	log gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	ErrNotImplemented = errors.New("Not Implemented")
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register layout module
	registerFlags(gopi.RegisterModule(gopi.Module{
		Name: "default/layout",
		Type: gopi.MODULE_TYPE_LAYOUT,
		New:  newLayout,
	}))
}

func registerFlags(flags *util.Flags) {
	/* no flags for the layout module */
}

func newLayout(config *gopi.AppConfig, logger gopi.Logger) (gopi.Driver, error) {
	var err gopi.Error
	if layout, ok := gopi.Open2(Config{}, logger, &err).(gopi.Layout); !ok {
		return nil, err
	} else {
		return layout, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config Config) Open(logger gopi.Logger) (gopi.Driver, error) {
	this := new(driver)
	this.log = logger

	this.log.Debug2("gopi.sys.default.layout.Open()")

	return this, nil
}

func (this *driver) Close() error {
	this.log.Debug2("gopi.sys.default.layout.Close()")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *driver) View() gopi.View {
	return nil
}

// CalculateLayout performs the layout using the root node as the correct size
func (this *driver) CalculateLayout() bool {
	return false
}

// CalculateLayoutWithSize calculates layout with a new root size
func (this *driver) CalculateLayoutWithSize(w, h float32) bool {
	return false
}

// Encode return XML encoded version of the layout
func (this *driver) Encode(w io.Writer, indent gopi.EncodeIndentOptions) error {
	return ErrNotImplemented
}
