package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	err := filepath.WalkDir(`a:\Project\gogpui\components`, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		oldMeasure1 := []byte("if uictx.MeasureOnly {\n\t\treturn w, h\n\t}")
		newMeasure1 := []byte("if uictx.MeasureOnly || !uictx.IsVisible(x, y, w, h) {\n\t\treturn w, h\n\t}")

		oldMeasure2 := []byte("if uictx.MeasureOnly {\n\t\treturn size.W, size.H\n\t}")
		newMeasure2 := []byte("if uictx.MeasureOnly || !uictx.IsVisible(x, y, size.W, size.H) {\n\t\treturn size.W, size.H\n\t}")

		oldMeasure3 := []byte("if uictx.MeasureOnly {\n\t\treturn fw, fh\n\t}")
		newMeasure3 := []byte("if uictx.MeasureOnly || !uictx.IsVisible(x, y, fw, fh) {\n\t\treturn fw, fh\n\t}")

		modified := bytes.ReplaceAll(content, oldMeasure1, newMeasure1)
		modified = bytes.ReplaceAll(modified, oldMeasure2, newMeasure2)
		modified = bytes.ReplaceAll(modified, oldMeasure3, newMeasure3)
		
		// Add LogicOnly wrapper replaces
		modified = bytes.ReplaceAll(modified, []byte("ctx.Fill()"), []byte("uictx.Fill()"))
		modified = bytes.ReplaceAll(modified, []byte("ctx.Stroke()"), []byte("uictx.Stroke()"))
		modified = bytes.ReplaceAll(modified, []byte("ctx.DrawString("), []byte("uictx.DrawString("))
		modified = bytes.ReplaceAll(modified, []byte("ctx.DrawStringAnchored("), []byte("uictx.DrawStringAnchored("))
		modified = bytes.ReplaceAll(modified, []byte("ctx.DrawRoundedRectangle("), []byte("uictx.DrawRoundedRectangle("))

		if !bytes.Equal(content, modified) {
			err = os.WriteFile(path, modified, 0644)
			if err != nil {
				fmt.Println("Error writing", path, err)
			} else {
				fmt.Println("Updated", path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking", err)
	}
}
