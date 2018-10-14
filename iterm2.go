package mkrg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"
)

type iterm2 struct {
	height, width int
	until         time.Time
}

func newIterm2(height, width int, until time.Time) *iterm2 {
	return &iterm2{height, width, until}
}

func (ui *iterm2) output(graph graph, ms metricsByName) error {
	buf := new(bytes.Buffer)
	printImage(buf, graph, ms, ui.height*24, ui.width*16,
		ui.until.Add(-time.Duration(ui.width*2)*time.Minute), ui.until)
	fmt.Print("\x1b]1337;File=inline=1;preserveAspectRatio=1:")
	fmt.Print(base64.StdEncoding.EncodeToString(buf.Bytes()))
	fmt.Print("\x07\n")
	return nil
}

func (ui *iterm2) cleanup() error {
	return nil
}
