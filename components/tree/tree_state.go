package tree

// TreeEntry は TreeItem のフラットな文字列表現と深さ(Depth)を持つ。
type TreeEntry struct {
	Item  *TreeItem
	Depth int
}

// NewTreeEntry は TreeEntry を作成する。
func NewTreeEntry(item *TreeItem, depth int) TreeEntry {
	return TreeEntry{
		Item:  item,
		Depth: depth,
	}
}

// TreeEntryState はレンダリング時に使用される、エントリのインタラクション状態。
type TreeEntryState struct {
	Selected     bool
	RightClicked bool
}

// TreeState は仮想化されたツリーの振る舞いとインタラクション状態を保持する。
// （gpui-component では Entity に包まれてアプリ全体で共有されるが、
// gogpui ではユーザーが手元で保持するか、グローバルに持つ想定）
type TreeState struct {
	items          []*TreeItem
	Entries        []TreeEntry
	SelectedIndex  *int
	RightClickedIx *int
}

// NewTreeState は TreeState を作成する。
func NewTreeState() *TreeState {
	return &TreeState{
		items:          make([]*TreeItem, 0),
		Entries:        make([]TreeEntry, 0),
		SelectedIndex:  nil,
		RightClickedIx: nil,
	}
}

// Items はツリーのルートアイテムを設定し、自身を返す。
func (s *TreeState) Items(items ...*TreeItem) *TreeState {
	s.SetItems(items)
	return s
}

// SetItems はルートアイテムを置き換え、選択状態をリセットし、内部のエントリリストを再構築する。
func (s *TreeState) SetItems(items []*TreeItem) {
	s.items = items
	s.SelectedIndex = nil
	s.RightClickedIx = nil
	s.rebuildEntries()
}

// SetSelectedIndex は選択インデックスを設定する。
func (s *TreeState) SetSelectedIndex(ix *int) {
	s.SelectedIndex = ix
}

// SetSelectedItem は指定したアイテムを選択状態にする。
// 見つからない場合は祖先を展開し、見えるようにしてから選択する。
func (s *TreeState) SetSelectedItem(item *TreeItem) {
	if item != nil {
		ix := s.indexOf(item.ID())
		if ix == nil {
			s.expandAncestors(item.ID())
			ix = s.indexOf(item.ID())
		}
		s.SelectedIndex = ix
	} else {
		s.SelectedIndex = nil
	}
}

// SelectedItem は現在選択されている TreeItem を返す。
func (s *TreeState) SelectedItem() *TreeItem {
	if s.SelectedIndex != nil && *s.SelectedIndex < len(s.Entries) {
		return s.Entries[*s.SelectedIndex].Item
	}
	return nil
}

// indexOf は指定した ID を持つアイテムのインデックスを返す。
func (s *TreeState) indexOf(id string) *int {
	for i, entry := range s.Entries {
		if entry.Item.ID() == id {
			ix := i
			return &ix
		}
	}
	return nil
}

// expandAncestors はターゲットIDのアイテムが見えるように祖先をすべて展開する。
func (s *TreeState) expandAncestors(targetID string) {
	// Root elements are in s.items. 
	// We need to find the target item starting from roots.
	var ancestors []*TreeItem
	for _, item := range s.items {
		if path := item.Ancestors(targetID); path != nil {
			ancestors = path
			// path ends with the direct parent of targetID
			break
		}
	}

	if len(ancestors) == 0 {
		return
	}

	for i := len(ancestors) - 1; i >= 0; i-- {
		ancestor := ancestors[i]
		if !ancestor.IsExpanded() {
			ancestor.Expanded(true)
		}
	}
	s.rebuildEntries()
}

// ToggleExpand は指定されたインデックスのアイテムがフォルダであれば展開状態をトグルする。
func (s *TreeState) ToggleExpand(ix int) {
	if ix < 0 || ix >= len(s.Entries) {
		return
	}
	entry := s.Entries[ix]
	if !entry.Item.IsFolder() {
		return
	}

	expanded := !entry.Item.IsExpanded()
	entry.Item.Expanded(expanded)
	s.RightClickedIx = nil
	s.rebuildEntries()
}

func (s *TreeState) rebuildEntries() {
	s.Entries = make([]TreeEntry, 0)
	for _, item := range s.items {
		s.addEntry(item, 0)
	}
}

func (s *TreeState) addEntry(item *TreeItem, depth int) {
	s.Entries = append(s.Entries, NewTreeEntry(item, depth))
	if item.IsExpanded() {
		for _, child := range item.children {
			s.addEntry(child, depth+1)
		}
	}
}

