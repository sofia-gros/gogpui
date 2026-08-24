package layout

import (
	"image"

	"github.com/gogpu/gg"
	"github.com/sofiagros/gogpui/core/context"
)

// CachedWidget は子コンポーネントの描画結果をオフスクリーンバッファにキャッシュするウィジェットです。
// 内容が静的（アニメーションや頻繁なホバー・クリック状態変化がない）な巨大ツリーを
// 包むことで、毎フレームの再パス描画を回避し劇的なパフォーマンス向上をもたらします。
type CachedWidget struct {
	id         string
	child      Widget
}

// CachedWidgetState holds the actual cache across frames
type CachedWidgetState struct {
	CachedCtx  *gg.Context
	CachedImg  *gg.ImageBuf
	LastW      float64
	LastH      float64
	IsDirty    bool
}

// NewCachedWidget creates a new CachedWidget wrapping the given child.
func NewCachedWidget(id string, child Widget) *CachedWidget {
	return &CachedWidget{
		id:    id,
		child: child,
	}
}

// MarkDirty marks the cache as invalid, forcing a redraw of the child on the next frame.
// Note: This needs UIContext if called externally, but internally it's handled.

func (w *CachedWidget) Render(uictx *context.UIContext, startX, startY float64) (float64, float64) {
	if uictx.MeasureOnly {
		return w.child.Render(uictx, startX, startY)
	}

	state := uictx.GetState(w.id)
	var cache *CachedWidgetState
	if state.CustomData == nil {
		cache = &CachedWidgetState{IsDirty: true}
		state.CustomData = cache
	} else {
		cache = state.CustomData.(*CachedWidgetState)
	}

	// 1. Measure the child to know how big the cache needs to be
	uictx.MeasureOnly = true
	childW, childH := w.child.Render(uictx, 0, 0)
	uictx.MeasureOnly = false

	if childW <= 0 || childH <= 0 {
		return 0, 0
	}

	// 2. Re-create or resize the cache context if necessary
	if cache.CachedCtx == nil || cache.LastW != childW || cache.LastH != childH {
		cache.CachedCtx = gg.NewContext(int(childW), int(childH))
		cache.LastW = childW
		cache.LastH = childH
		cache.IsDirty = true
	}

	// 2.5 Check if any damage rects intersect with this widget
	if !cache.IsDirty && len(uictx.DamageRects) > 0 {
		widgetRect := image.Rect(int(startX), int(startY), int(startX+childW), int(startY+childH))
		for _, dmg := range uictx.DamageRects {
			if widgetRect.Overlaps(dmg) {
				cache.IsDirty = true
				break
			}
		}
	}

	if uictx.LogicOnly {
		// Logic pass: just run the child's logic and return
		originalMouse := uictx.Mouse
		uictx.Mouse.X -= startX
		uictx.Mouse.Y -= startY

		w.child.Render(uictx, 0, 0)

		uictx.Mouse = originalMouse
		return childW, childH
	}

	// 3. Redraw the child to the cached context if dirty
	if cache.IsDirty {
		// Temporarily swap the UIContext's GG context and offset mouse coordinates
		originalGG := uictx.GG
		uictx.GG = cache.CachedCtx

		originalMouse := uictx.Mouse
		uictx.Mouse.X -= startX
		uictx.Mouse.Y -= startY

		// Clear the cache background (transparent by default)
		cache.CachedCtx.Clear() 
		
		// Render the child at (0, 0) inside the cache
		w.child.Render(uictx, 0, 0)

		// Restore original context and mouse
		uictx.GG = originalGG
		uictx.Mouse = originalMouse
		
		// Update the image buffer
		cache.CachedImg = gg.ImageBufFromImage(cache.CachedCtx.Image())
		cache.IsDirty = false
	}

	// 4. Draw the cached image onto the main context
	if cache.CachedImg != nil {
		uictx.GG.DrawImage(cache.CachedImg, startX, startY)
	}

	return childW, childH
}
