package layout

import rl "github.com/gen2brain/raylib-go/raylib"

type WidgetCallback func(Widget, rl.Rectangle)

type Widget interface {
	SizeHint() rl.Vector2
	Arrange(bounds rl.Rectangle, callback WidgetCallback)
	GetId() string
}

var marginSize float32 = 5.0

type Id string

func (id Id) GetId() string { return string(id) }

// # Box
type Box struct {
	sizeHint rl.Vector2
	Id
}

func (widget Box) SizeHint() rl.Vector2 {
	return widget.sizeHint
}

func (widget *Box) Arrange(bounds rl.Rectangle, callback WidgetCallback) {
	callback(widget, bounds)
}

func NewBox(id string, width, height float32) *Box {
	return &Box{
		sizeHint: rl.Vector2{X: width, Y: height},
		Id:       Id(id),
	}
}

// # VStack

type VStack struct {
	children []Widget
	Id
}

func (widget VStack) SizeHint() rl.Vector2 {
	var ret rl.Vector2
	for _, child := range widget.children {
		hint := child.SizeHint()
		ret.X = max(ret.X, hint.X)
		ret.Y += hint.Y
	}

	ret.Y += float32(len(widget.children)-1) * marginSize
	return ret
}

func (widget *VStack) Arrange(bounds rl.Rectangle, callback WidgetCallback) {
	width := bounds.Width
	x := bounds.X
	y := bounds.Y
	lastY := y

	for _, child := range widget.children {
		childHint := child.SizeHint()
		childRect := rl.Rectangle{
			X:      x,
			Y:      y,
			Height: childHint.Y,
			Width:  width,
		}
		child.Arrange(childRect, callback)
		y += childRect.Height
		lastY = y
		y += marginSize
	}

	rect := bounds
	rect.Height = lastY - bounds.Y
	callback(widget, rect)
}

func NewVStack(id string, children ...Widget) *VStack {
	return &VStack{
		children: children,
		Id:       Id(id),
	}
}

// # HFlex

type HFlex struct {
	children []Widget
	Id
}

func (widget HFlex) SizeHint() rl.Vector2 {
	var ret rl.Vector2
	for _, child := range widget.children {
		hint := child.SizeHint()
		ret.Y = max(ret.Y, hint.Y)
		ret.X += hint.X
	}

	ret.X += float32(len(widget.children)-1) * marginSize
	return ret
}

func (widget *HFlex) Arrange(bounds rl.Rectangle, callback WidgetCallback) {
	availableWidth := bounds.Width - marginSize*float32(len(widget.children)-1)

	height := bounds.Height
	x := bounds.X
	y := bounds.Y

	var totalWidthHint float32
	for _, child := range widget.children {
		childHint := child.SizeHint()
		totalWidthHint += childHint.X
	}

	for _, child := range widget.children {
		childHint := child.SizeHint()
		childRect := rl.Rectangle{
			X:      x,
			Y:      y,
			Height: height,
			Width:  availableWidth * childHint.X / totalWidthHint,
		}
		child.Arrange(childRect, callback)
		x += childRect.Width
		x += marginSize
	}

	callback(widget, bounds)
}

func NewHFlex(id string, children ...Widget) *HFlex {
	return &HFlex{
		children: children,
		Id:       Id(id),
	}
}

// Wrapper

const (
	UP_IDX    = 0
	RIGHT_IDX = 1
	DOWN_IDX  = 2
	LEFT_IDX  = 3
)

type Wrapper struct {
	child   Widget
	padding [4]float32
	Id
}

func (widget Wrapper) SizeHint() rl.Vector2 {
	ret := widget.child.SizeHint()
	ret.Y += widget.padding[UP_IDX] + widget.padding[DOWN_IDX]
	ret.X += widget.padding[RIGHT_IDX] + widget.padding[LEFT_IDX]
	return ret
}

func (widget *Wrapper) Arrange(bounds rl.Rectangle, callback WidgetCallback) {
	childBounds := bounds
	childBounds.X += widget.padding[LEFT_IDX]
	childBounds.Y += widget.padding[UP_IDX]
	childBounds.Width -= widget.padding[LEFT_IDX] + widget.padding[RIGHT_IDX]
	childBounds.Height -= widget.padding[UP_IDX] + widget.padding[DOWN_IDX]
	widget.child.Arrange(childBounds, callback)

	callback(widget, bounds)
}

func NewWrapper(id string, child Widget, paddings ...float32) *Wrapper {
	ret := Wrapper{
		child: child,
		Id:    Id(id),
	}

	for i := range 4 {
		ret.padding[i] = paddings[i%len(paddings)]
	}

	return &ret
}
