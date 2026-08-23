package avatar

import (
	"hash/fnv"
	"image/color"
	"strings"

	"github.com/sofiagros/gogpui/core/context"
)

// Size はアバターのサイズバリアントを表す。
type Size int

const (
	SizeMedium Size = iota
	SizeSmall
	SizeXSmall
	SizeLarge
)

// Avatar はユーザーアバター要素。
type Avatar struct {
	src       string
	name      string
	shortName string
	size      Size
}

// New は新しい Avatar 要素を作成する。
func New() *Avatar {
	return &Avatar{
		size: SizeMedium,
	}
}

// Src はアバターの画像ソースパスを設定する。
// （注: 現在の gogpui は画像ロード未サポートのため、プロパティ保持のみ）
func (a *Avatar) Src(src string) *Avatar {
	a.src = src
	return a
}

// Name はアバターのユーザー名を設定する。Srcがない場合はプレースホルダーとして使用される。
func (a *Avatar) Name(name string) *Avatar {
	a.name = name
	a.shortName = extractTextInitials(name)
	return a
}

// Size はアバターのサイズを設定する。
func (a *Avatar) Size(size Size) *Avatar {
	a.size = size
	return a
}

// extractTextInitials は Rust版の extract_text_initials と等価なロジック
func extractTextInitials(text string) string {
	words := strings.Split(text, " ")
	var result string
	count := 0
	for _, w := range words {
		if w == "" {
			continue
		}
		runes := []rune(w)
		if len(runes) > 0 {
			result += string(runes[0])
			count++
			if count == 2 {
				break
			}
		}
	}

	if len([]rune(result)) == 1 {
		runes := []rune(text)
		if len(runes) > 2 {
			runes = runes[:2]
		}
		result = string(runes)
	}

	return strings.ToUpper(result)
}

func hashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	var fR, fG, fB float64
	if s == 0 {
		fR, fG, fB = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q
		fR = hueToRGB(p, q, h+1.0/3.0)
		fG = hueToRGB(p, q, h)
		fB = hueToRGB(p, q, h-1.0/3.0)
	}
	return uint8(fR * 255), uint8(fG * 255), uint8(fB * 255)
}

// Render は指定位置に Avatar を描画し、占有サイズ (幅, 高さ) を返す。
func (a *Avatar) Render(uictx *context.UIContext, x, y float64) (float64, float64) {
	ctx := uictx.GG

	var sizePx float64
	var fontPx float64
	switch a.size {
	case SizeXSmall:
		sizePx = 16.0
		fontPx = 8.0
	case SizeSmall:
		sizePx = 24.0
		fontPx = 10.0
	case SizeLarge:
		sizePx = 40.0
		fontPx = 16.0
	case SizeMedium:
		fallthrough
	default:
		sizePx = 32.0
		fontPx = 14.0
	}

	if uictx.MeasureOnly {
		return sizePx, sizePx
	}

	radius := sizePx / 2.0
	cx := x + radius
	cy := y + radius

	if a.name != "" {
		// ハッシュ計算 (COLOR_COUNT = 360 / 15 = 24)
		colorIx := hashString(a.shortName) % 24
		hue := float64(colorIx*15) / 360.0

		// 色生成 (S=0.8, L=0.55 は固定)
		r, g, b := hslToRGB(hue, 0.8, 0.55)

		// 背景 (opacity 0.2)
		bgColor := color.NRGBA{R: r, G: g, B: b, A: uint8(255 * 0.2)}
		// テキスト (opacity 1.0)
		textColor := color.RGBA{R: r, G: g, B: b, A: 255}

		// 背景描画
		ctx.SetColor(bgColor)
		ctx.DrawCircle(cx, cy, radius)
		_ = ctx.Fill()

		// テキスト描画
		ctx.SetColor(textColor)
		
		scaleRatio := fontPx / float64(uictx.Theme.FontSize)
		ctx.Push()
		ctx.Translate(cx, cy)
		ctx.Scale(scaleRatio, scaleRatio)
		ctx.Translate(-cx, -cy)
		ctx.DrawStringAnchored(a.shortName, cx, cy, 0.5, 0.5)
		ctx.Pop()

	} else {
		// placeholder
		bgColor := uictx.Theme.Colors.Secondary
		textColor := uictx.Theme.Colors.Background
		borderColor := uictx.Theme.Colors.Border

		ctx.SetColor(bgColor)
		ctx.DrawCircle(cx, cy, radius)
		_ = ctx.FillPreserve()

		ctx.SetColor(borderColor)
		ctx.SetLineWidth(1.0)
		_ = ctx.Stroke()

		// Icon::User の代わりに "U" を描画する
		ctx.SetColor(textColor)
		scaleRatio := fontPx / float64(uictx.Theme.FontSize)
		ctx.Push()
		ctx.Translate(cx, cy)
		ctx.Scale(scaleRatio, scaleRatio)
		ctx.Translate(-cx, -cy)
		ctx.DrawStringAnchored("U", cx, cy, 0.5, 0.5)
		ctx.Pop()
	}

	return sizePx, sizePx
}
