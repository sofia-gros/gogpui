package appbar

import (
	"image"
	"image/color"

	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/layout"
	"github.com/sofiagros/gogpui/core/theme"
)

// AppBar はフレームレスウィンドウのためのトップレベルコンテナです。
// 背景色を描画し、高さを固定します。
type AppBar struct {
	content layout.Widget
	height  float64
}

// New は新しいAppBarを作成します。
func New() *AppBar {
	return &AppBar{
		height: 32.0,
	}
}

// Content はAppBar内に表示するウィジェット（タイトルやFlexなど）を設定します。
func (a *AppBar) Content(w layout.Widget) *AppBar {
	a.content = w
	return a
}

// Height はAppBarの固定高さを設定します。
func (a *AppBar) Height(h float64) *AppBar {
	a.height = h
	return a
}

// Render はAppBarを描画し、必要に応じてOSのドラッグ領域を登録します。
func (a *AppBar) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	if !uictx.MeasureOnly {
		// 背景の描画 (ウィンドウタイトルバーと同等の色)
		uictx.GG.SetColor(uictx.Theme.Colors.Background) // または少し暗い色
		uictx.GG.DrawRectangle(x, y, uictx.WindowWidth, a.height)
		uictx.GG.Fill()
	}

	if a.content != nil {
		if gw, ok := a.content.(layout.GrowWidget); ok {
			gw.SetAllocatedSize(uictx.WindowWidth, a.height)
		}
		a.content.Render(uictx, x, y)
	}

	return uictx.WindowWidth, a.height
}

// DragZone はラップしたウィジェットの領域をOSのドラッグ可能領域として登録します。
// これにより、このウィジェット上のドラッグ操作でウィンドウ全体を移動できるようになります。
type DragZone struct {
	content layout.Widget
}

// NewDragZone はドラッグ可能領域ウィジェットを作成します。
func NewDragZone(content layout.Widget) *DragZone {
	return &DragZone{content: content}
}

// GrowFactor は内部のWidgetがGrowWidgetであればその値を継承します。
func (d *DragZone) GrowFactor() float64 {
	if gw, ok := d.content.(layout.GrowWidget); ok {
		return gw.GrowFactor()
	}
	return 0
}

// SetAllocatedSize は内部のWidgetに割り当てサイズを伝播します。
func (d *DragZone) SetAllocatedSize(w, h float64) {
	if gw, ok := d.content.(layout.GrowWidget); ok {
		gw.SetAllocatedSize(w, h)
	}
}

// Render は内部を描画し、その領域を WindowControl の DragRegion として登録します。
func (d *DragZone) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	cw, ch := d.content.Render(uictx, x, y)
	if !uictx.MeasureOnly {
		rect := image.Rect(int(x), int(y), int(x+cw), int(y+ch))
		if uictx.Window != nil {
			uictx.Window.AddDragRegion(rect)
		}
	}
	return cw, ch
}

// WindowControls は最小化、最大化、閉じるの3つのボタンを持つ Flex を返します。
type WindowControls struct {
	flex *layout.Flex
}

// NewWindowControls はウィンドウ操作ボタンのセットを作成します。
// uictx.Window を操作するため、クリック時のコールバックで uictx を利用する専用ウィジェットを内部で構築します。
func NewWindowControls() *WindowControls {
	return &WindowControls{}
}

type windowControlButton struct {
	action string
	label  string
}

func (b *windowControlButton) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	w, h := 46.0, 32.0 // Windows standard-ish size

	if uictx.MeasureOnly {
		return w, h
	}

	_, _, isClicked := uictx.ProcessInteraction("appbar-btn-"+b.action, x, y, w, h, false)

	if isClicked {
		if uictx.Window != nil {
			switch b.action {
			case "min":
				uictx.Window.Minimize()
			case "max":
				uictx.Window.Maximize()
			case "close":
				uictx.Window.Close()
			}
		}
	}

	state := uictx.GetState("appbar-btn-" + b.action)
	th := uictx.Theme

	// Draw Background
	if state.HoverRatio > 0 {
		var bg color.Color
		if b.action == "close" {
			bg = color.RGBA{R: 232, G: 17, B: 35, A: uint8(255 * state.HoverRatio)}
		} else {
			if th.Mode == theme.DarkMode {
				bg = color.RGBA{R: 255, G: 255, B: 255, A: uint8(25 * state.HoverRatio)}
			} else {
				bg = color.RGBA{R: 0, G: 0, B: 0, A: uint8(15 * state.HoverRatio)}
			}
		}
		uictx.GG.SetColor(bg)
		uictx.GG.DrawRectangle(x, y, w, h)
		uictx.GG.Fill()
	}

	// Draw Icon
	fg := th.Colors.Foreground
	if b.action == "close" && state.HoverRatio > 0.5 {
		fg = color.White
	}
	uictx.GG.SetColor(fg)
	uictx.GG.SetLineWidth(1.2)
	
	cx, cy := x+w/2, y+h/2
	switch b.action {
	case "min":
		uictx.GG.DrawLine(cx-5, cy, cx+5, cy)
		uictx.GG.Stroke()
	case "max":
		uictx.GG.DrawRectangle(cx-4, cy-4, 8, 8)
		uictx.GG.Stroke()
	case "close":
		uictx.GG.DrawLine(cx-4, cy-4, cx+4, cy+4)
		uictx.GG.DrawLine(cx-4, cy+4, cx+4, cy-4)
		uictx.GG.Stroke()
	}

	return w, h
}

func (wc *WindowControls) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	if wc.flex == nil {
		wc.flex = layout.NewFlex().Direction(layout.Row).Gap(0).Add(
			&windowControlButton{action: "min", label: "-"},
			&windowControlButton{action: "max", label: "□"},
			&windowControlButton{action: "close", label: "x"},
		)
	}
	return wc.flex.Render(uictx, x, y)
}
