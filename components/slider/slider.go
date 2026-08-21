package slider

import (
	"image/color"

	"github.com/sofiagros/gogpui/core/context"
)

type Slider struct {
	id       string
	value    float64 // 0.0 to 1.0
	disabled bool
	onChange func(val float64)
}

func New(id string) *Slider {
	return &Slider{
		id: id,
		value: 0.5,
	}
}

func (s *Slider) Id(id string) *Slider {
	s.id = id
	return s
}

func (s *Slider) Value(val float64) *Slider {
	if val < 0 {
		val = 0
	} else if val > 1 {
		val = 1
	}
	s.value = val
	return s
}

func (s *Slider) Disabled(disabled bool) *Slider {
	s.disabled = disabled
	return s
}

func (s *Slider) OnChange(f func(val float64)) *Slider {
	s.onChange = f
	return s
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

func (s *Slider) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	ctx := uictx.GG
	th := uictx.Theme

	// Standard slider size in shadcn/gpui
	// Width is flexible, let's make it 200 by default for the layout
	trackW := 200.0
	trackH := 8.0
	thumbSize := 16.0
	radius := trackH / 2.0

	totalW := trackW
	totalH := thumbSize
	if trackH > totalH {
		totalH = trackH
	}

	if uictx.MeasureOnly {
		return totalW, totalH
	}

	state := uictx.GetState(s.id)

	// Center track vertically
	trackY := y + (totalH - trackH) / 2.0
	thumbY := y + (totalH - thumbSize) / 2.0

	hovered := !s.disabled && uictx.HitTest(x, thumbY, trackW, thumbSize)
	
	// Drag interaction
	if hovered && uictx.Mouse.LeftPressed {
		state.IsDragging = true
	}
	if !uictx.Mouse.LeftDown {
		state.IsDragging = false
	}

	currentVal := s.value
	if state.IsDragging {
		// Calculate new value based on mouse X
		relativeX := uictx.Mouse.X - x
		newVal := relativeX / trackW
		if newVal < 0 {
			newVal = 0
		} else if newVal > 1 {
			newVal = 1
		}
		
		if newVal != currentVal {
			currentVal = newVal
			if s.onChange != nil {
				s.onChange(currentVal)
			}
		}
	}

	// Update hover and active animation
	targetHover := 0.0
	if hovered {
		targetHover = 1.0
	}
	state.HoverRatio = uictx.Animate(state.HoverRatio, targetHover, 10.0)

	targetActive := 0.0
	if state.IsDragging {
		targetActive = 1.0
	}
	state.ActiveRatio = uictx.Animate(state.ActiveRatio, targetActive, 15.0)

	// Layout colors
	// Switch track (background) is usually Secondary or Muted
	uncheckedBg := th.Colors.Secondary
	checkedBg := th.Colors.Primary
	
	// Opacity helper
	withAlpha := func(c color.Color, alpha uint8) color.Color {
		r, g, b, _ := c.RGBA()
		return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: alpha}
	}

	disabledUncheckedBg := withAlpha(uncheckedBg, 127)
	disabledCheckedBg := withAlpha(checkedBg, 127)

	bg1 := checkedBg
	bg2 := uncheckedBg

	if s.disabled {
		bg1 = disabledCheckedBg
		bg2 = disabledUncheckedBg
	} else {
		// subtle hover effect on the filled track
		if state.HoverRatio > 0 {
			bg1 = mixColor(bg1, color.Black, state.HoverRatio*0.1)
		}
	}

	// Draw Track (Unfilled part)
	ctx.DrawRoundedRectangle(x, trackY, trackW, trackH, radius)
	ctx.SetColor(bg2)
	ctx.Fill()

	// Draw Track (Filled part)
	filledW := trackW * currentVal
	if filledW > 0 {
		ctx.DrawRoundedRectangle(x, trackY, filledW, trackH, radius)
		ctx.SetColor(bg1)
		ctx.Fill()
	}

	// Thumb animation
	actualThumbSize := thumbSize
	if !s.disabled {
		if state.HoverRatio > 0 {
			actualThumbSize += state.HoverRatio * 2.0
		}
		if state.ActiveRatio > 0 {
			actualThumbSize += state.ActiveRatio * 2.0
		}
	}

	// Draw Thumb
	thumbX := x + filledW - actualThumbSize/2.0
	if thumbX < x {
		thumbX = x
	} else if thumbX > x+trackW-actualThumbSize {
		thumbX = x + trackW - actualThumbSize
	}

	thumbBg := th.Colors.Background
	thumbBorder := th.Colors.Primary
	
	if s.disabled {
		thumbBorder = withAlpha(thumbBorder, 127)
	} else {
		if state.HoverRatio > 0 {
			thumbBorder = mixColor(thumbBorder, color.Black, state.HoverRatio*0.1)
		}
	}

	ctx.DrawCircle(thumbX+actualThumbSize/2.0, thumbY+actualThumbSize/2.0, actualThumbSize/2.0)
	ctx.SetColor(thumbBg)
	ctx.FillPreserve()
	ctx.SetColor(thumbBorder)
	ctx.SetLineWidth(2.0)
	ctx.Stroke()

	return totalW, totalH
}
