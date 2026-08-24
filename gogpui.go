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
	_ "embed"
	"fmt"
	"image"
	"log"
	"sync/atomic"
	"time"

	"github.com/gogpu/gg"
	_ "github.com/gogpu/gg/gpu" // Enable GPU acceleration
	"github.com/gogpu/gg/integration/ggcanvas"
	"github.com/gogpu/gg/text"
	"github.com/gogpu/gogpu"
	"github.com/gogpu/gogpu/input"
	"github.com/gogpu/gpucontext"
	"github.com/sofiagros/gogpui/core/context"
	"github.com/sofiagros/gogpui/core/theme"
)

//go:embed assets/fonts/Inter-Regular.ttf
var defaultFontData []byte

// Options はアプリケーションウィンドウの設定を保持する。
type Options struct {
	// Title はウィンドウタイトルバーに表示される文字列。
	Title string
	// Width / Height はウィンドウの論理ピクセルサイズ。デフォルトは 800x600。
	Width, Height int
	// FontPath はフォントファイルへのパス。
	// 未指定の場合はライブラリに組み込まれたデフォルトフォント(Inter)を使用する。
	FontPath string
	// Frameless はOSのウィンドウフレーム（タイトルバーや境界線）を非表示にするかどうかを指定する。
	Frameless bool
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

	// サブピクセルテキストメトリクスを有効化する。
	// コンテキスト生成前に呼ぶことが必須。
	text.SetGlobalSubpixelCache(text.NewSubpixelCache(text.HighQualitySubpixelConfig()))

	return &App{opts: opts}
}

type gogpuWindowControl struct {
	app   *gogpu.App
	uictx *context.UIContext
}

func (wc *gogpuWindowControl) Minimize() {
	if chrome, ok := any(wc.app).(gpucontext.WindowChrome); ok {
		chrome.Minimize()
	}
}

func (wc *gogpuWindowControl) Maximize() {
	if chrome, ok := any(wc.app).(gpucontext.WindowChrome); ok {
		chrome.Maximize()
	}
}

func (wc *gogpuWindowControl) IsMaximized() bool {
	if chrome, ok := any(wc.app).(gpucontext.WindowChrome); ok {
		return chrome.IsMaximized()
	}
	return false
}

func (wc *gogpuWindowControl) Close() {
	if chrome, ok := any(wc.app).(gpucontext.WindowChrome); ok {
		chrome.Close()
	}
}

func (wc *gogpuWindowControl) AddDragRegion(rect image.Rectangle) {
	wc.uictx.AddDragRegion(rect)
}

