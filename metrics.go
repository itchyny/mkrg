package mkrg

import "github.com/mackerelio/mackerel-client-go"

type metricsByName map[string][]mackerel.MetricValue

func (ms metricsByName) Add(metricName string, metricValues []mackerel.MetricValue) {
	ms[metricName] = metricValues
}

func (ms metricsByName) MaxValue() float64 {
	maxValue := 0.0
	for _, metrics := range ms {
		for _, m := range metrics {
			v := m.Value.(float64)
			if v > maxValue {
				maxValue = v
			}
		}
	}
	return maxValue
}
