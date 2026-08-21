package badge

import (
	"image/color"
	"testing"
)

func TestBadge_New(t *testing.T) {
	b := New()
	if b.variant != Number {
		t.Errorf("Expected default variant Number, got %v", b.variant)
	}

	b.Dot()
	if b.variant != Dot {
		t.Errorf("Expected Dot variant, got %v", b.variant)
	}

	b.Count(5).Max(10)
	if b.count != 5 || b.max != 10 {
		t.Errorf("Expected count 5 and max 10, got %d and %d", b.count, b.max)
	}

	b.Color(color.White)
	if b.color == nil {
		t.Errorf("Expected color to be set")
	}

	b.Size(Large)
	if b.size != Large {
		t.Errorf("Expected Large size, got %v", b.size)
	}
}
