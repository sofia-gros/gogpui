package theme

// SetMode switches the theme between light and dark modes.
func (t *Theme) SetMode(mode ThemeMode) {
	t.Mode = mode
	if mode == DarkMode {
		t.Colors = DarkThemeColors()
	} else {
		t.Colors = LightThemeColors()
	}
}

// IsDark returns true if the current theme is dark mode.
func (t *Theme) IsDark() bool {
	return t.Mode == DarkMode
}
