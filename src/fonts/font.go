package fonts

import (
	"golang.org/x/image/font"
	// "golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/opentype"
)

var (
	BaseFontSize = 5
	GoMonoFace   []font.Face
	DPI          = 200
	Hinting      = font.HintingFull
)

func init() {
	f, err := opentype.Parse(gomonobold.TTF)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		face, err := opentype.NewFace(f, &opentype.FaceOptions{
			Size:    float64(i * BaseFontSize),
			DPI:     float64(DPI),
			Hinting: font.HintingFull,
		})

		if err != nil {
			panic(err)
		}
		GoMonoFace = append(GoMonoFace, face)
	}
}
