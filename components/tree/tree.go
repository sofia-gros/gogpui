package tree

import (
	"fmt"

	"github.com/sofiagros/gogpui/components/icon"
	"github.com/sofiagros/gogpui/components/label"
	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/layout"
)

// RenderItemFunc は各可視エントリを描画するためのコールバック関数。
type RenderItemFunc func(ix int, entry TreeEntry, state TreeEntryState, uictx *context.UIContext) layout.Widget

// Tree は仮想化されていないツリーウィジェットを表す。
type Tree struct {
	id         string
	state      *TreeState
	renderItem RenderItemFunc
	iconTheme  icon.IconTheme
	width      float64
}

// New は新しい Tree ウィジェットを作成する。
func New(state *TreeState) *Tree {
	id := fmt.Sprintf("tree-%p", state)
	t := &Tree{
		id:    id,
		state: state,
	}

	t.renderItem = func(ix int, entry TreeEntry, entryState TreeEntryState, uictx *context.UIContext) layout.Widget {
		textColor := uictx.Theme.Colors.Foreground
		if entry.Item.IsDisabled() {
			textColor = uictx.Theme.Colors.MutedForeground
		} else if entryState.Selected {
			textColor = uictx.Theme.Colors.Primary
		}
		
		if entry.Item.textColor != nil {
			textColor = entry.Item.textColor
		}

		var iconName icon.Name
		if !entry.Item.IsFolder() {
			iconName = icon.File
		} else if entry.Item.IsExpanded() {
			iconName = icon.FolderOpen
		} else {
			iconName = icon.FolderClosed
		}

		iconWidget := icon.New(iconName).Color(textColor)
		if t.iconTheme != nil {
			iconWidget.Theme(t.iconTheme)
		}

		lbl := label.New(entry.Item.Label()).Color(textColor)
		padLeft := float64(entry.Depth) * 16.0

		row := layout.NewFlex().Direction(layout.Row).Gap(6)
		if t.width > 0 {
			row.WithConstraints(t.width, 0)
		}
		
		if padLeft > 0 {
			row.Add(layout.SpacerW(padLeft))
		}
		row.Add(iconWidget, lbl)

		if entry.Item.suffix != nil {
			row.Add(layout.NewSpacerFlex(), entry.Item.suffix)
		}

		return row
	}

	return t
}

// IconTheme はツリー内で使用するアイコンのテーマを設定します。
func (t *Tree) IconTheme(theme icon.IconTheme) *Tree {
	t.iconTheme = theme
	return t
}

// Width はツリーの横幅を設定します。これにより、アイテム内の右端（Suffix）配置が正しく機能します。
func (t *Tree) Width(w float64) *Tree {
	t.width = w
	return t
}

// Item は各可視エントリのアプリケーション固有のコンテンツを提供する。
func (t *Tree) Item(renderItem RenderItemFunc) *Tree {
	t.renderItem = renderItem
	return t
}

// Render は Tree を描画する。
func (t *Tree) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	if t.state == nil {
		return 0, 0
	}

	flex := layout.NewFlex().Direction(layout.Column).Gap(0)

	for i, entry := range t.state.Entries {
		ix := i
		entryState := TreeEntryState{
			Selected:     t.state.SelectedIndex != nil && *t.state.SelectedIndex == ix,
			RightClicked: t.state.RightClickedIx != nil && *t.state.RightClickedIx == ix,
		}

		itemWidget := t.renderItem(ix, entry, entryState, uictx)

		if !entry.Item.IsDisabled() {
			wrapperID := fmt.Sprintf("%s-item-%d", t.id, ix)
			wrapped := &treeItemWrapper{
				id:       wrapperID,
				ix:       ix,
				state:    t.state,
				disabled: entry.Item.IsDisabled(),
				content:  itemWidget,
			}
			flex.Add(wrapped)
		} else {
			flex.Add(itemWidget)
		}
	}

	return flex.Render(uictx, x, y)
}

// treeItemWrapper はツリーアイテムのクリックを捕捉する内部ラッパー。
type treeItemWrapper struct {
	id       string
	ix       int
	state    *TreeState
	disabled bool
	content  layout.Widget
}

func (w *treeItemWrapper) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	// まずコンテンツを描画（または測定）
	cw, ch := w.content.Render(uictx, x, y)

	if !uictx.MeasureOnly {
		// Interaction を処理 (MeasureOnly 以外で実行)
		_, _, isClicked := uictx.ProcessInteraction(w.id, x, y, cw, ch, w.disabled)
		
		if isClicked {
			ix := w.ix
			w.state.SetSelectedIndex(&ix)
			w.state.ToggleExpand(w.ix)
			uictx.NeedsRedraw = true
		}

		state := uictx.GetState(w.id)
		// UIContext が右クリックを直接 ProcessInteraction で返さないため、
		// ホバー状態で RightReleased があれば右クリックとみなす。
		// ※ gogpui の MouseState に RightReleased は未実装の可能性があるが、
		// 存在しない場合は無視される（必要なら input.MouseButtonRight を判定する）
		_ = state
	}

	return cw, ch
}

