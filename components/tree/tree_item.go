package tree

// TreeItemState は TreeItem の展開・無効状態を保持する。
type TreeItemState struct {
	expanded bool
	disabled bool
}

// TreeItem は安定した ID、表示ラベル、子要素、および状態を持つツリーアイテムを表す。
type TreeItem struct {
	id       string
	label    string
	children []*TreeItem
	state    *TreeItemState
}

// NewItem は新しい TreeItem を作成する。
func NewItem(id, label string) *TreeItem {
	return &TreeItem{
		id:       id,
		label:    label,
		children: make([]*TreeItem, 0),
		state: &TreeItemState{
			expanded: false,
			disabled: false,
		},
	}
}

// Child は子要素を追加し、自身を返す。
func (i *TreeItem) Child(child *TreeItem) *TreeItem {
	i.children = append(i.children, child)
	return i
}

// Children は複数の子要素を追加し、自身を返す。
func (i *TreeItem) Children(children ...*TreeItem) *TreeItem {
	i.children = append(i.children, children...)
	return i
}

// Expanded は展開状態を設定し、自身を返す。
func (i *TreeItem) Expanded(expanded bool) *TreeItem {
	i.state.expanded = expanded
	return i
}

// Disabled は無効状態を設定し、自身を返す。
func (i *TreeItem) Disabled(disabled bool) *TreeItem {
	i.state.disabled = disabled
	return i
}

// IsFolder は子要素を持っているかどうかを返す。
func (i *TreeItem) IsFolder() bool {
	return len(i.children) > 0
}

// IsDisabled は無効状態かどうかを返す。
func (i *TreeItem) IsDisabled() bool {
	return i.state.disabled
}

// IsExpanded は展開状態かどうかを返す。
func (i *TreeItem) IsExpanded() bool {
	return i.state.expanded
}

// ID はアイテムの ID を返す。
func (i *TreeItem) ID() string {
	return i.id
}

// Label はアイテムのラベルを返す。
func (i *TreeItem) Label() string {
	return i.label
}

// Ancestors は自身をルートとして、指定した ID を持つアイテムの祖先パスを検索して返す。
// 返されるリストは、検索対象アイテムの直接の親から順にルートへ向かう祖先のリスト。
func (i *TreeItem) Ancestors(targetID string) []*TreeItem {
	if i.id == targetID {
		return []*TreeItem{}
	}

	for _, child := range i.children {
		if path := child.Ancestors(targetID); path != nil {
			return append(path, i)
		}
	}

	return nil
}

