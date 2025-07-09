package main

import (
	"fmt"
	"os"
	"github.com/ajstarks/svgo" // Import the svgo library
)

func main() {
	data := []int{10, 33, 73, 64} // Example data

	const (
		width     = 500
		height    = 300
		barWidth  = 60
		barMargin = 20	
		chartPadding = 50 	
	)
	
	chartWidth := (len(data) * barWidth) + ((len(data) - 1) * barMargin)
	chartHeight := height - (2 * chartPadding)

	file, err := os.Create("barchart_svg.svg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	s := svg.New(file)
	s.Start(width, height)

	s.Rect(chartPadding, chartPadding, chartWidth, chartHeight, "fill:lightgray")

	xPos := chartPadding
	for i, value := range data {
		barHeight := (value * chartHeight) / 100

		barY := chartPadding + chartHeight - barHeight

		barColor := fmt.Sprintf("fill:rgb(%d,%d,%d)", i*50, 200-(i*20), 255)
		s.Rect(xPos, barY, barWidth, barHeight, barColor)

		textY := barY - 10
		s.Text(xPos+(barWidth/2), textY, fmt.Sprintf("%d", value), "text-anchor:middle;font-size:16px;fill:black")

		xPos += barWidth + barMargin
	}
	s.Text(width/2, chartPadding/2, "My Awesome Bar Chart", "text-anchor:middle;font-size:24px;font-weight:bold;fill:black")

	s.End()
	println("barchart_svg.svg created!")
}
