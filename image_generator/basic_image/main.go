package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	data := []int{10, 33, 73, 64}

	const (
		width  = 400
		height = 200
		barWidth = 50
		barGap   = 20
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.White)
		}
	}

	barColor := color.RGBA{R: 0, G: 0, B: 255, A: 255}

	xOffset := barGap
	for _, value := range data {
		barHeight := (value * height) / 100

		barX1 := xOffset
		barY1 := height - barHeight
		barX2 := xOffset + barWidth
		barY2 := height

		for x := barX1; x < barX2; x++ {
			for y := barY1; y < barY2; y++ {
				if x >= 0 && x < width && y >= 0 && y < height {
					img.Set(x, y, barColor)
				}
			}
		}
		xOffset += barWidth + barGap
	}

	file, err := os.Create("barchart_pixel.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

	println("barchart_pixel.png created!")
}