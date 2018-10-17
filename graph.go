package mkrg

import (
	"regexp"
	"strings"
)

type graph struct {
	name    string
	metrics []metric
}

type metric struct {
	name    string
	stacked bool
}

var systemGraphs = []graph{
	{
		name: "loadavg",
		metrics: []metric{
			{"loadavg1", false},
			{"loadavg5", false},
			{"loadavg15", false},
		},
	},
	{
		name: "cpu",
		metrics: []metric{
			{"cpu.user.percentage", true},
			{"cpu.nice.percentage", true},
			{"cpu.system.percentage", true},
			{"cpu.irq.percentage", true},
			{"cpu.softirq.percentage", true},
			{"cpu.iowait.percentage", true},
			{"cpu.steal.percentage", true},
			{"cpu.guest.percentage", true},
			{"cpu.idle.percentage", true},
		},
	},
	{
		name: "memory",
		metrics: []metric{
			{"memory.used", true},
			{"memory.mem_available", true},
			{"memory.buffers", true},
			{"memory.cached", true},
			{"memory.total", false},
			{"memory.free", true},
			{"memory.pagefile_total", false},
			{"memory.swap_used", false},
			{"memory.swap_cached", false},
			{"memory.pagefile_free", false},
			{"memory.swap_total", false},
		},
	},
	{
		name: "disk",
		metrics: []metric{
			{"disk.#.reads.delta", false},
			{"disk.#.writes.delta", false},
		},
	},
	{
		name: "interface",
		metrics: []metric{
			{"interface.#.rxBytes.delta", false},
			{"interface.#.txBytes.delta", false},
		},
	},
	{
		name: "filesystem",
		metrics: []metric{
			{"filesystem.#.used", false},
			{"filesystem.#.size", false},
		},
	},
}

func metricNamePattern(name string) *regexp.Regexp {
	return regexp.MustCompile(
		`\A` + strings.Replace(name, "#", `([-a-zA-Z0-9_]+)`, -1) + `\z`,
	)
}
