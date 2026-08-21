package layout

import "github.com/sofiagros/gogpui/core/context"

// Padding は単一の子 Widget を上下左右のインセットで包む Widget。
// GPUI の div().p(px(N)) / div().px(px(N)).py(px(N)) 等に相当する。
//
// 使用例:
//
//	layout.NewPadding(btn).All(16).Render(uictx, x, y)
//	layout.NewPadding(btn).Horizontal(8).Vertical(4).Render(uictx, x, y)
//	layout.NewPadding(btn).Top(2).Bottom(2).Left(8).Right(8).Render(uictx, x, y)
type Padding struct {
	child          Widget
	top, bottom    float64
	left, right    float64
}

// NewPadding は子 Widget を受け取り、Padding を返す。
func NewPadding(child Widget) *Padding {
	return &Padding{child: child}
}

// All はすべての方向に同じパディングを設定する。
func (p *Padding) All(v float64) *Padding {
	p.top, p.bottom, p.left, p.right = v, v, v, v
	return p
}

// Horizontal は左右のパディングを設定する。
func (p *Padding) Horizontal(v float64) *Padding {
	p.left, p.right = v, v
	return p
}

// Vertical は上下のパディングを設定する。
func (p *Padding) Vertical(v float64) *Padding {
	p.top, p.bottom = v, v
	return p
}

// Top は上パディングを設定する。
func (p *Padding) Top(v float64) *Padding {
	p.top = v
	return p
}

// Bottom は下パディングを設定する。
func (p *Padding) Bottom(v float64) *Padding {
	p.bottom = v
	return p
}

// Left は左パディングを設定する。
func (p *Padding) Left(v float64) *Padding {
	p.left = v
	return p
}

// Right は右パディングを設定する。
func (p *Padding) Right(v float64) *Padding {
	p.right = v
	return p
}

// Render は Widget インターフェースを実装する。
// 子 Widget を (x + left, y + top) に配置し、
// このウィジェット自身のサイズとして (child.w + left + right, child.h + top + bottom) を返す。
func (p *Padding) Render(uictx *context.UIContext, x, y float64) (w, h float64) {
	cw, ch := p.child.Render(uictx, x+p.left, y+p.top)
	return cw + p.left + p.right, ch + p.top + p.bottom
}
