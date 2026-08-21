package layout

import (
	"math"

	"github.com/sofiagros/gogpui/core/context"
)

// Direction はフレックスレイアウトの主軸方向を定義する。
type Direction int

const (
	// Row は水平方向（左から右）を示す。
	Row Direction = iota
	// Column は垂直方向（上から下）を示す。
	Column
)

// JustifyContent は主軸方向のアイテム配置方法を定義する。
type JustifyContent int

const (
	// JustifyStart はアイテムを主軸の先頭に配置する。
	JustifyStart JustifyContent = iota
	// JustifyCenter はアイテムを主軸の中央に配置する（WithConstraints 必須）。
	JustifyCenter
	// JustifyEnd はアイテムを主軸の末尾に配置する（WithConstraints 必須）。
	JustifyEnd
	// JustifySpaceBetween はアイテム間に均等なスペースを配置する（WithConstraints 必須）。
	JustifySpaceBetween
	// JustifySpaceAround はアイテムの周囲に均等なスペースを配置する（WithConstraints 必須）。
	JustifySpaceAround
)

// AlignItems は交差軸方向のアイテム整列方法を定義する。
type AlignItems int

const (
	// AlignStart はアイテムを交差軸の先頭に整列する。
	AlignStart AlignItems = iota
	// AlignCenter はアイテムを交差軸の中央に整列する。
	AlignCenter
	// AlignEnd はアイテムを交差軸の末尾に整列する。
	AlignEnd
	// AlignStretch はアイテムを交差軸方向に引き伸ばす（外部サイズ制約が必要）。
	AlignStretch
)

// WrapMode はフレックスアイテムの折り返しモードを定義する。
type WrapMode int

const (
	// NoWrap は折り返しを行わない（デフォルト）。
	NoWrap WrapMode = iota
	// Wrap は主軸方向にアイテムが収まらない場合、前方向へ折り返す。
	Wrap
	// WrapReverse は主軸方向にアイテムが収まらない場合、逆方向へ折り返す。
	WrapReverse
)

// Widget はすべての UI コンポーネントの標準インターフェースである。
type Widget interface {
	Render(uictx *context.UIContext, x, y float64) (w, h float64)
}

// GrowWidget は Flex 内で残りスペースを受け取れる Widget のインターフェース。
// GPUI の div().flex_grow() に相当する。
// Flex は Render を呼ぶ前に SetAllocatedSize で割り当てサイズを通知する。
type GrowWidget interface {
	Widget
	// GrowFactor は CSS flex-grow に相当する係数を返す（通常 1.0）。
	GrowFactor() float64
	// SetAllocatedSize は Flex から割り当てられた (w, h) を受け取る。
	// Flex は Render 呼び出しの直前にこのメソッドを呼ぶ。
	SetAllocatedSize(w, h float64)
}

// Flex はフレックスボックスレイアウトコンテナを表す。
type Flex struct {
	dir      Direction
	justify  JustifyContent
	align    AlignItems
	gap      float64
	children []Widget
	wrapMode WrapMode
	maxW     float64 // 0 = 制約なし
	maxH     float64 // 0 = 制約なし
}

// NewFlex はデフォルト設定（Row 方向）で新しい Flex コンテナを作成する。
func NewFlex() *Flex {
	return &Flex{
		dir:      Row,
		justify:  JustifyStart,
		align:    AlignStart,
		gap:      0,
		children: []Widget{},
		wrapMode: NoWrap,
		maxW:     0,
		maxH:     0,
	}
}

// Direction は主軸方向（Row または Column）を設定する。
func (f *Flex) Direction(dir Direction) *Flex {
	f.dir = dir
	return f
}

// Justify は主軸方向の配置方法を設定する。
func (f *Flex) Justify(j JustifyContent) *Flex {
	f.justify = j
	return f
}

// Align は交差軸方向の整列方法を設定する。
func (f *Flex) Align(a AlignItems) *Flex {
	f.align = a
	return f
}

// Gap はアイテム間のスペースを設定する。
func (f *Flex) Gap(gap float64) *Flex {
	f.gap = gap
	return f
}

// WrapContent は折り返しモードを設定する。
func (f *Flex) WrapContent(mode WrapMode) *Flex {
	f.wrapMode = mode
	return f
}

// WithConstraints は親コンテナの最大幅・最大高さを設定する。
// JustifyCenter / JustifyEnd / JustifySpaceBetween / JustifySpaceAround を使う場合は
// このメソッドで制約を指定すること。0 は制約なしを意味する。
func (f *Flex) WithConstraints(maxW, maxH float64) *Flex {
	f.maxW = maxW
	f.maxH = maxH
	return f
}

