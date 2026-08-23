package breadcrumb

import (
	"image/color"

	"github.com/sofiagros/gogpui/core/context"
	"github.com/gogpu/gg"
)

// BreadcrumbItem はパンくずリストの各項目を表す。
type BreadcrumbItem struct {
	id       string
	label    string
	disabled bool
	onClick  func()
}

// NewItem は新しい BreadcrumbItem を作成する。
func NewItem(id, label string) *BreadcrumbItem {
	return &BreadcrumbItem{
		id:    id,
		label: label,
	}
}

// Disabled は項目の有効/無効を設定する。
func (i *BreadcrumbItem) Disabled(d bool) *BreadcrumbItem {
	i.disabled = d
	return i
}

// OnClick はクリック時のコールバックを設定する。
func (i *BreadcrumbItem) OnClick(f func()) *BreadcrumbItem {
	i.onClick = f
	return i
}

// Breadcrumb はナビゲーション用のパンくずリストコンポーネント。
type Breadcrumb struct {
	items []*BreadcrumbItem
}

// New は新しい Breadcrumb を作成する。
func New() *Breadcrumb {
	return &Breadcrumb{
		items: make([]*BreadcrumbItem, 0),
	}
}

// Add はアイテムを追加する。
func (b *Breadcrumb) Add(items ...*BreadcrumbItem) *Breadcrumb {
	b.items = append(b.items, items...)
	return b
}

func drawChevronRight(ctx *gg.Context, x, y, size float64, c color.Color) {
	ctx.SetColor(c)
	ctx.SetLineWidth(1.5)
	
	halfW := size * 0.25
	halfH := size * 0.4

	ctx.MoveTo(x-halfW, y-halfH)
	ctx.LineTo(x+halfW, y)
	ctx.LineTo(x-halfW, y+halfH)
	ctx.Stroke()
}

// Render は指定位置に Breadcrumb を描画し、占有サイズ (幅, 高さ) を返す。
func (b *Breadcrumb) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	ctx := uictx.GG
	th := uictx.Theme

	gap := 6.0
	iconSize := 14.0 // size_3p5 in gpui is roughly 14px

	var totalW float64
	var totalH float64

	// 計測パスと描画パスを統合するための計算
	for i, item := range b.items {
		// Label measurement
		lblW, lblH := ctx.MeasureString(item.label)
		
		if lblH > totalH {
			totalH = lblH
		}
		if iconSize > totalH {
			totalH = iconSize
		}

		totalW += lblW

		// セパレーターの幅を追加（最後の要素以外）
		if i < len(b.items)-1 {
			totalW += gap + iconSize + gap
		}
	}

	if uictx.MeasureOnly {
		return totalW, totalH
	}

	currentX := x
	centerY := y + totalH/2.0

	for i, item := range b.items {
		isLast := (i == len(b.items)-1)
		
		lblW, _ := ctx.MeasureString(item.label)

		// インタラクション処理
		// isLast の場合でも onClick があればクリック可能とするのが一般的だが、
		// Rust版は常に is_last でテキストカラーが変わる。
		isClickable := !item.disabled && item.onClick != nil
		if isClickable {
			_, _, isClicked := uictx.ProcessInteraction(item.id, currentX, y, lblW, totalH, item.disabled)
			if isClicked {
				item.onClick()
			}
		}

		// テキスト色の決定
		// Rust版の優先度: Default -> muted, isLast -> foreground, disabled -> muted
		textColor := th.Colors.MutedForeground
		if isLast {
			textColor = th.Colors.Foreground
		}
		if item.disabled {
			textColor = th.Colors.MutedForeground
		}
		
		// ホバー時の視覚的フィードバック (gogpu 独自の UX 向上のため軽く追加)
		if isClickable {
			state := uictx.GetState(item.id)
			if state.HoverRatio > 0 && !item.disabled {
				// ホバー時は少し明るく/濃くするなどの処理が可能だが、今回はRust版に忠実に従うため
				// 必須ではない。カーソルポインターは gogpu 内部で設定される。
			}
		}

		// テキストの描画
		ctx.SetColor(textColor)
		// y座標はテキストのベースラインではなく中心(Anchored)として描画
		ctx.DrawStringAnchored(item.label, currentX, centerY, 0.0, 0.5)
		currentX += lblW

		// セパレーターの描画
		if !isLast {
			currentX += gap
			drawChevronRight(ctx, currentX+iconSize/2.0, centerY, iconSize, th.Colors.MutedForeground)
			currentX += iconSize + gap
		}
	}

	// Wait, if lblH is used, make sure totalH doesn't shrink to 0 if there are no items
	if len(b.items) == 0 {
		return 0, 0
	}

	return totalW, totalH
}
