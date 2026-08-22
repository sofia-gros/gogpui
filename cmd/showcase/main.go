package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/gogpu/gg"
	"github.com/gogpu/gg/integration/ggcanvas"
	"github.com/gogpu/gg/text"
	"github.com/gogpu/gogpu"
	"github.com/gogpu/gogpu/input"
	"github.com/sofiagros/gogpui/components/badge"
	"github.com/sofiagros/gogpui/components/button"
	"github.com/sofiagros/gogpui/components/checkbox"
	"github.com/sofiagros/gogpui/components/label"
	"github.com/sofiagros/gogpui/components/radio"
	"github.com/sofiagros/gogpui/components/slider"
	switch_comp "github.com/sofiagros/gogpui/components/switch"
	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/layout"
	"github.com/sofiagros/gogpui/core/theme"
)

const (
	winW = 1200
	winH = 860
)

func main() {
	// サブピクセルテキストメトリクスをグローバルに有効化する。
	// これにより各文字の advance width がサブピクセル精度で山計され、文字間隔が平均になる。
	// tester.go と同様にコンテキスト生成前に呼ぶ必要がある。
	text.SetGlobalSubpixelCache(text.NewSubpixelCache(text.HighQualitySubpixelConfig()))

	config := gogpu.DefaultConfig().
		WithTitle("gogpui Component Showcase").
		WithSize(winW, winH)

	app := gogpu.NewApp(config)

	var canvas *ggcanvas.Canvas
	var lastScale float64
	uictx := context.NewUIContext()

	// フォントソースを一度だけロードする（毎フレーム再読み込みを防ぐ）。
	// 毎フレーム再ロードすると MSDF アトラスが毎回初期化され、文字間隔が定になる。不安
	fontPath := filepath.Join("assets", "fonts", "Inter-Regular.ttf")
	fontSource, fontLoadErr := text.NewFontSourceFromFile(fontPath)
	if fontLoadErr != nil {
		log.Printf("Warning: Failed to load font from %s: %v", fontPath, fontLoadErr)
	}

	app.OnSurfaceAvailable(func() {
		lastScale = app.ScaleFactor()
		canvas = ggcanvas.MustNewWithScale(app.GPUContextProvider(), winW, winH, lastScale)
	})

	var lastMx, lastMy float32
	var lastLDown bool
	var pendingLeftPressed bool
	var pendingLeftReleased bool

	app.OnUpdate(func(dt float64) {
		in := app.Input()
		if in != nil {
			mx, my := in.Mouse().Position()
			ld := in.Mouse().Pressed(input.MouseButtonLeft)

			if in.Mouse().JustPressed(input.MouseButtonLeft) {
				pendingLeftPressed = true
			}
			if in.Mouse().JustReleased(input.MouseButtonLeft) {
				pendingLeftReleased = true
			}

			if mx != lastMx || my != lastMy || ld != lastLDown || pendingLeftPressed || pendingLeftReleased {
				lastMx, lastMy, lastLDown = mx, my, ld
				app.RequestRedraw()
			}
		}
	})

	lastTime := time.Now()

	// --- コンポーネントの状態 ---
	var chk1Val bool = false
	var chk2Val bool = true
	var sw1Val bool = false
	var sw2Val bool = true
	var rd1Val bool = false
	var rd2Val bool = true
	var sliderVal1 float64 = 0.3
	var sliderVal2 float64 = 0.7

	app.OnDraw(func(dc *gogpu.Context) {
		if canvas == nil {
			return
		}

		now := time.Now()
		dt := now.Sub(lastTime).Seconds()
		lastTime = now

		ctx := canvas.Context()
		th := theme.DefaultTheme()
		if app.DarkMode() {
			th.Mode = theme.DarkMode
			th.Colors = theme.DarkThemeColors()
		}
		in := app.Input()

		// 背景
		ctx.SetColor(th.Colors.Background)
		ctx.Clear()

		// フォントをセット（ソースは起動時に一度だけロード済み）。
		if fontSource != nil {
			ctx.SetFont(fontSource.Face(float64(th.FontSize),
				text.WithHinting(text.HintingNone),
				text.WithFeatures(text.NewFontFeature("kern", 1), text.NewFontFeature("liga", 1)),
			))
		}
		ctx.SetTextMode(gg.TextModeVector)

		uictx.Update(ctx, th, dt, in, lastScale)
		if pendingLeftPressed {
			uictx.Mouse.LeftPressed = true
			pendingLeftPressed = false
		}
		if pendingLeftReleased {
			uictx.Mouse.LeftReleased = true
			pendingLeftReleased = false
		}

		// --- レイアウト定数 ---
		const (
			padX    = 24.0  // 左右マージン
			padY    = 14.0  // 縦マージン
			colGap  = 16.0  // 列間ギャップ
			rowGap  = 16.0  // 行間ギャップ
			titleH  = 44.0  // ヘッダー高さ
			cardPad = 14.0  // カード内パディング
			secH    = 22.0  // セクションタイトル行高さ
			cardH1  = 182.0 // 1行目カード高さ
			cardH2  = 200.0 // 2行目カード高さ
		)
		numCols := 4.0
		colW := (float64(winW) - padX*2 - colGap*(numCols-1)) / numCols // ≈ 277px

		// ----- ヘッダー -----
		ctx.SetColor(th.Colors.Foreground)
		ctx.DrawStringAnchored("gogpui · Component Showcase", float64(winW)/2, titleH/2, 0.5, 0.5)

		// ヘッダー下の区切り線
		ctx.SetColor(th.Colors.Border)
		ctx.DrawRectangle(padX, titleH-1, float64(winW)-padX*2, 1)
		ctx.Fill()

		// ----- カード描画ヘルパー -----
		// カード背景と左上のセクションタイトルを描画する
		drawCard := func(x, y, w, h float64, title string) {
			// カード背景
			ctx.SetColor(th.Colors.Muted)
			ctx.DrawRoundedRectangle(x, y, w, h, 8)
			ctx.Fill()
			// セクション名
			ctx.SetColor(th.Colors.MutedForeground)
			ctx.DrawString(title, x+cardPad, y+cardPad+float64(th.FontSize)*0.85)
		}

		// ----- Row 1: Button / Checkbox / Switch / Radio -----
		row1Y := titleH + padY

		// Button
		col0x := padX
		drawCard(col0x, row1Y, colW, cardH1, "Button")
		layout.NewFlex().Direction(layout.Column).Gap(8).
			Add(
				button.New("sc-btn-primary").Label("Primary Button").Primary(),
				button.New("sc-btn-danger").Label("Danger Button").Danger(),
				button.New("sc-btn-link").Label("Link Button").Link(),
			).
			Render(uictx, col0x+cardPad, row1Y+secH+cardPad)

		// Checkbox
		col1x := padX + (colW + colGap)
		drawCard(col1x, row1Y, colW, cardH1, "Checkbox")
		layout.NewFlex().Direction(layout.Column).Gap(10).
			Add(
				checkbox.New("sc-chk1").Label("Unchecked").Checked(chk1Val).OnChange(func(v bool) { chk1Val = v }),
				checkbox.New("sc-chk2").Label("Checked").Checked(chk2Val).OnChange(func(v bool) { chk2Val = v }),
				checkbox.New("sc-chk3").Label("Disabled").Checked(false).Disabled(true),
			).
			Render(uictx, col1x+cardPad, row1Y+secH+cardPad)

		// Switch
		col2x := padX + (colW+colGap)*2
		drawCard(col2x, row1Y, colW, cardH1, "Switch")
		layout.NewFlex().Direction(layout.Column).Gap(10).
			Add(
				switch_comp.New("sc-sw1").Label("Off state").Checked(sw1Val).OnChange(func(v bool) { sw1Val = v }),
				switch_comp.New("sc-sw2").Label("On state").Checked(sw2Val).OnChange(func(v bool) { sw2Val = v }),
				switch_comp.New("sc-sw3").Label("Disabled").Checked(false).Disabled(true),
			).
			Render(uictx, col2x+cardPad, row1Y+secH+cardPad)

		// Radio
		col3x := padX + (colW+colGap)*3
		drawCard(col3x, row1Y, colW, cardH1, "Radio")
		layout.NewFlex().Direction(layout.Column).Gap(10).
			Add(
				radio.New("sc-rd1").Label("Unselected").Checked(rd1Val).OnChange(func(v bool) { rd1Val = v }),
				radio.New("sc-rd2").Label("Selected").Checked(rd2Val).OnChange(func(v bool) { rd2Val = v }),
				radio.New("sc-rd3").Label("Disabled").Checked(false).Disabled(true),
			).
			Render(uictx, col3x+cardPad, row1Y+secH+cardPad)

		// ----- Row 2: Slider / Label / Badge / Layout Demo -----
		row2Y := row1Y + cardH1 + rowGap

		// Slider
		drawCard(col0x, row2Y, colW, cardH2, "Slider")
		layout.NewFlex().Direction(layout.Column).Gap(14).
			Add(
				slider.New("sc-sl1").Value(sliderVal1).OnChange(func(v float64) { sliderVal1 = v }),
				slider.New("sc-sl2").Value(sliderVal2).OnChange(func(v float64) { sliderVal2 = v }),
				slider.New("sc-sl3").Value(0.5).Disabled(true),
			).
			Render(uictx, col0x+cardPad, row2Y+secH+cardPad)

		// Label
		drawCard(col1x, row2Y, colW, cardH2, "Label")
		layout.NewFlex().Direction(layout.Column).Gap(10).
			Add(
				label.New("Normal label"),
				label.New("With secondary").Secondary("extra info"),
				label.New("Highlighted text").Highlights("igh", false),
				label.New("Masked content").Masked(true),
			).
			Render(uictx, col1x+cardPad, row2Y+secH+cardPad)

		// Badge
		drawCard(col2x, row2Y, colW, cardH2, "Badge")
		layout.NewFlex().Direction(layout.Column).Gap(12).
			Add(
				badge.New().Count(5).Add(label.New("Messages")),
				badge.New().Count(120).Max(99).Add(label.New("Capped at 99+")),
				badge.New().Dot().Add(label.New("Dot variant")),
				badge.New().Count(3).Color(th.Colors.Info).Add(label.New("Info color")),
			).
			Render(uictx, col2x+cardPad, row2Y+secH+cardPad)

		// Layout Demo (Flex + Grid + Spacer + Padding)
		drawCard(col3x, row2Y, colW, cardH2, "Layout")
		innerW := colW - cardPad*2
		layout.NewFlex().Direction(layout.Column).Gap(10).
			Add(
				// Spacer: ラベルを左右に分ける
				layout.NewFlex().Direction(layout.Row).
					WithConstraints(innerW, 0).
					Add(
						label.New("Left"),
						layout.NewSpacerFlex(),
						label.New("Right"),
					),
				// Padding: ボタンに水平パディングを追加
				layout.NewPadding(
					button.New("sc-pad-btn").Label("Padded").Primary(),
				).Horizontal(10).Vertical(2),
				// Grid: 3列グリッドにボタン
				layout.NewGrid().Cols(3).Gap(6).
					Add(
						button.New("sc-g1").Label("A").Primary(),
						button.New("sc-g2").Label("B").Danger(),
						button.New("sc-g3").Label("C").Link(),
					),
			).
			Render(uictx, col3x+cardPad, row2Y+secH+cardPad)

		// ----- フッター -----
		row3Y := row2Y + cardH2 + rowGap
		statusLine := fmt.Sprintf(
			"Components: Button · Checkbox · Switch · Radio · Slider · Label · Badge     "+
				"Layout: Flex · Grid · Spacer · Padding",
		)
		ctx.SetColor(th.Colors.MutedForeground)
		ctx.DrawStringAnchored(statusLine, float64(winW)/2, row3Y+10, 0.5, 0.5)

		// canvas → ウィンドウへ描画
		canvas.MarkDirty()
		if renderErr := canvas.RenderTo(dc.AsTextureDrawer()); renderErr != nil {
			log.Printf("Failed to render canvas to screen: %v", renderErr)
		}

		if uictx.NeedsRedraw {
			app.RequestRedraw()
		}
	})

	log.Println("Starting gogpu window...")
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
