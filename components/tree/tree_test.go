package tree

import (
	"image/color"
	"testing"

	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/layout"
	gpuitesting "github.com/sofiagros/gogpui/testing"
)

func TestTreeItem_Ancestors(t *testing.T) {
	leaf := NewItem("leaf", "Leaf")
	branch := NewItem("branch", "Branch").Child(leaf)
	root := NewItem("root", "Root").Child(branch)

	leaf.Disabled(true).Expanded(true)
	if !leaf.IsDisabled() {
		t.Error("expected leaf to be disabled")
	}
	if !leaf.IsExpanded() {
		t.Error("expected leaf to be expanded")
	}

	ancestors := root.Ancestors("leaf")
	if len(ancestors) != 2 {
		t.Fatalf("expected 2 ancestors, got %d", len(ancestors))
	}
	if ancestors[0].ID() != "branch" {
		t.Errorf("expected branch, got %s", ancestors[0].ID())
	}
	if ancestors[1].ID() != "root" {
		t.Errorf("expected root, got %s", ancestors[1].ID())
	}
}

func TestTreeState_FlatEntries(t *testing.T) {
	items := []*TreeItem{
		NewItem("src", "src").
			Expanded(true).
			Child(NewItem("src/lib.go", "lib.go")),
		NewItem("README.md", "README.md"),
	}

	state := NewTreeState().Items(items...)

	if len(state.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(state.Entries))
	}
	if state.Entries[1].Depth != 1 {
		t.Errorf("expected depth 1 for lib.go, got %d", state.Entries[1].Depth)
	}

	ix1 := 1
	state.SetSelectedIndex(&ix1)
	
	state.SetItems([]*TreeItem{NewItem("go.mod", "go.mod")})
	if state.SelectedIndex != nil {
		t.Error("expected selected index to be reset")
	}
	if len(state.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(state.Entries))
	}
}

func TestTreeState_SelectingHiddenItemExpandsAncestors(t *testing.T) {
	target := NewItem("src/ui/tree.go", "tree.go")
	root := NewItem("src", "src").Child(
		NewItem("src/ui", "ui").Child(target),
	)

	state := NewTreeState().Items(root)
	state.SetSelectedItem(target)

	if len(state.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(state.Entries))
	}
	if state.SelectedItem().ID() != "src/ui/tree.go" {
		t.Errorf("expected src/ui/tree.go, got %s", state.SelectedItem().ID())
	}
}

func TestTreeState_ToggleFolder(t *testing.T) {
	root := NewItem("src", "src").Child(NewItem("src/lib.go", "lib.go"))
	state := NewTreeState().Items(root)

	state.ToggleExpand(0)
	if len(state.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(state.Entries))
	}

	state.ToggleExpand(0)
	if len(state.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(state.Entries))
	}
}

func TestTree_GoldenImages(t *testing.T) {
	tester := gpuitesting.NewTester()

	root := NewItem("src", "src").Expanded(true).Child(
		NewItem("src/ui", "ui").Expanded(false).Child(
			NewItem("src/ui/tree.go", "tree.go"),
		),
	)
	state := NewTreeState().Items(root)
	
	// mock list item for tree rendering
	renderItem := func(ix int, entry TreeEntry, s TreeEntryState, uictx *context.UIContext) layout.Widget {
		// Just a simple box
		bg := color.RGBA{255, 255, 255, 255}
		if s.Selected {
			bg = color.RGBA{200, 200, 255, 255}
		}
		
		padding := float64(entry.Depth) * 20.0
		
		return layout.NewFlex().Direction(layout.Row).Add(
			&mockTextWidget{label: entry.Item.Label(), bg: bg, padding: padding},
		)
	}

	treeWidget := New(state).Item(renderItem)

	score, err := tester.AssertGoldenImage(treeWidget, "Tree", "Default")
	if err != nil {
		t.Fatalf("AssertGoldenImage failed: %v", err)
	}
	if score > 0.0001 {
		t.Errorf("Golden image mismatch, score: %f", score)
	}
}

type mockTextWidget struct {
	label   string
	bg      color.Color
	padding float64
}

func (w *mockTextWidget) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	ctx := uictx.GG
	width := 200.0
	height := 24.0

	if !uictx.MeasureOnly {
		ctx.SetColor(w.bg)
		ctx.DrawRectangle(x, y, width, height)
		uictx.Fill()

		ctx.SetColor(color.Black)
		uictx.DrawString(w.label, x+w.padding+4, y+height/2+4)
	}

	return width, height
}

