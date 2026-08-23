package scroll

import (
	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/layout"
)

// Scroll はコンテンツをY軸方向にスクロール可能にするコンテナです。
type Scroll struct {
	id     string
	width  float64
	height float64
	child  layout.Widget
}

// New は新しい Scroll コンポーネントを作成します。
func New(id string) *Scroll {
	return &Scroll{
		id:     id,
		width:  0, // 0 means expand or unspecified
		height: 0,
	}
}

// Child はスクロール領域内のコンテンツを設定します。
func (s *Scroll) Child(w layout.Widget) *Scroll {
	s.child = w
	return s
}

// Size はスクロールビューアの表示サイズを固定で指定します。
func (s *Scroll) Size(w, h float64) *Scroll {
	s.width = w
	s.height = h
	return s
}

// GrowFactor は layout.GrowWidget インターフェースを満たし、
// サイズが未指定の場合は親コンテナの中で伸縮（flex-1）するようにします。
func (s *Scroll) GrowFactor() float64 { return 1.0 }

// SetAllocatedSize は Flex 等から割り当てられた領域を受け取ります。
func (s *Scroll) SetAllocatedSize(w, h float64) {
	s.width = w
	s.height = h
}

// Render は Scroll コンテナを描画します。
func (s *Scroll) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	if s.child == nil {
		return s.width, s.height
	}

	state := uictx.GetState(s.id)

	// Measure step: find the total height of the content.
	// We want to give the child the same width, but infinite height to measure.
	// We temporarily disable rendering and constraints on height.
	originalMeasure := uictx.MeasureOnly
	uictx.MeasureOnly = true
	
	// Set the context width if necessary (some widgets might need it)
	// But our interface just returns dimensions.
	_, contentH := s.child.Render(uictx, x, y)
	
	uictx.MeasureOnly = originalMeasure

	if uictx.MeasureOnly {
		return s.width, s.height
	}

	// Calculate max offset
	maxOffset := contentH - s.height
	if maxOffset < 0 {
		maxOffset = 0
	}

	// Handle scroll events
	// 50.0 is a scroll multiplier (speed)
	isHovered := uictx.HitTest(x, y, s.width, s.height)
	if isHovered && uictx.Mouse.ScrollY != 0 {
		state.OffsetY -= uictx.Mouse.ScrollY * 50.0
		uictx.NeedsRedraw = true
	}

	// Clamp offset
	if state.OffsetY < 0 {
		state.OffsetY = 0
	}
	if state.OffsetY > maxOffset {
		state.OffsetY = maxOffset
	}

	ctx := uictx.GG

	// Draw clipping region
	ctx.Push()
	ctx.DrawRectangle(x, y, s.width, s.height)
	ctx.Clip()

	// Translate context and adjust mouse for hit testing
	ctx.Translate(0, -state.OffsetY)
	uictx.Mouse.Y += state.OffsetY

	// Draw child
	s.child.Render(uictx, x, y)

	// Restore mouse and context
	uictx.Mouse.Y -= state.OffsetY
	ctx.Pop()

	// Draw scrollbar if content exceeds height
	if maxOffset > 0 {
		barW := 6.0
		barX := x + s.width - barW - 2.0
		
		// Ratio of visible height to total height
		ratio := s.height / contentH
		thumbH := s.height * ratio
		if thumbH < 20 {
			thumbH = 20
		}
		
		// Position thumb
		availH := s.height - thumbH
		scrollRatio := state.OffsetY / maxOffset
		thumbY := y + (availH * scrollRatio)

		ctx.SetColor(uictx.Theme.Colors.Border)
		ctx.DrawRoundedRectangle(barX, thumbY, barW, thumbH, 3)
		_ = ctx.Fill()
	}

	return s.width, s.height
}
