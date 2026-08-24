package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	path := "cmd/showcase/main.go"
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}
	str := string(content)
	
	// Remove Scroll: scroll.New("main-scroll").Size(float64(wW), float64(wH)).Child(contentWidget).Render(uictx, 0, 0)
	// Replace with: contentWidget.Render(uictx, 0, 0)
	
	str = strings.ReplaceAll(str, "scroll.New(\"main-scroll\").Size(float64(wW), float64(wH)).Child(contentWidget).Render(uictx, 0, 0)", "contentWidget.Render(uictx, 0, 0)")
	
	os.WriteFile(path, []byte(str), 0644)
	fmt.Println("Removed Scroll")
}
