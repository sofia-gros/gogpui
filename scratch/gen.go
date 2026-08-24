package main

import (
	"os"
)

func main() {
	content := package main

import (
	"fmt"

	gogpui "github.com/sofiagros/gogpui"
	"github.com/sofiagros/gogpui/components/avatar"
	"github.com/sofiagros/gogpui/components/badge"
	"github.com/sofiagros/gogpui/components/breadcrumb"
	"github.com/sofiagros/gogpui/components/button"
	"github.com/sofiagros/gogpui/components/checkbox"
	"github.com/sofiagros/gogpui/components/collapsible"
	"github.com/sofiagros/gogpui/components/label"
	"github.com/sofiagros/gogpui/components/progress"
	"github.com/sofiagros/gogpui/components/radio"
	"github.com/sofiagros/gogpui/components/separator"
	"github.com/sofiagros/gogpui/components/scroll"
	"github.com/sofiagros/gogpui/components/debugoverlay"
	"github.com/sofiagros/gogpui/components/tree"
	"github.com/sofiagros/gogpui/components/skeleton"
	"github.com/sofiagros/gogpui/components/slider"
	switch_comp "github.com/sofiagros/gogpui/components/switch"
	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/layout"
	"github.com/sofiagros/gogpui/core/theme"
)

const (
	winW = 1200
	winH = 1000
)

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
	var chk1Val bool = false
	var chk2Val bool = true
	var sw1Val bool = false
	var sw2Val bool = true
	var rd1Val bool = false
	var rd2Val bool = true
	var sliderVal1 float64 = 0.3
	var sliderVal2 float64 = 0.7
	isSkelOpen := false

	demoTreeState := tree.NewTreeState().Items(
		tree.NewItem("root1", "RootProject").Expanded(true).Child(
			tree.NewItem("folder1", "Folder").Child(
				tree.NewItem("file1", "main.go").Suffix(label.New("M").Color(theme.DefaultTheme().Colors.Warning)),
			),
		).Child(
			tree.NewItem("folder2", "Folder").Child(
				tree.NewItem("folder3", "Folder").Child(
					tree.NewItem("file2", "a.go").Suffix(label.New("M").Color(theme.DefaultTheme().Colors.Warning)),
				),
			),
		),
		tree.NewItem("root2", "Root 2 (Disabled)").Disabled(true),
	)

	app := gogpui.New(gogpui.Options{
		Title:  "gogpui Component Showcase",
		Width:  winW,
		Height: winH,
	})

	overlay := debugoverlay.New()

	app.Run(func(uictx *context.UIContext) {
		th := uictx.Theme
		ctx := uictx.GG
		wW := uictx.WindowWidth
		wH := uictx.WindowHeight

		contentWidget := &layout.CustomWidget{
			RenderFunc: func(uictx *context.UIContext, startX, startY float64) (float64, float64) {
				const (
					padX    = 24.0
					padY    = 14.0
					colGap  = 16.0
					rowGap  = 16.0
					titleH  = 44.0
					cardPad = 14.0
					secH    = 22.0
					cardH1  = 182.0
					cardH2  = 200.0
					cardH3  = 140.0
					cardH4  = 200.0
				)
				wW := float64(800)
				if uictx.WindowWidth > 0 {
					wW = uictx.WindowWidth
				}
				numCols := 4.0
				colW := (wW - padX*2 - colGap*(numCols-1)) / numCols

				if !uictx.MeasureOnly {
					ctx.SetColor(th.Colors.Foreground)
					ctx.DrawStringAnchored("gogpui Component Showcase", startX+wW/2, startY+titleH/2, 0.5, 0.5)
					separator.Horizontal().Length(wW - padX*2).Render(uictx, startX+padX, startY+titleH-1)
				}

				row1Y := startY + titleH + padY
				col0x := startX + padX
				
				layout.NewCachedWidget("btn-card", &CardWidget{title: "Button", w: colW, h: cardH1, child: layout.NewFlex().Direction(layout.Column).Gap(8).
					Add(
						button.New("sc-btn-primary").Label("Primary Button").Primary(),
						button.New("sc-btn-danger").Label("Danger Button").Danger(),
						button.New("sc-btn-link").Label("Link Button").Link(),
					)}).Render(uictx, col0x, row1Y)

				col1x := startX + padX + colW + colGap
				
				layout.NewCachedWidget("chk-card", &CardWidget{title: "Checkbox", w: colW, h: cardH1, child: layout.NewFlex().Direction(layout.Column).Gap(10).
					Add(
						checkbox.New("sc-chk1").Label("Unchecked").Checked(chk1Val).OnChange(func(v bool) { chk1Val = v }),
						checkbox.New("sc-chk2").Label("Checked").Checked(chk2Val).OnChange(func(v bool) { chk2Val = v }),
						checkbox.New("sc-chk3").Label("Disabled").Checked(false).Disabled(true),
					)}).Render(uictx, col1x, row1Y)

				col2x := startX + padX + (colW+colGap)*2
				layout.NewCachedWidget("sw-card", &CardWidget{title: "Switch", w: colW, h: cardH1, child: layout.NewFlex().Direction(layout.Column).Gap(10).
					Add(
						switch_comp.New("sc-sw1").Label("Off state").Checked(sw1Val).OnChange(func(v bool) { sw1Val = v }),
						switch_comp.New("sc-sw2").Label("On state").Checked(sw2Val).OnChange(func(v bool) { sw2Val = v }),
						switch_comp.New("sc-sw3").Label("Disabled").Checked(false).Disabled(true),
					)}).Render(uictx, col2x, row1Y)

				col3x := startX + padX + (colW+colGap)*3
				layout.NewCachedWidget("rd-card", &CardWidget{title: "Radio", w: colW, h: cardH1, child: layout.NewFlex().Direction(layout.Column).Gap(10).
					Add(
						radio.New("sc-rd1").Label("Unselected").Checked(rd1Val).OnChange(func(v bool) { rd1Val = v }),
						radio.New("sc-rd2").Label("Selected").Checked(rd2Val).OnChange(func(v bool) { rd2Val = v }),
						radio.New("sc-rd3").Label("Disabled").Checked(false).Disabled(true),
					)}).Render(uictx, col3x, row1Y)

				row2Y := row1Y + cardH1 + rowGap
				layout.NewCachedWidget("sl-card", &CardWidget{title: "Slider", w: colW, h: cardH2, child: layout.NewFlex().Direction(layout.Column).Gap(14).
					Add(
						slider.New("sc-sl1").Value(sliderVal1).OnChange(func(v float64) { sliderVal1 = v }),
						slider.New("sc-sl2").Value(sliderVal2).OnChange(func(v float64) { sliderVal2 = v }),
						slider.New("sc-sl3").Value(0.5).Disabled(true),
					)}).Render(uictx, col0x, row2Y)

				layout.NewCachedWidget("lbl-card", &CardWidget{title: "Label", w: colW, h: cardH2, child: layout.NewFlex().Direction(layout.Column).Gap(10).
					Add(
						label.New("Normal label"),
						label.New("With secondary").Secondary("extra info"),
						label.New("Highlighted text").Highlights("igh", false),
						label.New("Masked content").Masked(true),
					)}).Render(uictx, col1x, row2Y)

				layout.NewCachedWidget("bdg-card", &CardWidget{title: "Badge", w: colW, h: cardH2, child: layout.NewFlex().Direction(layout.Column).Gap(12).
					Add(
						badge.New().Count(5).Add(label.New("Messages")),
						badge.New().Count(120).Max(99).Add(label.New("Capped at 99+")),
						badge.New().Dot().Add(label.New("Dot variant")),
						badge.New().Count(3).Color(th.Colors.Info).Add(label.New("Info color")),
					)}).Render(uictx, col2x, row2Y)

				innerW := colW - cardPad*2
				layout.NewCachedWidget("sep-card", &CardWidget{title: "Separator & Layout", w: colW, h: cardH2, child: layout.NewFlex().Direction(layout.Column).Gap(10).
					Add(
						layout.NewFlex().Direction(layout.Row).
							WithConstraints(innerW, 0).
							Add(label.New("Left"), layout.NewSpacerFlex(), label.New("Right")),
						separator.Horizontal().Length(innerW),
						separator.HorizontalDashed().Label("OR").Length(innerW),
						layout.NewGrid().Cols(3).Gap(6).
							Add(
								button.New("sc-g1").Label("A").Primary(),
								button.New("sc-g2").Label("B").Danger(),
								button.New("sc-g3").Label("C").Link(),
							),
					)}).Render(uictx, col3x, row2Y)

				row3Y := row2Y + cardH2 + rowGap
				progW := colW - cardPad*2
				
				(&CardWidget{title: "Progress (Linear)", w: colW, h: cardH3, child: layout.NewFlex().Direction(layout.Column).Gap(10).
					Add(
						progress.New("sc-prog-det").Value(sliderVal1 * 100).Width(progW),
						progress.New("sc-prog-loading").Loading(true).Width(progW),
						progress.New("sc-prog-sm").Value(75).Size(progress.SizeSmall).Width(progW).Color(th.Colors.Info),
					)}).Render(uictx, col0x, row3Y)

				(&CardWidget{title: "Progress (Circle)", w: colW, h: cardH3, child: layout.NewFlex().Direction(layout.Row).Gap(10).
					Add(
						progress.NewCircle("sc-circ-det").Value(sliderVal1 * 100),
						progress.NewCircle("sc-circ-load").Loading(true),
						progress.NewCircle("sc-circ-lg").Value(85).Size(progress.SizeLarge).Color(th.Colors.Success),
						progress.NewCircle("sc-circ-xs").Value(40).Size(progress.SizeXSmall).Color(th.Colors.Danger),
					)}).Render(uictx, col1x, row3Y)

				(&CardWidget{title: "Skeleton", w: colW, h: cardH3, child: layout.NewFlex().Direction(layout.Column).Gap(10).
					Add(
						skeleton.New("sc-skel-1").Width(progW),
						skeleton.New("sc-skel-2").Width(progW).Height(24),
						skeleton.New("sc-skel-3").Width(progW).Secondary(true),
					)}).Render(uictx, col2x, row3Y)

				layout.NewCachedWidget("bc-card", &CardWidget{title: "Breadcrumb", w: colW, h: cardH3, child: layout.NewFlex().Direction(layout.Column).Gap(16).
					Add(
						breadcrumb.New().Add(
							breadcrumb.NewItem("sc-bc-1", "Home").OnClick(func() { fmt.Println("Breadcrumb Home clicked") }),
							breadcrumb.NewItem("sc-bc-2", "Components"),
						),
						breadcrumb.New().Add(
							breadcrumb.NewItem("sc-bc-3", "Home").OnClick(func() {}),
							breadcrumb.NewItem("sc-bc-4", "Docs").OnClick(func() {}),
							breadcrumb.NewItem("sc-bc-5", "Layout"),
						),
					)}).Render(uictx, col3x, row3Y)

				row4Y := row3Y + cardH3 + rowGap
				layout.NewCachedWidget("av-card", &CardWidget{title: "Avatar", w: colW, h: cardH4, child: layout.NewFlex().Direction(layout.Row).Gap(15).
					Add(
						avatar.New().Name("Jason Lee").Size(avatar.SizeLarge),
						avatar.New().Name("foo bar"),
						avatar.New().Name("huacnlee").Size(avatar.SizeSmall),
						avatar.New().Name("small a").Size(avatar.SizeXSmall),
						avatar.New(),
					)}).Render(uictx, col0x, row4Y)

				toggleBtn := button.New("toggle-col").Ghost().Label("Toggle Details").OnClick(func() { isSkelOpen = !isSkelOpen })
				skelContent := layout.NewFlex().Direction(layout.Column).Gap(5).Add(
					skeleton.New("skel1").Width(150).Height(12),
					skeleton.New("skel2").Width(120).Height(12),
					skeleton.New("skel3").Width(180).Height(12),
				)
				(&CardWidget{title: "Collapsible", w: colW, h: cardH4, child: collapsible.New().Open(isSkelOpen).Trigger(toggleBtn).Content(skelContent)}).Render(uictx, col1x, row4Y)

				layout.NewCachedWidget("tree-card", &CardWidget{title: "Tree", w: colW, h: cardH4, child: tree.New(demoTreeState).Width(colW - cardPad*2)}).Render(uictx, col2x, row4Y)

				rowFooterY := row4Y + cardH4 + rowGap
				if !uictx.MeasureOnly {
					statusLine := "Components: Button  Checkbox  Switch  Radio  Slider  Label  Badge  Separator  Progress  Skeleton  Breadcrumb  Avatar  Collapsible  Tree     Layout: Flex  Grid  Spacer  Padding"
					ctx.SetColor(th.Colors.MutedForeground)
					ctx.DrawStringAnchored(statusLine, startX+wW/2, startY+rowFooterY+10, 0.5, 0.5)
				}
				
				return wW, rowFooterY + 30
			},
		}

		scroll.New("main-scroll").Size(float64(wW), float64(wH)).Child(contentWidget).Render(uictx, 0, 0)

		overlay.Render(uictx, float64(wW)-200.0, 10.0)
	})
}

	os.WriteFile("cmd/showcase/main.go", []byte(content), 0644)
}
