package mkrg

import (
	"fmt"
	"regexp"
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
	metricNamesMap, err := app.getMetricNamesMap()
	if err != nil {
		return err
	}
	now := time.Now().Round(time.Minute)
	from := now.Add(-2 * time.Hour)
	height, width := 80, 120
	var column int
	maxColumn := 2
	lines := make([]string, height/4+1)
	for _, graph := range systemGraphs {
		var metricNames []string
		for _, metric := range graph.metrics {
			metricNames = append(metricNames, filterMetricNames(metricNamesMap, metric.name)...)
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
		v := newViewer(graph, height, width)
		for i, l := range v.GetLines(ms, from) {
			lines[i] += l
			if column < maxColumn-1 {
				lines[i] += "    "
			}
		}
		if column == maxColumn-1 {
			for i := range lines {
				fmt.Println(lines[i])
				lines[i] = ""
			}
			column = 0
		} else {
			column += 1
		}
	}
	if column > 0 {
		for i := range lines {
			fmt.Println(lines[i])
		}
	}
	return nil
}

func (app *app) getMetricNamesMap() (map[string]bool, error) {
	metricNames, err := app.client.ListHostMetricNames(app.hostID)
	if err != nil {
		return nil, err
	}
	metricNamesMap := make(map[string]bool, len(metricNames))
	for _, metricName := range metricNames {
		metricNamesMap[metricName] = true
	}
	return metricNamesMap, nil
}

func filterMetricNames(metricNamesMap map[string]bool, name string) []string {
	if metricNamesMap[name] {
		return []string{name}
	}
	namePattern := regexp.MustCompile(
		"^" + strings.Replace(name, "#", `[-a-zA-Z0-9_]+`, -1) + "$",
	)
	var metricNames []string
	for metricName := range metricNamesMap {
		if namePattern.MatchString(metricName) {
			metricNames = append(metricNames, metricName)
		}
	}
	return metricNames
}
