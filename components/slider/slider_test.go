package slider

import (
	"testing"

	testutil "github.com/sofiagros/gogpui/testing"
)

func TestSlider_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	tests := []struct {
		name      string
		component *Slider
	}{
		{
			name:      "Default",
			component: New("slider-default"),
		},
		{
			name:      "Value25",
			component: New("slider-25").Value(0.25),
		},
		{
			name:      "Value75",
			component: New("slider-75").Value(0.75),
		},
		{
			name:      "Disabled",
			component: New("slider-disabled").Value(0.5).Disabled(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tester.SimulateHover(0, 0)
			tester.StepFrame()
			score, err := tester.AssertGoldenImage(tt.component, "Slider", tt.name)
			if err != nil {
				t.Fatalf("Failed to assert golden image: %v", err)
			}
			if score > 0.05 {
				t.Errorf("Golden image mismatch for %s, score: %f", tt.name, score)
			}
		})
	}
}

func TestSlider_Dragging(t *testing.T) {
	tester := testutil.NewTester()
	var currentVal float64
	fired := false

	slider := New("slider-drag").Value(0.5).OnChange(func(val float64) {
		fired = true
		currentVal = val
	})

	// Render once to measure (width is 200, height is 16)
	slider.Render(tester.UI, 10, 10)

	// Simulate drag from middle (x=110, y=18) to right (x=160, y=18)
	tester.SimulateClick(110, 18)
	slider.Render(tester.UI, 10, 10)

	tester.SimulateHover(160, 18) // Mouse moves while left down (LeftPressed was false, LeftDown is true, wait tester.SimulateHover resets LeftDown?)
	// Let's manually set LeftDown = true
	tester.UI.Mouse.LeftDown = true
	tester.UI.Mouse.LeftPressed = false
	slider.Render(tester.UI, 10, 10)

	if !fired {
		t.Error("Expected Slider onChange to be fired")
	}
	
	expectedVal := 150.0 / 200.0 // (160 - 10) / 200 = 150 / 200 = 0.75
	if currentVal != expectedVal {
		t.Errorf("Expected Slider value to be %f, got %f", expectedVal, currentVal)
	}

	tester.SimulateClickRelease(160, 18)
	slider.Render(tester.UI, 10, 10)
}
