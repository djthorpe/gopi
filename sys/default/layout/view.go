package layout /* import "github.com/djthorpe/gopi/sys/default/layout" */

import (
	"github.com/djthorpe/gopi/third_party/flex"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type view struct {
	root  bool
	tag   uint
	class string
	node  *flex.Node
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *view) SetMargin(top, right, bottom, left string) {

}

func (this *view) SetPadding(top, right, bottom, left string) {

}

func (this *view) SetSize(width, height string) {

}
