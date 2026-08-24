package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	path := "core/layout/cached.go"
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	str := string(content)
	
	// We need to inject the LogicOnly check before IsDirty
	
	findStr := 	// 2.5 Check if any damage rects intersect with this widget
	
	injectStr := 	if uictx.LogicOnly {
		// Logic pass: just run the child's logic and return
		originalMouse := uictx.Mouse
		uictx.Mouse.X -= startX
		uictx.Mouse.Y -= startY

		w.child.Render(uictx, 0, 0)

		uictx.Mouse = originalMouse
		return childW, childH
	}

 + findStr

	str = strings.Replace(str, findStr, injectStr, 1)
	
	// Also remove the else block that does the logic again
	findElse := 	} else {
		// Even if not dirty, we MUST run the child's logic so it can process hover/click!
		// But we don't need to redraw it to the cache.
		originalMouse := uictx.Mouse
		uictx.Mouse.X -= startX
		uictx.Mouse.Y -= startY

		uictx.LogicOnly = true
		w.child.Render(uictx, 0, 0)
		uictx.LogicOnly = false
		
		uictx.Mouse = originalMouse
	}
	
	str = strings.Replace(str, findElse, "", 1)
	
	err = os.WriteFile(path, []byte(str), 0644)
	if err != nil {
		fmt.Println("Error writing:", err)
	} else {
		fmt.Println("Updated cached.go")
	}
}
