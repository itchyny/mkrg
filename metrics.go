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

func (ms metricsByName) Stack(graph graph) {
	stackedValue := make(map[int64]float64)
	for _, metric := range graph.metrics {
		if !metric.stacked {
			continue
		}
		if metrics, ok := ms[metric.name]; ok {
			for i, m := range metrics {
				w := metrics[i].Value.(float64)
				if v, ok := stackedValue[m.Time]; ok {
					stackedValue[m.Time] = v + w
					metrics[i].Value = v + w
				} else {
					stackedValue[m.Time] = w
				}
			}
		}
	}
}
