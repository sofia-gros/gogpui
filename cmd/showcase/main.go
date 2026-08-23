package main

import (
	"fmt"

	gogpui "github.com/sofiagros/gogpui"
	"github.com/sofiagros/gogpui/components/badge"
	"github.com/sofiagros/gogpui/components/button"
	"github.com/sofiagros/gogpui/components/checkbox"
	"github.com/sofiagros/gogpui/components/label"
	"github.com/sofiagros/gogpui/components/radio"
	"github.com/sofiagros/gogpui/components/slider"
	switch_comp "github.com/sofiagros/gogpui/components/switch"
	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/layout"
)

const (
	winW = 1200
	winH = 860
)

func main() {
	// --- コンポーネントの状態 ---
	var chk1Val bool = false
	var chk2Val bool = true
	var sw1Val bool = false
	var sw2Val bool = true
	var rd1Val bool = false
	var rd2Val bool = true
	var sliderVal1 float64 = 0.3
	var sliderVal2 float64 = 0.7

	app := gogpui.New(gogpui.Options{
		Title:  "gogpui Component Showcase",
		Width:  winW,
		Height: winH,
	})

	app.Run(func(uictx *context.UIContext) {
		th := uictx.Theme
		ctx := uictx.GG

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
		// ウィンドウサイズは毎フレーム UIContext から取得しリサイズに自動追従する。
		wW := uictx.WindowWidth
		wH := uictx.WindowHeight
		_ = wH // 現在未使用だが将来的なレイアウト制約用
		numCols := 4.0
		colW := (wW - padX*2 - colGap*(numCols-1)) / numCols

		// ----- ヘッダー -----
		ctx.SetColor(th.Colors.Foreground)
		ctx.DrawStringAnchored("gogpui Component Showcase", wW/2, titleH/2, 0.5, 0.5)

		// ヘッダー下の区切り線
		ctx.SetColor(th.Colors.Border)
		ctx.DrawRectangle(padX, titleH-1, wW-padX*2, 1)
		ctx.Fill()

		// ----- カード描画ヘルパー -----
		drawCard := func(x, y, w, h float64, title string) {
			ctx.SetColor(th.Colors.Muted)
			ctx.DrawRoundedRectangle(x, y, w, h, 8)
			ctx.Fill()
			ctx.SetColor(th.Colors.MutedForeground)
			ctx.DrawString(title, x+cardPad, y+cardPad+float64(th.FontSize)*0.85)
		}

		// ----- Row 1: Button / Checkbox / Switch / Radio -----
		row1Y := titleH + padY

		col0x := padX
		drawCard(col0x, row1Y, colW, cardH1, "Button")
		layout.NewFlex().Direction(layout.Column).Gap(8).
			Add(
				button.New("sc-btn-primary").Label("Primary Button").Primary(),
				button.New("sc-btn-danger").Label("Danger Button").Danger(),
				button.New("sc-btn-link").Label("Link Button").Link(),
			).
			Render(uictx, col0x+cardPad, row1Y+secH+cardPad)

		col1x := padX + (colW + colGap)
		drawCard(col1x, row1Y, colW, cardH1, "Checkbox")
		layout.NewFlex().Direction(layout.Column).Gap(10).
			Add(
				checkbox.New("sc-chk1").Label("Unchecked").Checked(chk1Val).OnChange(func(v bool) { chk1Val = v }),
				checkbox.New("sc-chk2").Label("Checked").Checked(chk2Val).OnChange(func(v bool) { chk2Val = v }),
				checkbox.New("sc-chk3").Label("Disabled").Checked(false).Disabled(true),
			).
			Render(uictx, col1x+cardPad, row1Y+secH+cardPad)

		col2x := padX + (colW+colGap)*2
		drawCard(col2x, row1Y, colW, cardH1, "Switch")
		layout.NewFlex().Direction(layout.Column).Gap(10).
			Add(
				switch_comp.New("sc-sw1").Label("Off state").Checked(sw1Val).OnChange(func(v bool) { sw1Val = v }),
				switch_comp.New("sc-sw2").Label("On state").Checked(sw2Val).OnChange(func(v bool) { sw2Val = v }),
				switch_comp.New("sc-sw3").Label("Disabled").Checked(false).Disabled(true),
			).
			Render(uictx, col2x+cardPad, row1Y+secH+cardPad)

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

		drawCard(col0x, row2Y, colW, cardH2, "Slider")
		layout.NewFlex().Direction(layout.Column).Gap(14).
			Add(
				slider.New("sc-sl1").Value(sliderVal1).OnChange(func(v float64) { sliderVal1 = v }),
				slider.New("sc-sl2").Value(sliderVal2).OnChange(func(v float64) { sliderVal2 = v }),
				slider.New("sc-sl3").Value(0.5).Disabled(true),
			).
			Render(uictx, col0x+cardPad, row2Y+secH+cardPad)

		drawCard(col1x, row2Y, colW, cardH2, "Label")
		layout.NewFlex().Direction(layout.Column).Gap(10).
			Add(
				label.New("Normal label"),
				label.New("With secondary").Secondary("extra info"),
				label.New("Highlighted text").Highlights("igh", false),
				label.New("Masked content").Masked(true),
			).
			Render(uictx, col1x+cardPad, row2Y+secH+cardPad)

		drawCard(col2x, row2Y, colW, cardH2, "Badge")
		layout.NewFlex().Direction(layout.Column).Gap(12).
			Add(
				badge.New().Count(5).Add(label.New("Messages")),
				badge.New().Count(120).Max(99).Add(label.New("Capped at 99+")),
				badge.New().Dot().Add(label.New("Dot variant")),
				badge.New().Count(3).Color(th.Colors.Info).Add(label.New("Info color")),
			).
			Render(uictx, col2x+cardPad, row2Y+secH+cardPad)

		drawCard(col3x, row2Y, colW, cardH2, "Layout")
		innerW := colW - cardPad*2
		layout.NewFlex().Direction(layout.Column).Gap(10).
			Add(
				layout.NewFlex().Direction(layout.Row).
					WithConstraints(innerW, 0).
					Add(label.New("Left"), layout.NewSpacerFlex(), label.New("Right")),
				layout.NewPadding(
					button.New("sc-pad-btn").Label("Padded").Primary(),
				).Horizontal(10).Vertical(2),
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
			"Components: Button  Checkbox  Switch  Radio  Slider  Label  Badge" +
				"     Layout: Flex  Grid  Spacer  Padding",
		)
		ctx.SetColor(th.Colors.MutedForeground)
		ctx.DrawStringAnchored(statusLine, wW/2, row3Y+10, 0.5, 0.5)
	})
}
