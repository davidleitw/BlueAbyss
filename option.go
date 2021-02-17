package BlueAbyss

import "time"

type Options struct {
	ExpiryDuration time.Duration

	MaxWorkerNumber int

	NonBlocking bool
}

type option func(opts *Options)

// load option for options.
func loadOptions(options ...option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

func SetUpExpiryDuration(expiry time.Duration) option {
	return func(opts *Options) {
		opts.ExpiryDuration = expiry
	}
}

func SetUpMaxWorkerNumber(maxWorkerNum int) option {
	return func(opts *Options) {
		opts.MaxWorkerNumber = maxWorkerNum
	}
}

func SetUpNonBlocking(nonb bool) option {
	return func(opts *Options) {
		opts.NonBlocking = nonb
	}
}
