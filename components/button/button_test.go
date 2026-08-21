package button

import (
	"testing"
	testutil "github.com/sofiagros/gogpui/testing"
)

func TestButton_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	tests := []struct {
		name      string
		component *Button
	}{
		{
			name:      "Default",
			component: New("btn-default").Label("Primary Button"),
		},
		{
			name:      "Primary",
			component: New("btn-primary").Label("Primary").Primary(),
		},
		{
			name:      "Danger",
			component: New("btn-danger").Label("Danger").Danger(),
		},
		{
			name:      "Disabled",
			component: New("btn-disabled").Label("Disabled").Primary().Disabled(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Normal state
			tester.SimulateHover(0, 0)
			tester.StepFrame()
			score, err := tester.AssertGoldenImage(tt.component, "Button", tt.name+"_Normal")
			if err != nil {
				t.Fatalf("Failed to assert golden image: %v", err)
			}
			t.Logf("Score for %s_Normal: %f", tt.name, score)
			if score > 0.05 { // Allow slight sub-pixel anti-aliasing differences
				t.Errorf("Golden image mismatch for %s Normal, score: %f", tt.name, score)
			}

			// Hover state
			tester.SimulateHover(20, 20)
			for i := 0; i < 10; i++ { // Run a few frames to let animation reach 1.0
				tester.StepFrame()
				tester.Ctx.Clear()
				tt.component.Render(tester.UI, 10, 10)
			}
			score, err = tester.AssertGoldenImage(tt.component, "Button", tt.name+"_Hovered")
			if err != nil {
				t.Fatalf("Failed to assert golden image: %v", err)
			}
			if score > 0.05 {
				t.Errorf("Golden image mismatch for %s Hovered, score: %f", tt.name, score)
			}
		})
	}
}

func TestButton_OnClick(t *testing.T) {
	tester := testutil.NewTester()
	clicked := false
	btn := New("btn-click").Label("Click Me").OnClick(func() {
		clicked = true
	})

	// Hover and click
	tester.SimulateHover(20, 20)
	btn.Render(tester.UI, 10, 10)

	tester.SimulateClick(20, 20)
	btn.Render(tester.UI, 10, 10)

	tester.SimulateClickRelease(20, 20)
	btn.Render(tester.UI, 10, 10)

	if !clicked {
		t.Error("Expected button onClick to be fired")
	}

	// Test Disabled click
	clicked = false
	btn.Disabled(true)
	tester.SimulateHover(20, 20)
	btn.Render(tester.UI, 10, 10)

	tester.SimulateClick(20, 20)
	btn.Render(tester.UI, 10, 10)

	tester.SimulateClickRelease(20, 20)
	btn.Render(tester.UI, 10, 10)

	if clicked {
		t.Error("Expected button onClick NOT to be fired when disabled")
	}
}
