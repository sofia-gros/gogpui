package theme

import (
	"context"
	"testing"
)

func TestThemeMode(t *testing.T) {
	theme := DefaultTheme()
	
	if theme.IsDark() {
		t.Errorf("expected default theme to be light")
	}

	theme.SetMode(DarkMode)
	if !theme.IsDark() {
		t.Errorf("expected theme to be dark after SetMode(DarkMode)")
	}

	if theme.Colors.Background != DarkThemeColors().Background {
		t.Errorf("expected colors to switch to dark mode colors")
	}
}

func TestThemeFromContext(t *testing.T) {
	ctx := context.Background()
	defaultTh := ThemeFrom(ctx)

	if defaultTh == nil {
		t.Errorf("ThemeFrom should return a default theme if none exists in context")
	}

	customTheme := DefaultTheme()
	customTheme.FontSize = 20.0
	ctx = WithTheme(ctx, customTheme)

	retrieved := ThemeFrom(ctx)
	if retrieved.FontSize != 20.0 {
		t.Errorf("expected to retrieve custom theme from context, got %v", retrieved.FontSize)
	}
}
