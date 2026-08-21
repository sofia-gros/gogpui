package layout

import (
	"testing"

	"github.com/sofiagros/gogpui/core/context"
)

// --- SpacerFixed テスト ---

// TestSpacerW は SpacerW が水平方向の固定サイズを正しく返すことを確認する。
func TestSpacerW(t *testing.T) {
	uictx := context.NewUIContext()
	s := SpacerW(20)
	w, h := s.Render(uictx, 0, 0)
	if w != 20 || h != 0 {
		t.Errorf("SpacerW(20): expected (20, 0), got (%v, %v)", w, h)
	}
}

// TestSpacerH は SpacerH が垂直方向の固定サイズを正しく返すことを確認する。
func TestSpacerH(t *testing.T) {
	uictx := context.NewUIContext()
	s := SpacerH(15)
	w, h := s.Render(uictx, 0, 0)
	if w != 0 || h != 15 {
		t.Errorf("SpacerH(15): expected (0, 15), got (%v, %v)", w, h)
	}
}

// TestSpacerFixed は NewSpacerFixed が指定したサイズを正しく返すことを確認する。
func TestSpacerFixed(t *testing.T) {
	uictx := context.NewUIContext()
	s := NewSpacerFixed(30, 10)
	w, h := s.Render(uictx, 0, 0)
	if w != 30 || h != 10 {
		t.Errorf("SpacerFixed(30, 10): expected (30, 10), got (%v, %v)", w, h)
	}
}

// --- SpacerFlex テスト ---

// TestSpacerFlex_IsGrow は SpacerFlex が GrowWidget を実装し GrowFactor=1.0 を返すことを確認する。
func TestSpacerFlex_IsGrow(t *testing.T) {
	s := NewSpacerFlex()
	var _ GrowWidget = s // コンパイル時チェック

	if s.GrowFactor() != 1.0 {
		t.Errorf("SpacerFlex default GrowFactor: expected 1.0, got %v", s.GrowFactor())
	}
}

// TestSpacerFlex_WithGrowFactor は WithGrowFactor でカスタム係数が設定できることを確認する。
func TestSpacerFlex_WithGrowFactor(t *testing.T) {
	s := NewSpacerFlex().WithGrowFactor(2.0)
	if s.GrowFactor() != 2.0 {
		t.Errorf("SpacerFlex WithGrowFactor: expected 2.0, got %v", s.GrowFactor())
	}
}

// TestSpacerFlex_AllocatedSize は SetAllocatedSize 後に Render が正しいサイズを返すことを確認する。
func TestSpacerFlex_AllocatedSize(t *testing.T) {
	uictx := context.NewUIContext()
	s := NewSpacerFlex()

	// 割り当て前: (0, 0)
	w, h := s.Render(uictx, 0, 0)
	if w != 0 || h != 0 {
		t.Errorf("SpacerFlex before alloc: expected (0,0), got (%v,%v)", w, h)
	}

	// 割り当て後: SetAllocatedSize で幅 100 を設定
	s.SetAllocatedSize(100, 0)
	w, h = s.Render(uictx, 0, 0)
	if w != 100 || h != 0 {
		t.Errorf("SpacerFlex after alloc: expected (100,0), got (%v,%v)", w, h)
	}
}

// TestFlex_GrowSpacer_Row は Row Flex + SpacerFlex が残りスペースを正しく埋めることを確認する。
func TestFlex_GrowSpacer_Row(t *testing.T) {
	uictx := context.NewUIContext()

	w1 := &mockWidget{w: 100, h: 40}
	spacer := NewSpacerFlex()
	w2 := &mockWidget{w: 80, h: 40}

	// 親幅 400, gap 0
	// w1=100, spacer=残り(400-100-80=220), w2=80
	flex := NewFlex().
		Direction(Row).
		WithConstraints(400, 0).
		Add(w1, spacer, w2)

	flex.Render(uictx, 0, 0)

	// w1: x=0
	if w1.renderX != 0 {
		t.Errorf("GrowRow: w1.x expected 0, got %v", w1.renderX)
	}
	// spacer は 220px 割り当て → w2: x=100+220=320
	if w2.renderX != 320 {
		t.Errorf("GrowRow: w2.x expected 320, got %v", w2.renderX)
	}
	// spacer のサイズが 220 になっていることを確認
	sw, _ := spacer.Render(uictx, 0, 0)
	if sw != 220 {
		t.Errorf("GrowRow: spacer.w expected 220, got %v", sw)
	}
}

// TestFlex_GrowSpacer_Column は Column Flex + SpacerFlex が残りスペースを正しく埋めることを確認する。
func TestFlex_GrowSpacer_Column(t *testing.T) {
	uictx := context.NewUIContext()

	w1 := &mockWidget{w: 80, h: 50}
	spacer := NewSpacerFlex()
	w2 := &mockWidget{w: 80, h: 30}

	// 親高さ 300, gap 0
	// w1=50, spacer=残り(300-50-30=220), w2=30
	flex := NewFlex().
		Direction(Column).
		WithConstraints(0, 300).
		Add(w1, spacer, w2)

	flex.Render(uictx, 0, 0)

	// w1: y=0
	if w1.renderY != 0 {
		t.Errorf("GrowCol: w1.y expected 0, got %v", w1.renderY)
	}
	// w2: y=50+220=270
	if w2.renderY != 270 {
		t.Errorf("GrowCol: w2.y expected 270, got %v", w2.renderY)
	}
}

// TestFlex_MultipleGrow は grow factor が異なる複数 SpacerFlex の配分が正しいことを確認する。
func TestFlex_MultipleGrow(t *testing.T) {
	uictx := context.NewUIContext()

	// 親幅 400, アイテムなし, spacer1(factor=1) + spacer2(factor=3)
	// remaining = 400, total factor = 4
	// spacer1 = 400 * 1/4 = 100, spacer2 = 400 * 3/4 = 300
	spacer1 := NewSpacerFlex().WithGrowFactor(1)
	spacer2 := NewSpacerFlex().WithGrowFactor(3)

	flex := NewFlex().
		Direction(Row).
		WithConstraints(400, 0).
		Add(spacer1, spacer2)

	totalW, _ := flex.Render(uictx, 0, 0)

	if totalW != 400 {
		t.Errorf("MultiGrow: totalW expected 400, got %v", totalW)
	}

	sw1, _ := spacer1.Render(uictx, 0, 0)
	sw2, _ := spacer2.Render(uictx, 0, 0)

	if sw1 != 100 {
		t.Errorf("MultiGrow: spacer1 expected 100, got %v", sw1)
	}
	if sw2 != 300 {
		t.Errorf("MultiGrow: spacer2 expected 300, got %v", sw2)
	}
}

// TestFlex_GrowWithGap は gap がある場合でも grow 配分が正しいことを確認する。
func TestFlex_GrowWithGap(t *testing.T) {
	uictx := context.NewUIContext()

	w1 := &mockWidget{w: 100, h: 40}
	spacer := NewSpacerFlex()
	w2 := &mockWidget{w: 100, h: 40}

	// 親幅 400, gap=10, 子3つ
	// gaps = 10 * 2 = 20
	// nonGrow = 100 + 100 = 200
	// remaining = 400 - 200 - 20 = 180
	// spacer = 180
	// w2.x = 100 + 10 + 180 + 10 = 300
	flex := NewFlex().
		Direction(Row).
		Gap(10).
		WithConstraints(400, 0).
		Add(w1, spacer, w2)

	flex.Render(uictx, 0, 0)

	if w2.renderX != 300 {
		t.Errorf("GrowWithGap: w2.x expected 300, got %v", w2.renderX)
	}
}
