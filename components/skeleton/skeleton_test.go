package skeleton

import (
	"testing"

	testutil "github.com/sofiagros/gogpui/testing"
)

func TestSkeleton_New(t *testing.T) {
	s := New("test-skel").Width(300).Height(24).Secondary(true)

	if s.id != "test-skel" {
		t.Errorf("Expected id 'test-skel', got %v", s.id)
	}
	if s.width != 300 {
		t.Errorf("Expected width 300, got %v", s.width)
	}
	if s.height != 24 {
		t.Errorf("Expected height 24, got %v", s.height)
	}
	if !s.secondary {
		t.Errorf("Expected secondary true, got %v", s.secondary)
	}
}

func TestSkeleton_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	tests := []struct {
		name       string
		setup      func() *Skeleton
		advanceSec float64
	}{
		{
			name: "Primary_Start", // 0s
			setup: func() *Skeleton {
				return New("test-1").Width(200).Height(16)
			},
			advanceSec: 0.0,
		},
		{
			name: "Primary_Peak", // 1s (ToggleRatio = 0.5 -> delta = 1.0)
			setup: func() *Skeleton {
				return New("test-1").Width(200).Height(16)
			},
			advanceSec: 1.0,
		},
		{
			name: "Primary_End", // 2s (ToggleRatio = 1.0 -> delta = 0.0)
			setup: func() *Skeleton {
				return New("test-1").Width(200).Height(16)
			},
			advanceSec: 2.0,
		},
		{
			name: "Secondary_Peak", // 1s with Secondary = true
			setup: func() *Skeleton {
				return New("test-2").Width(200).Height(16).Secondary(true)
			},
			advanceSec: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tester.Ctx.Clear()
			s := tt.setup()
			
			// Set state manually
			state := tester.UI.GetState(s.id)
			state.ToggleRatio = (tt.advanceSec / 2.0)
			tester.UI.DeltaTime = 0 // prevent further advance in Render

			s.Render(tester.UI, 10, 10)

			score, err := tester.AssertGoldenImage(s, "Skeleton", tt.name)
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
