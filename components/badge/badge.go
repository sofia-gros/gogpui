package badge

import (
	"image/color"
	"strconv"

	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/layout"
)

type Variant int

const (
	Number Variant = iota
	Dot
)

type Size int

const (
	Small Size = iota
	Medium
	Large
)

type Badge struct {
	child   layout.Widget
	variant Variant
	count   int
	max     int
	color   *color.Color
	size    Size
}

func New() *Badge {
	return &Badge{
		variant: Number,
		max:     99,
		size:    Medium,
	}
}

func (b *Badge) Add(child layout.Widget) *Badge {
	b.child = child
	return b
}

func (b *Badge) Dot() *Badge {
	b.variant = Dot
	return b
}

func (b *Badge) Count(count int) *Badge {
	b.count = count
	return b
}

func (b *Badge) Max(max int) *Badge {
	b.max = max
	return b
}

func (b *Badge) Color(c color.Color) *Badge {
	b.color = &c
	return b
}

func (b *Badge) Size(s Size) *Badge {
	b.size = s
	return b
}

func (b *Badge) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	ctx := uictx.GG
	th := uictx.Theme

	var childW, childH float64
	if b.child != nil {
		childW, childH = b.child.Render(uictx, x, y)
	}

	visible := true
	if b.variant == Number && b.count <= 0 {
		visible = false
	}

	var badgeW, badgeH float64
	var textStr string
	if visible {
		if b.variant == Dot {
			badgeW, badgeH = 6.0, 6.0
		} else {
			textStr = strconv.Itoa(b.count)
			if b.count > b.max {
				textStr = strconv.Itoa(b.max) + "+"
			}
			w, h := uictx.MeasureText(textStr)
			if h == 0 {
				_, h = uictx.MeasureText("0")
			}
			
			px, py := 4.0, 2.0
			badgeW = w + px*2
			badgeH = h + py*2
			
			// Minimum size to make it look like a circle if single digit
			if badgeW < badgeH {
				badgeW = badgeH
			}
		}
	}

	if uictx.MeasureOnly {
		if b.child == nil && visible {
			return badgeW, badgeH
		}
		return childW, childH
	}

	if !visible {
		return childW, childH
	}

	// Calculate badge position
	var bx, by float64
	if b.child != nil {
		// Top right corner of child, overlapping slightly
		bx = x + childW - badgeW/2.0
		by = y - badgeH/2.0
	} else {
		// Standalone
		bx = x
		by = y
	}

	// Draw badge background
	bgColor := th.Colors.Danger
	if b.color != nil {
		bgColor = *b.color
	}

	ctx.SetColor(bgColor)
	radius := badgeH / 2.0
	uictx.DrawRoundedRectangle(bx, by, badgeW, badgeH, radius)
	uictx.Fill()

	// Draw badge text
	if b.variant == Number {
		ctx.SetColor(color.White)
		uictx.DrawStringAnchored(textStr, bx+badgeW/2.0, by+badgeH/2.0, 0.5, 0.5)
	}

	if b.child == nil {
		return badgeW, badgeH
	}
	return childW, childH
}
