package mkrg

import (
	"fmt"
	"time"
)

type tui struct {
	height, width, column, maxColumn int
	until                            time.Time
	lines                            []string
}

func newTui(height, width, maxColumn int, until time.Time) *tui {
	return &tui{height, width, 0, maxColumn, until, make([]string, height)}
}

func (ui *tui) output(graph graph, ms metricsByName) error {
	v := newViewer(graph, ui.height, ui.width)
	for i, l := range v.GetLines(ms, ui.until) {
		ui.lines[i] += l
		if ui.column < ui.maxColumn-1 {
			ui.lines[i] += "    "
		}
	}
	if ui.column == ui.maxColumn-1 {
		for i := range ui.lines {
			fmt.Println(ui.lines[i])
			ui.lines[i] = ""
		}
		ui.column = 0
	} else {
		ui.column++
	}
	return nil
}

func (ui *tui) cleanup() error {
	if ui.column > 0 {
		for _, l := range ui.lines {
			fmt.Println(l)
		}
	}
	return nil
}
