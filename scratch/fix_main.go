package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	path := "cmd/showcase/main.go"
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	oldDraw := []byte("if uictx.MeasureOnly { return }\n\t\t\t\t\tctx.SetColor(th.Colors.Background)")
	newDraw := []byte("if uictx.MeasureOnly || !uictx.IsVisible(x, y, w, h) { return }\n\t\t\t\t\tctx.SetColor(th.Colors.Background)")
	content = bytes.ReplaceAll(content, oldDraw, newDraw)
    
    oldTitle1 := []byte("if !uictx.MeasureOnly {\n\t\t\t\t\tctx.SetColor(th.Colors.Foreground)\n\t\t\t\t\tctx.DrawStringAnchored(")
    newTitle1 := []byte("if !uictx.MeasureOnly && uictx.IsVisible(startX+wW/2, startY+titleH/2, 1, 1) {\n\t\t\t\t\tctx.SetColor(th.Colors.Foreground)\n\t\t\t\t\tctx.DrawStringAnchored(")
    content = bytes.ReplaceAll(content, oldTitle1, newTitle1)
    
    oldTitle2 := []byte("if !uictx.MeasureOnly {\n\t\t\t\t\tstatusLine := \"Components")
    newTitle2 := []byte("if !uictx.MeasureOnly && uictx.IsVisible(startX+wW/2, startY+rowFooterY+10, 1, 1) {\n\t\t\t\t\tstatusLine := \"Components")
    content = bytes.ReplaceAll(content, oldTitle2, newTitle2)

	err = os.WriteFile(path, content, 0644)
	if err != nil {
		fmt.Println("Error writing:", err)
	} else {
		fmt.Println("Updated main.go")
	}
}
