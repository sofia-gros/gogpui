package theme

import (
	"context"
)

// ThemeMode represents whether the theme is light or dark.
type ThemeMode int

const (
	LightMode ThemeMode = iota
	DarkMode
)

// Theme contains the global application theme including colors, fonts, and sizes.
type Theme struct {
	Mode   ThemeMode
	Colors ThemeColor

	// Basic configurations
	FontSize     float32
	MonoFontSize float32
	Radius       float32
	RadiusLg     float32
	Shadow       bool
	FocusRing    bool
}

// DefaultTheme returns a default initialized light theme.
func DefaultTheme() *Theme {
	return &Theme{
		Mode:         LightMode,
		Colors:       LightThemeColors(),
		FontSize:     16.0,
		MonoFontSize: 13.0,
		Radius:       6.0,
		RadiusLg:     8.0,
		Shadow:       true,
		FocusRing:    true,
	}
}

// themeContextKey is the type for the context key.
type themeContextKey struct{}

// WithTheme attaches a theme to a context.
func WithTheme(ctx context.Context, theme *Theme) context.Context {
	return context.WithValue(ctx, themeContextKey{}, theme)
}

// ThemeFrom retrieves the theme from a context.
// If it's a *gg.Context, it tries to extract it from its internal UserData/Context.
// For now we accept standard context.Context. In actual use, you'll wrap or pass gg.Context appropriately.
func ThemeFrom(ctx context.Context) *Theme {
	if theme, ok := ctx.Value(themeContextKey{}).(*Theme); ok {
		return theme
	}
	// Fallback to default if none is found
	return DefaultTheme()
}

// ggContextKey is used to store theme in gg.Context UserData if gg supports it,
// otherwise this relies on Go standard contexts.
