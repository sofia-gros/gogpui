package layout

import (
	"testing"

	"github.com/sofiagros/gogpui/core/context"
)

// mockWidget is a dummy widget that returns fixed dimensions and records where it was rendered.
type mockWidget struct {
	w, h   float64
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

func TestFlexLayout_Row(t *testing.T) {
	uictx := context.NewUIContext()
	
	w1 := &mockWidget{w: 100, h: 40}
	w2 := &mockWidget{w: 150, h: 50}

	flex := NewFlex().Direction(Row).Gap(10).Add(w1, w2)
	
	totalW, totalH := flex.Render(uictx, 10, 20)

	// Measure total dimensions
	if totalW != 260 {
		t.Errorf("Expected total width 260 (100 + 10 + 150), got %v", totalW)
	}
	if totalH != 50 {
		t.Errorf("Expected total height 50 (max of 40, 50), got %v", totalH)
	}

	// Verify coordinates for child 1
	if w1.renderX != 10 || w1.renderY != 20 {
		t.Errorf("Expected w1 at (10, 20), got (%v, %v)", w1.renderX, w1.renderY)
	}

	// Verify coordinates for child 2
	if w2.renderX != 120 || w2.renderY != 20 {
		t.Errorf("Expected w2 at (120, 20), got (%v, %v)", w2.renderX, w2.renderY)
	}
}

func TestFlexLayout_ColumnCenter(t *testing.T) {
	uictx := context.NewUIContext()
	
	w1 := &mockWidget{w: 100, h: 40}
	w2 := &mockWidget{w: 150, h: 50}

	// Column direction, cross-axis Center
	flex := NewFlex().Direction(Column).Align(AlignCenter).Gap(5).Add(w1, w2)
	
	totalW, totalH := flex.Render(uictx, 0, 0)

	if totalW != 150 {
		t.Errorf("Expected total width 150, got %v", totalW)
	}
	if totalH != 95 {
		t.Errorf("Expected total height 95 (40 + 5 + 50), got %v", totalH)
	}

	// Cross axis (X for column) center alignment:
	// Max width is 150.
	// w1 (100) -> (150 - 100)/2 = 25
	// w2 (150) -> (150 - 150)/2 = 0
	
	if w1.renderX != 25 || w1.renderY != 0 {
		t.Errorf("Expected w1 at (25, 0), got (%v, %v)", w1.renderX, w1.renderY)
	}

	if w2.renderX != 0 || w2.renderY != 45 {
		t.Errorf("Expected w2 at (0, 45), got (%v, %v)", w2.renderX, w2.renderY)
	}
}
