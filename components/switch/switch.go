package switch_comp

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

type Switch struct {
	id       string
	label    string
	checked  bool
	disabled bool
	size     Size
	onChange func(checked bool)
}

func New(id string) *Switch {
	return &Switch{
		id: id,
		size: SizeMedium,
	}
}

func (s *Switch) Id(id string) *Switch {
	s.id = id
	return s
}

func (s *Switch) Label(label string) *Switch {
	s.label = label
	return s
}

func (s *Switch) Checked(checked bool) *Switch {
	s.checked = checked
	return s
}

func (s *Switch) Disabled(disabled bool) *Switch {
	s.disabled = disabled
	return s
}

func (s *Switch) Size(size Size) *Switch {
	s.size = size
	return s
}

func (s *Switch) OnChange(f func(checked bool)) *Switch {
	s.onChange = f
	return s
}

func (s *Switch) getTrackSize() (float64, float64) {
	switch s.size {
	case SizeXSmall, SizeSmall:
		return 28.0, 16.0
	case SizeLarge:
		return 44.0, 24.0
	default: // SizeMedium
		return 36.0, 20.0
	}
}

func (s *Switch) getThumbSize() float64 {
	switch s.size {
	case SizeXSmall, SizeSmall:
		return 12.0
	case SizeLarge:
		return 20.0
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

func (s *Switch) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	ctx := uictx.GG
	th := uictx.Theme

	trackW, trackH := s.getTrackSize()
	thumbSize := s.getThumbSize()
	inset := 2.0
	radius := trackH / 2.0

	// Calculate text width if label exists
	var textW, textH float64
	if s.label != "" {
		textW, textH = ctx.MeasureString(s.label)
	}

	gap := 8.0
	totalW := trackW
	if textW > 0 {
		totalW += gap + textW
	}
	totalH := trackH
	if textH > trackH {
		totalH = textH
	}

	if uictx.MeasureOnly {
		return totalW, totalH
	}

	// Process interaction
	_, _, isClicked := uictx.ProcessInteraction(s.id, x, y, totalW, totalH, s.disabled)

	if isClicked && s.onChange != nil {
		s.onChange(!s.checked)
	}

	// Center items vertically
	centerY := y + totalH/2
	trackY := centerY - trackH/2

	// Smooth checked transition
	state := uictx.GetState(s.id)
	
	targetToggle := 0.0
	if s.checked {
		targetToggle = 1.0
	}
	state.ToggleRatio = uictx.Animate(state.ToggleRatio, targetToggle, 10.0)

	// Layout colors
	// Switch track (background)
	uncheckedBg := th.Colors.Input
	checkedBg := th.Colors.Primary
	
	// Opacity helper
	withAlpha := func(c color.Color, alpha uint8) color.Color {
		r, g, b, _ := c.RGBA()
		return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: alpha}
	}

	disabledCheckedBg := withAlpha(checkedBg, 127)

	bg := uncheckedBg
	if s.disabled {
		if s.checked {
			bg = disabledCheckedBg
		}
	} else {
		if s.checked {
			bg = checkedBg
		}
		// Hover effect (brighten slightly)
		if state.HoverRatio > 0 {
			targetHoverColor := color.White
			if s.checked {
				targetHoverColor = color.Black
			}
			bg = mixColor(bg, targetHoverColor, state.HoverRatio*0.1)
		}
	}

	// Active (click) effect - thumb widens slightly
	actualThumbW := thumbSize
	if !s.disabled && state.ActiveRatio > 0 {
		actualThumbW += state.ActiveRatio * (trackW * 0.1) // widen by 10% of track width
	}

	// Draw Track
	ctx.DrawRoundedRectangle(x, trackY, trackW, trackH, radius)
	ctx.SetColor(bg)
	ctx.Fill()

	// Switch thumb
	thumbBg := th.Colors.Background
	if s.disabled {
		thumbBg = withAlpha(thumbBg, 89) // opacity 0.35
	}
	
	// Calculate thumb position
	// offset 0 when unchecked, (trackW - actualThumbW - inset*2) when checked
	thumbMaxLeft := trackW - actualThumbW - inset*2
	thumbOffset := state.ToggleRatio * thumbMaxLeft
	thumbX := x + inset + thumbOffset
	thumbY := trackY + inset

	ctx.DrawRoundedRectangle(thumbX, thumbY, actualThumbW, thumbSize, thumbSize/2.0)
	ctx.SetColor(thumbBg)
	ctx.Fill()

	// Draw Label
	if s.label != "" {
		textColor := th.Colors.Foreground
		if s.disabled {
			textColor = th.Colors.MutedForeground
		}
		ctx.SetColor(textColor)
		textX := x + trackW + gap
		ctx.DrawStringAnchored(s.label, textX, centerY, 0.0, 0.5)
	}

	return totalW, totalH
}
