package waveshare

// References:
// https://github.com/waveshare/e-Paper/blob/master/RaspberryPi_JetsonNano/python/lib/waveshare_epd/epd7in5_HD.py

/*
package main

import (
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"math"
	"os"

	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

var (
	flagWidth  = flag.Uint("width", 0, "Width of destination image")
	flagHeight = flag.Uint("height", 0, "Height of destination of image")
	flagRotate = flag.Int("rotate", 0, "Rotation of image around center")
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Require input filenames")
		os.Exit(-1)
	}
	for _, path := range flag.Args() {
		if err := Process(path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
	}
}

func Process(path string) error {
	fh, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fh.Close()

	img, _, err := image.Decode(fh)
	if err != nil {
		return err
	}

	if out, err := Transform(img); err != nil {
		return err
	} else if fh, err := ioutil.TempFile("", "*.jpg"); err != nil {
		return err
	} else {
		defer fh.Close()
		if err := png.Encode(fh, out); err != nil {
			return err
		}
		fmt.Println(fh.Name())
	}

	// Return success
	return nil
}

func GetDestBounds(src image.Image) image.Rectangle {
	w := float64(*flagWidth)
	h := float64(*flagHeight)
	if w == 0 {
		if h == 0 {
			h = float64(src.Bounds().Dy())
		}
		w = float64(h) * float64(src.Bounds().Dx()) / float64(src.Bounds().Dy())
	}
	if h == 0 {
		h = float64(w) * float64(src.Bounds().Dy()) / float64(src.Bounds().Dx())
	}
	return image.Rectangle{image.ZP, image.Point{int(w), int(h)}}
}

func GetRotation(src image.Image) (float64, float64, float64) {
	centerx := float64(src.Bounds().Min.X+src.Bounds().Max.X) / 2
	centery := float64(src.Bounds().Min.Y+src.Bounds().Max.Y) / 2
	theta := math.Pi * float64(*flagRotate) / 180.0
	return theta, float64(centerx), float64(centery)
}

func Transform(src image.Image) (image.Image, error) {
	// Resize image into 'scaled'
	scaled := image.NewRGBA(GetDestBounds(src))
	transform := NewAffineTransform().Scale(
		float64(scaled.Bounds().Dx())/float64(src.Bounds().Dx()),
		float64(scaled.Bounds().Dy())/float64(src.Bounds().Dy()),
	)
	draw.ApproxBiLinear.Transform(scaled, f64.Aff3(transform), src, src.Bounds(), draw.Over, nil)

	// Rotate image
	if theta, x, y := GetRotation(scaled); theta != 0 {
		fmt.Fprintln(os.Stderr, theta, x, y)
		rotated := image.NewRGBA(scaled.Bounds())
		transform := NewAffineTransform().Rotate(theta, x, y)
		draw.ApproxBiLinear.Transform(rotated, f64.Aff3(transform), scaled, scaled.Bounds(), draw.Over, nil)
		scaled = rotated
	}

	// Convert to BW using dithering
	dst := image.NewPaletted(GetDestBounds(src), []color.Color{
		color.Gray{Y: 255},
		color.Gray{Y: 0},
	})
	draw.FloydSteinberg.Draw(dst, dst.Bounds(), scaled, image.ZP)

	return dst, nil
}

*/
