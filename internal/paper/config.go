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
