package mkrg

type graph struct {
	name    string
	metrics []metric
}

type metric struct {
	name    string
	stacked bool
}

var systemGraphs = []graph{
	graph{
		name: "loadavg",
		metrics: []metric{
			metric{"loadavg1", false},
			metric{"loadavg5", false},
			metric{"loadavg15", false},
		},
	},
	graph{
		name: "cpu",
		metrics: []metric{
			metric{"cpu.user.percentage", true},
			metric{"cpu.nice.percentage", true},
			metric{"cpu.system.percentage", true},
			metric{"cpu.irq.percentage", true},
			metric{"cpu.softirq.percentage", true},
			metric{"cpu.iowait.percentage", true},
			metric{"cpu.steal.percentage", true},
			metric{"cpu.guest.percentage", true},
			metric{"cpu.idle.percentage", true},
		},
	},
	graph{
		name: "memory",
		metrics: []metric{
			metric{"memory.used", true},
			metric{"memory.mem_available", true},
			metric{"memory.buffers", true},
			metric{"memory.cached", true},
			metric{"memory.total", false},
			metric{"memory.free", true},
			metric{"memory.pagefile_total", false},
			metric{"memory.swap_used", false},
			metric{"memory.swap_cached", false},
			metric{"memory.pagefile_free", false},
			metric{"memory.swap_total", false},
		},
	},
	graph{
		name: "disk",
		metrics: []metric{
			metric{"disk.#.reads.delta", false},
			metric{"disk.#.writes.delta", false},
		},
	},
	graph{
		name: "interface",
		metrics: []metric{
			metric{"interface.#.rxBytes.delta", false},
			metric{"interface.#.txBytes.delta", false},
		},
	},
	graph{
		name: "filesystem",
		metrics: []metric{
			metric{"filesystem.#.used", false},
			metric{"filesystem.#.size", false},
		},
	},
}
