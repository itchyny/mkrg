package mkrg

import (
	"fmt"
	"math"
	"time"

	"github.com/mackerelio/mackerel-client-go"
)

type app struct {
	client *mackerel.Client
	hostID string
}

func NewApp(client *mackerel.Client, hostID string) *app {
	return &app{
		client: client,
		hostID: hostID,
	}
}

func (app *app) Run() error {
	now := time.Now().Round(time.Minute)
	from := now.Add(-3 * time.Hour)
	metricNames := []string{"loadavg1", "loadavg5", "loadavg15"}
	metricsByName := make(map[string][]mackerel.MetricValue, len(metricNames))
	for _, metricName := range metricNames {
		metrics, err := app.client.FetchHostMetricValues(app.hostID, metricName, from.Unix(), now.Unix())
		if err != nil {
			return err
		}
		metricsByName[metricName] = metrics
	}
	dots := make([][]int, 100)
	for i := range dots {
		dots[i] = make([]int, 180)
	}
	maxValue := 0.0
	for _, metrics := range metricsByName {
		for _, m := range metrics {
			v := m.Value.(float64)
			if v > maxValue {
				maxValue = v
			}
		}
	}
	if maxValue == math.MaxFloat64 || maxValue <= 0 {
		maxValue = 1.0
	} else {
		maxValue *= 1.1
	}
	for _, metrics := range metricsByName {
		for _, m := range metrics {
			x := (m.Time - from.Unix()) / 60
			y := int(m.Value.(float64) / maxValue * 100)
			dots[y][x] = 1
		}
	}
	for i := 100 - 4; i >= 0; i -= 4 {
		for j := 0; j < 180; j += 2 {
			b := (0x2800 | dots[i+3][j] | dots[i+2][j]<<1 | dots[i+1][j]<<2 | dots[i+3][j+1]<<3 |
				dots[i+2][j+1]<<4 | dots[i+1][j+1]<<5 | dots[i][j]<<6 | dots[i][j+1]<<7)
			fmt.Printf("%c", rune(b))
		}
		fmt.Printf("\n")
	}
	return nil
}
