package main

import (
	"image/png"
	"os"
	
	"github.com/gogpu/gg"
)

func main() {
	ctx := gg.NewContext(100, 100)
	ctx.SetHexColor("#ff0000")
	ctx.DrawRectangle(0, 0, 100, 100)
	ctx.Fill()
	
	f, _ := os.Create("test.png")
	png.Encode(f, ctx.Image())
	f.Close()
}
