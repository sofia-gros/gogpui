package statusbar

import (
	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/layout"
)

// StatusBar はウィンドウ下部に配置される水平ステータスバーコンポーネント。
type StatusBar struct {
	width    float64
	left     []layout.Widget
	right    []layout.Widget
	children []layout.Widget
}

// New は新しい StatusBar を作成する。
func New() *StatusBar {
	return &StatusBar{
		width: 200, // デフォルト幅（通常は親やウィンドウサイズから設定される）
	}
}

// Width はステータスバーの横幅を設定する。
func (s *StatusBar) Width(w float64) *StatusBar {
	s.width = w
	return s
}

// Left は左側セクションに要素を追加する。
func (s *StatusBar) Left(child layout.Widget) *StatusBar {
	s.left = append(s.left, child)
	return s
}

// Right は右側セクションに要素を追加する。
func (s *StatusBar) Right(child layout.Widget) *StatusBar {
	s.right = append(s.right, child)
	return s
}

// Child は中央セクションに要素を追加する。
func (s *StatusBar) Child(child layout.Widget) *StatusBar {
	s.children = append(s.children, child)
	return s
}

// growContainer は StatusBar 内部で flex-1 として振る舞い、
// 与えられた幅の中で内部の Flex を配置するためのカスタム Widget です。
type growContainer struct {
	allocW, allocH float64
	inner          *layout.Flex
}

func (g *growContainer) GrowFactor() float64 { return 1.0 }

func (g *growContainer) SetAllocatedSize(w, h float64) {
	g.allocW = w
	g.allocH = h
}

func (g *growContainer) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	if uictx.MeasureOnly {
		return g.allocW, g.allocH
	}
	g.inner.WithConstraints(g.allocW, g.allocH)
	g.inner.Render(uictx, x, y)
	return g.allocW, g.allocH
}

// Render は指定位置に StatusBar を描画し、占有サイズ (幅, 高さ) を返す。
func (s *StatusBar) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	th := uictx.Theme
	ctx := uictx.GG

	// px_2 (8px), py_1 (4px), 内部の高さ目安を 16px とする -> 全体高さ 24px
	barH := 24.0
	px := 8.0
	py := 4.0

	if uictx.MeasureOnly {
		return s.width, barH
	}

	// 背景描画
	ctx.SetColor(th.Colors.Background) // Status bar BG (通常は少し暗いかBackgroundと同色)
	ctx.DrawRectangle(x, y, s.width, barH)
	_ = ctx.Fill()

	// 上部ボーダー描画 (border_t_1)
	ctx.SetColor(th.Colors.Border)
	ctx.DrawRectangle(x, y, s.width, 1.0)
	_ = ctx.Fill()

	// メインのコンテナ (横幅から左右のパディングを引いた領域)
	mainW := s.width - (px * 2)
	mainX := x + px
	mainY := y + py

	mainFlex := layout.NewFlex().Direction(layout.Row).Align(layout.AlignCenter).WithConstraints(mainW, barH-py*2)

	hasLeft := len(s.left) > 0
	hasRight := len(s.right) > 0

	// 1. Left Region
	if hasLeft {
		leftFlex := layout.NewFlex().Direction(layout.Row).Align(layout.AlignCenter).Gap(8)
		for _, child := range s.left {
			leftFlex.Add(child)
		}
		mainFlex.Add(leftFlex)
	}

	// 2. Center Region (flex-1)
	centerFlex := layout.NewFlex().Direction(layout.Row).Align(layout.AlignCenter).Gap(8)
	if hasLeft && hasRight {
		centerFlex.Justify(layout.JustifyCenter)
	} else if hasLeft && !hasRight {
		centerFlex.Justify(layout.JustifyEnd)
	} else {
		centerFlex.Justify(layout.JustifyStart)
	}
	for _, child := range s.children {
		centerFlex.Add(child)
	}
	mainFlex.Add(&growContainer{inner: centerFlex})

	// 3. Right Region
	if hasRight {
		rightFlex := layout.NewFlex().Direction(layout.Row).Align(layout.AlignCenter).Gap(8)
		for _, child := range s.right {
			rightFlex.Add(child)
		}
		mainFlex.Add(rightFlex)
	}

	mainFlex.Render(uictx, mainX, mainY)

	return s.width, barH
}
