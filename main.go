package main

import (
	_ "embed"
	"fmt"

	rgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nietaki/stingray/internal/layout"
	"github.com/nietaki/stingray/internal/paper"
	"github.com/nietaki/stingray/internal/state"
)

// //go:embed assets/shaders/raymarching.fs
// var raymarchingShaderText string

const (
	screenWidth  = 1280
	screenHeight = 720
	panelWidth   = 220
)

var (
	// Custom GUI font loading
	// font rl.Font = rl.LoadFontEx("fonts/rainyhearts16.ttf", 12, nil, 0)

	exitWindow bool = false

	// paper
	paperCam rl.Camera2D

	paperPixelDimensions paper.PaperDimensions
	paperConfig          *state.StateManager[paper.Config]
)

func main() {
	paperConfig = state.NewStateManager[paper.Config]()
	paperConfig.LoadOrDefault(paper.DefaultConfig())

	paperPixelDimensions = paperConfig.Data.PixelDims()

	startingOffset := rl.Vector2{X: screenWidth / 2, Y: screenHeight / 2}
	startingTarget := rl.Vector2{
		X: float32(paperPixelDimensions[0]) / 2,
		Y: float32(paperPixelDimensions[1]) / 2}
	paperCam = rl.NewCamera2D(startingOffset, startingTarget, 0.0, 1.0)

	// rl.SetConfigFlags(rl.FlagWindowUndecorated | rl.FlagWindowMousePassthrough)
	rl.SetConfigFlags(rl.FlagWindowAlwaysRun)

	rl.InitWindow(screenWidth, screenHeight, "stingray - control layout experiments")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	paperImage := rl.GenImagePerlinNoise(int(paperPixelDimensions[0]), int(paperPixelDimensions[1]), 0, 0, 10.0)
	// paperImage := rl.LoadImage("assets/images/lisek.png")
	// if !rl.IsImageValid(paperImage) {
	// 	panic("image invalid")
	// }
	paperTexture := rl.LoadTextureFromImage(paperImage)
	if !rl.IsTextureValid(paperTexture) {
		panic("texture invalid")
	}

	// rl.GuiSetFont(font)

	// panel := gui.NewPanel(rl.NewRectangle(screenWidth-panelWidth, 0, panelWidth, screenHeight))
	layoutRoot :=
		layout.NewVStack("root",
			layout.Group("padding",
				layout.Group("paperGroup",
					layout.NewVStack("paperStack",
						layout.NewHFlex("",
							layout.Label("paperSizeLabel"),
							layout.Control("paperSize"),
						),
						layout.NewHFlex("",
							layout.Checkbox("paperOrientation"),
							layout.Control("paperOrientationLabel"),
						),
						layout.NewHFlex("",
							layout.Control("renderScaleLabel"),
							layout.Control("renderScale"),
							layout.Label("renderScaleHelper"),
						),
						layout.NewHFlex("",
							layout.Control("paperPixelsLabel"),
							layout.Label("paperPixelsWidth"),
							layout.Label("paperPixelsHeight"),
						),
						layout.NewHFlex("",
							layout.Label("paperReset"),
							layout.Label("paperApply"),
							layout.Label("paperSave"),
						),
					),
				),
			),
		)

	widgetRectangles := make(map[string]rl.Rectangle)
	var cb layout.WidgetCallback = func(widget layout.Widget, bounds rl.Rectangle) {
		widgetRectangles[widget.GetId()] = bounds
	}
	panelBounds := rl.NewRectangle(screenWidth-panelWidth-10, 10, panelWidth, screenHeight-20)
	layoutRoot.Arrange(panelBounds, cb)

	getRect := func(id string) rl.Rectangle {
		if rect, ok := widgetRectangles[id]; ok {
			return rect
		}
		panic(fmt.Sprintf("Widget with id %s not found", id))
	}

	// Main game loop
	for !exitWindow { // Detect window close button or ESC key
		// logic update

		// dims := paper.APaperSizeInPixels(paperSizeIdx, landscape, renderScale)
		// paperPixelDimensions[0] = dims[0]
		// paperPixelDimensions[1] = dims[1]

		wheelDiff := rl.GetMouseWheelMove()
		if wheelDiff != 0 {
			newTarget := rl.GetScreenToWorld2D(rl.GetMousePosition(), paperCam)
			paperCam.Target = newTarget
			paperCam.Zoom += float32(rl.GetMouseWheelMove()) * 0.05
			paperCam.Zoom = rl.Clamp(paperCam.Zoom, 0.1, 3.0)
			rl.SetMousePosition(screenWidth/2, screenHeight/2)
		}

		// Update
		exitWindow = rl.WindowShouldClose()
		if rl.IsKeyPressed(rl.KeyEscape) {
			exitWindow = true
		}

		// DRAWING
		rl.BeginDrawing()
		rl.ClearBackground(rl.GetColor(uint(rgui.GetStyle(rgui.DEFAULT, rgui.BACKGROUND_COLOR))))
		rl.BeginMode2D(paperCam)
		rl.DrawTexture(paperTexture, 0, 0, rl.White)
		rl.EndMode2D()

		// gui
		rl.DrawRectangleRec(getRect("root"), rl.RayWhite)
		rgui.GroupBox(getRect("paperGroup"), "Paper Settings")
		rgui.Label(getRect("paperSizeLabel"), "paper size")
		rgui.ComboBox(getRect("paperSize"), "A0;A1;A2;A3;A4;A5", &paperConfig.Data.SizeIdx)
		rgui.CheckBox(getRect("paperOrientation"), "", &paperConfig.Data.Landscape)
		rgui.Label(getRect("paperOrientationLabel"), "landscape?")
		rgui.Label(getRect("renderScaleLabel"), "render scale")
		rgui.ValueBox(getRect("renderScale"), "", &paperConfig.Data.RenderScale, 1, 30, true)
		rgui.Label(getRect("renderScaleHelper"), "px/mm")

		rgui.Label(getRect("paperPixelsLabel"), "paper pixels w/h")
		rgui.ValueBox(getRect("paperPixelsWidth"), "", &paperPixelDimensions[0], 0, 1000000, false)
		rgui.ValueBox(getRect("paperPixelsHeight"), "", &paperPixelDimensions[1], 0, 1000000, false)

		rgui.SetStyle(rgui.BUTTON, rgui.TEXT_ALIGNMENT, rgui.TEXT_ALIGN_CENTER)
		if rgui.Button(getRect("paperReset"), "Reset") {
			paperConfig.Data = paper.DefaultConfig()
		}
		if rgui.Button(getRect("paperApply"), "Apply") {
		}
		if rgui.Button(getRect("paperSave"), "Save") {
			err := paperConfig.Save()
			if err != nil {
				panic("could not save paper config")
			}
		}

		// STATUS_BAR
		mousePos := rl.GetMousePosition()
		mousePosText := fmt.Sprintf("Mouse Position: (%.0f, %.0f), camera zoom: %.3f", mousePos.X, mousePos.Y, paperCam.Zoom)
		rgui.StatusBar(rl.NewRectangle(0, float32(rl.GetScreenHeight())-20, float32(rl.GetScreenWidth()), 20), mousePosText)

		rl.EndDrawing()
	}

	// De-Initialization
	rl.CloseWindow() // Close window and OpenGL context
}
