package mkrg

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

var (
	borderColor = color.RGBA{0xff, 0xff, 0xff, 0x88}
	axisColor   = color.RGBA{0xff, 0xff, 0xff, 0xff}
	tickColor   = color.RGBA{0xff, 0xff, 0xff, 0xaa}
)

func printImage(img *image.RGBA, graph graph, ms metricsByName, height, width, leftMargin int, from, until time.Time) error {
	drawGraph(img, graph, ms, height, width, leftMargin, from, until)
	drawBorder(img, height, width, leftMargin)
	return nil
}

func drawGraph(img *image.RGBA, graph graph, ms metricsByName, height, width, leftMargin int, from, until time.Time) {
	graphLeftMargin, bottomMargin := 60, 30
	drawSeries(img, graph, ms, height-bottomMargin, width-graphLeftMargin, leftMargin+graphLeftMargin, from, until)
	drawAxis(img, height, width, leftMargin, graphLeftMargin, bottomMargin, from, until)
}

func drawSeries(img *image.RGBA, graph graph, ms metricsByName, height, width, leftMargin int, from, until time.Time) {
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
}

func drawAxis(img *image.RGBA, height, width, leftMargin, graphLeftMargin, bottomMargin int, from, until time.Time) {
	for i := graphLeftMargin; i < width; i++ {
		img.Set(leftMargin+i, height-bottomMargin-1, axisColor)
	}
	for i := 0; i < height-bottomMargin; i++ {
		img.Set(leftMargin+graphLeftMargin, i, axisColor)
	}
	stepX := 30 * time.Minute
	for t := from.Truncate(stepX).Add(stepX); t.Before(until); t = t.Add(stepX) {
		offset := int(float64(t.Sub(from)) / float64(until.Sub(from)) * float64(width-graphLeftMargin))
		for i := 0; i < height-bottomMargin; i++ {
			img.Set(leftMargin+graphLeftMargin+offset, i, tickColor)
		}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(axisColor),
			Face: inconsolata.Bold8x16,
			Dot:  fixed.P(leftMargin+graphLeftMargin+offset-17, height-bottomMargin+20),
		}
		d.DrawString(fmt.Sprintf("%2d:%02d", t.Hour(), t.Minute()))
	}
}

func drawBorder(img *image.RGBA, height, width, leftMargin int) {
	for i := 0; i < width; i++ {
		img.Set(leftMargin+i, 0, borderColor)
		img.Set(leftMargin+i, height-1, borderColor)
	}
	for i := 0; i < height; i++ {
		img.Set(leftMargin, i, borderColor)
		img.Set(leftMargin+width-1, i, borderColor)
	}
}
