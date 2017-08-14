package layout /* import "github.com/djthorpe/gopi/sys/default/layout" */

import (
	"fmt"
	"math"
	"strings"

	"github.com/djthorpe/gopi"
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
// GLOBAL VARIABLES

var (
	flexEdgeMap = map[gopi.ViewEdge]flex.Edge{
		gopi.VIEW_EDGE_ALL:    flex.EdgeAll,
		gopi.VIEW_EDGE_TOP:    flex.EdgeTop,
		gopi.VIEW_EDGE_BOTTOM: flex.EdgeBottom,
		gopi.VIEW_EDGE_LEFT:   flex.EdgeLeft,
		gopi.VIEW_EDGE_RIGHT:  flex.EdgeRight,
	}
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - GETTERS

// Return view tag or zero if not defined
func (this *view) Tag() uint {
	return this.tag
}

// Return view class
func (this *view) Class() string {
	return this.class
}

// Return position value (relative or absolute)
func (this *view) Positioning() gopi.ViewPositioning {
	switch this.node.Style.PositionType {
	case flex.PositionTypeRelative:
		return gopi.VIEW_POSITIONING_RELATIVE
	case flex.PositionTypeAbsolute:
		return gopi.VIEW_POSITIONING_ABSOLUTE
	default:
		panic("Invalid ViewPositioning value")
	}
}

// Return direction value
func (this *view) Direction() gopi.ViewDirection {
	switch this.node.Style.FlexDirection {
	case flex.FlexDirectionColumn:
		return gopi.VIEW_DIRECTION_COLUMN
	case flex.FlexDirectionColumnReverse:
		return gopi.VIEW_DIRECTION_COLUMN_REVERSE
	case flex.FlexDirectionRow:
		return gopi.VIEW_DIRECTION_ROW
	case flex.FlexDirectionRowReverse:
		return gopi.VIEW_DIRECTION_ROW_REVERSE
	default:
		panic("Invalid ViewDirection value")
	}
}

// Return justify value
func (this *view) Justify() gopi.ViewJustify {
	switch this.node.Style.JustifyContent {
	case flex.JustifyCenter:
		return gopi.VIEW_JUSTIFY_CENTER
	case flex.JustifyFlexEnd:
		return gopi.VIEW_JUSTIFY_FLEX_END
	case flex.JustifyFlexStart:
		return gopi.VIEW_JUSTIFY_FLEX_START
	case flex.JustifySpaceAround:
		return gopi.VIEW_JUSTIFY_SPACE_AROUND
	case flex.JustifySpaceBetween:
		return gopi.VIEW_JUSTIFY_SPACE_BETWEEN
	default:
		panic("Invalid ViewJustify value")
	}
}

// Return wrap value
func (this *view) Wrap() gopi.ViewWrap {
	switch this.node.Style.FlexWrap {
	case flex.WrapWrap:
		return gopi.VIEW_WRAP_ON
	case flex.WrapWrapReverse:
		return gopi.VIEW_WRAP_REVERSE
	case flex.WrapNoWrap:
		return gopi.VIEW_WRAP_OFF
	default:
		panic("Invalid ViewWrap value")
	}
}

// Return align value
func (this *view) Align() gopi.ViewAlign {
	switch this.node.Style.AlignContent {
	case flex.AlignAuto:
		return gopi.VIEW_ALIGN_AUTO
	case flex.AlignFlexStart:
		return gopi.VIEW_ALIGN_FLEX_START
	case flex.AlignCenter:
		return gopi.VIEW_ALIGN_CENTER
	case flex.AlignFlexEnd:
		return gopi.VIEW_ALIGN_FLEX_END
	case flex.AlignStretch:
		return gopi.VIEW_ALIGN_STRETCH
	case flex.AlignBaseline:
		return gopi.VIEW_ALIGN_BASELINE
	case flex.AlignSpaceAround:
		return gopi.VIEW_ALIGN_SPACE_AROUND
	case flex.AlignSpaceBetween:
		return gopi.VIEW_ALIGN_SPACE_BETWEEN
	default:
		panic("Invalid ViewAlign value")
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - SETTERS

func (this *view) SetDirection(value gopi.ViewDirection) {
	switch value {
	case gopi.VIEW_DIRECTION_COLUMN:
		this.node.StyleSetFlexDirection(flex.FlexDirectionColumn)
	case gopi.VIEW_DIRECTION_COLUMN_REVERSE:
		this.node.StyleSetFlexDirection(flex.FlexDirectionColumnReverse)
	case gopi.VIEW_DIRECTION_ROW:
		this.node.StyleSetFlexDirection(flex.FlexDirectionRow)
	case gopi.VIEW_DIRECTION_ROW_REVERSE:
		this.node.StyleSetFlexDirection(flex.FlexDirectionRowReverse)
	default:
		panic("Invalid ViewDirection value")
	}
}

func (this *view) SetJustify(value gopi.ViewJustify) {
	switch value {
	case gopi.VIEW_JUSTIFY_FLEX_START:
		this.node.StyleSetJustifyContent(flex.JustifyFlexStart)
	case gopi.VIEW_JUSTIFY_FLEX_END:
		this.node.StyleSetJustifyContent(flex.JustifyFlexEnd)
	case gopi.VIEW_JUSTIFY_CENTER:
		this.node.StyleSetJustifyContent(flex.JustifyCenter)
	case gopi.VIEW_JUSTIFY_SPACE_BETWEEN:
		this.node.StyleSetJustifyContent(flex.JustifySpaceBetween)
	case gopi.VIEW_JUSTIFY_SPACE_AROUND:
		this.node.StyleSetJustifyContent(flex.JustifySpaceAround)
	default:
		panic("Invalid ViewJustify value")
	}
}

func (this *view) SetWrap(value gopi.ViewWrap) {
	switch value {
	case gopi.VIEW_WRAP_OFF:
		this.node.StyleSetFlexWrap(flex.WrapNoWrap)
	case gopi.VIEW_WRAP_ON:
		this.node.StyleSetFlexWrap(flex.WrapWrap)
	case gopi.VIEW_WRAP_REVERSE:
		this.node.StyleSetFlexWrap(flex.WrapWrapReverse)
	default:
		panic("Invalid ViewWrap value")
	}
}

func (this *view) SetAlign(value gopi.ViewAlign) {
	switch value {
	case gopi.VIEW_ALIGN_AUTO:
		this.node.StyleSetAlignContent(flex.AlignAuto)
	case gopi.VIEW_ALIGN_FLEX_START:
		this.node.StyleSetAlignContent(flex.AlignFlexStart)
	case gopi.VIEW_ALIGN_CENTER:
		this.node.StyleSetAlignContent(flex.AlignCenter)
	case gopi.VIEW_ALIGN_FLEX_END:
		this.node.StyleSetAlignContent(flex.AlignFlexEnd)
	case gopi.VIEW_ALIGN_STRETCH:
		this.node.StyleSetAlignContent(flex.AlignStretch)
	case gopi.VIEW_ALIGN_BASELINE:
		this.node.StyleSetAlignContent(flex.AlignBaseline)
	case gopi.VIEW_ALIGN_SPACE_BETWEEN:
		this.node.StyleSetAlignContent(flex.AlignSpaceBetween)
	case gopi.VIEW_ALIGN_SPACE_AROUND:
		this.node.StyleSetAlignContent(flex.AlignSpaceAround)
	default:
		panic("Invalid ViewAlign value")
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - SET POSITION, MARGIN AND PADDING

func (this *view) SetPositionValue(value float32, edges ...gopi.ViewEdge) {
	for _, edge := range edges {
		if edge == gopi.VIEW_EDGE_ALL {
			this.SetPositionValue(value, gopi.VIEW_EDGE_ALL, gopi.VIEW_EDGE_TOP, gopi.VIEW_EDGE_RIGHT, gopi.VIEW_EDGE_BOTTOM, gopi.VIEW_EDGE_LEFT)
			return
		}
		this.node.StyleSetPosition(flexEdge(edge), value)
	}
}

func (this *view) SetPositionPercent(percent float32, edges ...gopi.ViewEdge) {
	for _, edge := range edges {
		if edge == gopi.VIEW_EDGE_ALL {
			this.SetPositionPercent(percent, gopi.VIEW_EDGE_ALL, gopi.VIEW_EDGE_TOP, gopi.VIEW_EDGE_RIGHT, gopi.VIEW_EDGE_BOTTOM, gopi.VIEW_EDGE_LEFT)
			return
		}
		this.node.StyleSetPositionPercent(flexEdge(edge), percent)
	}
}

func (this *view) SetMarginValue(value float32, edges ...gopi.ViewEdge) {
	for _, edge := range edges {
		if edge == gopi.VIEW_EDGE_ALL {
			this.SetMarginValue(value, gopi.VIEW_EDGE_ALL, gopi.VIEW_EDGE_TOP, gopi.VIEW_EDGE_RIGHT, gopi.VIEW_EDGE_BOTTOM, gopi.VIEW_EDGE_LEFT)
			return
		}
		this.node.StyleSetMargin(flexEdge(edge), value)
	}
}

func (this *view) SetMarginPercent(percent float32, edges ...gopi.ViewEdge) {
	for _, edge := range edges {
		if edge == gopi.VIEW_EDGE_ALL {
			this.SetMarginPercent(percent, gopi.VIEW_EDGE_ALL, gopi.VIEW_EDGE_TOP, gopi.VIEW_EDGE_RIGHT, gopi.VIEW_EDGE_BOTTOM, gopi.VIEW_EDGE_LEFT)
			return
		}
		this.node.StyleSetMarginPercent(flexEdge(edge), percent)
	}
}

func (this *view) SetMarginAuto(edges ...gopi.ViewEdge) {
	for _, edge := range edges {
		if edge == gopi.VIEW_EDGE_ALL {
			this.SetMarginAuto(gopi.VIEW_EDGE_ALL, gopi.VIEW_EDGE_TOP, gopi.VIEW_EDGE_RIGHT, gopi.VIEW_EDGE_BOTTOM, gopi.VIEW_EDGE_LEFT)
			return
		}
		this.node.StyleSetMarginAuto(flexEdge(edge))
	}
}

func (this *view) SetPaddingValue(value float32, edges ...gopi.ViewEdge) {
	for _, edge := range edges {
		if edge == gopi.VIEW_EDGE_ALL {
			this.SetPaddingValue(value, gopi.VIEW_EDGE_ALL, gopi.VIEW_EDGE_TOP, gopi.VIEW_EDGE_RIGHT, gopi.VIEW_EDGE_BOTTOM, gopi.VIEW_EDGE_LEFT)
			return
		}
		this.node.StyleSetPadding(flexEdge(edge), value)
	}
}

func (this *view) SetPaddingPercent(percent float32, edges ...gopi.ViewEdge) {
	for _, edge := range edges {
		if edge == gopi.VIEW_EDGE_ALL {
			this.SetPaddingPercent(percent, gopi.VIEW_EDGE_ALL, gopi.VIEW_EDGE_TOP, gopi.VIEW_EDGE_RIGHT, gopi.VIEW_EDGE_BOTTOM, gopi.VIEW_EDGE_LEFT)
			return
		}
		this.node.StyleSetPaddingPercent(flexEdge(edge), percent)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - GET POSITION, MARGIN AND PADDING

func (this *view) PositionString(edges ...gopi.ViewEdge) string {
	edges_string := make([]string, 0, len(edges))
	for _, edge := range edges {
		if edge == gopi.VIEW_EDGE_ALL {
			return this.PositionString(gopi.VIEW_EDGE_TOP, gopi.VIEW_EDGE_RIGHT, gopi.VIEW_EDGE_BOTTOM, gopi.VIEW_EDGE_LEFT)
		}
		edges_string = append(edges_string, flexEdgeString(this.node.StyleGetPosition(flexEdge(edge))))
	}
	return strings.Join(edges_string, " ")
}

func (this *view) MarginString(edges ...gopi.ViewEdge) string {
	edges_string := make([]string, 0, len(edges))
	for _, edge := range edges {
		if edge == gopi.VIEW_EDGE_ALL {
			return this.MarginString(gopi.VIEW_EDGE_TOP, gopi.VIEW_EDGE_RIGHT, gopi.VIEW_EDGE_BOTTOM, gopi.VIEW_EDGE_LEFT)
		}
		edges_string = append(edges_string, flexEdgeString(this.node.StyleGetMargin(flexEdge(edge))))
	}
	return strings.Join(edges_string, " ")
}

func (this *view) PaddingString(edges ...gopi.ViewEdge) string {
	edges_string := make([]string, 0, len(edges))
	for _, edge := range edges {
		if edge == gopi.VIEW_EDGE_ALL {
			return this.PaddingString(gopi.VIEW_EDGE_TOP, gopi.VIEW_EDGE_RIGHT, gopi.VIEW_EDGE_BOTTOM, gopi.VIEW_EDGE_LEFT)
		}
		edges_string = append(edges_string, flexEdgeString(this.node.StyleGetPadding(flexEdge(edge))))
	}
	return strings.Join(edges_string, " ")
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - DIRTY

func (this *view) IsDirty() bool {
	return this.node.IsDirty
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func flexEdge(edge gopi.ViewEdge) flex.Edge {
	if edge, exists := flexEdgeMap[edge]; exists {
		return edge
	}
	panic("Invalid ViewEdge value")
}

func flexEdgeString(value flex.Value) string {
	if value.Unit == flex.UnitAuto {
		return "auto"
	}
	if math.IsNaN(float64(value.Value)) {
		// Can't check for NaN using ==
		return "inherit"
	}
	if value.Unit == flex.UnitPercent {
		return fmt.Sprintf("%v%%", value.Value)
	}
	if value.Unit == flex.UnitPoint {
		return fmt.Sprintf("%v", value.Value)
	}
	panic(value.Value == gopi.EdgeUndefined)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *view) String() string {
	parts := make([]string, 0)
	parts = append(parts, fmt.Sprintf("positioning=%v", this.Positioning()))
	if this.Tag() != gopi.TagNone {
		parts = append(parts, fmt.Sprintf("tag=%v", this.Tag()))
	}
	if this.Class() != "" {
		parts = append(parts, fmt.Sprintf("class=%v", this.Class()))
	}
	switch this.Positioning() {
	case gopi.VIEW_POSITIONING_ABSOLUTE:
		parts = append(parts, fmt.Sprintf("position=\"%v\"", this.PositionString(gopi.VIEW_EDGE_ALL)))
	case gopi.VIEW_POSITIONING_RELATIVE:
		parts = append(parts, fmt.Sprintf("direction=\"%v\"", this.Direction()))
		parts = append(parts, fmt.Sprintf("justify=\"%v\"", this.Justify()))
		parts = append(parts, fmt.Sprintf("wrap=\"%v\"", this.Wrap()))
		parts = append(parts, fmt.Sprintf("align=\"%v\"", this.Align()))
	}
	parts = append(parts, fmt.Sprintf("margin=\"%v\"", this.MarginString(gopi.VIEW_EDGE_ALL)))
	parts = append(parts, fmt.Sprintf("padding=\"%v\"", this.PaddingString(gopi.VIEW_EDGE_ALL)))
	return fmt.Sprintf("gopi.View{ %v }", strings.Join(parts, " "))
}
