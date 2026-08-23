package progress

import (
	"image/color"
	"math"

	"github.com/sofiagros/gogpui/core/context"
)

// easeInOut は Rust GPUI の ease_in_out と一致する 2次イージング関数。
func easeInOut(t float64) float64 {
	if t < 0.5 {
		return 2.0 * t * t
	}
	return -1.0 + (4.0-2.0*t)*t
}

// Size はプログレスバーおよびサークルのサイズバリアントを表す。
type Size int

const (
	// SizeXSmall は極小サイズ (バー高さ: 4px, サークル: 8px)。
	SizeXSmall Size = iota
	// SizeSmall は小サイズ (バー高さ: 6px, サークル: 12px)。
	SizeSmall
	// SizeMedium は中サイズ (バー高さ: 8px, サークル: 16px)。
	SizeMedium
	// SizeLarge は大サイズ (バー高さ: 10px, サークル: 20px)。
	SizeLarge
)

// Progress は線形の水平プログレスバーコンポーネント。
type Progress struct {
	id      string
	value   float64
	size    Size
	color   *color.Color
	width   float64
	loading bool
}

// New は新しい水平プログレスバーを作成する。
func New(id string) *Progress {
	return &Progress{
		id:    id,
		size:  SizeMedium,
		width: 200.0,
	}
}

// Value は進捗率 (0.0 〜 100.0) を設定する。
func (p *Progress) Value(val float64) *Progress {
	if val < 0 {
		val = 0
	} else if val > 100 {
		val = 100
	}
	p.value = val
	return p
}

// Size はプログレスバーのサイズを設定する。
func (p *Progress) Size(s Size) *Progress {
	p.size = s
	return p
}

// Color はプログレスバーのカラーを設定する。
func (p *Progress) Color(c color.Color) *Progress {
	p.color = &c
	return p
}

// Width はプログレスバーの幅を設定する。
func (p *Progress) Width(w float64) *Progress {
	p.width = w
	return p
}

// Loading は不確定 (ローディング) アニメーションを有効にする。
func (p *Progress) Loading(loading bool) *Progress {
	p.loading = loading
	return p
}

// Render は指定位置に水平プログレスバーを描画し、占有サイズ (幅, 高さ) を返す。
func (p *Progress) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	ctx := uictx.GG
	th := uictx.Theme

	var barH, radius float64
	switch p.size {
	case SizeXSmall:
		barH, radius = 4.0, 2.0
	case SizeSmall:
		barH, radius = 6.0, 3.0
	case SizeLarge:
		barH, radius = 10.0, 5.0
	case SizeMedium:
		fallthrough
	default:
		barH, radius = 8.0, 4.0
	}

	w := p.width
	if w <= 0 {
		w = 200.0
	}

	if uictx.MeasureOnly {
		return w, barH
	}

	primaryColor := th.Colors.Primary
	if p.color != nil {
		primaryColor = *p.color
	}

	pr, pg, pb, _ := primaryColor.RGBA()
	// 注意: 不透明度0.2の背景を描画するためには、Premultiplied alphaのRGBAではなく、NRGBAを使用する
	trackColor := color.NRGBA{
		R: uint8(pr >> 8),
		G: uint8(pg >> 8),
		B: uint8(pb >> 8),
		A: uint8(255 / 5),
	}
	ctx.SetColor(trackColor)
	ctx.DrawRoundedRectangle(x, y, w, barH, radius)
	_ = ctx.Fill()

	// インジケーター描画
	ctx.SetColor(primaryColor)
	if p.loading {
		uictx.NeedsRedraw = true
		state := uictx.GetState(p.id)
		state.ToggleRatio += uictx.DeltaTime * 1.0 // 1秒周期
		if state.ToggleRatio >= 1.0 {
			state.ToggleRatio = math.Mod(state.ToggleRatio, 1.0)
		}
		delta := state.ToggleRatio

		start := 0.0
		if delta > 0.5 {
			start = easeInOut((delta - 0.5) / 0.5)
		}
		end := 1.0 - easeInOut(1.0-delta)

		startX := x + start*w
		endX := x + end*w
		indicatorW := endX - startX
		if indicatorW > 0 {
			ctx.Push()
			ctx.DrawRoundedRectangle(x, y, w, barH, radius)
			ctx.Clip()

			// loading中は常に両端角丸 (rounded_r_noneしない)
			ctx.DrawRoundedRectangle(startX, y, indicatorW, barH, radius)
			_ = ctx.Fill()
			ctx.Pop()
		}
	} else {
		state := uictx.GetState(p.id)
		state.ValueRatio = uictx.Animate(state.ValueRatio, p.value, 150.0)
		indicatorW := (state.ValueRatio / 100.0) * w
		if indicatorW > 0 {
			if indicatorW > w {
				indicatorW = w
			}
			ctx.Push()
			ctx.DrawRoundedRectangle(x, y, w, barH, radius)
			ctx.Clip()

			if state.ValueRatio >= 100.0 {
				// 100%時は両端角丸
				ctx.DrawRoundedRectangle(x, y, indicatorW, barH, radius)
			} else {
				// 100%未満の場合は右端を角丸にしない (rounded_r_none)
				ctx.DrawRectangle(x, y, indicatorW, barH)
			}
			_ = ctx.Fill()
			ctx.Pop()
		}
	}

	return w, barH
}

