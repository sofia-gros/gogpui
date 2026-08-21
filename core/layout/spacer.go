package layout

import "github.com/sofiagros/gogpui/core/context"

// SpacerFixed は固定サイズのスペーサー Widget。
// Flex 内で一定の余白を確保するために使用する。
// GPUI の div().w(px(N)) や div().h(px(N)) に相当する。
type SpacerFixed struct {
	w, h float64
}

// SpacerW は幅 w の水平スペーサーを返す。
// Row Flex の中で水平方向の余白として使用する。
func SpacerW(w float64) *SpacerFixed {
	return &SpacerFixed{w: w, h: 0}
}

// SpacerH は高さ h の垂直スペーサーを返す。
// Column Flex の中で垂直方向の余白として使用する。
func SpacerH(h float64) *SpacerFixed {
	return &SpacerFixed{w: 0, h: h}
}

// SpacerFixed は幅 w・高さ h を両方指定した固定サイズスペーサーを返す。
func NewSpacerFixed(w, h float64) *SpacerFixed {
	return &SpacerFixed{w: w, h: h}
}

// Render は Widget インターフェースを実装し、固定サイズを返す。
func (s *SpacerFixed) Render(_ *context.UIContext, _, _ float64) (float64, float64) {
	return s.w, s.h
}

// --- SpacerFlex ---

// SpacerFlex は Flex 内で残りスペースをすべて埋める GrowWidget。
// CSS の flex-grow: 1 / GPUI の div().flex_grow() に相当する。
// WithConstraints で親サイズを指定した Flex の中でのみ有効に機能する。
type SpacerFlex struct {
	factor float64
	allocW float64
	allocH float64
}

// NewSpacerFlex はデフォルトの grow factor (1.0) で SpacerFlex を返す。
func NewSpacerFlex() *SpacerFlex {
	return &SpacerFlex{factor: 1.0}
}

// WithGrowFactor は grow 係数を変更した SpacerFlex を返す。
// 複数の SpacerFlex を同一 Flex に配置する場合、係数に比例してスペースが配分される。
func (s *SpacerFlex) WithGrowFactor(factor float64) *SpacerFlex {
	s.factor = factor
	return s
}

// GrowFactor は GrowWidget インターフェースを実装し、grow 係数を返す。
func (s *SpacerFlex) GrowFactor() float64 {
	return s.factor
}

// SetAllocatedSize は GrowWidget インターフェースを実装し、
// Flex から割り当てられたサイズを受け取る。
func (s *SpacerFlex) SetAllocatedSize(w, h float64) {
	s.allocW = w
	s.allocH = h
}

// Render は Widget インターフェースを実装する。
// Flex から SetAllocatedSize が呼ばれた後は割り当てサイズを返す。
// 呼ばれていない場合は (0, 0) を返す（計測パスで利用される）。
func (s *SpacerFlex) Render(_ *context.UIContext, _, _ float64) (float64, float64) {
	return s.allocW, s.allocH
}
