package layout

import "github.com/sofiagros/gogpui/core/context"

// CustomWidget is a helper to wrap custom drawing logic into a Widget.
type CustomWidget struct {
	RenderFunc func(uictx *context.UIContext, x, y float64) (float64, float64)
}

func (c *CustomWidget) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	if c.RenderFunc != nil {
		return c.RenderFunc(uictx, x, y)
	}
	return 0, 0
}
