package mkrg

import (
	"os"
	"sync"
	"time"

	"github.com/mackerelio/mackerel-client-go"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sync/errgroup"
)

// App ...
type App struct {
	client  *mackerel.Client
	hostID  string
	fetcher *fetcher
}

// NewApp creates a new app.
func NewApp(client *mackerel.Client, hostID string) *App {
	return &App{
		client:  client,
		hostID:  hostID,
		fetcher: newFetcher(client),
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
	eg, mu := errgroup.Group{}, new(sync.Mutex)
	for _, graph := range systemGraphs {
		graph := graph
		var metricNames []string
		for _, metric := range graph.metrics {
			metricNames = append(metricNames, filterMetricNames(metricNamesMap, metric.name)...)
		}
		if len(metricNames) == 0 {
			continue
		}
		eg.Go(func() error {
			ms, err := app.fetchMetrics(graph, metricNames, from, until)
			if err != nil {
				return err
			}
			mu.Lock()
			defer mu.Unlock()
			return ui.output(graph, ms)
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	mu.Lock()
	defer mu.Unlock()
	return ui.cleanup()
}

func (app *App) fetchMetrics(graph graph, metricNames []string, from, until time.Time) (metricsByName, error) {
	eg, mu := errgroup.Group{}, new(sync.Mutex)
	ms := make(metricsByName, len(metricNames))
	for _, metricName := range metricNames {
		metricName := metricName
		eg.Go(func() error {
			res := <-app.fetcher.fetchMetric(app.hostID, metricName, from, until)
			if res.err != nil {
				return res.err
			}
			mu.Lock()
			defer mu.Unlock()
			ms.Add(res.metricName, res.metrics)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	ms.AddMemorySwapUsed()
	ms.Stack(graph)
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
	namePattern := metricNamePattern(name)
	var metricNames []string
	for metricName := range metricNamesMap {
		if namePattern.MatchString(metricName) {
			metricNames = append(metricNames, metricName)
		}
	}
	return metricNames
}
