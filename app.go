package mkrg

import (
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mackerelio/mackerel-client-go"
	"golang.org/x/crypto/ssh/terminal"
)

// App ...
type App struct {
	client *mackerel.Client
	hostID string
}

// NewApp creates a new app.
func NewApp(client *mackerel.Client, hostID string) *App {
	return &App{
		client: client,
		hostID: hostID,
	}
}

// Run the app.
func (app *App) Run() error {
	metricNamesMap, err := app.getMetricNamesMap()
	if err != nil {
		return err
	}
	termWidth, _, err := terminal.GetSize(0)
	if err != nil {
		return err
	}
	var maxColumn int
	if termWidth > 160 {
		maxColumn = 3
	} else if termWidth > 80 {
		maxColumn = 2
	} else {
		maxColumn = 1
	}
	width := (termWidth+4)/maxColumn - 4
	height := width / 8 * 3
	until := time.Now().Round(time.Minute)
	from := until.Add(-time.Duration(width*3) * time.Minute)
	var ui ui
	if os.Getenv("TERM_PROGRAM") == "iTerm.app" && os.Getenv("MKRG_VIEWER") == "" ||
		os.Getenv("MKRG_VIEWER") == "iTerm2" {
		ui = newIterm2UI(height, width, maxColumn, from, until)
	} else if os.Getenv("MKRG_VIEWER") == "Sixel" {
		ui = newSixel(height, width, maxColumn, from, until)
	} else {
		from = until.Add(-time.Duration(width*2) * time.Minute)
		ui = newTui(height, width, maxColumn, until)
	}
	for _, graph := range systemGraphs {
		var metricNames []string
		for _, metric := range graph.metrics {
			metricNames = append(metricNames, filterMetricNames(metricNamesMap, metric.name)...)
		}
		if len(metricNames) == 0 {
			continue
		}
		ms, err := app.fetchMetrics(metricNames, from, until)
		if err != nil {
			return err
		}
		ms.AddMemorySwapUsed()
		ms.Stack(graph)
		if err := ui.output(graph, ms); err != nil {
			return err
		}
	}
	return ui.cleanup()
}

func (app *App) fetchMetrics(metricNames []string, from, until time.Time) (metricsByName, error) {
	ms := make(metricsByName, len(metricNames))
	for _, metricName := range metricNames {
		metrics, err := app.client.FetchHostMetricValues(app.hostID, metricName, from.Unix(), until.Unix())
		if err != nil {
			return nil, err
		}
		ms.Add(metricName, metrics)
	}
	return ms, nil
}

func (app *App) getMetricNamesMap() (map[string]bool, error) {
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
	if name == "memory.swap_used" {
		if metricNamesMap["memory.swap_total"] && metricNamesMap["memory.swap_free"] {
			return []string{"memory.swap_free"}
		}
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
