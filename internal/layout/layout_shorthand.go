package layout

var controlHeight float32 = 20.0
var controlWidth float32 = 100.0
var labelWidth float32 = 50.0

func Control(id string) *Box {
	return NewBox(id, controlWidth, controlHeight)
}

func Label(id string) *Box {
	return NewBox(id, labelWidth, controlHeight)
}

func Group(id string, child Widget) *Wrapper {
	return NewWrapper(id, child, 10, 5, 5, 5)
}
