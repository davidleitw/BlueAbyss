package BlueAbyss

import "time"

type Options struct {
	ExpiryDuration  time.Duration
	PreAlloc        bool
	MaxWorkerNumber int
	NonBlocking     bool
	PanicHandler    func(interface{})
	BackupRate      float32
	BackupWqSize    int32
}

type option func(opts *Options)

// load option for options.
func loadOptions(options ...option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	if opts.ExpiryDuration == 0 {
		opts.ExpiryDuration = DefaultAbyssExpiryTime
	}
	if opts.BackupRate <= 0.0 {
		opts.BackupRate = DefaultAbyssBackupRate
	}
	return opts
}

func SetUpExpiryDuration(expiry time.Duration) option {
	return func(opts *Options) {
		opts.ExpiryDuration = expiry
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

func SetUpNonBlocking(nonblock bool) option {
	return func(opts *Options) {
		opts.NonBlocking = nonblock
	}
}

func SetUpPanicHandler(panicHandle func(interface{})) option {
	return func(opts *Options) {
		opts.PanicHandler = panicHandle
	}
}

func SetUpBackupRate(rate float32) option {
	return func(opts *Options) {
		opts.BackupRate = rate
	}
}

func SetUpBackupWqSize(size int32) option {
	return func(opts *Options) {
		opts.BackupWqSize = size
	}
}
