package main

import (
	"image"
	"image/color"

	gogpui "github.com/sofiagros/gogpui"
	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/layout"
	"github.com/gogpu/gg"
)

func main() {
	app := gogpui.New(gogpui.Options{Width: 800, Height: 600})
	
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	buf := gg.ImageBufFromImage(img)

	app.Run(func(uictx *context.UIContext) {
		uictx.GG.DrawImage(buf, 100, 100)
		
		contentWidget := &layout.CustomWidget{
			RenderFunc: func(uictx *context.UIContext, startX, startY float64) (float64, float64) {
				return 800, 600
			},
		}
		contentWidget.Render(uictx, 0, 0)
	})
}
