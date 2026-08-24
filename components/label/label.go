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
	textColor  *color.Color
}

func New(text string) *Label {
	return &Label{
		text: text,
	}
}

// Color sets a custom color for the label text.
func (l *Label) Color(c color.Color) *Label {
	l.textColor = &c
	return l
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

	baseColor := uictx.Theme.Colors.Foreground
	if l.textColor != nil {
		baseColor = *l.textColor
	}

	if l.masked {
		maskStr := strings.Repeat("•", len([]rune(l.text)))
		segments = append(segments, TextSegment{Text: maskStr, Color: baseColor})
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

		start := 0
		if l.isPrefix {
			if strings.HasPrefix(txtLower, hlLower) {
				segments = append(segments, TextSegment{Text: l.text[:len(hlLower)], Color: uictx.Theme.Colors.Info})
				if len(l.text) > len(hlLower) {
					segments = append(segments, TextSegment{Text: l.text[len(hlLower):], Color: baseColor})
				}
			} else {
				segments = append(segments, TextSegment{Text: l.text, Color: baseColor})
			}
		} else {
			// Full matching (all occurrences)
			for start < len(l.text) {
				idx := strings.Index(txtLower[start:], hlLower)
				if idx == -1 {
					if start < len(l.text) {
						segments = append(segments, TextSegment{Text: l.text[start:], Color: baseColor})
					}
					break
				}
				realIdx := start + idx
				if realIdx > start {
					segments = append(segments, TextSegment{Text: l.text[start:realIdx], Color: baseColor})
				}
				segments = append(segments, TextSegment{Text: l.text[realIdx : realIdx+len(hlLower)], Color: uictx.Theme.Colors.Info})
				start = realIdx + len(hlLower)
			}
		}
	} else {
		segments = append(segments, TextSegment{Text: l.text, Color: baseColor})
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
		w, h := uictx.MeasureText(seg.Text)
		totalW += w
		if h > totalH {
			totalH = h
		}
	}

	// Make totalH at least font size to keep height stable
	if totalH == 0 {
		_, totalH = uictx.MeasureText("Ag") 
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
		w, _ := uictx.MeasureText(seg.Text)
		
		// Draw anchored at left-center to align horizontally while keeping vertical middle
		uictx.DrawStringAnchored(seg.Text, currentX, y+totalH/2.0, 0.0, 0.5)
		currentX += w
	}

	return totalW, totalH
}
