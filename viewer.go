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
	h, w := (v.height-3)*4, (v.width-6)*2
	dots := make([][]int, h)
	for i := range dots {
		dots[i] = make([]int, w)
	}
	maxValue := math.Max(ms.MaxValue(), 1.0) * 1.1
	tick := getTick(maxValue)
	format, scale := formatAxisY(tick, maxValue)
	from := until.Add(-time.Duration(w/2) * time.Minute)
	for _, metrics := range ms {
		prevPrevTime, prevTime, nextTime, prevX, prevY := int64(0), int64(0), int64(0), -1.0, 0.0
		for i, m := range metrics {
			x := float64((m.Time - from.Unix()) / 30)
			y := m.Value.(float64) / maxValue * float64(h)
			if 0 <= x {
				if i < len(metrics)-1 {
					nextTime = metrics[i+1].Time
				}
				start, step := 0.0, math.Min(1.0/math.Sqrt((x-prevX)*(x-prevX)+(y-prevY)*(y-prevY)), 1.0)
				if prevX < 0 || prevTime < m.Time-3*60 && (prevTime-3*60 < prevPrevTime || prevPrevTime == 0 && nextTime-3*60 < m.Time) {
					start = 1.0
				}
				prevPrevTime, prevTime = prevTime, m.Time
				for p := start; p <= 1.0; p += step {
					dots[int(prevY*(1.0-p)+y*p)][int(prevX*(1.0-p)+x*p)] = 1
				}
			}
			prevX, prevY = x, y
		}
	}
	lines := make([]string, v.height)
	leftPadding := int(math.Max(float64((v.width-len(v.graph.name)+1)/2), 0))
	lines[0] = strings.Repeat(" ", leftPadding) + v.graph.name
	lines[0] += strings.Repeat(" ", int(math.Max(float64(v.width-len(lines[0])), 0)))
	for y := 0.0; y < maxValue; y += tick {
		if y > 0.0 {
			posY := v.height - int(math.Round(y/maxValue*float64(v.height-3))) - 2
			lines[posY] = fmt.Sprintf("%4s +", fmt.Sprintf(format, y/scale))
		}
	}
	line := make([]rune, w/2)
	for i := h - 4; i >= 0; i -= 4 {
		for j := 0; j < w; j += 2 {
			line[j/2] = rune(0x2800 | dots[i+3][j] | dots[i+2][j]<<1 | dots[i+1][j]<<2 | dots[i+3][j+1]<<3 |
				dots[i+2][j+1]<<4 | dots[i+1][j+1]<<5 | dots[i][j]<<6 | dots[i][j+1]<<7)
		}
		y := (h - i) / 4
		if lines[y] == "" {
			lines[y] = "     |"
		}
		lines[y] += string(line)
	}
	axisX := []rune("   0 +" + strings.Repeat("-", v.width-6))
	stepX := 30 * time.Minute
	var axisXLabels string
	for t := from.Truncate(stepX); !until.Before(t); t = t.Add(stepX) {
		offset := int(float64(t.Sub(from))/float64(until.Sub(from))*float64(v.width-6)) + 6
		if offset < 5 || len(axisX) <= offset {
			continue
		}
		axisX[offset] = '+'
		axisXLabels += strings.Repeat(" ", int(math.Max(float64(offset-len(axisXLabels)-2), 0)))
		axisXLabels += fmt.Sprintf("%1d:%02d", t.Hour(), t.Minute())
	}
	axisXLabels += strings.Repeat(" ", int(math.Max(float64(v.width-len(axisXLabels)), 0)))
	lines[v.height-2] += string(axisX)
	lines[v.height-1] = axisXLabels[:v.width]
	return lines
}
