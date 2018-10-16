package mkrg

import (
	"time"

	"github.com/mackerelio/mackerel-client-go"
)

type fetcher struct {
	client *mackerel.Client
	sem    chan struct{}
}

type metricAPIResult struct {
	metricName string
	metrics    []mackerel.MetricValue
	err        error
}

func newFetcher(client *mackerel.Client) *fetcher {
	return &fetcher{client, make(chan struct{}, 5)}
}

func (f *fetcher) fetchMetric(hostID, metricName string, from, until time.Time) <-chan metricAPIResult {
	ch := make(chan metricAPIResult)
	f.sem <- struct{}{}
	go func() {
		metrics, err := f.client.FetchHostMetricValues(hostID, metricName, from.Unix(), until.Unix())
		<-f.sem
		ch <- metricAPIResult{metricName, metrics, err}
	}()
	return ch
}
