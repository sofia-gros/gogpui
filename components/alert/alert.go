package alert

import (
	"image/color"

	"github.com/sofiagros/gogpui/components/button"
	"github.com/sofiagros/gogpui/components/label"
	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/layout"
)

// AlertVariant は Alert の種類を表す。
type AlertVariant int

const (
	Default AlertVariant = iota
	Info
	Success
	Warning
	Error
)

// Alert はユーザーにメッセージや警告を表示するコンポーネント。
type Alert struct {
	id       string
	variant  AlertVariant
	title    string
	message  string
	onClose  func()
	visible  bool
	width    float64
}

// New は新しい Alert を作成する。
func New(id string, message string) *Alert {
	return &Alert{
		id:      id,
		variant: Default,
		message: message,
		visible: true,
		width:   300,
	}
}

// NewInfo は Info バリアントの Alert を作成する。
func NewInfo(id string, message string) *Alert {
	return New(id, message).Variant(Info)
}

// NewSuccess は Success バリアントの Alert を作成する。
func NewSuccess(id string, message string) *Alert {
	return New(id, message).Variant(Success)
}

// NewWarning は Warning バリアントの Alert を作成する。
func NewWarning(id string, message string) *Alert {
	return New(id, message).Variant(Warning)
}

// NewError は Error バリアントの Alert を作成する。
func NewError(id string, message string) *Alert {
	return New(id, message).Variant(Error)
}

// Variant は Alert の種類を設定する。
func (a *Alert) Variant(v AlertVariant) *Alert {
	a.variant = v
	return a
}

// Title は Alert のタイトルを設定する。
func (a *Alert) Title(t string) *Alert {
	a.title = t
	return a
}

// Width は Alert の幅を設定する。
func (a *Alert) Width(w float64) *Alert {
	a.width = w
	return a
}

// OnClose は閉じるボタンがクリックされた時のコールバックを設定する。
func (a *Alert) OnClose(f func()) *Alert {
	a.onClose = f
	return a
}

func paleColor(c color.Color, alpha uint8) color.Color {
	r, g, b, _ := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: alpha}
}

// Render は指定位置に Alert を描画し、占有サイズ (幅, 高さ) を返す。
func (a *Alert) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	if !a.visible {
		return 0, 0
	}

	th := uictx.Theme
	ctx := uictx.GG

	var fgColor color.Color
	switch a.variant {
	case Info:
		fgColor = th.Colors.Info
	case Success:
		fgColor = th.Colors.Success
	case Warning:
		fgColor = th.Colors.Warning
	case Error:
		fgColor = th.Colors.Danger
	default:
		fgColor = th.Colors.Foreground
	}

	bgColor := paleColor(fgColor, 25)  // ~10% opacity
	borderColor := paleColor(fgColor, 76) // ~30% opacity

	// 内側のコンテンツをレイアウト計算するための一時コンテナ
	contentFlex := layout.NewFlex().Direction(layout.Column).Gap(4)
	if a.title != "" {
		contentFlex.Add(label.New(a.title).Color(fgColor))
	}
	contentFlex.Add(label.New(a.message).Color(th.Colors.MutedForeground))

	headerFlex := layout.NewFlex().Direction(layout.Row).Gap(8).Align(layout.AlignStart)
	
	// アイコン代わりの四角形など（文字で代用）
	var iconText string
	switch a.variant {
	case Info:
		iconText = "ℹ"
	case Success:
		iconText = "✔"
	case Warning:
		iconText = "⚠"
	case Error:
		iconText = "✖"
	default:
		iconText = "ℹ"
	}
	headerFlex.Add(label.New(iconText).Color(fgColor))
	headerFlex.Add(contentFlex)

	// 閉じるボタン
	mainFlex := layout.NewFlex().Direction(layout.Row).Align(layout.AlignStart).Gap(8).WithConstraints(a.width-16, 0)
	mainFlex.Add(headerFlex)

	if a.onClose != nil {
		closeBtn := button.New(a.id + "-close").
			Ghost().
			Label("×").
			OnClick(func() {
				if a.onClose != nil {
					a.onClose()
				}
			})
		// 右側に寄せるためにダミースペーサー等は使わず、シンプルに追加
		mainFlex.Add(closeBtn)
	}

	// 高さを事前計測
	originalMeasure := uictx.MeasureOnly
	uictx.MeasureOnly = true
	_, contentH := mainFlex.Render(uictx, x, y)
	uictx.MeasureOnly = originalMeasure

	pad := 12.0
	totalH := contentH + pad*2

	if uictx.MeasureOnly {
		return a.width, totalH
	}

	// 背景描画
	ctx.SetColor(bgColor)
	ctx.DrawRoundedRectangle(x, y, a.width, totalH, 6)
	_ = ctx.Fill()

	// 枠線描画
	ctx.SetColor(borderColor)
	ctx.SetLineWidth(1)
	ctx.DrawRoundedRectangle(x, y, a.width, totalH, 6)
	_ = ctx.Stroke()

	// コンテンツ描画
	mainFlex.Render(uictx, x+pad, y+pad)

	return a.width, totalH
}
