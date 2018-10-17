package mkrg

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type viewer struct {
	graph         graph
	height, width int
}

func newViewer(graph graph, height, width int) *viewer {
	return &viewer{graph, height, width}
}

func (v *viewer) GetLines(ms metricsByName, until time.Time) []string {
	h, w := (v.height-3)*4, (v.width-1)*2
	dots := make([][]int, h)
	for i := range dots {
		dots[i] = make([]int, w)
	}
	maxValue := math.Max(ms.MaxValue(), 1.0) * 1.1
	from := until.Add(-time.Duration(w) * time.Minute)
	for _, metrics := range ms {
		for _, m := range metrics {
			x := int((m.Time - from.Unix()) / 60)
			if 0 <= x && x < w {
				y := int(m.Value.(float64) / maxValue * float64(h))
				dots[y][x] = 1
			}
		}
	}
	lines := make([]string, v.height)
	leftPadding := int(math.Max(float64((v.width-len(v.graph.name)+1)/2), 0))
	lines[0] = strings.Repeat(" ", leftPadding) + v.graph.name
	lines[0] += strings.Repeat(" ", int(math.Max(float64(v.width-len(lines[0])+1), 0)))
	line := make([]rune, w/2)
	for i := h - 4; i >= 0; i -= 4 {
		for j := 0; j < w; j += 2 {
			line[j/2] = rune(0x2800 | dots[i+3][j] | dots[i+2][j]<<1 | dots[i+1][j]<<2 | dots[i+3][j+1]<<3 |
				dots[i+2][j+1]<<4 | dots[i+1][j+1]<<5 | dots[i][j]<<6 | dots[i][j+1]<<7)
		}
		lines[(h-i)/4] = "|" + string(line)
	}
	axisX := []rune("+" + strings.Repeat("-", v.width-1))
	stepX := 30 * time.Minute
	var axisXLabels string
	for t := from.Truncate(stepX).Add(stepX); !until.Before(t); t = t.Add(stepX) {
		offset := int(float64(t.Sub(from)) / float64(until.Sub(from)) * float64(v.width))
		axisX[offset] = '+'
		axisXLabels += strings.Repeat(" ", int(math.Max(float64(offset-len(axisXLabels)-2), 0)))
		axisXLabels += fmt.Sprintf("%1d:%02d", t.Hour(), t.Minute())
		if offset < 2 {
			axisXLabels = axisXLabels[2-offset:]
		}
	}
	axisXLabels += strings.Repeat(" ", int(math.Max(float64(v.width-len(axisXLabels)), 0)))
	lines[v.height-2], lines[v.height-1] = string(axisX), axisXLabels[:v.width]
	return lines
}
