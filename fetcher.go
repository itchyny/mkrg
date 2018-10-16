package mkrg

import (
	"time"

	"github.com/mackerelio/mackerel-client-go"
)

type fetcher struct {
	client *mackerel.Client
	sem    chan struct{}
}

func newFetcher(client *mackerel.Client) *fetcher {
	return &fetcher{client, make(chan struct{}, 5)}
}

func (f *fetcher) fetchMetric(hostID, metricName string, from, until time.Time) ([]mackerel.MetricValue, error) {
	f.sem <- struct{}{}
	metrics, err := f.client.FetchHostMetricValues(hostID, metricName, from.Unix(), until.Unix())
	<-f.sem
	return metrics, err
}
