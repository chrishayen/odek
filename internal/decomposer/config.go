package decomposer

// Config controls how a decomposition run is executed: how many parallel
// initial attempts, how deep to recurse, and how many total runes to cap at.
type Config struct {
	ParallelInitial int
	MaxDepth        int
	RuneCap         int
	Recurse         bool
}

// ConfigForEffort maps an effort level (1-5) to a run config. Levels
// outside 1-5 fall through to the level-3 default.
func ConfigForEffort(level int) Config {
	switch level {
	case 1:
		return Config{ParallelInitial: 1, MaxDepth: 0, RuneCap: 10, Recurse: false}
	case 2:
		return Config{ParallelInitial: 1, MaxDepth: 10, RuneCap: 25, Recurse: true}
	case 3:
		return Config{ParallelInitial: 3, MaxDepth: 10, RuneCap: 50, Recurse: true}
	case 4:
		return Config{ParallelInitial: 5, MaxDepth: 10, RuneCap: 100, Recurse: true}
	case 5:
		return Config{ParallelInitial: 5, MaxDepth: 10, RuneCap: 200, Recurse: true}
	}
	return Config{ParallelInitial: 3, MaxDepth: 10, RuneCap: 50, Recurse: true}
}
