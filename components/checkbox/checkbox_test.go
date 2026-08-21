package checkbox

import (
	"testing"

	testutil "github.com/sofiagros/gogpui/testing"
)

func TestCheckbox_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	tests := []struct {
		name      string
		component *Checkbox
	}{
		{
			name:      "Default",
			component: New("chk-default").Label("Default Checkbox"),
		},
		{
			name:      "Checked",
			component: New("chk-checked").Label("Checked Checkbox").Checked(true),
		},
		{
			name:      "Disabled",
			component: New("chk-disabled").Label("Disabled Checkbox").Disabled(true),
		},
		{
			name:      "DisabledChecked",
			component: New("chk-disabled-checked").Label("Disabled Checked").Disabled(true).Checked(true),
		},
		{
			name:      "Small",
			component: New("chk-small").Label("Small Checkbox").Size(SizeSmall),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Normal state
			tester.SimulateHover(0, 0)
			tester.StepFrame()
			score, err := tester.AssertGoldenImage(tt.component, "Checkbox", tt.name+"_Normal")
			if err != nil {
				t.Fatalf("Failed to assert golden image: %v", err)
			}
			if score > 0.05 { // Allow slight sub-pixel anti-aliasing differences
				t.Errorf("Golden image mismatch for %s Normal, score: %f", tt.name, score)
			}

			// Hover state (only relevant if not disabled, but we can capture it anyway)
			tester.SimulateHover(10, 10)
			for i := 0; i < 10; i++ { // Run a few frames to let animation reach 1.0
				tester.StepFrame()
				tester.Ctx.Clear()
				tt.component.Render(tester.UI, 10, 10)
			}
			score, err = tester.AssertGoldenImage(tt.component, "Checkbox", tt.name+"_Hovered")
			if err != nil {
				t.Fatalf("Failed to assert golden image: %v", err)
			}
			if score > 0.05 {
				t.Errorf("Golden image mismatch for %s Hovered, score: %f", tt.name, score)
			}
		})
	}
}

func TestCheckbox_OnChange(t *testing.T) {
	tester := testutil.NewTester()
	changedTo := false
	fired := false

	chk := New("chk-change").Label("Change Me").OnChange(func(checked bool) {
		fired = true
		changedTo = checked
	})

	// Hover and click
	tester.SimulateHover(10, 10)
	chk.Render(tester.UI, 10, 10)

	tester.SimulateClick(10, 10)
	chk.Render(tester.UI, 10, 10)

	tester.SimulateClickRelease(10, 10)
	chk.Render(tester.UI, 10, 10)

	if !fired {
		t.Error("Expected Checkbox onChange to be fired")
	}
	if !changedTo {
		t.Error("Expected Checkbox to pass true to onChange")
	}

	// Test Disabled click
	fired = false
	chk.Disabled(true)
	tester.SimulateHover(10, 10)
	chk.Render(tester.UI, 10, 10)

	tester.SimulateClick(10, 10)
	chk.Render(tester.UI, 10, 10)

	tester.SimulateClickRelease(10, 10)
	chk.Render(tester.UI, 10, 10)

	if fired {
		t.Error("Expected Checkbox onChange NOT to be fired when disabled")
	}
}
