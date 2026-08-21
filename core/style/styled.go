package style

import (
	"github.com/gogpu/gg"
	"github.com/sofiagros/gogpui/core/theme"
)

// ApplyRadius is a helper that configures the gogpu context to draw with the theme's corner radius.
func ApplyRadius(ctx *gg.Context, t *theme.Theme) {
	// In gogpu/gg, corner radius or paths might be set on the context.
	// This is a stub for where styling logic translates to gg calls.
	_ = ctx
	_ = t
}

// ApplyShadow configures the gogpu context to draw a shadow if enabled in the theme.
func ApplyShadow(ctx *gg.Context, t *theme.Theme) {
	if !t.Shadow {
		return
	}
	// Stub: apply shadow parameters to gg context.
}
