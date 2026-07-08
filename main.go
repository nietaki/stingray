package main

import (
	_ "embed"
	"fmt"

	rgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nietaki/stingray/internal/layout"
	// "github.com/nietaki/stingray/internal/gui"
)

// //go:embed assets/shaders/raymarching.fs
// var raymarchingShaderText string

func main() {
	const (
		screenWidth  = 1280
		screenHeight = 720
		panelWidth   = 220
	)

	// rl.SetConfigFlags(rl.FlagWindowUndecorated | rl.FlagWindowMousePassthrough)
	rl.SetConfigFlags(rl.FlagWindowAlwaysRun)

	rl.InitWindow(screenWidth, screenHeight, "stingray - control layout experiments")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	var (
		// Custom GUI font loading
		// font rl.Font = rl.LoadFontEx("fonts/rainyhearts16.ttf", 12, nil, 0)

		exitWindow bool = false

		// values
		paperSizeIdx int32 = 1
		paperAspect  int32 = 1
	)

	// rl.GuiSetFont(font)

	// panel := gui.NewPanel(rl.NewRectangle(screenWidth-panelWidth, 0, panelWidth, screenHeight))
	layoutRoot := layout.NewVStack(
		"root",
		layout.NewHFlex("",
			layout.Label("paperSizeLabel"),
			layout.Control("paperSize"),
		),

		layout.NewHFlex("",
			layout.Label("paperOrientationLabel"),
			layout.Control("paperOrientation"),
		),
	)

	widgetRectangles := make(map[string]rl.Rectangle)
	var cb layout.WidgetCallback
	cb = func(widget layout.Widget, bounds rl.Rectangle) {
		widgetRectangles[widget.GetId()] = bounds
	}
	panelBounds := rl.NewRectangle(screenWidth-panelWidth, 0, panelWidth, screenHeight)
	layoutRoot.Arrange(panelBounds, cb)

	getRect := func(id string) rl.Rectangle {
		if rect, ok := widgetRectangles[id]; ok {
			return rect
		}
		panic(fmt.Sprintf("Widget with id %s not found", id))
	}

	// Main game loop
	for !exitWindow { // Detect window close button or ESC key
		// Update
		exitWindow = rl.WindowShouldClose()
		if rl.IsKeyPressed(rl.KeyEscape) {
			exitWindow = true
		}

		// DRAWING
		rl.BeginDrawing()
		rl.ClearBackground(rl.GetColor(uint(rgui.GetStyle(rgui.DEFAULT, rgui.BACKGROUND_COLOR))))
		rgui.Label(getRect("paperSizeLabel"), "paper size")
		rgui.ComboBox(getRect("paperSize"), "A5;A4;A3", &paperSizeIdx)
		rgui.ComboBox(getRect("paperOrientation"), "portrait;landscape", &paperAspect)

		// STATUS_BAR
		mousePos := rl.GetMousePosition()
		mousePosText := fmt.Sprintf("Mouse Position: (%.0f, %.0f)", mousePos.X, mousePos.Y)
		rgui.StatusBar(rl.NewRectangle(0, float32(rl.GetScreenHeight())-20, float32(rl.GetScreenWidth()), 20), mousePosText)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.CloseWindow() // Close window and OpenGL context
}
