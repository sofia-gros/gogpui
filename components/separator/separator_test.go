package separator

import (
	"image/color"
	"testing"

	testutil "github.com/sofiagros/gogpui/testing"
)

// TestSeparator_New は Separator の初期化およびビルダーメソッドをテストする。
func TestSeparator_New(t *testing.T) {
	// デフォルト値
	s := New()
	if s.orientation != OrientationHorizontal {
		t.Errorf("Expected OrientationHorizontal orientation, got %v", s.orientation)
	}
	if s.style != StyleSolid {
		t.Errorf("Expected StyleSolid style, got %v", s.style)
	}

	// 各種コンストラクタとビルダー
	hDashed := HorizontalDashed().Label("OR").Length(300)
	if hDashed.orientation != OrientationHorizontal || hDashed.style != StyleDashed || hDashed.label != "OR" || hDashed.length != 300 {
		t.Errorf("HorizontalDashed properties mismatch")
	}

	vDashed := VerticalDashed().Color(color.RGBA{255, 0, 0, 255}).Length(50)
	if vDashed.orientation != OrientationVertical || vDashed.style != StyleDashed || vDashed.color == nil || vDashed.length != 50 {
		t.Errorf("VerticalDashed properties mismatch")
	}

	s.Solid()
	if s.style != StyleSolid {
		t.Errorf("Expected StyleSolid, got %v", s.style)
	}
}

// TestSeparator_GoldenImages は Separator の各バリアントの描画とゴールデンイメージ一致をテストする。
func TestSeparator_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	tests := []struct {
		name      string
		component *Separator
	}{
		{
			name:      "Horizontal_Solid",
			component: Horizontal().Length(200),
		},
		{
			name:      "Horizontal_Dashed",
			component: HorizontalDashed().Length(200),
		},
		{
			name:      "Horizontal_WithLabel",
			component: Horizontal().Label("CONTINUE WITH").Length(260),
		},
		{
			name:      "Vertical_Solid",
			component: Vertical().Length(80),
		},
		{
			name:      "Vertical_Dashed",
			component: VerticalDashed().Length(80),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tester.Ctx.Clear()
			tt.component.Render(tester.UI, 10, 10)
			score, err := tester.AssertGoldenImage(tt.component, "Separator", tt.name)
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
