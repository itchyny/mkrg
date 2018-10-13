package mkrg

import (
	"fmt"
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
	for _, graph := range systemGraphs {
		var metricNames []string
		for _, metric := range graph.metrics {
			metricNames = append(metricNames, metric.name)
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
		v := newViewer(graph, 100, 180)
		for _, l := range v.GetLines(ms, from) {
			fmt.Println(l)
		}
	}
	return nil
}
