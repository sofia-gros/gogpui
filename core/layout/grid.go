package layout

import (
	"math"

	"github.com/sofiagros/gogpui/core/context"
)

// Grid はグリッドレイアウトコンテナを表す。
// GPUI の grid_cols / grid_rows / gap_x / gap_y API に 1:1 対応する。
type Grid struct {
	cols     int // 0 = Rows から自動計算
	rows     int // 0 = Cols から自動計算
	gapX     float64
	gapY     float64
	children []Widget
}

// NewGrid はデフォルト設定で新しい Grid コンテナを作成する。
func NewGrid() *Grid {
	return &Grid{}
}

// Cols は列数を設定する（GPUI: grid_cols）。
// Rows が未設定の場合、行数は子要素数と列数から自動計算される。
func (g *Grid) Cols(cols int) *Grid {
	g.cols = cols
	return g
}

// Rows は行数を設定する（GPUI: grid_rows）。
// Cols が未設定の場合、列数は子要素数と行数から自動計算される。
func (g *Grid) Rows(rows int) *Grid {
	g.rows = rows
	return g
}

// Gap は水平・垂直両方のギャップを同じ値に設定する。
func (g *Grid) Gap(gap float64) *Grid {
	g.gapX = gap
	g.gapY = gap
	return g
}

// GapX は列間（水平）のギャップを設定する（GPUI: gap_x）。
func (g *Grid) GapX(gap float64) *Grid {
	g.gapX = gap
	return g
}

// GapY は行間（垂直）のギャップを設定する（GPUI: gap_y）。
func (g *Grid) GapY(gap float64) *Grid {
	g.gapY = gap
	return g
}

// Add はグリッドに 1 つ以上の子 Widget を追加する。
func (g *Grid) Add(children ...Widget) *Grid {
	g.children = append(g.children, children...)
	return g
}

// resolveColsRows は cols と rows を子要素数から解決する。
// Cols だけ指定 → rows = ceil(n / cols)
// Rows だけ指定 → cols = ceil(n / rows)
// 両方指定     → そのまま使用
// 両方未指定   → cols = n（1行グリッド）
func (g *Grid) resolveColsRows(n int) (cols, rows int) {
	cols, rows = g.cols, g.rows
	if n == 0 {
		return 1, 1
	}
	switch {
	case cols > 0 && rows > 0:
		// 両方明示済み
	case cols > 0:
		rows = (n + cols - 1) / cols
	case rows > 0:
		cols = (n + rows - 1) / rows
	default:
		cols = n
		rows = 1
	}
	return cols, rows
}

// Render は Widget インターフェースを実装し、子 Widget をグリッド状に配置・描画する。
func (g *Grid) Render(uictx *context.UIContext, x, y float64) (w, h float64) {
	n := len(g.children)
	if n == 0 {
		return 0, 0
	}

	cols, rows := g.resolveColsRows(n)

	// --- 1. 計測パス ---
	originalMeasure := uictx.MeasureOnly
	uictx.MeasureOnly = true

	boxes := make([]childBox, n)
	for i, child := range g.children {
		cw, ch := child.Render(uictx, 0, 0)
		boxes[i] = childBox{cw, ch}
	}

	uictx.MeasureOnly = originalMeasure

	// 列幅: 各列に属するアイテムの幅の最大値
	colWidths := make([]float64, cols)
	// 行高: 各行に属するアイテムの高さの最大値
	rowHeights := make([]float64, rows)

	for i, box := range boxes {
		col := i % cols
		row := i / cols
		if row >= rows {
			break // Cols/Rows 指定を超える要素は無視
		}
		if box.w > colWidths[col] {
			colWidths[col] = box.w
		}
		if box.h > rowHeights[row] {
			rowHeights[row] = box.h
		}
	}

	// グリッド全体の幅・高さ
	var totalW, totalH float64
	for i, cw := range colWidths {
		totalW += cw
		if i > 0 {
			totalW += g.gapX
		}
	}
	for i, rh := range rowHeights {
		totalH += rh
		if i > 0 {
			totalH += g.gapY
		}
	}

	if uictx.MeasureOnly {
		return totalW, totalH
	}

	// --- 2. 列・行のオフセットを事前計算 ---
	colOffsets := make([]float64, cols)
	cur := 0.0
	for i, cw := range colWidths {
		colOffsets[i] = cur
		cur += cw + g.gapX
	}

	rowOffsets := make([]float64, rows)
	cur = 0.0
	for i, rh := range rowHeights {
		rowOffsets[i] = cur
		cur += rh + g.gapY
	}

	// --- 3. 描画パス ---
	for i, child := range g.children {
		col := i % cols
		row := i / cols
		if row >= rows {
			break
		}
		cx := x + colOffsets[col]
		cy := y + rowOffsets[row]
		rx, ry := math.Round(cx), math.Round(cy)
		box := boxes[i]
		if uictx.IsVisible(rx, ry, box.w, box.h) {
			child.Render(uictx, rx, ry)
		}
	}

	return totalW, totalH
}
