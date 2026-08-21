package main

import (
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

func main() {
	const w, h = 800, 600

	config := gogpu.DefaultConfig().
		WithTitle("gogpui Components Showcase").
		WithSize(w, h)

	app := gogpu.NewApp(config)

	var canvas *ggcanvas.Canvas
	var lastScale float64
	uictx := context.NewUIContext()

	app.OnSurfaceAvailable(func() {
		lastScale = app.ScaleFactor()
		canvas = ggcanvas.MustNewWithScale(app.GPUContextProvider(), w, h, lastScale)
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

	var sliderVal1 float64 = 0.3
	var sliderVal2 float64 = 0.8

	var chk1Val bool = false
	var chk2Val bool = true
	var chk3Val bool = false

	var sw1Val bool = false
	var sw2Val bool = true
	var sw3Val bool = false

	var rd1Val bool = false
	var rd2Val bool = true
	var rd3Val bool = false

	app.OnDraw(func(dc *gogpu.Context) {
		if canvas == nil {
			return
		}

		// Calculate Delta Time
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

		// Background
		ctx.SetColor(th.Colors.Background)
		ctx.Clear()

		// Load font to enable text rendering
		fontPath := filepath.Join("assets", "fonts", "Inter-Regular.ttf")
		source, err := text.NewFontSourceFromFile(fontPath)
		if err == nil {
			// Force enable 'kern' and 'liga' OpenType features
			ctx.SetFont(source.Face(float64(th.FontSize),
				text.WithFeatures(text.NewFontFeature("kern", 1), text.NewFontFeature("liga", 1)),
			))
		} else {
			log.Printf("Warning: Failed to load font from %s: %v", fontPath, err)
		}

		ctx.SetTextMode(gg.TextModeMSDF)

		// Title
		ctx.SetColor(th.Colors.Foreground)
		ctx.DrawStringAnchored("gogpui Components Showcase", float64(w)/2, 15, 0.5, 0.5)

		// Update UI Context with current frame's input and delta time
		uictx.Update(ctx, th, dt, in, lastScale)

		// Override volatile input state with buffered state
		if pendingLeftPressed {
			uictx.Mouse.LeftPressed = true
			pendingLeftPressed = false
		}
		if pendingLeftReleased {
			uictx.Mouse.LeftReleased = true
			pendingLeftReleased = false
		}

		// Draw buttons using Flexbox layout
		btn1 := button.New("showcase-btn-primary").Label("Primary Button").Primary()
		btn2 := button.New("showcase-btn-danger").Label("Danger Button").Danger()
		btn3 := button.New("showcase-btn-link").Label("Link Button").Link()

		chk1 := checkbox.New("showcase-chk-default").Label("Default Checkbox").Checked(chk1Val).OnChange(func(val bool) { chk1Val = val; println(" chk1 clicked:", val) })
		chk2 := checkbox.New("showcase-chk-checked").Label("Checked Checkbox").Checked(chk2Val).OnChange(func(val bool) { chk2Val = val })
		chk3 := checkbox.New("showcase-chk-disabled").Label("Disabled Checkbox").Checked(chk3Val).Disabled(true).OnChange(func(val bool) { chk3Val = val })

		sw1 := switch_comp.New("showcase-sw-default").Label("Default Switch").Checked(sw1Val).OnChange(func(val bool) { sw1Val = val })
		sw2 := switch_comp.New("showcase-sw-checked").Label("Checked Switch").Checked(sw2Val).OnChange(func(val bool) { sw2Val = val })
		sw3 := switch_comp.New("showcase-sw-disabled").Label("Disabled Switch").Checked(sw3Val).Disabled(true).OnChange(func(val bool) { sw3Val = val })

		rd1 := radio.New("showcase-rd-default").Label("Default Radio").Checked(rd1Val).OnChange(func(val bool) { rd1Val = val })
		rd2 := radio.New("showcase-rd-checked").Label("Checked Radio").Checked(rd2Val).OnChange(func(val bool) { rd2Val = val })
		rd3 := radio.New("showcase-rd-disabled").Label("Disabled Radio").Checked(rd3Val).Disabled(true).OnChange(func(val bool) { rd3Val = val })

		sl1 := slider.New("showcase-sl-default").Value(sliderVal1).OnChange(func(val float64) { sliderVal1 = val })
		sl2 := slider.New("showcase-sl-checked").Value(sliderVal2).OnChange(func(val float64) { sliderVal2 = val })
		sl3 := slider.New("showcase-sl-disabled").Value(0.7).Disabled(true)

		lbl1 := label.New("Normal Label")
		lbl2 := label.New("Secondary Label").Secondary("with extra info")
		lbl3 := label.New("Highlighted Label").Highlights("light", false)
		lbl4 := label.New("Masked Secret").Masked(true)
		lbl5 := label.New("Button Label").Masked(true)

		bdg1 := badge.New().Count(5).Add(label.New("Messages"))
		bdg2 := badge.New().Count(120).Max(99).Add(label.New("Notifications"))
		bdg3 := badge.New().Dot().Add(label.New("Updates"))
		bdg4 := badge.New().Count(3).Color(th.Colors.Info).Add(label.New("Info Badge"))
		bdg5 := badge.New().Count(3).Color(th.Colors.Info).Add(button.New("showcase-btn-label").Label("Badge Button"))

		layout.NewFlex().
			Direction(layout.Column).
			Gap(10).
			Add(btn1, btn2, btn3).
			Render(uictx, 50, 50)

		layout.NewFlex().
			Direction(layout.Column).
			Gap(10).
			Add(chk1, chk2, chk3).
			Render(uictx, 300, 50)

		layout.NewFlex().
			Direction(layout.Column).
			Gap(10).
			Add(sw1, sw2, sw3).
			Render(uictx, 550, 50)

		layout.NewFlex().
			Direction(layout.Column).
			Gap(10).
			Add(rd1, rd2, rd3).
			Render(uictx, 50, 250)

		layout.NewFlex().
			Direction(layout.Column).
			Gap(10).
			Add(sl1, sl2, sl3).
			Render(uictx, 250, 250)

		layout.NewFlex().
			Direction(layout.Column).
			Gap(20).
			Add(
				layout.NewFlex().Direction(layout.Column).Gap(10).Add(lbl1, lbl2, lbl3, lbl4, lbl5),
				layout.NewFlex().Direction(layout.Column).Gap(20).Add(bdg1, bdg2, bdg3, bdg4, bdg5),
			).
			Render(uictx, 550, 200)

		// Render the gg canvas texture to the window
		canvas.MarkDirty()
		err = canvas.RenderTo(dc.AsTextureDrawer())
		if err != nil {
			log.Printf("Failed to render canvas to screen: %v", err)
		}

		// Request a redraw ONLY if an animation is in progress.
		// gogpu automatically calls OnDraw when window events (mouse, resize) occur!
		if uictx.NeedsRedraw {
			app.RequestRedraw()
		}
	})

	log.Println("Starting gogpu window...")
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
