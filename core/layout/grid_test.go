package layout

import (
	"testing"

	"github.com/sofiagros/gogpui/core/context"
)

// TestGrid_Cols は Cols(3) で 3列グリッドが正しく配置されることを確認する。
func TestGrid_Cols(t *testing.T) {
	uictx := context.NewUIContext()

	widgets := []*mockWidget{
		{w: 100, h: 40}, // (0,0)
		{w: 100, h: 40}, // (0,1)
		{w: 100, h: 40}, // (0,2)
		{w: 100, h: 40}, // (1,0)
		{w: 100, h: 40}, // (1,1)
	}

	// 3列、gapX=10, gapY=5
	// 列幅: すべて 100, 列オフセット: [0, 110, 220]
	// 行高: 行0=40, 行1=40, 行オフセット: [0, 45]
	g := NewGrid().Cols(3).GapX(10).GapY(5)
	for _, w := range widgets {
		g.Add(w)
	}

	totalW, totalH := g.Render(uictx, 0, 0)

	// 全体サイズの確認
	// 幅: 100 + 10 + 100 + 10 + 100 = 320
	if totalW != 320 {
		t.Errorf("Cols: expected totalW=320, got %v", totalW)
	}
	// 高さ: 40 + 5 + 40 = 85
	if totalH != 85 {
		t.Errorf("Cols: expected totalH=85, got %v", totalH)
	}

	// 行0の配置確認
	if widgets[0].renderX != 0 || widgets[0].renderY != 0 {
		t.Errorf("widget[0] expected (0,0), got (%v,%v)", widgets[0].renderX, widgets[0].renderY)
	}
	if widgets[1].renderX != 110 || widgets[1].renderY != 0 {
		t.Errorf("widget[1] expected (110,0), got (%v,%v)", widgets[1].renderX, widgets[1].renderY)
	}
	if widgets[2].renderX != 220 || widgets[2].renderY != 0 {
		t.Errorf("widget[2] expected (220,0), got (%v,%v)", widgets[2].renderX, widgets[2].renderY)
	}

	// 行1の配置確認
	if widgets[3].renderX != 0 || widgets[3].renderY != 45 {
		t.Errorf("widget[3] expected (0,45), got (%v,%v)", widgets[3].renderX, widgets[3].renderY)
	}
	if widgets[4].renderX != 110 || widgets[4].renderY != 45 {
		t.Errorf("widget[4] expected (110,45), got (%v,%v)", widgets[4].renderX, widgets[4].renderY)
	}
}

// TestGrid_Rows は Rows(2) で行数から列数が自動計算されることを確認する。
func TestGrid_Rows(t *testing.T) {
	uictx := context.NewUIContext()

	// 4アイテム、2行 → 列数は ceil(4/2)=2
	widgets := []*mockWidget{
		{w: 80, h: 30}, // col0, row0
		{w: 80, h: 30}, // col1, row0
		{w: 80, h: 30}, // col0, row1
		{w: 80, h: 30}, // col1, row1
	}

	g := NewGrid().Rows(2).Gap(8)
	for _, w := range widgets {
		g.Add(w)
	}

	totalW, totalH := g.Render(uictx, 10, 20)

	// 幅: 80 + 8 + 80 = 168
	if totalW != 168 {
		t.Errorf("Rows: expected totalW=168, got %v", totalW)
	}
	// 高さ: 30 + 8 + 30 = 68
	if totalH != 68 {
		t.Errorf("Rows: expected totalH=68, got %v", totalH)
	}

	// オフセット 10, 20 が加算されることを確認
	if widgets[0].renderX != 10 || widgets[0].renderY != 20 {
		t.Errorf("Rows: widget[0] expected (10,20), got (%v,%v)", widgets[0].renderX, widgets[0].renderY)
	}
	if widgets[1].renderX != 98 || widgets[1].renderY != 20 {
		// 10 + 80 + 8 = 98
		t.Errorf("Rows: widget[1] expected (98,20), got (%v,%v)", widgets[1].renderX, widgets[1].renderY)
	}
	if widgets[2].renderX != 10 || widgets[2].renderY != 58 {
		// y: 20 + 30 + 8 = 58
		t.Errorf("Rows: widget[2] expected (10,58), got (%v,%v)", widgets[2].renderX, widgets[2].renderY)
	}
}

