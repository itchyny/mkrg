package mkrg

type ui interface {
	output(graph, metricsByName) error
	cleanup() error
}
