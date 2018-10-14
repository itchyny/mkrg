package mkrg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"
)

type iterm2 struct {
	height, width int
	maxColumn     int
	from, until   time.Time
}

func newIterm2(height, width int, maxColumn int, from, until time.Time) *iterm2 {
	return &iterm2{height, width, maxColumn, from, until}
}

func (ui *iterm2) output(graph graph, ms metricsByName) error {
	buf := new(bytes.Buffer)
	printImage(buf, graph, ms, ui.height*24, ui.width*16, ui.from, ui.until)
	fmt.Printf("\x1b]1337;File=inline=1;preserveAspectRatio=1;width=%d%%:", 100.0/ui.maxColumn)
	fmt.Print(base64.StdEncoding.EncodeToString(buf.Bytes()))
	fmt.Print("\x07\n")
	return nil
}

func (ui *iterm2) cleanup() error {
	return nil
}
