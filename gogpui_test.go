package gogpui

import (
	"testing"

	"github.com/gogpu/gg"
	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/theme"
)

// TestNew は App の初期化とデフォルトオプションの適用をテストする。
func TestNew(t *testing.T) {
	// デフォルト値のテスト
	app := New(Options{})
	if app.opts.Width != 800 {
		t.Errorf("expected Width 800, got %d", app.opts.Width)
	}
	if app.opts.Height != 600 {
		t.Errorf("expected Height 600, got %d", app.opts.Height)
	}
	if app.opts.Title != "gogpui App" {
		t.Errorf("expected Title 'gogpui App', got %s", app.opts.Title)
	}
	if app.opts.FontPath != "" {
		t.Errorf("expected empty FontPath by default, got %s", app.opts.FontPath)
	}

	// カスタムオプションのテスト
	custom := New(Options{
		Title:    "Custom App",
		Width:    1024,
		Height:   768,
		FontPath: "custom/font.ttf",
	})
	if custom.opts.Width != 1024 {
		t.Errorf("expected Width 1024, got %d", custom.opts.Width)
	}
	if custom.opts.Height != 768 {
		t.Errorf("expected Height 768, got %d", custom.opts.Height)
	}
	if custom.opts.Title != "Custom App" {
		t.Errorf("expected Title 'Custom App', got %s", custom.opts.Title)
	}
	if custom.opts.FontPath != "custom/font.ttf" {
		t.Errorf("expected FontPath 'custom/font.ttf', got %s", custom.opts.FontPath)
	}
}

// TestUIContextLogicalDimensions は UIContext に渡される論理ウィンドウサイズの正確性をテストする。
func TestUIContextLogicalDimensions(t *testing.T) {
	uictx := context.NewUIContext()
	ctx := gg.NewContext(1920, 1080)
	th := theme.DefaultTheme()

	// リサイズ後の論理サイズ 1920x1080、スケール 1.5
	uictx.Update(ctx, th, 0.016, nil, 1.5, 1920, 1080)

	if uictx.WindowWidth != 1920 {
		t.Errorf("expected WindowWidth 1920, got %f", uictx.WindowWidth)
	}
	if uictx.WindowHeight != 1080 {
		t.Errorf("expected WindowHeight 1080, got %f", uictx.WindowHeight)
	}
}
