package separator

import (
	"image/color"

	"github.com/sofiagros/gogpui/core/context"
)

// Style はセパレーターの線の描画スタイルを表す。
type Style int

const (
	// StyleSolid は実線スタイル。
	StyleSolid Style = iota
	// StyleDashed は破線スタイル。
	StyleDashed
)

// Orientation はセパレーターの方向を表す。
type Orientation int

const (
	// OrientationHorizontal は水平方向のセパレーター。
	OrientationHorizontal Orientation = iota
	// OrientationVertical は垂直方向のセパレーター。
	OrientationVertical
)

// Separator はコンテンツを視覚的に分割する区切り線コンポーネント。
type Separator struct {
	orientation Orientation
	style       Style
	label       string
	color       *color.Color
	length      float64
}

// New はデフォルトの水平実線セパレーターを作成する。
func New() *Separator {
	return &Separator{
		orientation: OrientationHorizontal,
		style:       StyleSolid,
	}
}

// Horizontal は水平セパレーターを作成する。
func Horizontal() *Separator {
	return New()
}

// Vertical は垂直セパレーターを作成する。
func Vertical() *Separator {
	return &Separator{
		orientation: OrientationVertical,
		style:       StyleSolid,
	}
}

// HorizontalDashed は水平破線セパレーターを作成する。
func HorizontalDashed() *Separator {
	return Horizontal().Dashed()
}

// VerticalDashed は垂直破線セパレーターを作成する。
func VerticalDashed() *Separator {
	return Vertical().Dashed()
}

// Label はセパレーターの中央に表示するテキストラベルを設定する。
func (s *Separator) Label(label string) *Separator {
	s.label = label
	return s
}

// Color はセパレーターの描画色を設定する。
func (s *Separator) Color(c color.Color) *Separator {
	s.color = &c
	return s
}

// Dashed はセパレーターを破線スタイルに設定する。
func (s *Separator) Dashed() *Separator {
	s.style = StyleDashed
	return s
}

// Solid はセパレーターを実線スタイルに設定する。
func (s *Separator) Solid() *Separator {
	s.style = StyleSolid
	return s
}

// Style はセパレーターの描画スタイル（実線/破線）を設定する。
func (s *Separator) Style(style Style) *Separator {
	s.style = style
	return s
}

// Length はセパレーターの長さ（水平時は幅、垂直時は高さ）を設定する。
func (s *Separator) Length(length float64) *Separator {
	s.length = length
	return s
}

// Render は指定位置にセパレーターを描画し、占有サイズ (幅, 高さ) を返す。
func (s *Separator) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	ctx := uictx.GG
	th := uictx.Theme

	length := s.length
	if length <= 0 {
		if s.orientation == OrientationHorizontal {
			length = 200.0 // デフォルトの水平長
		} else {
			length = 40.0 // デフォルトの垂直長
		}
	}

	if s.orientation == OrientationHorizontal {
		var h float64 = 1.0
		var textW, textH float64
		var totalLabelW float64
		paddingX := 8.0

		if s.label != "" {
			textW, textH = ctx.MeasureString(s.label)
			totalLabelW = textW + paddingX*2
			h = textH
			if h < 16 {
				h = 16
			}
		}

		if uictx.MeasureOnly {
			return length, h
		}

		lineColor := th.Colors.Border
		if s.color != nil {
			lineColor = *s.color
		}

		if s.label != "" {
			midY := y + h/2.0

			// 左側の線
			leftLineW := (length - totalLabelW) / 2.0
			if leftLineW > 0 {
				ctx.SetColor(lineColor)
				ctx.SetLineWidth(1.0)
				if s.style == StyleDashed {
					ctx.SetDash(4, 2)
				}
				ctx.DrawLine(x, midY, x+leftLineW, midY)
				_ = ctx.Stroke()
				if s.style == StyleDashed {
					ctx.ClearDash()
				}
			}

			// 中央ラベル
			labelX := x + leftLineW + paddingX
			ctx.SetColor(th.Colors.MutedForeground)
			ctx.DrawString(s.label, labelX, y+(h+textH)/2.0-2.0)

			// 右側の線
			rightLineStartX := x + leftLineW + totalLabelW
			if rightLineStartX < x+length {
				ctx.SetColor(lineColor)
				ctx.SetLineWidth(1.0)
				if s.style == StyleDashed {
					ctx.SetDash(4, 2)
				}
				ctx.DrawLine(rightLineStartX, midY, x+length, midY)
				_ = ctx.Stroke()
				if s.style == StyleDashed {
					ctx.ClearDash()
				}
			}

			return length, h
		}

		// ラベルなし水平セパレーター
		ctx.SetColor(lineColor)
		ctx.SetLineWidth(1.0)
		if s.style == StyleDashed {
			ctx.SetDash(4, 2)
		}
		ctx.DrawLine(x, y+0.5, x+length, y+0.5)
		_ = ctx.Stroke()
		if s.style == StyleDashed {
			ctx.ClearDash()
		}
		return length, 1.0
	}

	if uictx.MeasureOnly {
		return 1.0, length
	}

	lineColor := th.Colors.Border
	if s.color != nil {
		lineColor = *s.color
	}

	// 垂直セパレーター
	ctx.SetColor(lineColor)
	ctx.SetLineWidth(1.0)
	if s.style == StyleDashed {
		ctx.SetDash(4, 2)
	}
	ctx.DrawLine(x+0.5, y, x+0.5, y+length)
	_ = ctx.Stroke()
	if s.style == StyleDashed {
		ctx.ClearDash()
	}
	return 1.0, length
}
