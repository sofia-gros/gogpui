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
		fmt.Println("Error:", err)
		return
	}
	
	// We will use string replacement to inject CardWidget and use it.
	str := string(content)
	
	cardWidgetCode := 
type CardWidget struct {
	title string
	w, h  float64
	child layout.Widget
}

func (c *CardWidget) Render(uictx *context.UIContext, x, y float64) (w, h float64) {
	if uictx.MeasureOnly {
		return c.w, c.h
	}
	ctx := uictx.GG
	th := uictx.Theme
	cardPad := 14.0
	secH := 22.0

	ctx.SetColor(th.Colors.Background)
	uictx.DrawRoundedRectangle(x, y, c.w, c.h, 8)
	uictx.Fill()
	ctx.SetColor(th.Colors.Border)
	ctx.SetLineWidth(1.0)
	uictx.DrawRoundedRectangle(x, y, c.w, c.h, 8)
	uictx.Stroke()
	ctx.SetColor(th.Colors.MutedForeground)
	uictx.DrawString(c.title, x+cardPad, y+cardPad+float64(th.FontSize)*0.85)

	c.child.Render(uictx, x+cardPad, y+secH+cardPad)
	return c.w, c.h
}

func main() {

	str = strings.Replace(str, "func main() {", cardWidgetCode, 1)

	// Remove drawCard definition
	startDrawCard := strings.Index(str, "drawCard := func(")
	endDrawCard := strings.Index(str[startDrawCard:], "}") + startDrawCard + 1
	str = str[:startDrawCard] + str[endDrawCard:]

	// Replace card draw calls
	// Example:
	// drawCard(col0x, row1Y, colW, cardH1, "Button")
	// layout.NewFlex().Direction(layout.Column).Gap(8).
	//		Add( ... ).
	//		Render(uictx, col0x+cardPad, row1Y+secH+cardPad)
	
	// We will write a small function to do regex replacement or manual replacement
	// Since regex is complex in Go, we'll write a manual loop
	
	lines := strings.Split(str, "\n")
	var out []string
	
	i := 0
	for i < len(lines) {
		line := lines[i]
		if strings.Contains(line, "drawCard(") {
			// Extract parameters
			parts := strings.Split(line, "drawCard(")
			args := strings.Split(parts[1], ")")[0]
			argParts := strings.Split(args, ", ")
			cx := argParts[0]
			cy := argParts[1]
			cw := argParts[2]
			ch := argParts[3]
			title := argParts[4]
			
			// Extract the flex/child rendering that follows
			i++
			var childLines []string
			childLines = append(childLines, strings.TrimSpace(lines[i])) // e.g. "layout.NewFlex().Direction(layout.Column).Gap(8)."
			
			// Some are treeWidget or breadcrumb directly, but mostly layout.NewFlex()
			isTree := strings.Contains(childLines[0], "treeWidget :=")
			if isTree {
				childLines[0] = "tree.New(demoTreeState).Width(colW - cardPad*2)"
				i++ // skip treeWidget.Render...
			} else if strings.Contains(childLines[0], "toggleBtn :=") {
				childLines[0] = strings.TrimSpace(lines[i]) + "\n" + strings.TrimSpace(lines[i+1]) + "\n" + strings.TrimSpace(lines[i+2]) + "\n" + strings.TrimSpace(lines[i+3]) + "\n" + strings.TrimSpace(lines[i+4]) + "\n" + strings.TrimSpace(lines[i+5]) + "\n" + strings.TrimSpace(lines[i+6])
				i += 7 // skip collapsible block
				// Fix Render call
				childLines[0] = strings.Replace(childLines[0], ".Render(uictx, col1x+cardPad, row4Y+secH+cardPad)", "", 1)
			} else {
				for {
					i++
					l := lines[i]
					if strings.Contains(l, ".Render(uictx") {
						break
					}
					childLines = append(childLines, l)
				}
			}
			
			childStr := strings.Join(childLines, "\n")
			
			isDynamic := strings.Contains(title, "Progress") || strings.Contains(title, "Skeleton") || strings.Contains(title, "Collapsible")
			
			cardDef := fmt.Sprintf("&CardWidget{title: %s, w: %s, h: %s, child: %s}", title, cw, ch, childStr)
			
			if isDynamic {
				out = append(out, fmt.Sprintf("\t\t\t\t%s.Render(uictx, %s, %s)", cardDef, cx, cy))
			} else {
				out = append(out, fmt.Sprintf("\t\t\t\tlayout.NewCachedWidget(%s, %s).Render(uictx, %s, %s)", title, cardDef, cx, cy))
			}
		} else {
			out = append(out, line)
		}
		i++
	}
	
	err = os.WriteFile(path, []byte(strings.Join(out, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing:", err)
	} else {
		fmt.Println("Updated main.go")
	}
}
