package progress

import (
	"image/color"
	"testing"

	testutil "github.com/sofiagros/gogpui/testing"
)

// TestProgress_New は Progress および ProgressCircle の初期化と設定をテストする。
func TestProgress_New(t *testing.T) {
	p := New("test-prog").Value(45.0).Size(SizeLarge).Width(300).Color(color.RGBA{0, 255, 0, 255})
	if p.id != "test-prog" || p.value != 45.0 || p.size != SizeLarge || p.width != 300 || p.color == nil {
		t.Errorf("Progress properties mismatch")
	}

	p.Value(150.0)
	if p.value != 100.0 {
		t.Errorf("Expected clamped value 100.0, got %f", p.value)
	}
	p.Value(-20.0)
	if p.value != 0.0 {
		t.Errorf("Expected clamped value 0.0, got %f", p.value)
	}

	c := NewCircle("test-circle").Value(75.0).Size(SizeSmall).Radius(20).StrokeWidth(3).Loading(true)
	if c.id != "test-circle" || c.value != 75.0 || c.size != SizeSmall || c.radius != 20 || c.strokeWidth != 3 || !c.loading {
		t.Errorf("ProgressCircle properties mismatch")
	}
}

// TestProgress_GoldenImages は Progress の各状態におけるゴールデンイメージテストを行う。
func TestProgress_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	tests := []struct {
		name      string
		component testutil.RenderComponent
	}{
		{
			name:      "Progress_0Percent",
			component: New("prog-0").Value(0).Width(200),
		},
		{
			name:      "Progress_50Percent",
			component: New("prog-50").Value(50).Width(200),
		},
		{
			name:      "Progress_100Percent",
			component: New("prog-100").Value(100).Width(200),
		},
		{
			name:      "Progress_SizeLarge",
			component: New("prog-large").Value(70).Size(SizeLarge).Width(200),
		},
		{
			name:      "Progress_SizeXSmall",
			component: New("prog-xsmall").Value(70).Size(SizeXSmall).Width(200),
		},
		{
			name:      "Circle_0Percent",
			component: NewCircle("circle-0").Value(0),
		},
		{
			name:      "Circle_50Percent",
			component: NewCircle("circle-50").Value(50),
		},
		{
			name:      "Circle_100Percent",
			component: NewCircle("circle-100").Value(100),
		},
		{
			name:      "Circle_SizeLarge",
			component: NewCircle("circle-large").Value(60).Size(SizeLarge),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// アニメーション完了まで数フレーム進める
			for i := 0; i < 15; i++ {
				tester.StepFrame()
				tester.Ctx.Clear()
				tt.component.Render(tester.UI, 10, 10)
			}
			score, err := tester.AssertGoldenImage(tt.component, "Progress", tt.name)
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
