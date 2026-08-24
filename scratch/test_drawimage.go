package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/gogpu/gg"
)

func main() {
	// Create a software context
	ctx := gg.NewContext(100, 100)
	ctx.SetColor(color.RGBA{255, 0, 0, 255})
	ctx.DrawRectangle(0, 0, 100, 100)
	ctx.Fill()

	img := ctx.Image()
	buf := gg.ImageBufFromImage(img)

	fmt.Println("Image bounds:", img.Bounds())
	fmt.Println("Buf:", buf != nil)
}
