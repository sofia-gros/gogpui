package context

import (
	"github.com/gogpu/gg"
	"github.com/gogpu/gogpu/input"
	"github.com/sofiagros/gogpui/core/theme"
)

// MouseState abstracts the mouse input for UI components.
type MouseState struct {
	X, Y        float64
	LeftDown    bool
	LeftPressed bool // True only on the frame it was pressed
	LeftReleased bool // True only on the frame it was released
}

// WidgetState holds persistent state for a widget across frames.
type WidgetState struct {
	IsHovered    bool
	IsActive     bool
	HoverRatio   float64 // 0.0 to 1.0 for animations
	ActiveRatio  float64 // 0.0 to 1.0 for animations
	ToggleRatio  float64 // 0.0 to 1.0 for toggle animations
	ValueRatio   float64 // 0.0 to 1.0 (or custom range) for progress/slider animations
	IsDragging   bool
}

// UIContext wraps gg.Context and provides immediate-mode UI state and input.
type UIContext struct {
	GG           *gg.Context
	Theme        *theme.Theme
	Mouse        MouseState
	DeltaTime    float64
	Scale        float64
	// WindowWidth は現在のウィンドウ論理幅（ピクセル）。
	// リサイズ後に毎フレーム更新される。
	WindowWidth  float64
	// WindowHeight は現在のウィンドウ論理高さ（ピクセル）。
	// リサイズ後に毎フレーム更新される。
	WindowHeight float64
	States       map[string]*WidgetState
	MeasureOnly  bool
	NeedsRedraw  bool
}

// NewUIContext creates a new UIContext.
func NewUIContext() *UIContext {
	return &UIContext{
		States: make(map[string]*WidgetState),
		Scale:  1.0,
	}
}

// Update はフレーム開始時に UIContext を最新状態に更新する。
// windowW / windowH は論理ピクセル単位のウィンドウサイズ（物理サイズ ÷ スケールファクタ）。
func (c *UIContext) Update(ggCtx *gg.Context, th *theme.Theme, dt float64, in *input.State, scale float64, windowW, windowH float64) {
	c.GG = ggCtx
	c.Theme = th
	c.DeltaTime = dt
	c.Scale = scale
	c.WindowWidth = windowW
	c.WindowHeight = windowH
	c.MeasureOnly = false // デフォルトは描画モード
	c.NeedsRedraw = false // フレーム開始時にリセット

	if in != nil {
		mx, my := in.Mouse().Position()
		// ウィンドウのスケールファクタで割り、論理座標に変換する
		c.Mouse.X = float64(mx) / scale
		c.Mouse.Y = float64(my) / scale
		c.Mouse.LeftDown = in.Mouse().Pressed(input.MouseButtonLeft)
		c.Mouse.LeftPressed = in.Mouse().JustPressed(input.MouseButtonLeft)
		c.Mouse.LeftReleased = in.Mouse().JustReleased(input.MouseButtonLeft)
	} else {
		// tester.go など入力をモックで注入する場合は MouseState を直接更新する
	}
}

// GetState retrieves or creates the persistent state for a given widget ID.
func (c *UIContext) GetState(id string) *WidgetState {
	if s, ok := c.States[id]; ok {
		return s
	}
	s := &WidgetState{}
	c.States[id] = s
	return s
}

// HitTest checks if the mouse is currently over the given rectangular bounds.
func (c *UIContext) HitTest(x, y, w, h float64) bool {
	return c.Mouse.X >= x && c.Mouse.X <= x+w && c.Mouse.Y >= y && c.Mouse.Y <= y+h
}

// Animate approaches a target value from current value at a given speed based on DeltaTime.
// If the value is changing, it flags NeedsRedraw to ensure the animation completes.
func (c *UIContext) Animate(current, target, speed float64) float64 {
	if current < target {
		current += c.DeltaTime * speed
		if current > target {
			current = target
		}
		c.NeedsRedraw = true
	} else if current > target {
		current -= c.DeltaTime * speed
		if current < target {
			current = target
		}
		c.NeedsRedraw = true
	}
	return current
}

// ProcessInteraction handles standard hover and active (click) interactions for a widget.
// Returns (isHovered, isActive, isClicked)
func (c *UIContext) ProcessInteraction(id string, x, y, w, h float64, disabled bool) (bool, bool, bool) {
	state := c.GetState(id)
	
	// Skip processing if we are just measuring
	if c.MeasureOnly {
		return state.IsHovered, state.IsActive, false
	}

	isHovered := !disabled && c.HitTest(x, y, w, h)
	isActive := isHovered && c.Mouse.LeftDown
	isClicked := false

	if isHovered && c.Mouse.LeftReleased {
		if state.IsActive { // Was active on this widget
			isClicked = true
		}
	}

	// Update state
	state.IsHovered = isHovered
	state.IsActive = isActive

	// Update animations
	targetHover := 0.0
	if isHovered {
		targetHover = 1.0
	}
	state.HoverRatio = c.Animate(state.HoverRatio, targetHover, 10.0)

	targetActive := 0.0
	if isActive {
		targetActive = 1.0
	}
	state.ActiveRatio = c.Animate(state.ActiveRatio, targetActive, 15.0)

	return isHovered, isActive, isClicked
}
