package layout

import (
	"testing"

	"github.com/sofiagros/gogpui/core/context"
)

// TestPadding_All は All(16) で全方向にパディングが入ることを確認する。
func TestPadding_All(t *testing.T) {
	uictx := context.NewUIContext()
	child := &mockWidget{w: 100, h: 50}

	p := NewPadding(child).All(16)
	totalW, totalH := p.Render(uictx, 0, 0)

	// 全体サイズ: 100+16+16=132, 50+16+16=82
	if totalW != 132 {
		t.Errorf("Padding All: expected totalW=132, got %v", totalW)
	}
	if totalH != 82 {
		t.Errorf("Padding All: expected totalH=82, got %v", totalH)
	}
	// 子は (16, 16) にオフセットされる
	if child.renderX != 16 || child.renderY != 16 {
		t.Errorf("Padding All: child expected at (16,16), got (%v,%v)", child.renderX, child.renderY)
	}
}

// TestPadding_IndividualSides は各辺個別指定が正しく子 Widget に反映されることを確認する。
func TestPadding_IndividualSides(t *testing.T) {
	uictx := context.NewUIContext()
	child := &mockWidget{w: 100, h: 50}

	// top=2, bottom=8, left=4, right=12
	p := NewPadding(child).Top(2).Bottom(8).Left(4).Right(12)
	p.Render(uictx, 10, 20)

	// 子の座標: x=10+4=14, y=20+2=22
	if child.renderX != 14 {
		t.Errorf("IndividualSides: child.x expected 14, got %v", child.renderX)
	}
	if child.renderY != 22 {
		t.Errorf("IndividualSides: child.y expected 22, got %v", child.renderY)
	}
}

// TestPadding_TotalSize は Padding のサイズが child + padding の合計になることを確認する。
func TestPadding_TotalSize(t *testing.T) {
	uictx := context.NewUIContext()
	child := &mockWidget{w: 60, h: 30}

	// horizontal=8, vertical=4
	p := NewPadding(child).Horizontal(8).Vertical(4)
	totalW, totalH := p.Render(uictx, 0, 0)

	// w: 60+8+8=76, h: 30+4+4=38
	if totalW != 76 {
		t.Errorf("TotalSize: expected totalW=76, got %v", totalW)
	}
	if totalH != 38 {
		t.Errorf("TotalSize: expected totalH=38, got %v", totalH)
	}
}

// TestPadding_MeasureOnly は MeasureOnly=true のとき子座標を変更しないことを確認する。
func TestPadding_MeasureOnly(t *testing.T) {
	uictx := context.NewUIContext()
	uictx.MeasureOnly = true

	child := &mockWidget{w: 100, h: 50}
	p := NewPadding(child).All(10)

	totalW, totalH := p.Render(uictx, 5, 5)

	// サイズは正しく返る
	if totalW != 120 || totalH != 70 {
		t.Errorf("MeasureOnly: expected (120,70), got (%v,%v)", totalW, totalH)
	}
	// MeasureOnly=true のとき mockWidget は座標を記録しない (renderX/Y = 0 のまま)
	// ※ mockWidget は MeasureOnly=true のとき記録しない実装になっている
	if child.renderX != 0 || child.renderY != 0 {
		t.Errorf("MeasureOnly: child coords should not be set, got (%v,%v)", child.renderX, child.renderY)
	}
}

// TestPadding_InFlex は Flex + Padding の組み合わせで正しく配置されることを確認する。
func TestPadding_InFlex(t *testing.T) {
	uictx := context.NewUIContext()

	child1 := &mockWidget{w: 80, h: 40}
	child2 := &mockWidget{w: 80, h: 40}

	// child1 を Padding(left=10, right=10) で包む → 実効幅 = 80+10+10 = 100
	padded := NewPadding(child1).Left(10).Right(10)

	flex := NewFlex().Direction(Row).Gap(0).Add(padded, child2)
	flex.Render(uictx, 0, 0)

	// child1 は Padding の中で x=10 にオフセット
	if child1.renderX != 10 {
		t.Errorf("InFlex: child1.x expected 10, got %v", child1.renderX)
	}
	// child2 は Padding の後に配置される: x=100
	if child2.renderX != 100 {
		t.Errorf("InFlex: child2.x expected 100, got %v", child2.renderX)
	}
}

// TestPadding_ZeroInset は All(0) の場合に子の座標が変わらないことを確認する。
func TestPadding_ZeroInset(t *testing.T) {
	uictx := context.NewUIContext()
	child := &mockWidget{w: 50, h: 20}

	p := NewPadding(child).All(0)
	p.Render(uictx, 5, 10)

	if child.renderX != 5 || child.renderY != 10 {
		t.Errorf("ZeroInset: expected child at (5,10), got (%v,%v)", child.renderX, child.renderY)
	}
}