// Run はウィンドウを起動し、毎フレーム onDraw を呼び出す。
// onDraw は描画準備済みの *context.UIContext を受け取る。
// gogpu / ggcanvas / gg などのインポートは不要。
func (a *App) Run(onDraw func(*context.UIContext)) {
	config := gogpu.DefaultConfig().
		WithTitle(a.opts.Title).
		WithSize(a.opts.Width, a.opts.Height).
		WithFrameless(a.opts.Frameless)

	app := gogpu.NewApp(config)

	// フォントソースを一度だけ読み込む。
	// 毎フレーム再読み込みすると MSDF/Vector アトラスが毎回初期化され文字間隔が不安定になる。
	var fontSource *text.FontSource
	var fontLoadErr error
	if a.opts.FontPath == "" {
		fontSource, fontLoadErr = text.NewFontSource(defaultFontData)
	} else {
		fontSource, fontLoadErr = text.NewFontSourceFromFile(a.opts.FontPath)
	}
	if fontLoadErr != nil {
		log.Printf("gogpui: failed to load font: %v", fontLoadErr)
	}

	var canvas *ggcanvas.Canvas

	// targetW / targetH はリサイズイベント（メインスレッド）から
	// 描画処理（レンダースレッド）へサイズを安全に渡すためのアトミック変数。
	// canvas.Resize() は必ず OnDraw（レンダースレッド）内で呼ぶ。
	var targetW, targetH atomic.Int32
	var targetScale atomic.Value // float64 を保持

	// レンダースレッド側で管理するローカル変数（OnDraw 内のみで参照）
	var currentW, currentH int
	var lastScale float64

	var cachedFace text.Face
	var lastFontSize float32

	uictx := context.NewUIContext()
	uictx.Window = &gogpuWindowControl{app: app, uictx: uictx}

	if chrome, ok := any(app).(gpucontext.WindowChrome); ok {
		chrome.SetHitTestCallback(func(x, y float64) gpucontext.HitTestResult {
			regions := uictx.GetDragRegions()
			pt := image.Point{X: int(x), Y: int(y)}
			for _, r := range regions {
				if pt.In(r) {
					return gpucontext.HitTestCaption
				}
			}
			return gpucontext.HitTestClient
		})
	}

	app.OnSurfaceAvailable(func() {
		scale := app.ScaleFactor()
		w, h := app.Size()
		if w <= 0 || h <= 0 {
			// OnSurfaceAvailable 時にまだサイズが取れない場合はオプション初期値を使用する
			w, h = a.opts.Width, a.opts.Height
		}
		targetW.Store(int32(w))
		targetH.Store(int32(h))
		targetScale.Store(scale)
		if canvas == nil {
			canvas = ggcanvas.MustNewWithScale(app.GPUContextProvider(), w, h, scale)
		}

		app.RequestRedraw() // サーフェース有効化後に必ず再描画する
	})

	// --- ウィンドウリサイズ対応 ---
	// OnResize はメインスレッドで呼ばれる。
	// アトミック変数に新サイズを書き込み、OnDraw 側で canvas.Resize() を適用する。
	app.OnResize(func(w, h int) {
		targetW.Store(int32(w))
		targetH.Store(int32(h))
		targetScale.Store(app.ScaleFactor())
		app.RequestRedraw()
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
	var frameCount int
	var lastFPSLog time.Time = time.Now()

	app.OnDraw(func(dc *gogpu.Context) {
		if canvas == nil {
			return
		}

		// --- 最小化チェック ---
		// 最小化時はサイズが 0x0 になる。GPUサーフェース再構成と canvas.Resize() で
		// エラー（hal: surface width and height must be non-zero）になるため、
		// サイズが有効になるまで描画・スワップチェーン取得を完全にスキップする。
		if isMin, clientW, clientH := checkWindowMinimized(a.opts.Title); isMin {
			targetW.Store(0)
			targetH.Store(0)
			currentW, currentH = 0, 0
			return
		} else if clientW > 0 && clientH > 0 {
			targetW.Store(int32(clientW))
			targetH.Store(int32(clientH))
		}

		winW, winH := app.Size()
		if winW <= 0 || winH <= 0 {
			targetW.Store(0)
			targetH.Store(0)
			currentW, currentH = 0, 0
			return
		}

		tgtW := int(targetW.Load())
		tgtH := int(targetH.Load())
		if tgtW <= 0 || tgtH <= 0 {
			tgtW, tgtH = winW, winH
			targetW.Store(int32(tgtW))
			targetH.Store(int32(tgtH))
		}
		newScale := app.ScaleFactor()
		targetScale.Store(newScale)

		// --- リサイズ処理（レンダースレッド内で安全に適用）---
		isResized := false
		if tgtW != currentW || tgtH != currentH || newScale != lastScale {
			currentW, currentH = tgtW, tgtH
			lastScale = newScale
			if err := canvas.Resize(currentW, currentH); err != nil {
				log.Printf("gogpui: canvas resize error: %v", err)
			}
			canvas.SetDeviceScale(lastScale)
			isResized = true
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

		// ダメージがあるか、初回フレームかを確認
		hasDamage := uictx.NeedsRedraw || len(uictx.DamageRects) > 0 || isResized || frameCount == 0

		if hasDamage {
			// 背景クリア（フル再描画）
			// TODO: 本当は DamageRects の範囲だけクリップしてクリアしたいが、まずはフル再描画で進める
			ctx.SetColor(th.Colors.Background)
			ctx.Clear()
			uictx.LogicOnly = false
		} else {
			// ダメージがなければ、描画をスキップしてロジックだけ実行
			uictx.LogicOnly = true
		}

		// 次のフレームでダメージ処理を行うため、この時点で現在のダメージ領域をクリア
		uictx.DamageRects = uictx.DamageRects[:0]
		uictx.NeedsRedraw = false

		// フォントセット（ソースは起動時に一度だけロード済み）
		if fontSource != nil {
			if cachedFace == nil || th.FontSize != lastFontSize {
				cachedFace = fontSource.Face(float64(th.FontSize),
					text.WithHinting(text.HintingNone),
					text.WithFeatures(text.NewFontFeature("kern", 1), text.NewFontFeature("liga", 1)),
				)
				lastFontSize = th.FontSize
			}
			ctx.SetFont(cachedFace)
		}
		// TextModeGlyphMask: GPUアトラスを用いたUIテキストの高速・高品質レンダリング
		ctx.SetTextMode(gg.TextModeGlyphMask)

		uictx.Update(ctx, th, dt, in, lastScale,
			float64(currentW), float64(currentH))

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
		var onDrawDur, renderDur time.Duration
		if onDraw != nil {
			startOnDraw := time.Now()
			onDraw(uictx)
			onDrawDur = time.Since(startOnDraw)
		}

		// canvas を GPU に転送する。
		// ダメージがあった場合のみアップロードを行う。
		if onDraw != nil && !uictx.LogicOnly {
			canvas.MarkDirty()
		}
		
		startRender := time.Now()
		if err := canvas.Render(dc.RenderTarget()); err != nil {
			fmt.Println("gogpui: canvas.Render error:", err)
			// 「gogpu: surface frame not available」などの過渡的エラー時は次フレームで確実に再試行する。
			app.RequestRedraw()
		}
		renderDur = time.Since(startRender)

		// アニメーション実行中は gogpu に VSync フルレート描画を要求する。
		// 次のフレームで描画が必要とマークされたら、すぐにRedraw要求
		if len(uictx.DamageRects) > 0 || uictx.NeedsRedraw {
			app.RequestRedraw()
		}

		frameCount++
		if time.Since(lastFPSLog) >= time.Second {
			fps := float64(frameCount) / time.Since(lastFPSLog).Seconds()
			fmt.Printf("FPS: %.1f, dt: %.1fms (onDraw: %.1fms, canvas.Render: %.1fms)\n", fps, dt*1000, float64(onDrawDur.Microseconds())/1000.0, float64(renderDur.Microseconds())/1000.0)
			frameCount = 0
			lastFPSLog = time.Now()
		}
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
