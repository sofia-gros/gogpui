package collapsible

import (
	"image/color"
	"testing"

	"github.com/sofiagros/gogpui/components/button"
	"github.com/sofiagros/gogpui/components/label"
	testutil "github.com/sofiagros/gogpui/testing"
)

func TestCollapsible_API(t *testing.T) {
	c := New().
		Open(true).
		Trigger(button.New("btn")).
		Content(label.New("Content"))

	if !c.open {
		t.Errorf("expected open to be true")
	}
	if c.trigger == nil {
		t.Errorf("expected trigger to be set")
	}
	if c.content == nil {
		t.Errorf("expected content to be set")
	}
}

func TestCollapsible_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	cases := []struct {
		name string
		open bool
	}{
		{"Closed", false},
		{"Open", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tester.Ctx.SetColor(color.RGBA{240, 240, 240, 255})
			tester.Ctx.Clear()

			c := New().
				Open(tc.open).
				Trigger(button.New("btn").Label("Toggle")).
				Content(label.New("Expanded Content Details..."))

			c.Render(tester.UI, 10, 10)

			score, err := tester.AssertGoldenImage(c, "Collapsible", tc.name)
			if err != nil {
				t.Fatalf("Failed to assert golden image: %v", err)
			}
			t.Logf("Score for %s: %f", tc.name, score)
		})
	}
}
