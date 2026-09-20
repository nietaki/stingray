package layout

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// TODO: change the callback to accept the bounds localization function
type PanelDrawingCallback func(*Panel)

type Panel struct {
	Bounds          rl.Rectangle
	LayoutRoot      Widget
	DrawingCallback PanelDrawingCallback

	widgetRectangles map[string]rl.Rectangle
}

func NewPanel(bounds rl.Rectangle, layoutRoot Widget, drawingCallback PanelDrawingCallback) *Panel {
	return &Panel{
		Bounds:           bounds,
		LayoutRoot:       layoutRoot,
		DrawingCallback:  drawingCallback,
		widgetRectangles: make(map[string]rl.Rectangle),
	}
}

func (panel *Panel) Rect(id string) rl.Rectangle {
	if rect, ok := panel.widgetRectangles[id]; ok {
		return rect
	}
	panic(fmt.Sprintf("Widget with id %s not found", id))
}

func (panel *Panel) Arrange() {
	var cb WidgetCallback = func(widget Widget, bounds rl.Rectangle) {
		panel.widgetRectangles[widget.GetId()] = bounds
	}

	panel.LayoutRoot.Arrange(panel.Bounds, cb)
}

func (panel *Panel) Draw() {
	panel.DrawingCallback(panel)
}