// ProgressCircle は円形の進捗インジケーターコンポーネント。
type ProgressCircle struct {
	id          string
	value       float64
	size        Size
	radius      float64
	strokeWidth float64
	color       *color.Color
	loading     bool
}

// NewCircle は新しい円形プログレスインジケーターを作成する。
func NewCircle(id string) *ProgressCircle {
	return &ProgressCircle{
		id:   id,
		size: SizeMedium,
	}
}

// Value は進捗率 (0.0 〜 100.0) を設定する。
func (c *ProgressCircle) Value(val float64) *ProgressCircle {
	if val < 0 {
		val = 0
	} else if val > 100 {
		val = 100
	}
	c.value = val
	return c
}

// Size はサークルのサイズバリアントを設定する。
func (c *ProgressCircle) Size(s Size) *ProgressCircle {
	c.size = s
	return c
}

// Radius はカスタム半径を設定する。
func (c *ProgressCircle) Radius(r float64) *ProgressCircle {
	c.radius = r
	return c
}

// StrokeWidth は線の太さを設定する。
func (c *ProgressCircle) StrokeWidth(w float64) *ProgressCircle {
	c.strokeWidth = w
	return c
}

// Color は描画カラーを設定する。
func (c *ProgressCircle) Color(col color.Color) *ProgressCircle {
	c.color = &col
	return c
}

// Loading は回転ローディングアニメーションを有効にする。
func (c *ProgressCircle) Loading(loading bool) *ProgressCircle {
	c.loading = loading
	return c
}

// Render は指定位置に円形プログレスを描画し、占有サイズ (幅, 高さ) を返す。
func (c *ProgressCircle) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	ctx := uictx.GG
	th := uictx.Theme

	var sizePx float64
	switch c.size {
	case SizeXSmall:
		sizePx = 8.0
	case SizeSmall:
		sizePx = 12.0
	case SizeLarge:
		sizePx = 20.0
	case SizeMedium:
		fallthrough
	default:
		sizePx = 16.0
	}
	if c.radius > 0 {
		sizePx = c.radius * 2.0
	}

	stroke := sizePx * 0.15
	if stroke > 5.0 {
		stroke = 5.0
	}
	if stroke < 1.5 {
		stroke = 1.5
	}
	if c.strokeWidth > 0 {
		stroke = c.strokeWidth
	}

	actualRadius := (sizePx - stroke) / 2.0
	cx := x + sizePx/2.0
	cy := y + sizePx/2.0

	if uictx.MeasureOnly {
		return sizePx, sizePx
	}

	primaryColor := th.Colors.Primary
	if c.color != nil {
		primaryColor = *c.color
	}

	pr, pg, pb, _ := primaryColor.RGBA()
	trackColor := color.NRGBA{
		R: uint8(pr >> 8),
		G: uint8(pg >> 8),
		B: uint8(pb >> 8),
		A: uint8(255 / 5),
	}
	
	// DrawArc を用いて円を描画する
	ctx.SetColor(trackColor)
	ctx.SetLineWidth(stroke)
	ctx.DrawArc(cx, cy, actualRadius, 0, math.Pi*2)
	_ = ctx.Stroke()

	// インジケーター円弧描画
	ctx.SetColor(primaryColor)
	ctx.SetLineWidth(stroke)

	if c.loading {
		uictx.NeedsRedraw = true
		state := uictx.GetState(c.id)
		state.ToggleRatio += uictx.DeltaTime * 1.0 // 1秒周期
		if state.ToggleRatio >= 1.0 {
			state.ToggleRatio = math.Mod(state.ToggleRatio, 1.0)
		}
		delta := state.ToggleRatio

		end := easeInOut(delta) * 100.0
		start := 0.0
		if delta > 0.5 {
			start = easeInOut((delta - 0.5) / 0.5) * 100.0
		}

		startAngle := -math.Pi/2.0 + (start/100.0)*math.Pi*2.0
		endAngle := -math.Pi/2.0 + (end/100.0)*math.Pi*2.0
		if endAngle > startAngle {
			ctx.DrawArc(cx, cy, actualRadius, startAngle, endAngle)
			_ = ctx.Stroke()
		}
	} else {
		state := uictx.GetState(c.id)
		state.ValueRatio = uictx.Animate(state.ValueRatio, c.value, 150.0)
		if state.ValueRatio > 0 {
			startAngle := -math.Pi / 2.0
			endAngle := startAngle + (state.ValueRatio/100.0)*math.Pi*2.0
			ctx.DrawArc(cx, cy, actualRadius, startAngle, endAngle)
			_ = ctx.Stroke()
		}
	}

	return sizePx, sizePx
}

