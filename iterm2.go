package mkrg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"time"
)

type iterm2 struct {
	height, width     int
	column, maxColumn int
	from, until       time.Time
	img               *image.RGBA
}

func newIterm2(height, width int, maxColumn int, from, until time.Time) *iterm2 {
	return &iterm2{height, width, 0, maxColumn, from, until, nil}
}

func (ui *iterm2) output(graph graph, ms metricsByName) error {
	imgHeight, imgWidth, padding := ui.height*24, ui.width*16, ui.width/5
	if ui.column == 0 {
		ui.img = image.NewRGBA(image.Rect(0, 0, (imgWidth+padding)*ui.maxColumn-padding, imgHeight+padding))
	}
	printImage(ui.img, graph, ms, imgHeight, imgWidth, (imgWidth+padding)*ui.column, ui.from, ui.until)
	if ui.column == ui.maxColumn-1 {
		if err := ui.cleanup(); err != nil {
			return err
		}
		ui.column = 0
	} else {
		ui.column++
	}
	return nil
}

func (ui *iterm2) cleanup() error {
	if ui.column > 0 {
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, ui.img); err != nil {
			return err
		}
		fmt.Print("\x1b]1337;File=inline=1;preserveAspectRatio=1;width=100%:")
		fmt.Print(base64.StdEncoding.EncodeToString(buf.Bytes()))
		fmt.Print("\x07\n")
	}
	return nil
}
