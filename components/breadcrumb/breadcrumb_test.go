package breadcrumb

import (
	"testing"

	testutil "github.com/sofiagros/gogpui/testing"
)

func TestBreadcrumb_NewItem(t *testing.T) {
	item := NewItem("item1", "Home").Disabled(true).OnClick(func() {})

	if item.id != "item1" {
		t.Errorf("Expected id 'item1', got %v", item.id)
	}
	if item.label != "Home" {
		t.Errorf("Expected label 'Home', got %v", item.label)
	}
	if !item.disabled {
		t.Errorf("Expected disabled true, got %v", item.disabled)
	}
	if item.onClick == nil {
		t.Errorf("Expected onClick to be set")
	}
}

func TestBreadcrumb_New(t *testing.T) {
	b := New().Add(
		NewItem("b1", "Docs"),
		NewItem("b2", "Components"),
	)
	if len(b.items) != 2 {
		t.Errorf("Expected 2 items, got %v", len(b.items))
	}
}

func TestBreadcrumb_GoldenImages(t *testing.T) {
	tester := testutil.NewTester()

	tests := []struct {
		name  string
		setup func() *Breadcrumb
		hover string // hover する item の id
	}{
		{
			name: "SingleItem",
			setup: func() *Breadcrumb {
				return New().Add(NewItem("b1", "Home"))
			},
		},
		{
			name: "MultipleItems",
			setup: func() *Breadcrumb {
				return New().Add(
					NewItem("b1", "Home"),
					NewItem("b2", "Docs"),
					NewItem("b3", "Components"),
				)
			},
		},
		{
			name: "WithDisabled",
			setup: func() *Breadcrumb {
				return New().Add(
					NewItem("b1", "Home"),
					NewItem("b2", "Settings").Disabled(true),
				)
			},
		},
		{
			name: "HoveredItem",
			setup: func() *Breadcrumb {
				return New().Add(
					NewItem("b1", "Home").OnClick(func() {}),
					NewItem("b2", "Docs"),
				)
			},
			hover: "b1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tester.Ctx.Clear()
			b := tt.setup()

			if tt.hover != "" {
				// hover のためのステートセット
				state := tester.UI.GetState(tt.hover)
				state.HoverRatio = 1.0
			}

			// 描画
			b.Render(tester.UI, 20, 20)

			score, err := tester.AssertGoldenImage(b, "Breadcrumb", tt.name)
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
