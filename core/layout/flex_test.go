package layout

import (
	"testing"

	"github.com/sofiagros/gogpui/core/context"
)

// mockWidget は固定サイズを返し、描画座標を記録するダミーウィジェットである。
type mockWidget struct {
	w, h    float64
	renderX float64
	renderY float64
}

func (m *mockWidget) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	if !uictx.MeasureOnly {
		m.renderX = x
		m.renderY = y
	}
	return m.w, m.h
}

// --- 既存テスト ---

func TestFlexLayout_Row(t *testing.T) {
	uictx := context.NewUIContext()

	w1 := &mockWidget{w: 100, h: 40}
	w2 := &mockWidget{w: 150, h: 50}

	flex := NewFlex().Direction(Row).Gap(10).Add(w1, w2)

	totalW, totalH := flex.Render(uictx, 10, 20)

	if totalW != 260 {
		t.Errorf("Expected total width 260 (100 + 10 + 150), got %v", totalW)
	}
	if totalH != 50 {
		t.Errorf("Expected total height 50 (max of 40, 50), got %v", totalH)
	}
	if w1.renderX != 10 || w1.renderY != 20 {
		t.Errorf("Expected w1 at (10, 20), got (%v, %v)", w1.renderX, w1.renderY)
	}
	if w2.renderX != 120 || w2.renderY != 20 {
		t.Errorf("Expected w2 at (120, 20), got (%v, %v)", w2.renderX, w2.renderY)
	}
}

func TestFlexLayout_ColumnCenter(t *testing.T) {
	uictx := context.NewUIContext()

	w1 := &mockWidget{w: 100, h: 40}
	w2 := &mockWidget{w: 150, h: 50}

	flex := NewFlex().Direction(Column).Align(AlignCenter).Gap(5).Add(w1, w2)

	totalW, totalH := flex.Render(uictx, 0, 0)

	if totalW != 150 {
		t.Errorf("Expected total width 150, got %v", totalW)
	}
	if totalH != 95 {
		t.Errorf("Expected total height 95 (40 + 5 + 50), got %v", totalH)
	}
	// w1 (100): (150-100)/2 = 25
	if w1.renderX != 25 || w1.renderY != 0 {
		t.Errorf("Expected w1 at (25, 0), got (%v, %v)", w1.renderX, w1.renderY)
	}
	// w2 (150): (150-150)/2 = 0
	if w2.renderX != 0 || w2.renderY != 45 {
		t.Errorf("Expected w2 at (0, 45), got (%v, %v)", w2.renderX, w2.renderY)
	}
}

// --- JustifyContent テスト ---

// TestFlexLayout_JustifyEnd は Row + JustifyEnd で、アイテムが右端に寄ることを確認する。
func TestFlexLayout_JustifyEnd(t *testing.T) {
	uictx := context.NewUIContext()

	w1 := &mockWidget{w: 100, h: 40}
	w2 := &mockWidget{w: 80, h: 40}

	// 親幅 400、アイテム合計 100 + 10(gap) + 80 = 190
	// cursor = 400 - 190 = 210
	flex := NewFlex().
		Direction(Row).
		Gap(10).
		Justify(JustifyEnd).
		WithConstraints(400, 0).
		Add(w1, w2)

	flex.Render(uictx, 0, 0)

	// w1: 0 + 210 = 210
	if w1.renderX != 210 {
		t.Errorf("JustifyEnd: expected w1.x=210, got %v", w1.renderX)
	}
	// w2: 210 + 100 + 10 = 320
	if w2.renderX != 320 {
		t.Errorf("JustifyEnd: expected w2.x=320, got %v", w2.renderX)
	}
}

// TestFlexLayout_JustifyCenter は Row + JustifyCenter で、アイテムが中央に寄ることを確認する。
func TestFlexLayout_JustifyCenter(t *testing.T) {
	uictx := context.NewUIContext()

	w1 := &mockWidget{w: 100, h: 40}
	w2 := &mockWidget{w: 80, h: 40}

	// 親幅 400、アイテム合計 190
	// cursor = (400 - 190) / 2 = 105
	flex := NewFlex().
		Direction(Row).
		Gap(10).
		Justify(JustifyCenter).
		WithConstraints(400, 0).
		Add(w1, w2)

	flex.Render(uictx, 0, 0)

	if w1.renderX != 105 {
		t.Errorf("JustifyCenter: expected w1.x=105, got %v", w1.renderX)
	}
	// w2: 105 + 100 + 10 = 215
	if w2.renderX != 215 {
		t.Errorf("JustifyCenter: expected w2.x=215, got %v", w2.renderX)
	}
}

