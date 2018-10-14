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
		color.RGBA{0x63, 0xba, 0xc6, 0xff},
		color.RGBA{0xcc, 0x99, 0x00, 0xff},
		color.RGBA{0x81, 0x71, 0xb3, 0xff},
		color.RGBA{0x80, 0x9e, 0x10, 0xff},
		color.RGBA{0xb2, 0x66, 0x32, 0xff},
		color.RGBA{0x36, 0x99, 0x7d, 0xff},
		color.RGBA{0xb7, 0x95, 0x69, 0xff},
		color.RGBA{0x32, 0x6e, 0xc6, 0xff},
		color.RGBA{0x9c, 0x91, 0x00, 0xff},
		color.RGBA{0x53, 0x7c, 0x48, 0xff},
		color.RGBA{0xc9, 0x5b, 0x75, 0xff},
		color.RGBA{0x00, 0x5c, 0x9b, 0xff},
		color.RGBA{0x96, 0x75, 0x5a, 0xff},
		color.RGBA{0x67, 0xb0, 0x7d, 0xff},
		color.RGBA{0x5f, 0x83, 0xb8, 0xff},
		color.RGBA{0xa3, 0xa3, 0xe2, 0xff},
		color.RGBA{0x83, 0x9b, 0x4d, 0xff},
		color.RGBA{0xba, 0x55, 0x9b, 0xff},
		color.RGBA{0x3a, 0x8c, 0x86, 0xff},
		color.RGBA{0xb5, 0x83, 0x13, 0xff},
		color.RGBA{0x9e, 0x7f, 0x68, 0xff},
		color.RGBA{0x56, 0x54, 0xaf, 0xff},
	}
)

type Image struct {
	img                   draw.Image
	topMargin, leftMargin int
}

func (img *Image) Set(x, y int, c color.Color) {
	img.img.Set(x+img.leftMargin, y+img.topMargin, c)
}
func (img *Image) ColorModel() color.Model {
	return img.img.ColorModel()
}
func (img *Image) Bounds() image.Rectangle {
	return img.img.Bounds()
}
func (img *Image) At(x, y int) color.Color {
	return img.img.At(x+img.leftMargin+img.topMargin, y)
}

func printImage(img draw.Image, graph graph, ms metricsByName, height, width int, from, until time.Time) error {
	drawGraph(img, graph, ms, height, width, from, until)
	drawBorder(img, height, width)
	drawTitle(img, width, graph.name)
	return nil
}

func drawGraph(img draw.Image, graph graph, ms metricsByName, height, width int, from, until time.Time) {
	graphLeftMargin, bottomMargin := 48, 30
	maxValue := math.Max(ms.MaxValue(), 1.0) * 1.1
	drawAxisX(img, height, width, graphLeftMargin, bottomMargin, from, until)
	drawAxisY(img, height, width, graphLeftMargin, bottomMargin, from, until, maxValue)
	drawSeries(&Image{img, 0, graphLeftMargin}, graph, ms, height-bottomMargin, width-graphLeftMargin, from, until, maxValue)
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
	prevX, prevY := -1, 0
	for i, metricName := range ms.MetricNames() {
		metrics, seriesColor := ms[metricName], seriesColors[i%len(seriesColors)]
		for _, m := range metrics {
			x := int(m.Time-from.Unix()) * width / int(until.Sub(from)/time.Second)
			y := int(m.Value.(float64) / maxValue * float64(height))
			if 0 <= x && 0 <= prevX && prevX < x {
				step := int(math.Max(math.Sqrt(float64((x-prevX)*(x-prevX)+(y-prevY)*(y-prevY)))/2.0, 5.0))
				for i := 1; i <= step; i++ {
					imgSet(int(float64(prevX*(step-i)+x*i)/float64(step)), int((float64(prevY*(step-i)+y*i))/float64(step)), seriesColor)
				}
			}
			prevX, prevY = x, y
		}
		prevX, prevY = -1, 0
	}
}

func drawAxisX(img draw.Image, height, width, graphLeftMargin, bottomMargin int, from, until time.Time) {
	for i := 0; i < height-bottomMargin; i++ {
		img.Set(graphLeftMargin, i, axisColor)
	}
	stepX := 30 * time.Minute
	for t := from.Truncate(stepX).Add(stepX); t.Before(until); t = t.Add(stepX) {
		offset := int(float64(t.Sub(from)) / float64(until.Sub(from)) * float64(width-graphLeftMargin))
		for i := 0; i < height-bottomMargin; i++ {
			img.Set(graphLeftMargin+offset, i, tickColor)
		}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(axisColor),
			Face: inconsolata.Bold8x16,
			Dot:  fixed.P(graphLeftMargin+offset-17, height-bottomMargin+20),
		}
		d.DrawString(fmt.Sprintf("%2d:%02d", t.Hour(), t.Minute()))
	}
}

func drawAxisY(img draw.Image, height, width, graphLeftMargin, bottomMargin int, from, until time.Time, maxValue float64) {
	for i := graphLeftMargin; i < width; i++ {
		img.Set(i, height-bottomMargin-1, axisColor)
	}
	tick := math.Pow10(int(math.Floor(math.Log10(maxValue / 5.0))))
	if maxValue/tick > 12 {
		tick *= 5
	} else if maxValue/tick > 6 {
		tick *= 2
	}
	format, scale := formatAxisY(tick, maxValue)
	for y := 0.0; y < maxValue; y += tick {
		posY := height - bottomMargin - int(y/maxValue*float64(height-bottomMargin))
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
		} else {
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
