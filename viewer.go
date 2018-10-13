package mkrg

import (
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
	h, w := (v.height-2)*4, (v.width-1)*2
	dots := make([][]int, h)
	for i := range dots {
		dots[i] = make([]int, w)
	}
	stackedValue := make(map[int64]float64)
	for _, metric := range v.graph.metrics {
		if !metric.stacked {
			continue
		}
		if metrics, ok := ms[metric.name]; ok {
			for i, m := range metrics {
				w := metrics[i].Value.(float64)
				if v, ok := stackedValue[m.Time]; ok {
					stackedValue[m.Time] = v + w
					metrics[i].Value = v + w
				} else {
					stackedValue[m.Time] = w
				}
			}
		}
	}
	maxValue := math.Max(ms.MaxValue(), 1.0) * 1.1
	from := until.Unix() - int64(w)*60
	for _, metrics := range ms {
		for _, m := range metrics {
			x := int((m.Time - from) / 60)
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
	lines[v.height-1] = "+" + strings.Repeat("-", v.width-1)
	return lines
}
