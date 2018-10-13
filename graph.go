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
			metric{"cpu.system.percentage", true},
			metric{"cpu.idle.percentage", true},
		},
	},
}
