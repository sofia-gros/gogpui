package skeleton

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

// Skeleton は読み込み中の状態を示すプレースホルダー要素。
type Skeleton struct {
	id        string
	secondary bool
	width     float64
	height    float64
}

// New は新しい Skeleton 要素を作成する。
func New(id string) *Skeleton {
	return &Skeleton{
		id:     id,
		width:  200.0, // default placeholder width
		height: 16.0,  // default placeholder height (h_4 in gpui = 1.0rem = 16px)
	}
}

// Secondary はセカンダリカラーの使用を設定する。
func (s *Skeleton) Secondary(secondary bool) *Skeleton {
	s.secondary = secondary
	return s
}

// Width はプレースホルダーの幅を設定する。
func (s *Skeleton) Width(w float64) *Skeleton {
	s.width = w
	return s
}

// Height はプレースホルダーの高さを設定する。
func (s *Skeleton) Height(h float64) *Skeleton {
	s.height = h
	return s
}

// Render は指定位置に Skeleton を描画し、占有サイズ (幅, 高さ) を返す。
func (s *Skeleton) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	ctx := uictx.GG
	th := uictx.Theme

	w := s.width
	if w <= 0 {
		w = 200.0
	}
	h := s.height
	if h <= 0 {
		h = 16.0
	}

	if uictx.MeasureOnly {
		return w, h
	}

	// アニメーション更新 (周期: 2秒)
	uictx.NeedsRedraw = true
	state := uictx.GetState(s.id)
	state.ToggleRatio += uictx.DeltaTime * 0.5 // 2秒で 0.0 -> 1.0 に進行する
	if state.ToggleRatio >= 1.0 {
		state.ToggleRatio = math.Mod(state.ToggleRatio, 1.0)
	}

	// bounce(ease_in_out) の計算
	// 0.0〜0.5で 0.0->1.0 (1秒), 0.5〜1.0で 1.0->0.0 (1秒)
	t := state.ToggleRatio
	var delta float64
	if t < 0.5 {
		delta = t * 2.0
	} else {
		delta = (1.0 - t) * 2.0
	}
	eased := easeInOut(delta)

	// let v = 1.0 - delta * 0.5; this.opacity(v)
	alphaFactor := 1.0 - eased*0.5

	// 背景色の取得とアルファ計算
	// cx.theme().skeleton は th.Colors.Muted に相当
	bgColor := th.Colors.Muted
	pr, pg, pb, pa := bgColor.RGBA()
	
	// ベース不透明度
	baseOpacity := 1.0
	if s.secondary {
		baseOpacity = 0.5
	}
	
	finalAlpha := float64(pa>>8) / 255.0 * baseOpacity * alphaFactor
	if finalAlpha < 0 {
		finalAlpha = 0
	}
	if finalAlpha > 1.0 {
		finalAlpha = 1.0
	}

	// NRGBA を用いて正しく透過ブレンドを描画
	colorNRGBA := color.NRGBA{
		R: uint8(pr >> 8),
		G: uint8(pg >> 8),
		B: uint8(pb >> 8),
		A: uint8(finalAlpha * 255.0),
	}

	radius := float64(th.Radius)

	ctx.SetColor(colorNRGBA)
	uictx.DrawRoundedRectangle(x, y, w, h, radius)
	_ = uictx.Fill()

	return w, h
}
