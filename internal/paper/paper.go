package paper

// import rl "github.com/gen2brain/raylib-go/raylib"

type PaperDimensions [2]int32

func (d *PaperDimensions) Width() int32 {
	return d[0]
}

func (d *PaperDimensions) Height() int32 {
	return d[1]
}

var aSizes = [...]PaperDimensions{
	{841, 1189},
	{594, 841},
	{420, 594},
	{297, 420},
	{210, 297},
	{148, 210},
}

func APaperSizeInPixels(sizeIdx int32, landscape bool, renderScale int32) PaperDimensions {
	ret := aSizes[sizeIdx]

	ret[0] *= renderScale
	ret[1] *= renderScale

	if landscape {
		ret[0], ret[1] = ret[1], ret[0]
	}

	return ret
}
