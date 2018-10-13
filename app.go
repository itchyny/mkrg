package mkrg

import (
	"fmt"
	"math"
	"strings"
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
	metricsByName := make(metricsByName, len(metricNames))
	for _, metricName := range metricNames {
		metrics, err := app.client.FetchHostMetricValues(app.hostID, metricName, from.Unix(), now.Unix())
		if err != nil {
			return err
		}
		metricsByName.Add(metricName, metrics)
	}
	dots := make([][]int, 100)
	for i := range dots {
		dots[i] = make([]int, 180)
	}
	maxValue := math.Max(metricsByName.MaxValue(), 1.0) * 1.1
	for _, metrics := range metricsByName {
		for _, m := range metrics {
			x := (m.Time - from.Unix()) / 60
			y := int(m.Value.(float64) / maxValue * 100)
			dots[y][x] = 1
		}
	}
	line := make([]rune, 90)
	for i := 100 - 4; i >= 0; i -= 4 {
		for j := 0; j < 180; j += 2 {
			line[j/2] = rune(0x2800 | dots[i+3][j] | dots[i+2][j]<<1 | dots[i+1][j]<<2 | dots[i+3][j+1]<<3 |
				dots[i+2][j+1]<<4 | dots[i+1][j+1]<<5 | dots[i][j]<<6 | dots[i][j+1]<<7)
		}
		fmt.Printf("|%s\n", string(line))
	}
	fmt.Printf("+%s\n", strings.Repeat("-", 180/2))
	return nil
}
