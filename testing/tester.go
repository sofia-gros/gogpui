package testing

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/gogpu/gg"
	"github.com/gogpu/gg/text"
	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/theme"
)

// Tester provides a mock environment for UI component testing and golden image diffs.
type Tester struct {
	Width  int
	Height int
	Ctx    *gg.Context
	UI     *context.UIContext
}

// NewTester creates a new tester instance.
func NewTester() *Tester {
	// Enable subpixel rendering globally BEFORE creating context
	text.SetGlobalSubpixelCache(text.NewSubpixelCache(text.HighQualitySubpixelConfig()))

	w, h := 800, 600
	ctx := gg.NewContext(w, h)
	ctx.SetColor(color.White)
	ctx.Clear()

	// Find the project root by looking for go.mod
	rootDir := "."
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(filepath.Join(rootDir, "go.mod")); err == nil {
			break
		}
		rootDir = filepath.Join(rootDir, "..")
	}
	fontPath := filepath.Join(rootDir, "assets", "fonts", "Inter-Regular.ttf")

	// Load font to enable text rendering
	source, err := text.NewFontSourceFromFile(fontPath)
	if err == nil {
		// Disable hinting to prevent advance width integer rounding, enable kern/liga
		ctx.SetFont(source.Face(16.0, 
			text.WithHinting(text.HintingNone),
			text.WithFeatures(text.NewFontFeature("kern", 1), text.NewFontFeature("liga", 1)),
		))
	} else {
		fmt.Printf("Warning: Failed to load font from %s: %v\n", fontPath, err)
	}

	// UI テキスト向けのベクターモードを使用する。
	// TextModeMSDF はゲーム向けで、UI ラベルには TextModeVector が推奨（gg ドキュメント参照）。
	ctx.SetTextMode(gg.TextModeVector)

	uictx := context.NewUIContext()
	// windowW/windowH はテスト用キャンバスの論理サイズを渡す。
	uictx.Update(ctx, theme.DefaultTheme(), 0.016, nil, 1.0, float64(w), float64(h))

	return &Tester{
		Width:  w,
		Height: h,
		Ctx:    ctx,
		UI:     uictx,
	}
}

// RenderComponent is an interface for components that can be rendered.
type RenderComponent interface {
	Render(ctx *context.UIContext, x, y float64) (w, h float64)
}

// AssertGoldenImage renders the component to a PNG and compares it against a golden image.
func (t *Tester) AssertGoldenImage(component RenderComponent, componentName string, stateName string) (score float64, err error) {
	// 1. Render the component in an off-screen buffer
	t.Ctx.SetColor(color.White)
	t.Ctx.Clear()
	component.Render(t.UI, 10, 10)

	// Ensure the golden directory exists
	goldenDir := filepath.Join("testdata", "golden")
	if err := os.MkdirAll(goldenDir, 0755); err != nil {
		return 0, err
	}

	actualPath := filepath.Join(goldenDir, fmt.Sprintf("%s_%s_actual.png", componentName, stateName))
	goldenPath := filepath.Join(goldenDir, fmt.Sprintf("%s_%s_golden.png", componentName, stateName))

	// 2. Save the rendered image
	if err := t.Ctx.SavePNG(actualPath); err != nil {
		return 0, fmt.Errorf("failed to save actual image: %w", err)
	}

	// 3. Compare with golden image if it exists
	goldenFile, err := os.Open(goldenPath)
	if err != nil {
		// If golden image doesn't exist, we use the actual as golden for the first run
		fmt.Printf("Golden image not found for %s_%s, using actual as golden.\n", componentName, stateName)
		err = os.Rename(actualPath, goldenPath)
		if err != nil {
			return 0, fmt.Errorf("failed to create golden image: %w", err)
		}
		return 0.0, nil
	}
	defer goldenFile.Close()

	goldenImg, _, err := image.Decode(goldenFile)
	if err != nil {
		return 0, fmt.Errorf("failed to decode golden image: %w", err)
	}

	actualImg := t.Ctx.Image()
	score = calculateDiff(goldenImg, actualImg)
	
	// Clean up actual if match is perfect
	if score == 0.0 {
		os.Remove(actualPath)
	}
	return score, nil
}

// calculateDiff calculates a simple pixel diff score between two images.
// Returns 0.0 for identical images.
func calculateDiff(img1, img2 image.Image) float64 {
	b1 := img1.Bounds()
	b2 := img2.Bounds()
	if b1 != b2 {
		return 1.0 // Maximum diff if bounds don't match
	}

	var diff float64
	var count float64
	for y := b1.Min.Y; y < b1.Max.Y; y++ {
		for x := b1.Min.X; x < b1.Max.X; x++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()
			diff += math.Abs(float64(r1)-float64(r2)) / 65535.0
			diff += math.Abs(float64(g1)-float64(g2)) / 65535.0
			diff += math.Abs(float64(b1)-float64(b2)) / 65535.0
			diff += math.Abs(float64(a1)-float64(a2)) / 65535.0
			count += 4.0
		}
	}
	return diff / count
}

// SimulateHover simulates hovering the mouse over the widget coordinates.
func (t *Tester) SimulateHover(x, y float64) {
	t.UI.Mouse.X = x
	t.UI.Mouse.Y = y
	t.UI.Mouse.LeftDown = false
	t.UI.Mouse.LeftPressed = false
	t.UI.Mouse.LeftReleased = false
}

// SimulateClick simulates a mouse click on a given coordinate.
func (t *Tester) SimulateClick(x, y float64) {
	t.UI.Mouse.X = x
	t.UI.Mouse.Y = y
	t.UI.Mouse.LeftDown = true
	t.UI.Mouse.LeftPressed = true
	t.UI.Mouse.LeftReleased = false
}

// SimulateClickRelease simulates releasing the mouse.
func (t *Tester) SimulateClickRelease(x, y float64) {
	t.UI.Mouse.X = x
	t.UI.Mouse.Y = y
	t.UI.Mouse.LeftDown = false
	t.UI.Mouse.LeftPressed = false
	t.UI.Mouse.LeftReleased = true
}

// StepFrame advances the animation time by a fixed frame delta (e.g. 1/60s).
func (t *Tester) StepFrame() {
	t.UI.DeltaTime = 1.0 / 60.0
}
