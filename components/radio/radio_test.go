package radio

import (
	"testing"

	testutil "github.com/sofiagros/gogpui/testing"
)

func TestRadio_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	tests := []struct {
		name      string
		component *Radio
	}{
		{
			name:      "Default",
			component: New("rd-default").Label("Default Radio"),
		},
		{
			name:      "Checked",
			component: New("rd-checked").Label("Checked Radio").Checked(true),
		},
		{
			name:      "Disabled",
			component: New("rd-disabled").Label("Disabled Radio").Disabled(true),
		},
		{
			name:      "DisabledChecked",
			component: New("rd-disabled-checked").Label("Disabled Checked").Disabled(true).Checked(true),
		},
		{
			name:      "Small",
			component: New("rd-small").Label("Small Radio").Size(SizeSmall),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Normal state
			tester.SimulateHover(0, 0)
			tester.StepFrame()
			score, err := tester.AssertGoldenImage(tt.component, "Radio", tt.name+"_Normal")
			if err != nil {
				t.Fatalf("Failed to assert golden image: %v", err)
			}
			if score > 0.05 {
				t.Errorf("Golden image mismatch for %s Normal, score: %f", tt.name, score)
			}
		})
	}
}

func TestRadio_OnChange(t *testing.T) {
	tester := testutil.NewTester()
	changedTo := false
	fired := false

	rd := New("rd-change").Label("Change Me").OnChange(func(checked bool) {
		fired = true
		changedTo = checked
	})

	// Hover and click
	tester.SimulateHover(10, 10)
	rd.Render(tester.UI, 10, 10)

	tester.SimulateClick(10, 10)
	rd.Render(tester.UI, 10, 10)

	tester.SimulateClickRelease(10, 10)
	rd.Render(tester.UI, 10, 10)

	if !fired {
		t.Error("Expected Radio onChange to be fired")
	}
	if !changedTo {
		t.Error("Expected Radio to pass true to onChange")
	}

	// Test Disabled click
	fired = false
	rd.Disabled(true)
	tester.SimulateHover(10, 10)
	rd.Render(tester.UI, 10, 10)

	tester.SimulateClick(10, 10)
	rd.Render(tester.UI, 10, 10)

	tester.SimulateClickRelease(10, 10)
	rd.Render(tester.UI, 10, 10)

	if fired {
		t.Error("Expected Radio onChange NOT to be fired when disabled")
	}
}
