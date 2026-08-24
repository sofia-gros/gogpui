package radio

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

type Radio struct {
	id       string
	label    string
	checked  bool
	disabled bool
	size     Size
	onChange func(checked bool)
}

func New(id string) *Radio {
	return &Radio{
		id: id,
		size: SizeMedium,
	}
}

func (r *Radio) Id(id string) *Radio {
	r.id = id
	return r
}

func (r *Radio) Label(label string) *Radio {
	r.label = label
	return r
}

func (r *Radio) Checked(checked bool) *Radio {
	r.checked = checked
	return r
}

func (r *Radio) Disabled(disabled bool) *Radio {
	r.disabled = disabled
	return r
}

func (r *Radio) Size(size Size) *Radio {
	r.size = size
	return r
}

func (r *Radio) OnChange(f func(checked bool)) *Radio {
	r.onChange = f
	return r
}

func (r *Radio) getBoxSize() float64 {
	switch r.size {
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

	rr := uint8((float64(r1>>8)*(1-ratio) + float64(r2>>8)*ratio))
	gg := uint8((float64(g1>>8)*(1-ratio) + float64(g2>>8)*ratio))
	bb := uint8((float64(b1>>8)*(1-ratio) + float64(b2>>8)*ratio))
	aa := uint8((float64(a1>>8)*(1-ratio) + float64(a2>>8)*ratio))

	return color.RGBA{R: rr, G: gg, B: bb, A: aa}
}

func (r *Radio) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	ctx := uictx.GG
	th := uictx.Theme

	boxSize := r.getBoxSize()
	circleRadius := boxSize / 2.0

	// Calculate text width if label exists
	var textW, textH float64
	if r.label != "" {
		textW, textH = uictx.MeasureText(r.label)
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
	_, _, isClicked := uictx.ProcessInteraction(r.id, x, y, totalW, totalH, r.disabled)

	if isClicked && r.onChange != nil {
		r.onChange(!r.checked)
	}

	// Center items vertically
	centerY := y + totalH/2
	centerX := x + circleRadius

	// Layout colors
	uncheckedBorder := th.Colors.Input
	checkedColor := th.Colors.Primary
	
	// Opacity helper
	withAlpha := func(c color.Color, alpha uint8) color.Color {
		colR, colG, colB, _ := c.RGBA()
		return color.RGBA{R: uint8(colR >> 8), G: uint8(colG >> 8), B: uint8(colB >> 8), A: alpha}
	}

	disabledColor := withAlpha(uncheckedBorder, 127)
	disabledCheckedColor := withAlpha(checkedColor, 127)

	borderColor := uncheckedBorder
	
	state := uictx.GetState(r.id)

	if r.disabled {
		if r.checked {
			borderColor = disabledCheckedColor
		} else {
			borderColor = disabledColor
		}
	} else {
		if r.checked {
			borderColor = checkedColor
			// Hover effect when checked
			if state.HoverRatio > 0 {
				borderColor = mixColor(borderColor, color.Black, state.HoverRatio*0.1)
			}
		} else {
			// Hover effect when unchecked
			if state.HoverRatio > 0 {
				borderColor = mixColor(borderColor, th.Colors.Foreground, state.HoverRatio*0.3)
			}
		}
	}

	// Active (click) effect - slightly shift down
	yOffset := 0.0
	if !r.disabled && state.ActiveRatio > 0 {
		yOffset = state.ActiveRatio * 1.0
	}
	centerY += yOffset

	// 外枠を描く。ストローク幅をスケール対応にして物理 1.5px を保証する。
	// Windows 拡大表示(125%/150%)環境では 1 論理 px が見えないことがあるため。
	lineW := 1.5 / uictx.Scale
	if lineW < 1.0 {
		lineW = 1.0
	}
	ctx.DrawCircle(centerX, centerY, circleRadius)
	ctx.SetColor(th.Colors.Background)
	ctx.FillPreserve()
	ctx.SetColor(borderColor)
	ctx.SetLineWidth(lineW)
	uictx.Stroke()

	// Draw inner circle (dot) if checked
	if r.checked {
		dotColor := th.Colors.Primary
		if r.disabled {
			dotColor = withAlpha(dotColor, 127)
		}
		
		dotRadius := circleRadius * 0.6
		if !r.disabled && state.ActiveRatio > 0 {
			dotRadius *= (1.0 - state.ActiveRatio*0.2) // scale down dot on click
		}

		ctx.DrawCircle(centerX, centerY, dotRadius)
		ctx.SetColor(dotColor)
		uictx.Fill()
	}

	// Draw Label
	if r.label != "" {
		textColor := th.Colors.Foreground
		if r.disabled {
			textColor = th.Colors.MutedForeground
		}
		ctx.SetColor(textColor)
		textX := x + boxSize + gap
		uictx.DrawStringAnchored(r.label, textX, centerY-yOffset, 0.0, 0.5) // align label correctly if shifted
	}

	return totalW, totalH
}
