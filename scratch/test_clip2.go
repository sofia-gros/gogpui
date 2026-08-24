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
	
	// Remove unused imports and vars
	str = strings.ReplaceAll(str, "\"github.com/sofiagros/gogpui/components/scroll\"", "")
	str = strings.ReplaceAll(str, "wH := uictx.WindowHeight", "")
	
	os.WriteFile(path, []byte(str), 0644)
	fmt.Println("Fixed unused")
}
