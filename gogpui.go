// Package gogpui は GPU UI コンポーネントライブラリのトップレベル API を提供する。
//
// ユーザーは gogpu / ggcanvas / gg などの内部ライブラリを直接インポートする必要はない。
// 代わりに New と Run だけでウィンドウを起動できる。
//
// 使用例:
//
//	app := gogpui.New(gogpui.Options{Title: "My App", Width: 800, Height: 600})
//	app.Run(func(uictx *context.UIContext) {
//	    btn := button.New("ok").Primary().Label("Let's Go!")
//	    layout.NewFlex().Direction(layout.Column).Gap(8).Add(btn).Render(uictx, 50, 50)
//	})
package gogpui

import (
	"log"
	"path/filepath"
	"time"

	"github.com/gogpu/gg"
	"github.com/gogpu/gg/integration/ggcanvas"
	"github.com/gogpu/gg/text"
	"github.com/gogpu/gogpu"
	"github.com/gogpu/gogpu/input"
	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/theme"
)

// Options はアプリケーションウィンドウの設定を保持する。
type Options struct {
	// Title はウィンドウタイトルバーに表示される文字列。
	Title string
	// Width / Height はウィンドウの論理ピクセルサイズ。デフォルトは 800x600。
	Width, Height int
	// FontPath はフォントファイルへのパス。
	// 未指定の場合は "assets/fonts/Inter-Regular.ttf" を使用する。
	FontPath string
}

// App は gogpui アプリケーションを表す。
// 内部で gogpu / ggcanvas / gg/text の初期化・フレームループを管理する。
type App struct {
	opts Options
}

// New はオプションを受け取り、新しい App を返す。
// この関数はウィンドウを開かない。Run を呼び出すことでウィンドウが起動する。
func New(opts Options) *App {
	if opts.Width == 0 {
		opts.Width = 800
	}
	if opts.Height == 0 {
		opts.Height = 600
	}
	if opts.Title == "" {
		opts.Title = "gogpui App"
	}
	if opts.FontPath == "" {
		opts.FontPath = filepath.Join("assets", "fonts", "Inter-Regular.ttf")
	}

	// サブピクセルテキストメトリクスを有効化する。
	// コンテキスト生成前に呼ぶことが必須。
	text.SetGlobalSubpixelCache(text.NewSubpixelCache(text.HighQualitySubpixelConfig()))

	return &App{opts: opts}
}

// Run はウィンドウを起動し、毎フレーム onDraw を呼び出す。
// onDraw は描画準備済みの *context.UIContext を受け取る。
// gogpu / ggcanvas / gg などのインポートは不要。
func (a *App) Run(onDraw func(*context.UIContext)) {
	config := gogpu.DefaultConfig().
		WithTitle(a.opts.Title).
		WithSize(a.opts.Width, a.opts.Height)

	app := gogpu.NewApp(config)

	// フォントソースを一度だけ読み込む。
	// 毎フレーム再読み込みすると MSDF/Vector アトラスが毎回初期化され文字間隔が不安定になる。
	fontSource, fontLoadErr := text.NewFontSourceFromFile(a.opts.FontPath)
	if fontLoadErr != nil {
		log.Printf("gogpui: failed to load font from %s: %v", a.opts.FontPath, fontLoadErr)
	}

	var canvas *ggcanvas.Canvas
	var lastScale float64
	uictx := context.NewUIContext()

	app.OnSurfaceAvailable(func() {
		lastScale = app.ScaleFactor()
		canvas = ggcanvas.MustNewWithScale(app.GPUContextProvider(), a.opts.Width, a.opts.Height, lastScale)
	})

	// --- 入力バッファリング ---
	var lastMx, lastMy float32
	var lastLDown bool
	var pendingLeftPressed bool
	var pendingLeftReleased bool

	app.OnUpdate(func(_ float64) {
		in := app.Input()
		if in == nil {
			return
		}
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
	})

	lastTime := time.Now()

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

		// 背景クリア
		ctx.SetColor(th.Colors.Background)
		ctx.Clear()

		// フォントセット（ソースは起動時に一度だけロード済み）
		if fontSource != nil {
			ctx.SetFont(fontSource.Face(float64(th.FontSize),
				text.WithHinting(text.HintingNone),
				text.WithFeatures(text.NewFontFeature("kern", 1), text.NewFontFeature("liga", 1)),
			))
		}
		// TextModeVector: UI ラベル・品質重視テキスト向け（gg ドキュメント推奨）
		ctx.SetTextMode(gg.TextModeVector)

		uictx.Update(ctx, th, dt, in, lastScale)

		// 揮発性入力をバッファから反映
		if pendingLeftPressed {
			uictx.Mouse.LeftPressed = true
			pendingLeftPressed = false
		}
		if pendingLeftReleased {
			uictx.Mouse.LeftReleased = true
			pendingLeftReleased = false
		}

		// ユーザー定義の描画コールバックを呼び出す
		if onDraw != nil {
			onDraw(uictx)
		}

		// canvas → ウィンドウへ転送
		canvas.MarkDirty()
		if err := canvas.RenderTo(dc.AsTextureDrawer()); err != nil {
			log.Printf("gogpui: render error: %v", err)
		}

		if uictx.NeedsRedraw {
			app.RequestRedraw()
		}
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
