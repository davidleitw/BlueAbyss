package BlueAbyss

// pool state
const (
	// Opened state mean routine pool can accept task
	Abyss_OPENED int32 = iota

	// Closed state mean routine pool is closed, can't accept the task any more.
	Abyss_CLOSED
)

// pool default configuration
const (
	DefaultAbyssSize = 100000
)
