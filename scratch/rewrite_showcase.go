package main

import (
	"fmt"
	"os"
)

func main() {
	path := "cmd/showcase/main.go"
	content := "package main

import (
	\"fmt\"

	gogpui \"github.com/sofiagros/gogpui\"
	\"github.com/sofiagros/gogpui/components/avatar\"
	\"github.com/sofiagros/gogpui/components/badge\"
	\"github.com/sofiagros/gogpui/components/breadcrumb\"
	\"github.com/sofiagros/gogpui/components/button\"
	\"github.com/sofiagros/gogpui/components/checkbox\"
	\"github.com/sofiagros/gogpui/components/collapsible\"
	\"github.com/sofiagros/gogpui/components/debugoverlay\"
	\"github.com/sofiagros/gogpui/components/label\"
	\"github.com/sofiagros/gogpui/components/progress\"
	\"github.com/sofiagros/gogpui/components/radio\"
	\"github.com/sofiagros/gogpui/components/separator\"
	\"github.com/sofiagros/gogpui/components/skeleton\"
	\"github.com/sofiagros/gogpui/components/slider\"
	switch_comp \"github.com/sofiagros/gogpui/components/switch\"
	\"github.com/sofiagros/gogpui/components/tree\"
	\"github.com/sofiagros/gogpui/core/context\"
	\"github.com/sofiagros/gogpui/core/layout\"
	\"github.com/sofiagros/gogpui/core/theme\"
)

const (
	winW = 1200
	winH = 1000
)

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
		tree.NewItem(\"root1\", \"RootProject\").Expanded(true).Child(
			tree.NewItem(\"folder1\", \"Folder\").Child(
				tree.NewItem(\"file1\", \"main.go\").Suffix(label.New(\"M\").Color(theme.DefaultTheme().Colors.Warning)),
			),
		).Child(
			tree.NewItem(\"folder2\", \"Folder\").Child(
				tree.NewItem(\"folder3\", \"Folder\").Child(
					tree.NewItem(\"file2\", \"a.go\").Suffix(label.New(\"M\").Color(theme.DefaultTheme().Colors.Warning)),
				),
			),
		),
		tree.NewItem(\"root2\", \"Root 2 (Disabled)\").Disabled(true),
	)

	navTreeState := tree.NewTreeState().Items(
		tree.NewItem(\"forms\", \"Forms\").Expanded(true).Child(
			tree.NewItem(\"btn\", \"Button\"),
			tree.NewItem(\"chk\", \"Checkbox\"),
			tree.NewItem(\"sw\", \"Switch\"),
			tree.NewItem(\"rad\", \"Radio\"),
			tree.NewItem(\"sld\", \"Slider\"),
		),
		tree.NewItem(\"display\", \"Data Display\").Expanded(true).Child(
			tree.NewItem(\"lbl\", \"Label\"),
			tree.NewItem(\"bdg\", \"Badge\"),
			tree.NewItem(\"avt\", \"Avatar\"),
			tree.NewItem(\"tree\", \"Tree\"),
		),
		tree.NewItem(\"feedback\", \"Feedback\").Expanded(true).Child(
			tree.NewItem(\"prog\", \"Progress\"),
			tree.NewItem(\"skel\", \"Skeleton\"),
		),
		tree.NewItem(\"nav\", \"Navigation\").Expanded(true).Child(
			tree.NewItem(\"brd\", \"Breadcrumb\"),
			tree.NewItem(\"col\", \"Collapsible\"),
		),
		tree.NewItem(\"layout\", \"Layout\").Expanded(true).Child(
			tree.NewItem(\"sep\", \"Separator\"),
		),
	)

	app := gogpui.New(gogpui.Options{
		Title:  \"gogpui Component Showcase\",
		Width:  winW,
		Height: winH,
	})

	overlay := debugoverlay.New()

	app.Run(func(uictx *context.UIContext) {
		th := uictx.Theme
		ctx := uictx.GG
		wW := uictx.WindowWidth
		wH := uictx.WindowHeight
		if wW == 0 {
			wW = 1200
		}
		if wH == 0 {
			wH = 1000
		}

		selectedID := \"btn\"
		if navTreeState.SelectedIndex != nil && *navTreeState.SelectedIndex >= 0 && *navTreeState.SelectedIndex < len(navTreeState.Entries) {
			selectedID = navTreeState.Entries[*navTreeState.SelectedIndex].Item.ID()
		}

		contentWidget := &layout.CustomWidget{
			RenderFunc: func(uictx *context.UIContext, startX, startY float64) (float64, float64) {
				const (
					padX    = 24.0
					padY    = 14.0
					titleH  = 44.0
					cardPad = 14.0
					secH    = 22.0
				)

				// Header
				if !uictx.MeasureOnly {
					ctx.SetColor(th.Colors.Foreground)
					ctx.DrawStringAnchored(\"gogpui Component Showcase\", startX+wW/2, startY+titleH/2, 0.5, 0.5)
					separator.Horizontal().Length(wW-padX*2).Render(uictx, startX+padX, startY+titleH-1)
				}

				sidebarW := 250.0
				mainW := wW - sidebarW - padX*3

				sidebar := layout.NewFlex().Direction(layout.Column).Gap(10).Add(
					label.New(\"Components\").Color(th.Colors.MutedForeground),
					tree.New(navTreeState).Width(sidebarW),
				)

				drawCard := func(x, y, w, h float64, title string) {
					if uictx.MeasureOnly {
						return
					}
					ctx.SetColor(th.Colors.Background)
					ctx.DrawRoundedRectangle(x, y, w, h, 8)
					ctx.Fill()
					ctx.SetColor(th.Colors.Border)
					ctx.SetLineWidth(1.0)
					ctx.DrawRoundedRectangle(x, y, w, h, 8)
					ctx.Stroke()
					ctx.SetColor(th.Colors.MutedForeground)
					ctx.DrawString(title, x+cardPad, y+cardPad+float64(th.FontSize)*0.85)
				}

				var mainContent layout.Widget
				
				switch selectedID {
				case \"btn\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 182, \"Button\")
						layout.NewFlex().Direction(layout.Column).Gap(8).Add(
							button.New(\"sc-btn-primary\").Label(\"Primary Button\").Primary(),
							button.New(\"sc-btn-danger\").Label(\"Danger Button\").Danger(),
							button.New(\"sc-btn-link\").Label(\"Link Button\").Link(),
						).Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 182
					}}
				case \"chk\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 182, \"Checkbox\")
						layout.NewFlex().Direction(layout.Column).Gap(10).Add(
							checkbox.New(\"sc-chk1\").Label(\"Unchecked\").Checked(chk1Val).OnChange(func(v bool) { chk1Val = v }),
							checkbox.New(\"sc-chk2\").Label(\"Checked\").Checked(chk2Val).OnChange(func(v bool) { chk2Val = v }),
							checkbox.New(\"sc-chk3\").Label(\"Disabled\").Checked(false).Disabled(true),
						).Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 182
					}}
				case \"sw\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 182, \"Switch\")
						layout.NewFlex().Direction(layout.Column).Gap(10).Add(
							switch_comp.New(\"sc-sw1\").Label(\"Off state\").Checked(sw1Val).OnChange(func(v bool) { sw1Val = v }),
							switch_comp.New(\"sc-sw2\").Label(\"On state\").Checked(sw2Val).OnChange(func(v bool) { sw2Val = v }),
							switch_comp.New(\"sc-sw3\").Label(\"Disabled\").Checked(false).Disabled(true),
						).Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 182
					}}
				case \"rad\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 182, \"Radio\")
						layout.NewFlex().Direction(layout.Column).Gap(10).Add(
							radio.New(\"sc-rd1\").Label(\"Unselected\").Checked(rd1Val).OnChange(func(v bool) { rd1Val = v }),
							radio.New(\"sc-rd2\").Label(\"Selected\").Checked(rd2Val).OnChange(func(v bool) { rd2Val = v }),
							radio.New(\"sc-rd3\").Label(\"Disabled\").Checked(false).Disabled(true),
						).Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 182
					}}
				case \"sld\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 200, \"Slider\")
						layout.NewFlex().Direction(layout.Column).Gap(14).Add(
							slider.New(\"sc-sl1\").Value(sliderVal1).OnChange(func(v float64) { sliderVal1 = v }),
							slider.New(\"sc-sl2\").Value(sliderVal2).OnChange(func(v float64) { sliderVal2 = v }),
							slider.New(\"sc-sl3\").Value(0.5).Disabled(true),
						).Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 200
					}}
				case \"lbl\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 200, \"Label\")
						layout.NewFlex().Direction(layout.Column).Gap(10).Add(
							label.New(\"Normal label\"),
							label.New(\"With secondary\").Secondary(\"extra info\"),
							label.New(\"Highlighted text\").Highlights(\"igh\", false),
							label.New(\"Masked content\").Masked(true),
						).Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 200
					}}
				case \"bdg\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 200, \"Badge\")
						layout.NewFlex().Direction(layout.Column).Gap(12).Add(
							badge.New().Count(5).Add(label.New(\"Messages\")),
							badge.New().Count(120).Max(99).Add(label.New(\"Capped at 99+\")),
							badge.New().Dot().Add(label.New(\"Dot variant\")),
							badge.New().Count(3).Color(th.Colors.Info).Add(label.New(\"Info color\")),
						).Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 200
					}}
				case \"sep\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 200, \"Separator & Layout\")
						innerW := mainW - cardPad*2
						layout.NewFlex().Direction(layout.Column).Gap(10).Add(
							layout.NewFlex().Direction(layout.Row).WithConstraints(innerW, 0).Add(label.New(\"Left\"), layout.NewSpacerFlex(), label.New(\"Right\")),
							separator.Horizontal().Length(innerW),
							separator.HorizontalDashed().Label(\"OR\").Length(innerW),
							layout.NewGrid().Cols(3).Gap(6).Add(
								button.New(\"sc-g1\").Label(\"A\").Primary(),
								button.New(\"sc-g2\").Label(\"B\").Danger(),
								button.New(\"sc-g3\").Label(\"C\").Link(),
							),
						).Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 200
					}}
				case \"prog\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 200, \"Progress\")
						progW := mainW - cardPad*2
						layout.NewFlex().Direction(layout.Column).Gap(20).Add(
							progress.New(\"sc-prog-det\").Value(sliderVal1*100).Width(progW),
							progress.New(\"sc-prog-loading\").Loading(true).Width(progW),
							progress.New(\"sc-prog-sm\").Value(75).Size(progress.SizeSmall).Width(progW).Color(th.Colors.Info),
							layout.NewFlex().Direction(layout.Row).Gap(10).Add(
								progress.NewCircle(\"sc-circ-det\").Value(sliderVal1*100),
								progress.NewCircle(\"sc-circ-load\").Loading(true),
								progress.NewCircle(\"sc-circ-lg\").Value(85).Size(progress.SizeLarge).Color(th.Colors.Success),
								progress.NewCircle(\"sc-circ-xs\").Value(40).Size(progress.SizeXSmall).Color(th.Colors.Danger),
							),
						).Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 200
					}}
				case \"skel\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 140, \"Skeleton\")
						progW := mainW - cardPad*2
						layout.NewFlex().Direction(layout.Column).Gap(10).Add(
							skeleton.New(\"sc-skel-1\").Width(progW),
							skeleton.New(\"sc-skel-2\").Width(progW).Height(24),
							skeleton.New(\"sc-skel-3\").Width(progW).Secondary(true),
						).Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 140
					}}
				case \"brd\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 140, \"Breadcrumb\")
						layout.NewFlex().Direction(layout.Column).Gap(16).Add(
							breadcrumb.New().Add(
								breadcrumb.NewItem(\"sc-bc-1\", \"Home\").OnClick(func() { fmt.Println(\"Breadcrumb Home clicked\") }),
								breadcrumb.NewItem(\"sc-bc-2\", \"Components\"),
							),
							breadcrumb.New().Add(
								breadcrumb.NewItem(\"sc-bc-3\", \"Home\").OnClick(func() {}),
								breadcrumb.NewItem(\"sc-bc-4\", \"Docs\").OnClick(func() {}),
								breadcrumb.NewItem(\"sc-bc-5\", \"Layout\"),
							),
						).Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 140
					}}
				case \"avt\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 200, \"Avatar\")
						layout.NewFlex().Direction(layout.Row).Gap(15).Add(
							avatar.New().Name(\"Jason Lee\").Size(avatar.SizeLarge),
							avatar.New().Name(\"foo bar\"),
							avatar.New().Name(\"huacnlee\").Size(avatar.SizeSmall),
							avatar.New().Name(\"small a\").Size(avatar.SizeXSmall),
							avatar.New(),
						).Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 200
					}}
				case \"col\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 200, \"Collapsible\")
						toggleBtn := button.New(\"toggle-col\").Ghost().Label(\"Toggle Details\").OnClick(func() { isSkelOpen = !isSkelOpen })
						skelContent := layout.NewFlex().Direction(layout.Column).Gap(5).Add(
							skeleton.New(\"skel1\").Width(150).Height(12),
							skeleton.New(\"skel2\").Width(120).Height(12),
							skeleton.New(\"skel3\").Width(180).Height(12),
						)
						collapsible.New().Open(isSkelOpen).Trigger(toggleBtn).Content(skelContent).Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 200
					}}
				case \"tree\":
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 200, \"Tree\")
						treeWidget := tree.New(demoTreeState).Width(mainW - cardPad*2)
						treeWidget.Render(uictx, cx+cardPad, cy+secH+cardPad)
						return mainW, 200
					}}
				default:
					mainContent = &layout.CustomWidget{RenderFunc: func(uictx *context.UIContext, cx, cy float64) (float64, float64) {
						drawCard(cx, cy, mainW, 100, \"Select a component\")
						return mainW, 100
					}}
				}

				layout.NewFlex().Direction(layout.Row).Gap(padX).Add(
					sidebar,
					mainContent,
				).Render(uictx, startX+padX, startY+titleH+padY)

				return wW, wH
			},
		}

		contentWidget.Render(uictx, 0, 0)
		overlay.Render(uictx, float64(wW)-200.0, 10.0)
	})
}
"
	
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Showcase rewritten.")
	}
}
