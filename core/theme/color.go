package theme

import "image/color"

// ThemeColor contains the semantic color palette for the application.
// Based on shadcn/ui and GPUI's ThemeColor.
type ThemeColor struct {
	Background color.Color
	Foreground color.Color

	Primary           color.Color
	PrimaryForeground color.Color

	Secondary           color.Color
	SecondaryForeground color.Color

	Muted           color.Color
	MutedForeground color.Color

	Danger           color.Color
	DangerForeground color.Color

	Success           color.Color
	SuccessForeground color.Color

	Warning           color.Color
	WarningForeground color.Color

	Info           color.Color
	InfoForeground color.Color

	Border color.Color
	Input  color.Color
	Ring   color.Color
}

// LightThemeColors returns the default light theme color palette.
func LightThemeColors() ThemeColor {
	return ThemeColor{
		Background: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // hsl(0, 0%, 100%)
		Foreground: color.RGBA{R: 9, G: 9, B: 11, A: 255},     // hsl(240, 10%, 3.9%)

		Primary:           color.RGBA{R: 24, G: 24, B: 27, A: 255}, // hsl(240, 5.9%, 10%)
		PrimaryForeground: color.RGBA{R: 250, G: 250, B: 250, A: 255}, // hsl(0, 0%, 98%)

		Secondary:           color.RGBA{R: 244, G: 244, B: 245, A: 255}, // hsl(240, 4.8%, 95.9%)
		SecondaryForeground: color.RGBA{R: 24, G: 24, B: 27, A: 255}, // hsl(240, 5.9%, 10%)

		Muted:           color.RGBA{R: 244, G: 244, B: 245, A: 255}, // hsl(240, 4.8%, 95.9%)
		MutedForeground: color.RGBA{R: 113, G: 113, B: 122, A: 255}, // hsl(240, 3.8%, 46.1%)

		Danger:           color.RGBA{R: 239, G: 68, B: 68, A: 255}, // hsl(0, 84.2%, 60.2%)
		DangerForeground: color.RGBA{R: 250, G: 250, B: 250, A: 255}, // hsl(0, 0%, 98%)

		Success:           color.RGBA{R: 34, G: 197, B: 94, A: 255},  // green-500
		SuccessForeground: color.RGBA{R: 250, G: 250, B: 250, A: 255},

		Warning:           color.RGBA{R: 245, G: 158, B: 11, A: 255}, // amber-500
		WarningForeground: color.RGBA{R: 250, G: 250, B: 250, A: 255},

		Info:           color.RGBA{R: 59, G: 130, B: 246, A: 255},    // blue-500
		InfoForeground: color.RGBA{R: 250, G: 250, B: 250, A: 255},

		Border: color.RGBA{R: 228, G: 228, B: 231, A: 255}, // hsl(240, 5.9%, 90%)
		Input:  color.RGBA{R: 228, G: 228, B: 231, A: 255}, // hsl(240, 5.9%, 90%)
		Ring:   color.RGBA{R: 24, G: 24, B: 27, A: 255},    // hsl(240, 5.9%, 10%)
	}
}

// DarkThemeColors returns the default dark theme color palette.
func DarkThemeColors() ThemeColor {
	return ThemeColor{
		Background: color.RGBA{R: 9, G: 9, B: 11, A: 255},     // hsl(240, 10%, 3.9%)
		Foreground: color.RGBA{R: 250, G: 250, B: 250, A: 255}, // hsl(0, 0%, 98%)

		Primary:           color.RGBA{R: 250, G: 250, B: 250, A: 255}, // hsl(0, 0%, 98%)
		PrimaryForeground: color.RGBA{R: 24, G: 24, B: 27, A: 255}, // hsl(240, 5.9%, 10%)

		Secondary:           color.RGBA{R: 39, G: 39, B: 42, A: 255}, // hsl(240, 3.7%, 15.9%)
		SecondaryForeground: color.RGBA{R: 250, G: 250, B: 250, A: 255}, // hsl(0, 0%, 98%)

		Muted:           color.RGBA{R: 39, G: 39, B: 42, A: 255}, // hsl(240, 3.7%, 15.9%)
		MutedForeground: color.RGBA{R: 161, G: 161, B: 170, A: 255}, // hsl(240, 5%, 64.9%)

		Danger:           color.RGBA{R: 127, G: 29, B: 29, A: 255}, // hsl(0, 62.8%, 30.6%)
		DangerForeground: color.RGBA{R: 250, G: 250, B: 250, A: 255}, // hsl(0, 0%, 98%)

		Success:           color.RGBA{R: 21, G: 128, B: 61, A: 255},  // green-700
		SuccessForeground: color.RGBA{R: 250, G: 250, B: 250, A: 255},

		Warning:           color.RGBA{R: 180, G: 83, B: 9, A: 255},   // amber-700
		WarningForeground: color.RGBA{R: 250, G: 250, B: 250, A: 255},

		Info:           color.RGBA{R: 29, G: 78, B: 216, A: 255},     // blue-700
		InfoForeground: color.RGBA{R: 250, G: 250, B: 250, A: 255},

		// Border/Input: zinc-600 (#52525b) — ダークモードで視認可能なボーダー色。
		// 元の zinc-800 (#27272a) は背景 (#09090b) とのコントラストが低すぎてストロークが見えなかった。
		Border: color.RGBA{R: 82, G: 82, B: 91, A: 255},  // zinc-600 #52525b
		Input:  color.RGBA{R: 82, G: 82, B: 91, A: 255},  // zinc-600 #52525b
		Ring:   color.RGBA{R: 212, G: 212, B: 216, A: 255}, // hsl(240, 4.9%, 83.9%)
	}
}
