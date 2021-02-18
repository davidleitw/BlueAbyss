package BlueAbyss

import (
	"errors"
)

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

var (
	Abyss_CLOSE_ERROR  = errors.New("Pool is closed, can't summit tasks.")
	Abyss_SUMMIT_ERROR = errors.New("Get worker error.")
)
