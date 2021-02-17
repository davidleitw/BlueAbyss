package BlueAbyss

import "sync/atomic"

type Pool struct {
	capacity int32

	runningWorker int32

	state int32

	opts *Options
}

func NewPool(size int, options ...option) (*Pool, error) {
	opts := loadOptions(options...)

	if size <= 0 {
		size = DefaultAbyssSize
	}

	p := &Pool{
		capacity: int32(size),
		state:    Abyss_OPENED,
		opts:     opts,
	}

	return p, nil
}

func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

// return the number of running goroutine
func (p *Pool) Busy() int {
	return int(atomic.LoadInt32(&p.runningWorker))
}

func (p *Pool) Free() int {
	return p.Cap() - p.Busy()
}

func (p *Pool) runUp() {
	atomic.AddInt32(&p.runningWorker, 1)
}

func (p *Pool) runDown() {
	atomic.AddInt32(&p.runningWorker, -1)
}

func (p *Pool) closed() bool {
	if atomic.LoadInt32(&p.state) == Abyss_CLOSED {
		return true
	}
	return false
}

func (p *Pool) Release() {

}
