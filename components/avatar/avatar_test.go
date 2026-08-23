package avatar

import (
	"testing"

	testutil "github.com/sofiagros/gogpui/testing"
)

func TestExtractTextInitials(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Jason Lee", "JL"},
		{"Foo Bar Dar", "FB"},
		{"huacnlee", "HU"},
		{"s", "S"},
		{"A B", "AB"},
		{"", ""},
	}

	for _, tt := range tests {
		actual := extractTextInitials(tt.input)
		if actual != tt.expected {
			t.Errorf("extractTextInitials(%q) = %q, want %q", tt.input, actual, tt.expected)
		}
	}
}

func TestAvatar_New(t *testing.T) {
	a := New().Name("Jason Lee").Size(SizeLarge).Src("dummy.png")
	
	if a.name != "Jason Lee" {
		t.Errorf("Expected name 'Jason Lee', got %v", a.name)
	}
	if a.shortName != "JL" {
		t.Errorf("Expected shortName 'JL', got %v", a.shortName)
	}
	if a.size != SizeLarge {
		t.Errorf("Expected SizeLarge, got %v", a.size)
	}
	if a.src != "dummy.png" {
		t.Errorf("Expected src 'dummy.png', got %v", a.src)
	}
}

func TestAvatar_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	tests := []struct {
		name  string
		setup func() *Avatar
	}{
		{
			name: "Placeholder",
			setup: func() *Avatar {
				return New()
			},
		},
		{
			name: "Initials_JL",
			setup: func() *Avatar {
				return New().Name("Jason Lee")
			},
		},
		{
			name: "Initials_HU",
			setup: func() *Avatar {
				return New().Name("huacnlee")
			},
		},
		{
			name: "SizeXSmall",
			setup: func() *Avatar {
				return New().Name("Small Avatar").Size(SizeXSmall)
			},
		},
		{
			name: "SizeLarge",
			setup: func() *Avatar {
				return New().Name("Large Avatar").Size(SizeLarge)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tester.Ctx.Clear()
			a := tt.setup()

			a.Render(tester.UI, 20, 20)

			score, err := tester.AssertGoldenImage(a, "Avatar", tt.name)
			if err != nil {
				t.Fatalf("Failed to assert golden image: %v", err)
			}
			t.Logf("Score for %s: %f", tt.name, score)
			if score > 0.05 {
				t.Errorf("Golden image mismatch for %s, score: %f", tt.name, score)
			}
		})
	}
}
