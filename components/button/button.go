package button

import (
	"image/color"

	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/size"
)

// Variant represents the visual style variant of the button.
type Variant int

const (
	Default Variant = iota
	Primary
	Secondary
	Danger
	Info
	Success
	Warning
	Ghost
	Link
	Text
)

// Size represents the size of the button.
type Size int

const (
	Sm Size = iota
	Md
	Lg
)

// Button represents a clickable UI element.
type Button struct {
	id       string
	label    string
	variant  Variant
	size     Size
	disabled bool
	onClick  func()
}

// New creates a new Button instance.
func New(id string) *Button {
	return &Button{
		id:      id,
		variant: Default,
		size:    Md,
	}
}

// Label sets the text label of the button.
func (b *Button) Label(text string) *Button {
	b.label = text
	return b
}
func (b *Button) Primary() *Button   { b.variant = Primary; return b }
func (b *Button) Secondary() *Button { b.variant = Secondary; return b }
func (b *Button) Danger() *Button    { b.variant = Danger; return b }
func (b *Button) Warning() *Button   { b.variant = Warning; return b }
func (b *Button) Success() *Button   { b.variant = Success; return b }
func (b *Button) Info() *Button      { b.variant = Info; return b }
func (b *Button) Ghost() *Button     { b.variant = Ghost; return b }
func (b *Button) Link() *Button      { b.variant = Link; return b }
func (b *Button) Text() *Button      { b.variant = Text; return b }

// Size sets the size of the button.
func (b *Button) Size(s Size) *Button {
	b.size = s
	return b
}

// Disabled disables or enables the button.
func (b *Button) Disabled(disabled bool) *Button {
	b.disabled = disabled
	return b
}

// OnClick sets the click event handler.
func (b *Button) OnClick(handler func()) *Button {
	b.onClick = handler
	return b
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

// Render draws the button onto the given UIContext.
func (b *Button) Render(uictx *context.UIContext, x, y float64) (w, h float64) {
	ctx := uictx.GG
	th := uictx.Theme

	// Calculate layout properties based on size variant
	var px, py, minHeight float64
	switch b.size {
	case Sm:
		px, py, minHeight = float64(size.Md), float64(size.Xs), 36.0
	case Lg:
		px, py, minHeight = float64(size.Xl), float64(size.Sm), 44.0
	default:
		px, py, minHeight = float64(size.Lg), float64(size.Sm), 40.0
	}

	if b.variant == Link || b.variant == Text {
		px, py = 0, 0 // No padding for links and text
	}
	_ = py // Keep py for layout engine in the future

	// Measure text bounds
	textW, _ := uictx.MeasureText(b.label)

	// Calculate width and height
	w = textW + (px * 2)
	h = minHeight
	if b.variant == Link || b.variant == Text {
		h = 24.0 // smaller height for text/link
	}

	if uictx.MeasureOnly {
		return w, h
	}

	// Process interactions
	_, isActive, isClicked := uictx.ProcessInteraction(b.id, x, y, w, h, b.disabled)
	if isClicked && b.onClick != nil {
		b.onClick()
	}

	state := uictx.GetState(b.id)

	// Determine base and hover colors based on variant
	var bgBase, bgHover, fg color.Color

	switch b.variant {
	case Primary:
		bgBase, bgHover, fg = th.Colors.Primary, mixColor(th.Colors.Primary, color.White, 0.2), th.Colors.PrimaryForeground
	case Secondary:
		bgBase, bgHover, fg = th.Colors.Secondary, mixColor(th.Colors.Secondary, color.Black, 0.1), th.Colors.SecondaryForeground
	case Danger:
		bgBase, bgHover, fg = th.Colors.Danger, mixColor(th.Colors.Danger, color.Black, 0.1), th.Colors.DangerForeground
	case Success:
		bgBase, bgHover, fg = th.Colors.Success, mixColor(th.Colors.Success, color.Black, 0.1), th.Colors.SuccessForeground
	case Warning:
		bgBase, bgHover, fg = th.Colors.Warning, mixColor(th.Colors.Warning, color.Black, 0.1), th.Colors.WarningForeground
	case Info:
		bgBase, bgHover, fg = th.Colors.Info, mixColor(th.Colors.Info, color.Black, 0.1), th.Colors.InfoForeground
	case Ghost, Text:
		bgBase, bgHover, fg = color.Transparent, th.Colors.Muted, th.Colors.Foreground
	case Link:
		bgBase, bgHover, fg = color.Transparent, color.Transparent, th.Colors.Primary
	default:
		bgBase, bgHover, fg = th.Colors.Input, mixColor(th.Colors.Input, color.Black, 0.1), th.Colors.Foreground
	}

	if b.disabled {
		bgBase = th.Colors.Muted
		bgHover = th.Colors.Muted
		fg = th.Colors.MutedForeground
	}

	// Animate background color
	bg := mixColor(bgBase, bgHover, state.HoverRatio)

	// Animate push down
	var yOffset float64 = 0.0
	if isActive && !b.disabled {
		yOffset = state.ActiveRatio * 1.5 // Max 1.5px push down
	}

	// Draw button background
	ctx.SetColor(bg)
	uictx.DrawRoundedRectangle(x, y+yOffset, w, h, float64(th.Radius))
	uictx.Fill()

	if b.variant == Default {
		ctx.SetColor(th.Colors.Border)
		uictx.DrawRoundedRectangle(x, y+yOffset, w, h, float64(th.Radius))
		uictx.Stroke()
	}

	// Draw text
	ctx.SetColor(fg)
	uictx.DrawStringAnchored(b.label, x+w/2, y+h/2+yOffset, 0.5, 0.5)

	// Draw underline for link on hover
	if b.variant == Link && state.HoverRatio > 0 {
		ctx.SetColor(fg)
		lineWidth := 1.0
		ctx.SetLineWidth(lineWidth)
		ctx.DrawLine(x, y+h-4+yOffset, x+w, y+h-4+yOffset)
		uictx.Stroke()
	}

	return w, h
}
