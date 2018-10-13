package mkrg

import (
	"fmt"
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
	allMetricNames := []string{"loadavg1", "loadavg5", "loadavg15", "cpu.user.percentage", "cpu.system.percentage"}
	graphNames := []string{"loadavg", "cpu"}
	for _, graphName := range graphNames {
		var metricNames []string
		for _, metricName := range allMetricNames {
			if strings.HasPrefix(metricName, graphName) {
				metricNames = append(metricNames, metricName)
			}
		}
		if len(metricNames) == 0 {
			continue
		}
		ms := make(metricsByName, len(metricNames))
		for _, metricName := range metricNames {
			metrics, err := app.client.FetchHostMetricValues(app.hostID, metricName, from.Unix(), now.Unix())
			if err != nil {
				return err
			}
			ms.Add(metricName, metrics)
		}
		v := newViewer(graphName, 100, 180)
		for _, l := range v.GetLines(ms, from) {
			fmt.Println(l)
		}
	}
	return nil
}
