package label

import (
	"testing"
)

func TestLabel_New(t *testing.T) {
	lbl := New("Test")
	if lbl.text != "Test" {
		t.Errorf("Expected 'Test', got '%s'", lbl.text)
	}

	lbl.Secondary("Secondary text")
	if lbl.secondary != "Secondary text" {
		t.Errorf("Expected 'Secondary text', got '%s'", lbl.secondary)
	}

	lbl.Masked(true)
	if !lbl.masked {
		t.Errorf("Expected masked to be true")
	}

	lbl.Highlights("Test", true)
	if lbl.highlights != "Test" || !lbl.isPrefix {
		t.Errorf("Expected highlights to be 'Test' with prefix true")
	}
}
