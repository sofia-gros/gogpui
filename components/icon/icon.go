package icon

import (
	"image/color"

	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/layout"
)

// Name はアイコンの種類を表します。
type Name string

const (
	FolderOpen   Name = "folder-open"
	FolderClosed Name = "folder-closed"
	File         Name = "file"
)

// IconTheme はアイコンを実際に描画するロジックを提供するインターフェースです。
type IconTheme interface {
	// Render は指定されたアイコンを(cx, cy)位置に描画し、使用した幅と高さを返します。
	Render(name Name, uictx *context.UIContext, cx, cy, size float64, c color.Color) (float64, float64)
}

// DefaultTheme は組み込みの基本的な描画ロジックを用いたデフォルトテーマです。
type DefaultTheme struct{}

func (t DefaultTheme) Render(name Name, uictx *context.UIContext, cx, cy, size float64, c color.Color) (float64, float64) {
	if uictx.MeasureOnly {
		return size, size
	}

	ggCtx := uictx.GG
	ggCtx.SetColor(c)
	ggCtx.SetLineWidth(1.5)

	switch name {
	case File:
		ggCtx.DrawRectangle(cx+3, cy+2, 10, 12)
		ggCtx.Stroke()
		ggCtx.MoveTo(cx+5, cy+5)
		ggCtx.LineTo(cx+11, cy+5)
		ggCtx.MoveTo(cx+5, cy+8)
		ggCtx.LineTo(cx+11, cy+8)
		ggCtx.Stroke()
	case FolderOpen:
		ggCtx.MoveTo(cx+2, cy+3)
		ggCtx.LineTo(cx+7, cy+3)
		ggCtx.LineTo(cx+8, cy+5)
		ggCtx.LineTo(cx+14, cy+5)
		ggCtx.LineTo(cx+14, cy+13)
		ggCtx.LineTo(cx+2, cy+13)
		ggCtx.ClosePath()
		ggCtx.Stroke()
		// flap
		ggCtx.MoveTo(cx+2, cy+13)
		ggCtx.LineTo(cx+5, cy+7)
		ggCtx.LineTo(cx+16, cy+7)
		ggCtx.LineTo(cx+13, cy+13)
		ggCtx.ClosePath()
		ggCtx.Stroke()
	case FolderClosed:
		ggCtx.DrawRectangle(cx+2, cy+5, 12, 8)
		ggCtx.MoveTo(cx+2, cy+5)
		ggCtx.LineTo(cx+2, cy+3)
		ggCtx.LineTo(cx+6, cy+3)
		ggCtx.LineTo(cx+8, cy+5)
		ggCtx.Stroke()
	}

	return size, size
}

// Icon はアイコンを描画するためのコンポーネントです。
type Icon struct {
	name  Name
	size  float64
	color color.Color
	theme IconTheme
}

// New は新しいIconコンポーネントを作成します。デフォルトのサイズは16です。
func New(name Name) *Icon {
	return &Icon{
		name:  name,
		size:  16,
		theme: DefaultTheme{},
	}
}

// Size はアイコンのサイズを設定します。
func (i *Icon) Size(size float64) *Icon {
	i.size = size
	return i
}

// Color はアイコンの色を設定します。指定しない場合はテーマのForegroundなどが使われます。
func (i *Icon) Color(c color.Color) *Icon {
	i.color = c
	return i
}

// Theme は描画に使用するIconThemeを設定します。
func (i *Icon) Theme(theme IconTheme) *Icon {
	i.theme = theme
	return i
}

// Render は layout.Widget インターフェースを実装し、アイコンを描画します。
func (i *Icon) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	c := i.color
	if c == nil {
		c = uictx.Theme.Colors.Foreground
	}
	return i.theme.Render(i.name, uictx, x, y, i.size, c)
}

// layout.Widget を満たしているかチェック
var _ layout.Widget = (*Icon)(nil)