// Add はレイアウトに 1 つ以上の子 Widget を追加する。
func (f *Flex) Add(children ...Widget) *Flex {
	f.children = append(f.children, children...)
	return f
}

// Render は Widget インターフェースを実装し、子 Widget を配置・描画する。
func (f *Flex) Render(uictx *context.UIContext, x, y float64) (w, h float64) {
	if len(f.children) == 0 {
		return 0, 0
	}

	if f.wrapMode != NoWrap {
		return f.renderWrap(uictx, x, y)
	}
	return f.renderNoWrap(uictx, x, y)
}

// childBox は計測パスで得られた各子ウィジェットのサイズを保持する。
type childBox struct {
	w, h float64
}

// renderNoWrap は折り返しなしのフレックスレイアウトを描画する。
func (f *Flex) renderNoWrap(uictx *context.UIContext, x, y float64) (w, h float64) {
	// --- 1. 計測パス（GrowWidget は 0 として扱う）---
	originalMeasure := uictx.MeasureOnly
	uictx.MeasureOnly = true

	n := len(f.children)
	boxes := make([]childBox, n)
	var nonGrowMainTotal, crossMax float64

	for i, child := range f.children {
		// GrowWidget は計測時に 0 を返すことを期待するが、念のため 0 に固定する
		if _, ok := child.(GrowWidget); ok {
			boxes[i] = childBox{0, 0}
			continue
		}
		cw, ch := child.Render(uictx, 0, 0)
		boxes[i] = childBox{cw, ch}
		var main, cross float64
		if f.dir == Row {
			main, cross = cw, ch
		} else {
			main, cross = ch, cw
		}
		nonGrowMainTotal += main
		if cross > crossMax {
			crossMax = cross
		}
	}
	uictx.MeasureOnly = originalMeasure

	// --- 2. GrowWidget への残りスペース配分 ---
	var maxMain float64
	if f.dir == Row {
		maxMain = f.maxW
	} else {
		maxMain = f.maxH
	}

	// grow が有効なのは maxMain が制約されているときのみ
	if maxMain > 0 {
		gaps := f.gap * float64(n-1)
		remaining := maxMain - nonGrowMainTotal - gaps
		if remaining < 0 {
			remaining = 0
		}

		// grow 係数の合計を計算
		var totalFactor float64
		for _, child := range f.children {
			if gw, ok := child.(GrowWidget); ok {
				totalFactor += gw.GrowFactor()
			}
		}

		// 各 GrowWidget に割り当てサイズをセットし、boxes を更新
		if totalFactor > 0 {
			for i, child := range f.children {
				gw, ok := child.(GrowWidget)
				if !ok {
					continue
				}
				allocMain := remaining * gw.GrowFactor() / totalFactor
				var aw, ah float64
				if f.dir == Row {
					aw, ah = allocMain, 0
				} else {
					aw, ah = 0, allocMain
				}
				gw.SetAllocatedSize(aw, ah)
				boxes[i] = childBox{aw, ah}
			}
		}
	}

	// grow 後に交差軸の最大値を再計算
	for _, box := range boxes {
		var cross float64
		if f.dir == Row {
			cross = box.h
		} else {
			cross = box.w
		}
		if cross > crossMax {
			crossMax = cross
		}
	}

	// --- 3. 全体サイズの集計 ---
	var rawMainTotal float64
	for _, box := range boxes {
		if f.dir == Row {
			rawMainTotal += box.w
		} else {
			rawMainTotal += box.h
		}
	}
	mainTotal := rawMainTotal
	if n > 1 {
		mainTotal += f.gap * float64(n-1)
	}

	if f.dir == Row {
		w, h = mainTotal, crossMax
	} else {
		w, h = crossMax, mainTotal
	}

	// 親が計測中であれば早期リターン
	if uictx.MeasureOnly {
		return w, h
	}

	// maxMain の確定（grow 配分後の mainTotal を使う）
	if maxMain <= 0 {
		maxMain = mainTotal
	}

	// --- 4. JustifyContent によるカーソル・スペーシング計算 ---
	var cursor, spacing float64
	spacing = f.gap

	switch f.justify {
	case JustifyStart:
		cursor = 0
	case JustifyEnd:
		cursor = maxMain - mainTotal
		if cursor < 0 {
			cursor = 0
		}
	case JustifyCenter:
		cursor = (maxMain - mainTotal) / 2.0
		if cursor < 0 {
			cursor = 0
		}
	case JustifySpaceBetween:
		cursor = 0
		if n > 1 {
			spacing = (maxMain - rawMainTotal) / float64(n-1)
			if spacing < 0 {
				spacing = 0
			}
		}
	case JustifySpaceAround:
		if n > 0 {
			unit := (maxMain - rawMainTotal) / float64(n)
			if unit < 0 {
				unit = 0
			}
			spacing = unit
			cursor = unit / 2.0
		}
	}

	// --- 5. 描画パス ---
	for i, child := range f.children {
		box := boxes[i]
		var cx, cy float64

		if f.dir == Row {
			cx = x + cursor
			switch f.align {
			case AlignStart, AlignStretch:
				cy = y
			case AlignEnd:
				cy = y + (crossMax - box.h)
			case AlignCenter:
				cy = y + (crossMax-box.h)/2.0
			}
			cursor += box.w + spacing
		} else {
			cy = y + cursor
			switch f.align {
			case AlignStart, AlignStretch:
				cx = x
			case AlignEnd:
				cx = x + (crossMax - box.w)
			case AlignCenter:
				cx = x + (crossMax-box.w)/2.0
			}
			cursor += box.h + spacing
		}

		child.Render(uictx, math.Round(cx), math.Round(cy))
	}

	return w, h
}

