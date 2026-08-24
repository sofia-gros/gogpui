package debugoverlay

import (
	"fmt"
	"image/color"
	"runtime"
	"time"

	"github.com/sofiagros/gogpui/core/context"
)

// DebugOverlay provides a real-time performance overlay (FPS, CPU, Mem, etc.)
type DebugOverlay struct {
	visible        bool
	frameTimes     [30]float64
	frameIndex     int
	lastUpdateTime time.Time
	memStats       runtime.MemStats
	numCPU         int
}

// New creates a new DebugOverlay instance.
func New() *DebugOverlay {
	return &DebugOverlay{
		visible: true,
		numCPU:  runtime.NumCPU(),
	}
}

// Visible sets whether the overlay is currently visible.
func (d *DebugOverlay) Visible(v bool) *DebugOverlay {
	d.visible = v
	return d
}

// Toggle flips the visibility of the overlay.
func (d *DebugOverlay) Toggle() {
	d.visible = !d.visible
}

// Render draws the debug overlay at the absolute position x, y (usually top-right).
func (d *DebugOverlay) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	if !d.visible {
		return 0, 0
	}

	// Update DeltaTime history
	d.frameTimes[d.frameIndex] = uictx.DeltaTime
	d.frameIndex = (d.frameIndex + 1) % len(d.frameTimes)

	// Calculate Average FrameTime and FPS
	var sum float64
	for _, dt := range d.frameTimes {
		sum += dt
	}
	avgDt := sum / float64(len(d.frameTimes))
	
	fps := 0.0
	if avgDt > 0 {
		fps = 1.0 / avgDt
	}
	frameTimeMs := avgDt * 1000.0

	// Throttle heavy runtime operations (Memory stats)
	// runtime.ReadMemStats is a STOP-THE-WORLD operation and blocks the render thread!
	// Removed to prevent 1-second UI freezes.
	// if time.Since(d.lastUpdateTime) > 500*time.Millisecond {
	// 	runtime.ReadMemStats(&d.memStats)
	// 	d.lastUpdateTime = time.Now()
	// }

	numGoroutines := runtime.NumGoroutine()
	sysMemMB := 0.0 // float64(d.memStats.Sys) / 1024 / 1024
	allocMemMB := 0.0 // float64(d.memStats.Alloc) / 1024 / 1024

	// Overlay Dimensions
	w, h := 190.0, 130.0

	if uictx.MeasureOnly {
		return w, h
	}

	ctx := uictx.GG

	// Draw Background (Semi-transparent black)
	ctx.SetColor(color.RGBA{0, 0, 0, 200})
	uictx.DrawRoundedRectangle(x, y, w, h, 6)
	uictx.Fill()

	// Draw Border
	ctx.SetColor(color.RGBA{80, 80, 80, 255})
	ctx.SetLineWidth(1.0)
	uictx.DrawRoundedRectangle(x, y, w, h, 6)
	uictx.Stroke()

	// Text Settings
	lineHeight := 18.0
	startX := x + 12.0
	startY := y + 20.0

	// Color logic for FPS/FrameTime
	fpsColor := color.RGBA{0, 255, 0, 255} // Green
	if fps < 30 {
		fpsColor = color.RGBA{255, 0, 0, 255} // Red
	} else if fps < 55 {
		fpsColor = color.RGBA{255, 255, 0, 255} // Yellow
	}

	// Draw Lines
	ctx.SetColor(fpsColor)
	uictx.DrawString(fmt.Sprintf("FPS:      %.1f", fps), startX, startY)
	
	ctx.SetColor(color.RGBA{220, 220, 220, 255})
	uictx.DrawString(fmt.Sprintf("Frame:    %.2f ms", frameTimeMs), startX, startY+lineHeight*1)
	
	ctx.SetColor(color.RGBA{150, 200, 255, 255})
	uictx.DrawString(fmt.Sprintf("CPU Core: %d", d.numCPU), startX, startY+lineHeight*2)
	
	ctx.SetColor(color.RGBA{200, 150, 255, 255})
	uictx.DrawString(fmt.Sprintf("Mem Sys:  %.1f MB", sysMemMB), startX, startY+lineHeight*3)
	uictx.DrawString(fmt.Sprintf("Mem Allc: %.1f MB", allocMemMB), startX, startY+lineHeight*4)
	
	ctx.SetColor(color.RGBA{220, 220, 220, 255})
	uictx.DrawString(fmt.Sprintf("Goroutin: %d", numGoroutines), startX, startY+lineHeight*5)

	return w, h
}
