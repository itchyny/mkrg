package mkrg

import (
	"image"
	"os"
	"time"

	"github.com/mattn/go-sixel"
)

type sixelUI struct {
	height, width     int
	column, maxColumn int
	from, until       time.Time
	img               *image.RGBA
}

func newSixel(height, width, maxColumn int, from, until time.Time) *sixelUI {
	return &sixelUI{height, width, 0, maxColumn, from, until, nil}
}

// Needs improvements because I don't have checked the behavior.
func (ui *sixelUI) output(graph graph, ms metricsByName) error {
	imgHeight, imgWidth, padding := ui.height*20, ui.width*12, ui.width/5
	if ui.column == 0 {
		ui.img = image.NewRGBA(image.Rect(0, 0, (imgWidth+padding)*ui.maxColumn-padding, imgHeight+padding*2))
	}
	printImage(&imageWithMargins{ui.img, padding, (imgWidth + padding) * ui.column}, graph, ms, imgHeight, imgWidth, ui.from, ui.until)
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

func (ui *sixelUI) cleanup() error {
	if ui.column > 0 {
		if err := sixel.NewEncoder(os.Stdout).Encode(ui.img); err != nil {
			return err
		}
	}
	return nil
}
