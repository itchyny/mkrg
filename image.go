package mkrg

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

var (
	borderColor  = color.RGBA{0xff, 0xff, 0xff, 0x88}
	axisColor    = color.RGBA{0xff, 0xff, 0xff, 0xff}
	tickColor    = color.RGBA{0xff, 0xff, 0xff, 0xaa}
	seriesColors = []color.RGBA{
		{0x63, 0xba, 0xc6, 0xff},
		{0xcc, 0x99, 0x00, 0xff},
		{0x81, 0x71, 0xb3, 0xff},
		{0x80, 0x9e, 0x10, 0xff},
		{0xb2, 0x66, 0x32, 0xff},
		{0x36, 0x99, 0x7d, 0xff},
		{0xb7, 0x95, 0x69, 0xff},
		{0x32, 0x6e, 0xc6, 0xff},
		{0x9c, 0x91, 0x00, 0xff},
		{0x53, 0x7c, 0x48, 0xff},
		{0xc9, 0x5b, 0x75, 0xff},
		{0x00, 0x5c, 0x9b, 0xff},
		{0x96, 0x75, 0x5a, 0xff},
		{0x67, 0xb0, 0x7d, 0xff},
		{0x5f, 0x83, 0xb8, 0xff},
		{0xa3, 0xa3, 0xe2, 0xff},
		{0x83, 0x9b, 0x4d, 0xff},
		{0xba, 0x55, 0x9b, 0xff},
		{0x3a, 0x8c, 0x86, 0xff},
		{0xb5, 0x83, 0x13, 0xff},
		{0x9e, 0x7f, 0x68, 0xff},
		{0x56, 0x54, 0xaf, 0xff},
	}
)

type imageWithMargins struct {
	img                   draw.Image
	topMargin, leftMargin int
}

func (img *imageWithMargins) Set(x, y int, c color.Color) {
	img.img.Set(x+img.leftMargin, y+img.topMargin, c)
}
func (img *imageWithMargins) ColorModel() color.Model {
	return img.img.ColorModel()
}
func (img *imageWithMargins) Bounds() image.Rectangle {
	return img.img.Bounds()
}
func (img *imageWithMargins) At(x, y int) color.Color {
	return img.img.At(x+img.leftMargin+img.topMargin, y)
}

func printImage(img draw.Image, graph graph, ms metricsByName, height, width int, from, until time.Time) error {
	drawGraph(img, graph, ms, height, width, from, until)
	drawBorder(img, height, width)
	drawTitle(img, width, graph.name)
	return nil
}

func drawGraph(img draw.Image, graph graph, ms metricsByName, height, width int, from, until time.Time) {
	graphLeftMargin, bottomMargin := 48, 26
	maxValue := math.Max(ms.MaxValue(), 1.0) * 1.1
	drawAxisX(img, height-bottomMargin, width, graphLeftMargin, from, until)
	drawAxisY(img, height-bottomMargin, width, graphLeftMargin, from, until, maxValue)
	drawSeries(&imageWithMargins{img, 0, graphLeftMargin}, graph, ms, height-bottomMargin, width-graphLeftMargin, from, until, maxValue)
}

func drawSeries(img draw.Image, graph graph, ms metricsByName, height, width int, from, until time.Time, maxValue float64) {
	imgSet := func(x, y int, c color.RGBA) {
		pointSize := 2
		for i := 0; i < pointSize; i++ {
			for j := 0; j < pointSize; j++ {
				img.Set(x+i, height-(y+j), c)
			}
		}
	}
	for i, metricName := range ms.ListMetricNames(graph) {
		prevPrevTime, prevTime, prevX, prevY := int64(0), int64(0), -1.0, 0.0
		metrics, seriesColor := ms[metricName], seriesColors[i%len(seriesColors)]
		for _, m := range metrics {
			x := float64(m.Time-from.Unix()) * float64(width) / float64(until.Sub(from)/time.Second)
			y := m.Value.(float64) / maxValue * float64(height)
			if 0 <= x {
				start, step := 0.0, math.Min(2.0/math.Sqrt((x-prevX)*(x-prevX)+(y-prevY)*(y-prevY)), 0.2)
				if prevX < 0 || prevTime-3*60 < prevPrevTime && prevTime < m.Time-3*60 {
					start = 1.0
				}
				prevPrevTime, prevTime = prevTime, m.Time
				for p := start; p <= 1.0; p += step {
					imgSet(int(prevX*(1.0-p)+x*p), int(prevY*(1.0-p)+y*p), seriesColor)
				}
			}
			prevX, prevY = x, y
		}
	}
}

