package mkrg

import (
	"image"
	"image/color"
	"math"
	"time"
)

var (
	borderColor = color.RGBA{0xff, 0xff, 0xff, 0x88}
)

func printImage(img *image.RGBA, graph graph, ms metricsByName, height, width, leftMargin int, from, until time.Time) error {
	c := color.RGBA{0x63, 0xba, 0xc6, 0xff}
	maxValue := math.Max(ms.MaxValue(), 1.0) * 1.1
	imgSet := func(x, y int, c color.RGBA) {
		pointSize := 2
		for i := 0; i < pointSize; i++ {
			for j := 0; j < pointSize; j++ {
				img.Set(leftMargin+x+i, height-(y+j), c)
			}
		}
	}
	for i := 0; i < width; i++ {
		img.Set(leftMargin+i, 0, borderColor)
		img.Set(leftMargin+i, height-1, borderColor)
	}
	for i := 0; i < height; i++ {
		img.Set(leftMargin, i, borderColor)
		img.Set(leftMargin+width-1, i, borderColor)
	}
	prevX, prevY := -1, 0
	for _, metrics := range ms {
		for _, m := range metrics {
			x := int(m.Time-from.Unix()) * width / int(until.Sub(from)/time.Second)
			y := int(m.Value.(float64) / maxValue * float64(height))
			if 0 <= x && 0 <= prevX && prevX < x {
				step := int(math.Max(math.Sqrt(float64((x-prevX)*(x-prevX)+(y-prevY)*(y-prevY)))/2.0, 5.0))
				for i := 1; i <= step; i++ {
					imgSet(int(float64(prevX*(step-i)+x*i)/float64(step)), int((float64(prevY*(step-i)+y*i))/float64(step)), c)
				}
			}
			prevX, prevY = x, y
		}
		prevX, prevY = -1, 0
	}
	return nil
}
