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

func (v *viewer) GetLines(ms metricsByName, from time.Time) []string {
	dots := make([][]int, v.height)
	for i := range dots {
		dots[i] = make([]int, v.width)
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
	for _, metrics := range ms {
		for _, m := range metrics {
			x := (m.Time - from.Unix()) / 60
			y := int(m.Value.(float64) / maxValue * float64(v.height))
			dots[y][x] = 1
		}
	}
	lines := make([]string, v.height/4+1)
	lines[0] = strings.Repeat(" ", int(math.Max(float64((v.width/2-len(v.graph.name)+1)/2), 0))) + v.graph.name
	line := make([]rune, 90)
	for i := v.height - 4; i >= 0; i -= 4 {
		for j := 0; j < v.width; j += 2 {
			line[j/2] = rune(0x2800 | dots[i+3][j] | dots[i+2][j]<<1 | dots[i+1][j]<<2 | dots[i+3][j+1]<<3 |
				dots[i+2][j+1]<<4 | dots[i+1][j+1]<<5 | dots[i][j]<<6 | dots[i][j+1]<<7)
		}
		lines[(v.height-i)/4] = "|" + string(line)
	}
	lines[v.height/4] = "+" + strings.Repeat("-", v.width/2)
	return lines
}