func drawAxisX(img draw.Image, height, width, graphLeftMargin int, from, until time.Time) {
	for i := 0; i < height; i++ {
		img.Set(graphLeftMargin, i, axisColor)
	}
	stepX := 30 * time.Minute
	for t := from.Truncate(stepX).Add(stepX); !until.Before(t); t = t.Add(stepX) {
		offset := int(float64(t.Sub(from)) / float64(until.Sub(from)) * float64(width-graphLeftMargin))
		for i := 0; i < height; i++ {
			img.Set(graphLeftMargin+offset, i, tickColor)
		}
		diffX := -19
		if t.Hour() < 10 {
			diffX = -23
		}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(axisColor),
			Face: inconsolata.Bold8x16,
			Dot:  fixed.P(graphLeftMargin+offset+diffX, height+17),
		}
		d.DrawString(fmt.Sprintf("%2d:%02d", t.Hour(), t.Minute()))
	}
	for i := 0; i < 40; i++ {
		for j := 0; j < 20; j++ {
			img.Set(i+width, height+j+5, color.Alpha{0x00})
		}
	}
}

func drawAxisY(img draw.Image, height, width, graphLeftMargin int, from, until time.Time, maxValue float64) {
	for i := graphLeftMargin; i < width; i++ {
		img.Set(i, height-1, axisColor)
	}
	tick := math.Pow10(int(math.Floor(math.Log10(maxValue / 5.0))))
	if maxValue/tick > 12 {
		tick *= 5
	} else if maxValue/tick > 6 {
		tick *= 2
	}
	format, scale := formatAxisY(tick, maxValue)
	for y := 0.0; y < maxValue; y += tick {
		posY := height - int(y/maxValue*float64(height))
		for i := graphLeftMargin; 0.0 < y && i < width; i++ {
			img.Set(i, posY, tickColor)
		}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(axisColor),
			Face: inconsolata.Bold8x16,
			Dot:  fixed.P(8, posY+4),
		}
		if y == 0.0 {
			d.DrawString("   0")
		} else if y/maxValue < 0.985 {
			d.DrawString(fmt.Sprintf("%4s", fmt.Sprintf(format, y/scale)))
		}
	}
}

func formatAxisY(tick, maxValue float64) (string, float64) {
	var suffix string
	scale := 1.0
	if maxValue >= 1e12 {
		suffix, scale = "T", 1e12
	} else if maxValue >= 1e9 {
		suffix, scale = "G", 1e9
	} else if maxValue >= 1e6 {
		suffix, scale = "M", 1e6
	} else if maxValue >= 1e3 {
		suffix, scale = "K", 1e3
	}
	digits := int(math.Ceil(math.Max(1.0-math.Log10(maxValue/scale), 0.0)))
	return fmt.Sprintf("%%.%df%s", digits, suffix), scale
}

func drawBorder(img draw.Image, height, width int) {
	for i := 0; i < width; i++ {
		img.Set(i, 0, borderColor)
		img.Set(i, height-1, borderColor)
	}
	for i := 0; i < height; i++ {
		img.Set(0, i, borderColor)
		img.Set(width-1, i, borderColor)
	}
}

func drawTitle(img draw.Image, width int, title string) {
	x, y := width/2-len(title)*4, 20
	for i := -3; i < len(title)*8+3; i++ {
		for j := -4; j < 15; j++ {
			img.Set(x+i, y-j, color.Alpha{0x00})
		}
	}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(axisColor),
		Face: inconsolata.Bold8x16,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(title)
}
