package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

var transparent = color.RGBA{A: 255}

func main() {
	// Declare flags.
	var image1Path string
	flag.StringVar(&image1Path, "img1", "", "Path to image 1")
	var image2Path string
	flag.StringVar(&image2Path, "img2", "", "Path to image 1")
	var outName string
	flag.StringVar(&outName, "ou", "out.png", "Name of output file name.")
	flag.Parse()

	// Read first image.
	img1file, err := os.Open(image1Path)
	if err != nil {
		fmt.Printf("could not open file %s because %s", image1Path, err.Error())
	}
	defer img1file.Close()
	img1, _, err := image.Decode(img1file)
	if err != nil {
		panic(err.Error())
	}

	// Read second image.
	img2file, err := os.Open(image2Path)
	if err != nil {
		fmt.Printf("could not open file %s because %s", image2Path, err.Error())
	}
	defer img2file.Close()
	img2, _, err := image.Decode(img2file)
	if err != nil {
		panic(err.Error())
	}

	// Let's check necessary conditions before we start processing.
	size1 := img1.Bounds().Size()
	size2 := img2.Bounds().Size()
	if size1.X != size2.X || size1.Y != size2.Y {
		fmt.Printf("Images size should be identical. Image %s size: (%d,%d), image %s size: (%d,%d)", image1Path, size1.X, size1.Y, image2Path, size2.X, size2.Y)
		os.Exit(1)
	}

	minX := size1.X - 1
	minY := size1.Y - 1
	maxX := 0
	maxY := 0
	for x := 0; x < size1.X; x++ {
		for y := 0; y < size1.Y; y++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	img := image.NewRGBA(image.Rect(0, 0, maxX-minX+1, maxY-minY+1))
	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				img.Set(x-minX, y-minY, img1.At(x, y))
			} else {
				img.Set(x-minX, y-minY, transparent)
			}
		}
	}
	outFile, err := os.Create(outName)
	if err = png.Encode(outFile, img); err != nil {
		fmt.Printf("could not ecode output image because: %s", err.Error())
		os.Exit(1)
	}
	defer outFile.Close()
}
