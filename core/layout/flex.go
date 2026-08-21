package layout

import (
	"math"

	"github.com/sofiagros/gogpui/core/context"
)

// Direction defines the main axis for the flex layout.
type Direction int

const (
	Row Direction = iota
	Column
)

// JustifyContent defines how items are distributed along the main axis.
type JustifyContent int

const (
	JustifyStart JustifyContent = iota
	JustifyCenter
	JustifyEnd
	JustifySpaceBetween
	JustifySpaceAround
)

// AlignItems defines how items are aligned along the cross axis.
type AlignItems int

const (
	AlignStart AlignItems = iota
	AlignCenter
	AlignEnd
	AlignStretch // Note: Stretch requires components to accept external size constraints
)

// Widget is the standard interface for all UI components.
type Widget interface {
	Render(uictx *context.UIContext, x, y float64) (w, h float64)
}

// Flex represents a flexbox-style layout container.
type Flex struct {
	direction Direction
	justify   JustifyContent
	align     AlignItems
	gap       float64
	children  []Widget
}

// NewFlex creates a new Flex container with default row direction.
func NewFlex() *Flex {
	return &Flex{
		direction: Row,
		justify:   JustifyStart,
		align:     AlignStart,
		gap:       0,
		children:  []Widget{},
	}
}

// Direction sets the main axis (Row or Column).
func (f *Flex) Direction(dir Direction) *Flex {
	f.direction = dir
	return f
}

// Justify sets the justification along the main axis.
func (f *Flex) Justify(j JustifyContent) *Flex {
	f.justify = j
	return f
}

// Align sets the alignment along the cross axis.
func (f *Flex) Align(a AlignItems) *Flex {
	f.align = a
	return f
}

// Gap sets the spacing between children.
func (f *Flex) Gap(gap float64) *Flex {
	f.gap = gap
	return f
}

// Add appends one or more children to the layout.
func (f *Flex) Add(children ...Widget) *Flex {
	f.children = append(f.children, children...)
	return f
}

// Render implements the Widget interface for Flex.
func (f *Flex) Render(uictx *context.UIContext, x, y float64) (w, h float64) {
	if len(f.children) == 0 {
		return 0, 0
	}

	// 1. Measure Pass
	originalMeasure := uictx.MeasureOnly
	uictx.MeasureOnly = true

	type childBox struct {
		w, h float64
	}
	boxes := make([]childBox, len(f.children))
	var mainTotal, crossMax float64

	for i, child := range f.children {
		cw, ch := child.Render(uictx, 0, 0)
		boxes[i] = childBox{cw, ch}
		
		var main, cross float64
		if f.direction == Row {
			main, cross = cw, ch
		} else {
			main, cross = ch, cw
		}

		mainTotal += main
		if i > 0 {
			mainTotal += f.gap
		}
		if cross > crossMax {
			crossMax = cross
		}
	}

	uictx.MeasureOnly = originalMeasure

	// Total dimensions of this flex container
	if f.direction == Row {
		w, h = mainTotal, crossMax
	} else {
		w, h = crossMax, mainTotal
	}

	// Early return if parent is just measuring us
	if uictx.MeasureOnly {
		return w, h
	}

	// 2. Layout & Render Pass
	var cursor float64
	
	// Handle JustifyContent calculation
	spacing := f.gap
	
	switch f.justify {
	case JustifyStart:
		cursor = 0
	case JustifyEnd:
		// Not fully supported unless we have bounded parent w/h, using 0 for unbounded
		cursor = 0 
	case JustifyCenter:
		// Requires bounded parent, unsupported in unconstrained, using 0
		cursor = 0
	case JustifySpaceBetween:
		if len(f.children) > 1 {
			// Without fixed parent bounds, this behaves like start with gaps
			spacing = f.gap
		}
		cursor = 0
	case JustifySpaceAround:
		cursor = f.gap / 2.0
		spacing = f.gap
	}

	for i, child := range f.children {
		box := boxes[i]
		var cx, cy float64

		if f.direction == Row {
			cx = x + cursor
			
			// Cross axis alignment
			switch f.align {
			case AlignStart, AlignStretch:
				cy = y
			case AlignEnd:
				cy = y + (crossMax - box.h)
			case AlignCenter:
				cy = y + (crossMax - box.h) / 2.0
			}

			cursor += box.w + spacing
		} else {
			cy = y + cursor

			// Cross axis alignment
			switch f.align {
			case AlignStart, AlignStretch:
				cx = x
			case AlignEnd:
				cx = x + (crossMax - box.w)
			case AlignCenter:
				cx = x + (crossMax - box.w) / 2.0
			}

			cursor += box.h + spacing
		}

		child.Render(uictx, math.Round(cx), math.Round(cy))
	}

	return w, h
}
