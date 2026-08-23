package collapsible

import (
	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/layout"
)

// Collapsible は開閉可能なコンテナコンポーネント。
type Collapsible struct {
	open    bool
	trigger layout.Widget
	content layout.Widget
}

// New は新しい Collapsible を作成する。
func New() *Collapsible {
	return &Collapsible{
		open: false,
	}
}

// Open はコンテンツの開閉状態を設定する。
func (c *Collapsible) Open(open bool) *Collapsible {
	c.open = open
	return c
}

// Trigger は開閉のトリガーとなる要素を設定する。
func (c *Collapsible) Trigger(w layout.Widget) *Collapsible {
	c.trigger = w
	return c
}

// Content は展開時に表示される要素を設定する。
func (c *Collapsible) Content(w layout.Widget) *Collapsible {
	c.content = w
	return c
}

// Render は指定位置に Collapsible を描画し、占有サイズ (幅, 高さ) を返す。
func (c *Collapsible) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	// トリガーが未設定の場合は何も描画しない
	if c.trigger == nil {
		return 0, 0
	}

	flex := layout.NewFlex().Direction(layout.Column)
	flex.Add(c.trigger)

	if c.open && c.content != nil {
		flex.Add(c.content)
	}

	return flex.Render(uictx, x, y)
}