// TestGrid_GapXY は GapX と GapY が別々に設定できることを確認する。
func TestGrid_GapXY(t *testing.T) {
	uictx := context.NewUIContext()

	w1 := &mockWidget{w: 50, h: 20}
	w2 := &mockWidget{w: 50, h: 20}
	w3 := &mockWidget{w: 50, h: 20}
	w4 := &mockWidget{w: 50, h: 20}

	// 2列, gapX=20, gapY=10
	// 幅: 50+20+50=120, 高さ: 20+10+20=50
	g := NewGrid().Cols(2).GapX(20).GapY(10).Add(w1, w2, w3, w4)

	totalW, totalH := g.Render(uictx, 0, 0)

	if totalW != 120 {
		t.Errorf("GapXY: expected totalW=120, got %v", totalW)
	}
	if totalH != 50 {
		t.Errorf("GapXY: expected totalH=50, got %v", totalH)
	}
	// w3: col0, row1 → (0, 20+10) = (0, 30)
	if w3.renderX != 0 || w3.renderY != 30 {
		t.Errorf("GapXY: w3 expected (0,30), got (%v,%v)", w3.renderX, w3.renderY)
	}
	// w4: col1, row1 → (50+20, 30) = (70, 30)
	if w4.renderX != 70 || w4.renderY != 30 {
		t.Errorf("GapXY: w4 expected (70,30), got (%v,%v)", w4.renderX, w4.renderY)
	}
}

// TestGrid_MeasureOnly は MeasureOnly=true のとき座標を記録せず総サイズのみ返すことを確認する。
func TestGrid_MeasureOnly(t *testing.T) {
	uictx := context.NewUIContext()
	uictx.MeasureOnly = true

	w1 := &mockWidget{w: 60, h: 30}
	w2 := &mockWidget{w: 60, h: 30}

	g := NewGrid().Cols(2).Gap(5).Add(w1, w2)

	totalW, totalH := g.Render(uictx, 0, 0)

	// 幅: 60+5+60=125, 高さ: 30（1行）
	if totalW != 125 {
		t.Errorf("MeasureOnly: expected totalW=125, got %v", totalW)
	}
	if totalH != 30 {
		t.Errorf("MeasureOnly: expected totalH=30, got %v", totalH)
	}
	// 座標は記録されていないはず
	if w1.renderX != 0 || w1.renderY != 0 {
		t.Errorf("MeasureOnly: w1 coords should not be set, got (%v,%v)", w1.renderX, w1.renderY)
	}
}

// TestGrid_PartialLastRow はアイテム数が列数の倍数でない場合でも最終行が正しく配置されることを確認する。
func TestGrid_PartialLastRow(t *testing.T) {
	uictx := context.NewUIContext()

	// 3列、5アイテム: 最終行は1アイテムだけ
	widgets := []*mockWidget{
		{w: 100, h: 40},
		{w: 100, h: 40},
		{w: 100, h: 40},
		{w: 100, h: 40},
		{w: 100, h: 40}, // col0, row1 のみ
	}

	g := NewGrid().Cols(3).Gap(0)
	for _, w := range widgets {
		g.Add(w)
	}

	g.Render(uictx, 0, 0)

	// widget[4] は col=4%3=1, row=4/3=1 → (100, 40) に配置されるはず
	if widgets[4].renderX != 100 || widgets[4].renderY != 40 {
		t.Errorf("PartialLastRow: widget[4] expected (100,40), got (%v,%v)", widgets[4].renderX, widgets[4].renderY)
	}
}

// TestGrid_VariableColWidths は列ごとに幅が異なる場合でも正しく配置されることを確認する。
func TestGrid_VariableColWidths(t *testing.T) {
	uictx := context.NewUIContext()

	// 2列、列0の幅=60、列1の幅=120
	w1 := &mockWidget{w: 60, h: 30}
	w2 := &mockWidget{w: 120, h: 30}
	w3 := &mockWidget{w: 40, h: 30} // 列0、w=40 < 60 なので列幅は60のまま
	w4 := &mockWidget{w: 120, h: 30}

	g := NewGrid().Cols(2).GapX(10).GapY(0).Add(w1, w2, w3, w4)

	totalW, _ := g.Render(uictx, 0, 0)

	// 列幅: col0=max(60,40)=60, col1=max(120,120)=120
	// 幅: 60 + 10 + 120 = 190
	if totalW != 190 {
		t.Errorf("VarColWidths: expected totalW=190, got %v", totalW)
	}
	// w3 は col0, row1: x=0, y=30
	if w3.renderX != 0 {
		t.Errorf("VarColWidths: w3 expected x=0, got %v", w3.renderX)
	}
	// w4 は col1, row1: x=60+10=70, y=30
	if w4.renderX != 70 {
		t.Errorf("VarColWidths: w4 expected x=70, got %v", w4.renderX)
	}
}

// TestGrid_Empty は空グリッドが (0,0) を返すことを確認する。
func TestGrid_Empty(t *testing.T) {
	uictx := context.NewUIContext()

	g := NewGrid().Cols(3)
	w, h := g.Render(uictx, 0, 0)

	if w != 0 || h != 0 {
		t.Errorf("Empty: expected (0,0), got (%v,%v)", w, h)
	}
}
