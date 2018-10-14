package mkrg

import (
	"sort"

	"github.com/mackerelio/mackerel-client-go"
)

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

func (ms metricsByName) MetricNames() []string {
	metricNames := make([]string, 0, len(ms))
	for name := range ms {
		metricNames = append(metricNames, name)
	}
	sort.Strings(metricNames)
	return metricNames
}

func (ms metricsByName) AddMemorySwapUsed() {
	if totalMetrics, ok := ms["memory.swap_total"]; ok {
		if freeMetrics, ok := ms["memory.swap_free"]; ok {
			usedMetrics := make([]mackerel.MetricValue, 0, len(totalMetrics))
			for i, j := 0, 0; i < len(totalMetrics) && j < len(freeMetrics); i++ {
				for j < len(freeMetrics) && totalMetrics[i].Time > freeMetrics[j].Time {
					j++
				}
				if totalMetrics[i].Time == freeMetrics[j].Time {
					usedMetrics = append(usedMetrics, mackerel.MetricValue{
						Time:  totalMetrics[i].Time,
						Value: totalMetrics[i].Value.(float64) - freeMetrics[j].Value.(float64),
					})
				}
				for j < len(freeMetrics) && totalMetrics[i].Time >= freeMetrics[j].Time {
					j++
				}
			}
			delete(ms, "memory.swap_free")
			ms.Add("memory.swap_used", usedMetrics)
		}
	}
}
