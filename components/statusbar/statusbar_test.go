package statusbar

import (
	"image/color"
	"testing"

	"github.com/sofiagros/gogpui/components/label"
	testutil "github.com/sofiagros/gogpui/testing"
)

func TestStatusBar_API(t *testing.T) {
	s := New().
		Width(500).
		Left(label.New("Left")).
		Right(label.New("Right")).
		Child(label.New("Center"))

	if s.width != 500 {
		t.Errorf("expected width 500")
	}
	if len(s.left) != 1 {
		t.Errorf("expected 1 left child")
	}
	if len(s.right) != 1 {
		t.Errorf("expected 1 right child")
	}
	if len(s.children) != 1 {
		t.Errorf("expected 1 center child")
	}
}

func TestStatusBar_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	cases := []struct {
		name      string
		component *StatusBar
	}{
		{
			name: "AllRegions",
			component: New().Width(400).
				Left(label.New("Ln 1, Col 1")).
				Child(label.New("Building...")).
				Right(label.New("UTF-8")),
		},
		{
			name: "LeftOnly",
			component: New().Width(400).
				Left(label.New("Ln 1, Col 1")).
				Child(label.New("End Aligned Child")),
		},
		{
			name: "RightOnly",
			component: New().Width(400).
				Child(label.New("Start Aligned Child")).
				Right(label.New("UTF-8")),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tester.Ctx.SetColor(color.RGBA{240, 240, 240, 255})
			tester.Ctx.Clear()

			tc.component.Render(tester.UI, 10, 10)

			score, err := tester.AssertGoldenImage(tc.component, "StatusBar", tc.name)
			if err != nil {
				t.Fatalf("Failed to assert golden image: %v", err)
			}
			t.Logf("Score for %s: %f", tc.name, score)
		})
	}
}