// TestFlexLayout_JustifySpaceBetween は Row + JustifySpaceBetween で、
// アイテム間のスペースが均等になることを確認する。
func TestFlexLayout_JustifySpaceBetween(t *testing.T) {
	uictx := context.NewUIContext()

	w1 := &mockWidget{w: 80, h: 40}
	w2 := &mockWidget{w: 80, h: 40}
	w3 := &mockWidget{w: 80, h: 40}

	// 親幅 400、rawMainTotal = 80*3 = 240
	// spacing = (400 - 240) / (3-1) = 160 / 2 = 80
	// w1: 0, w2: 80+80=160, w3: 160+80+80=320
	flex := NewFlex().
		Direction(Row).
		Justify(JustifySpaceBetween).
		WithConstraints(400, 0).
		Add(w1, w2, w3)

	flex.Render(uictx, 0, 0)

	if w1.renderX != 0 {
		t.Errorf("SpaceBetween: expected w1.x=0, got %v", w1.renderX)
	}
	if w2.renderX != 160 {
		t.Errorf("SpaceBetween: expected w2.x=160, got %v", w2.renderX)
	}
	if w3.renderX != 320 {
		t.Errorf("SpaceBetween: expected w3.x=320, got %v", w3.renderX)
	}
}

// TestFlexLayout_JustifySpaceAround は Row + JustifySpaceAround で、
// アイテムの周囲のスペースが均等になることを確認する。
func TestFlexLayout_JustifySpaceAround(t *testing.T) {
	uictx := context.NewUIContext()

	w1 := &mockWidget{w: 80, h: 40}
	w2 := &mockWidget{w: 80, h: 40}
	w3 := &mockWidget{w: 80, h: 40}

	// 親幅 400、rawMainTotal = 240
	// unit = (400 - 240) / 3 = 160 / 3 ≈ 53.33...
	// cursor start = unit/2 ≈ 26.67 → Round → 27
	// w1: 27, w2: 27+80+53=160, w3: 160+80+53=293 → 各々 math.Round される
	flex := NewFlex().
		Direction(Row).
		Justify(JustifySpaceAround).
		WithConstraints(400, 0).
		Add(w1, w2, w3)

	flex.Render(uictx, 0, 0)

	unit := (400.0 - 240.0) / 3.0
	expectedX1 := unit / 2.0

	// math.Round が適用されるので誤差 1px を許容する
	diff := w1.renderX - expectedX1
	if diff < -1 || diff > 1 {
		t.Errorf("SpaceAround: expected w1.x≈%.1f, got %v", expectedX1, w1.renderX)
	}
}

// --- WrapMode テスト ---

// TestFlexLayout_Wrap は Wrap モードで折り返しが発生することを確認する。
func TestFlexLayout_Wrap(t *testing.T) {
	uictx := context.NewUIContext()

	w1 := &mockWidget{w: 150, h: 40}
	w2 := &mockWidget{w: 150, h: 40}
	w3 := &mockWidget{w: 150, h: 40}

	// 親幅 320: w1(150) + gap(10) + w2(150) = 310 < 320 → 同行
	//           w3 は 310 + 10 + 150 = 470 > 320 → 改行
	flex := NewFlex().
		Direction(Row).
		Gap(10).
		WrapContent(Wrap).
		WithConstraints(320, 0).
		Add(w1, w2, w3)

	_, totalH := flex.Render(uictx, 0, 0)

	// 2行: 40 + 10(gap) + 40 = 90
	if totalH != 90 {
		t.Errorf("Wrap: expected totalH=90, got %v", totalH)
	}
	// w1, w2 は Y=0 の行
	if w1.renderY != 0 {
		t.Errorf("Wrap: expected w1.y=0, got %v", w1.renderY)
	}
	if w2.renderY != 0 {
		t.Errorf("Wrap: expected w2.y=0, got %v", w2.renderY)
	}
	// w3 は Y=50 の行
	if w3.renderY != 50 {
		t.Errorf("Wrap: expected w3.y=50, got %v", w3.renderY)
	}
}

// TestFlexLayout_WrapReverse は WrapReverse モードで行の順序が逆になることを確認する。
func TestFlexLayout_WrapReverse(t *testing.T) {
	uictx := context.NewUIContext()

	w1 := &mockWidget{w: 150, h: 40}
	w2 := &mockWidget{w: 150, h: 40}
	w3 := &mockWidget{w: 150, h: 40}

	// Wrap と同じ分割だが行順序が逆: [w3] が Y=0、[w1, w2] が Y=50
	flex := NewFlex().
		Direction(Row).
		Gap(10).
		WrapContent(WrapReverse).
		WithConstraints(320, 0).
		Add(w1, w2, w3)

	flex.Render(uictx, 0, 0)

	// WrapReverse: 後の行が先に描画される
	if w3.renderY != 0 {
		t.Errorf("WrapReverse: expected w3.y=0, got %v", w3.renderY)
	}
	if w1.renderY != 50 {
		t.Errorf("WrapReverse: expected w1.y=50, got %v", w1.renderY)
	}
}

// TestFlexLayout_NoConstraint_JustifyCenter は制約なしの場合
// JustifyCenter が Start と等価（cursor=0）になることを確認する。
func TestFlexLayout_NoConstraint_JustifyCenter(t *testing.T) {
	uictx := context.NewUIContext()

	w1 := &mockWidget{w: 100, h: 40}

	flex := NewFlex().
		Direction(Row).
		Justify(JustifyCenter).
		Add(w1) // WithConstraints なし

	flex.Render(uictx, 0, 0)

	// maxMain = mainTotal のとき (maxMain - mainTotal) / 2 = 0
	if w1.renderX != 0 {
		t.Errorf("NoConstraint JustifyCenter: expected w1.x=0, got %v", w1.renderX)
	}
}