// lineInfo は Wrap レイアウト時の各行（または列）の情報を保持する。
type lineInfo struct {
	start, end int     // children スライスのインデックス範囲 [start, end)
	mainSize   float64 // 行の主軸合計サイズ（gap込み）
	crossSize  float64 // 行の交差軸最大サイズ
}

// renderWrap は折り返しありのフレックスレイアウトを描画する。
func (f *Flex) renderWrap(uictx *context.UIContext, x, y float64) (w, h float64) {
	if len(f.children) == 0 {
		return 0, 0
	}

	// 主軸の制約サイズを決定する
	var maxMain float64
	if f.dir == Row {
		maxMain = f.maxW
	} else {
		maxMain = f.maxH
	}

	// --- 1. 計測パス ---
	originalMeasure := uictx.MeasureOnly
	uictx.MeasureOnly = true

	allBoxes := make([]childBox, len(f.children))
	for i, child := range f.children {
		cw, ch := child.Render(uictx, 0, 0)
		allBoxes[i] = childBox{cw, ch}
	}
	uictx.MeasureOnly = originalMeasure

	// --- 2. 行（または列）に分割する ---
	var lines []lineInfo
	lineStart := 0
	var lineMain, lineCross float64

	for i, box := range allBoxes {
		var itemMain, itemCross float64
		if f.dir == Row {
			itemMain, itemCross = box.w, box.h
		} else {
			itemMain, itemCross = box.h, box.w
		}

		// 行に追加したときの主軸合計
		newMain := lineMain
		if i > lineStart {
			newMain += f.gap
		}
		newMain += itemMain

		// 制約を超えるかつ行が空でない場合は改行
		if maxMain > 0 && newMain > maxMain && i > lineStart {
			lines = append(lines, lineInfo{
				start:     lineStart,
				end:       i,
				mainSize:  lineMain,
				crossSize: lineCross,
			})
			lineStart = i
			lineMain = itemMain
			lineCross = itemCross
		} else {
			lineMain = newMain
			if itemCross > lineCross {
				lineCross = itemCross
			}
		}
	}
	// 最後の行を追加
	lines = append(lines, lineInfo{
		start:     lineStart,
		end:       len(f.children),
		mainSize:  lineMain,
		crossSize: lineCross,
	})

	// WrapReverse の場合は行順序を逆にする
	if f.wrapMode == WrapReverse {
		for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
			lines[i], lines[j] = lines[j], lines[i]
		}
	}

	// --- 3. 全体サイズを計算 ---
	var totalCross, totalMainMax float64
	for i, line := range lines {
		totalCross += line.crossSize
		if i > 0 {
			totalCross += f.gap
		}
		if line.mainSize > totalMainMax {
			totalMainMax = line.mainSize
		}
	}

	if f.dir == Row {
		w, h = totalMainMax, totalCross
	} else {
		w, h = totalCross, totalMainMax
	}

	if uictx.MeasureOnly {
		return w, h
	}

	// --- 4. 各行を描画する ---
	crossCursor := 0.0
	for _, line := range lines {
		rowChildren := f.children[line.start:line.end]
		sub := &Flex{
			dir:      f.dir,
			justify:  f.justify,
			align:    f.align,
			gap:      f.gap,
			children: rowChildren,
			wrapMode: NoWrap,
		}
		if f.dir == Row {
			sub.maxW = maxMain
			sub.maxH = line.crossSize
		} else {
			sub.maxW = line.crossSize
			sub.maxH = maxMain
		}

		var lx, ly float64
		if f.dir == Row {
			lx = x
			ly = y + crossCursor
		} else {
			lx = x + crossCursor
			ly = y
		}

		sub.renderNoWrap(uictx, lx, ly)
		crossCursor += line.crossSize + f.gap
	}

	return w, h
}
