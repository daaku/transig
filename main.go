package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
)

func run(inPath string, out io.Writer) error {
	const threshold = 175
	f, err := os.Open(inPath)
	if err != nil {
		return errors.Wrapf(err, "opening input file: %s", inPath)
	}
	defer f.Close()
	inImg, _, err := image.Decode(f)
	if err != nil {
		return errors.Wrapf(err, "decoding input file: %s", inPath)
	}
	b := inImg.Bounds()
	outImg := image.NewRGBA(b)
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			p := inImg.At(x, y)
			_, _, _, a := p.RGBA()
			if a == 0 {
				continue // skin transparent pixels
			}
			g := color.GrayModel.Convert(p).(color.Gray)
			if g.Y > threshold {
				outImg.Set(x, y, color.Transparent)
			} else {
				outImg.Set(x, y, color.Black)
			}
		}
	}
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	if err := encoder.Encode(out, outImg); err != nil {
		return errors.Wrap(err, "encoding output image")
	}
	return nil
}

func main() {
	log.SetOutput(os.Stderr)
	if err := run(os.Args[1], os.Stdout); err != nil {
		log.Fatalf("error: %+v\n", err)
	}
}
