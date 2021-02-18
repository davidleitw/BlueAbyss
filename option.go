package BlueAbyss

import "time"

type Options struct {
	ExpiryDuration  time.Duration
	UnlimitedSize   bool
	PreAlloc        bool
	MaxWorkerNumber int
	NonBlocking     bool
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

func SetUpUnlimitedSize(unlimited bool) option {
	return func(opts *Options) {
		opts.UnlimitedSize = unlimited
	}
}

func SetUpPreAlloc(prealloc bool) option {
	return func(opts *Options) {
		opts.PreAlloc = prealloc
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
