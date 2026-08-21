package label

import (
	"image/color"
	"strings"

	"github.com/sofiagros/gogpui/core/context"
)

type TextSegment struct {
	Text  string
	Color color.Color
}

type Label struct {
	text       string
	secondary  string
	masked     bool
	highlights string
	isPrefix   bool
}

func New(text string) *Label {
	return &Label{
		text: text,
	}
}

func (l *Label) Secondary(text string) *Label {
	l.secondary = text
	return l
}

func (l *Label) Masked(masked bool) *Label {
	l.masked = masked
	return l
}

func (l *Label) Highlights(text string, isPrefix bool) *Label {
	l.highlights = text
	l.isPrefix = isPrefix
	return l
}

// buildSegments splits the label and secondary text into colored segments
func (l *Label) buildSegments(uictx *context.UIContext) []TextSegment {
	var segments []TextSegment

	if l.masked {
		maskStr := strings.Repeat("•", len([]rune(l.text)))
		segments = append(segments, TextSegment{Text: maskStr, Color: uictx.Theme.Colors.Foreground})
		if l.secondary != "" {
			maskSec := strings.Repeat("•", len([]rune(l.secondary)))
			segments = append(segments, TextSegment{Text: " " + maskSec, Color: uictx.Theme.Colors.MutedForeground})
		}
		return segments
	}

	// Apply highlights to primary text
	if l.highlights != "" {
		hlLower := strings.ToLower(l.highlights)
		txtLower := strings.ToLower(l.text)

		if l.isPrefix {
			if strings.HasPrefix(txtLower, hlLower) {
				segments = append(segments, TextSegment{Text: l.text[:len(hlLower)], Color: uictx.Theme.Colors.Info})
				if len(l.text) > len(hlLower) {
					segments = append(segments, TextSegment{Text: l.text[len(hlLower):], Color: uictx.Theme.Colors.Foreground})
				}
			} else {
				segments = append(segments, TextSegment{Text: l.text, Color: uictx.Theme.Colors.Foreground})
			}
		} else {
			// Full matching (all occurrences)
			start := 0
			for {
				idx := strings.Index(txtLower[start:], hlLower)
				if idx == -1 {
					if start < len(l.text) {
						segments = append(segments, TextSegment{Text: l.text[start:], Color: uictx.Theme.Colors.Foreground})
					}
					break
				}
				realIdx := start + idx
				if realIdx > start {
					segments = append(segments, TextSegment{Text: l.text[start:realIdx], Color: uictx.Theme.Colors.Foreground})
				}
				segments = append(segments, TextSegment{Text: l.text[realIdx : realIdx+len(hlLower)], Color: uictx.Theme.Colors.Info})
				start = realIdx + len(hlLower)
			}
		}
	} else {
		segments = append(segments, TextSegment{Text: l.text, Color: uictx.Theme.Colors.Foreground})
	}

	// Add secondary text
	if l.secondary != "" {
		segments = append(segments, TextSegment{Text: " " + l.secondary, Color: uictx.Theme.Colors.MutedForeground})
	}

	return segments
}

func (l *Label) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	ctx := uictx.GG
	segments := l.buildSegments(uictx)

	var totalW float64
	var totalH float64

	// Measure pass
	for _, seg := range segments {
		if seg.Text == "" {
			continue
		}
		w, h := ctx.MeasureString(seg.Text)
		totalW += w
		if h > totalH {
			totalH = h
		}
	}

	// Make totalH at least font size to keep height stable
	if totalH == 0 {
		_, totalH = ctx.MeasureString("Ag") 
	}

	if uictx.MeasureOnly {
		return totalW, totalH
	}

	// Draw pass
	currentX := x
	for _, seg := range segments {
		if seg.Text == "" {
			continue
		}
		ctx.SetColor(seg.Color)
		w, _ := ctx.MeasureString(seg.Text)
		
		// Draw anchored at left-center to align horizontally while keeping vertical middle
		ctx.DrawStringAnchored(seg.Text, currentX, y+totalH/2.0, 0.0, 0.5)
		currentX += w
	}

	return totalW, totalH
}
