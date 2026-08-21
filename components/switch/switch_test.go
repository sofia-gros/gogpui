package switch_comp

import (
	"testing"

	testutil "github.com/sofiagros/gogpui/testing"
)

func TestSwitch_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	tests := []struct {
		name      string
		component *Switch
	}{
		{
			name:      "Default",
			component: New("sw-default").Label("Default Switch"),
		},
		{
			name:      "Checked",
			component: New("sw-checked").Label("Checked Switch").Checked(true),
		},
		{
			name:      "Disabled",
			component: New("sw-disabled").Label("Disabled Switch").Disabled(true),
		},
		{
			name:      "DisabledChecked",
			component: New("sw-disabled-checked").Label("Disabled Checked").Disabled(true).Checked(true),
		},
		{
			name:      "Small",
			component: New("sw-small").Label("Small Switch").Size(SizeSmall),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Normal state
			tester.SimulateHover(0, 0)
			// Step enough frames to let animations settle
			for i := 0; i < 20; i++ {
				tester.StepFrame()
				tester.Ctx.Clear()
				tt.component.Render(tester.UI, 10, 10)
			}
			
			score, err := tester.AssertGoldenImage(tt.component, "Switch", tt.name+"_Normal")
			if err != nil {
				t.Fatalf("Failed to assert golden image: %v", err)
			}
			if score > 0.05 {
				t.Errorf("Golden image mismatch for %s Normal, score: %f", tt.name, score)
			}
		})
	}
}

func TestSwitch_OnChange(t *testing.T) {
	tester := testutil.NewTester()
	changedTo := false
	fired := false

	sw := New("sw-change").Label("Change Me").OnChange(func(checked bool) {
		fired = true
		changedTo = checked
	})

	// Hover and click
	tester.SimulateHover(10, 10)
	sw.Render(tester.UI, 10, 10)

	tester.SimulateClick(10, 10)
	sw.Render(tester.UI, 10, 10)

	tester.SimulateClickRelease(10, 10)
	sw.Render(tester.UI, 10, 10)

	if !fired {
		t.Error("Expected Switch onChange to be fired")
	}
	if !changedTo {
		t.Error("Expected Switch to pass true to onChange")
	}

	// Test Disabled click
	fired = false
	sw.Disabled(true)
	tester.SimulateHover(10, 10)
	sw.Render(tester.UI, 10, 10)

	tester.SimulateClick(10, 10)
	sw.Render(tester.UI, 10, 10)

	tester.SimulateClickRelease(10, 10)
	sw.Render(tester.UI, 10, 10)

	if fired {
		t.Error("Expected Switch onChange NOT to be fired when disabled")
	}
}
