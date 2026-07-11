package paper

type Config struct {
	SizeIdx     int32
	Landscape   bool
	RenderScale int32
}

func DefaultConfig() Config {
	return Config{
		SizeIdx:     5,
		Landscape:   false,
		RenderScale: 10,
	}
}

func (conf Config) PixelDims() PaperDimensions {
	return APaperSizeInPixels(conf.SizeIdx, conf.Landscape, conf.RenderScale)
}
