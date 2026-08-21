package checkbox

import (
	"image/color"

	"github.com/sofiagros/gogpui/core/context"
)

type Size int

const (
	SizeMedium Size = iota
	SizeSmall
	SizeXSmall
	SizeLarge
)

type Checkbox struct {
	id       string
	label    string
	checked  bool
	disabled bool
	size     Size
	onChange func(checked bool)
}

func New(id string) *Checkbox {
	return &Checkbox{
		id: id,
		size: SizeMedium,
	}
}

func (c *Checkbox) Id(id string) *Checkbox {
	c.id = id
	return c
}

func (c *Checkbox) actualID() string {
	if c.id != "" {
		return c.id
	}
	if c.label != "" {
		return "checkbox:" + c.label
	}
	return "checkbox"
}

func (c *Checkbox) Label(label string) *Checkbox {
	c.label = label
	return c
}

func (c *Checkbox) Checked(checked bool) *Checkbox {
	c.checked = checked
	return c
}

func (c *Checkbox) Disabled(disabled bool) *Checkbox {
	c.disabled = disabled
	return c
}

func (c *Checkbox) Size(size Size) *Checkbox {
	c.size = size
	return c
}

func (c *Checkbox) OnChange(f func(checked bool)) *Checkbox {
	c.onChange = f
	return c
}

func (c *Checkbox) getBoxSize() float64 {
	switch c.size {
	case SizeXSmall:
		return 12.0
	case SizeSmall:
		return 14.0
	case SizeLarge:
		return 18.0
	default: // SizeMedium
		return 16.0
	}
}

// Helper to mix two colors (simple lerp)
func mixColor(c1, c2 color.Color, ratio float64) color.Color {
	if ratio <= 0 {
		return c1
	}
	if ratio >= 1 {
		return c2
	}
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	r := uint8((float64(r1>>8)*(1-ratio) + float64(r2>>8)*ratio))
	g := uint8((float64(g1>>8)*(1-ratio) + float64(g2>>8)*ratio))
	bl := uint8((float64(b1>>8)*(1-ratio) + float64(b2>>8)*ratio))
	a := uint8((float64(a1>>8)*(1-ratio) + float64(a2>>8)*ratio))

	return color.RGBA{R: r, G: g, B: bl, A: a}
}

func (c *Checkbox) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	ctx := uictx.GG
	th := uictx.Theme

	boxSize := c.getBoxSize()
	radius := float64(th.Radius) * 0.5
	if radius > 4.0 {
		radius = 4.0 // max radius
	}

	// Calculate text width if label exists
	var textW, textH float64
	if c.label != "" {
		textW, textH = ctx.MeasureString(c.label)
	}

	gap := 8.0
	totalW := boxSize
	if textW > 0 {
		totalW += gap + textW
	}
	totalH := boxSize
	if textH > boxSize {
		totalH = textH
	}

	if uictx.MeasureOnly {
		return totalW, totalH
	}

	// Process interaction
	_, _, isClicked := uictx.ProcessInteraction(c.id, x, y, totalW, totalH, c.disabled)

	if isClicked && c.onChange != nil {
		c.onChange(!c.checked)
	}

	// Center items vertically
	centerY := y + totalH/2
	boxY := centerY - boxSize/2

	// Layout colors
	uncheckedBorder := th.Colors.Input
	checkedColor := th.Colors.Primary
	
	// Opacity helper
	withAlpha := func(c color.Color, alpha uint8) color.Color {
		r, g, b, _ := c.RGBA()
		return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: alpha}
	}

	disabledColor := withAlpha(uncheckedBorder, 127)
	disabledCheckedColor := withAlpha(checkedColor, 127)

	bgColor := th.Colors.Background
	borderColor := uncheckedBorder

	state := uictx.GetState(c.id)

	if c.disabled {
		if c.checked {
			bgColor = disabledCheckedColor
			borderColor = disabledCheckedColor
		} else {
			borderColor = disabledColor
		}
	} else {
		if c.checked {
			bgColor = checkedColor
			borderColor = checkedColor
			
			// Hover effect when checked
			if state.HoverRatio > 0 {
				bgColor = mixColor(bgColor, color.Black, state.HoverRatio*0.1)
				borderColor = mixColor(borderColor, color.Black, state.HoverRatio*0.1)
			}
		} else {
			// Hover effect when unchecked (subtle border color change)
			if state.HoverRatio > 0 {
				borderColor = mixColor(borderColor, th.Colors.Foreground, state.HoverRatio*0.3)
			}
		}
	}

	// Active (click) effect - slightly shift down
	yOffset := 0.0
	if !c.disabled && state.ActiveRatio > 0 {
		yOffset = state.ActiveRatio * 1.0
	}
	boxY += yOffset
	centerY += yOffset

	// Draw Box
	ctx.DrawRoundedRectangle(x, boxY, boxSize, boxSize, radius)
	ctx.SetColor(bgColor)
	ctx.FillPreserve()
	ctx.SetColor(borderColor)
	ctx.SetLineWidth(1.0)
	ctx.Stroke()

	// Draw Checkmark
	if c.checked {
		iconColor := th.Colors.PrimaryForeground
		if c.disabled {
			iconColor = withAlpha(iconColor, 127)
		}
		ctx.SetColor(iconColor)
		ctx.SetLineWidth(2.0)
		ctx.MoveTo(x+boxSize*0.25, boxY+boxSize*0.55)
		ctx.LineTo(x+boxSize*0.45, boxY+boxSize*0.75)
		ctx.LineTo(x+boxSize*0.75, boxY+boxSize*0.3)
		ctx.Stroke()
	}

	// Draw Label
	if c.label != "" {
		textColor := th.Colors.Foreground
		if c.disabled {
			textColor = th.Colors.MutedForeground
		}
		ctx.SetColor(textColor)
		textX := x + boxSize + gap
		ctx.DrawStringAnchored(c.label, textX, centerY, 0.0, 0.5)
	}

	return totalW, totalH
}
