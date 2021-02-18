package BlueAbyss

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	capacity int32

	runningWorker int32

	wq workers

	state int32

	lock sync.Locker

	cache sync.Pool

	opts *Options
}

func NewPool(size int, options ...option) (*Pool, error) {
	opts := loadOptions(options...)

	if size <= 0 {
		size = DefaultAbyssSize
	}

	p := &Pool{
		capacity:      int32(size),
		runningWorker: int32(0),
		wq:            newWorkQueue(size),
		state:         Abyss_OPENED,
		lock:          &sync.Mutex{},
		opts:          opts,
	}

	p.cache = sync.Pool{
		New: func() interface{} {
			return &worker{
				owner:     p,
				taskQueue: make(chan func(), 1),
			}
		},
	}

	go p.background()
	return p, nil
}

func (p *Pool) background() {
	exp := p.opts.ExpiryDuration
	heartBeat := time.NewTicker(exp)
	defer heartBeat.Stop()

	for range heartBeat.C {
		if p.closed() {
			return
		}

		expirychan := p.wq.expiryChan()
		for w := range expirychan {
			fmt.Println(w)
		}
	}
}

func (p *Pool) Summit(t func()) error {
	if p.closed() {
		return Abyss_CLOSE_ERROR
	}

	w := p.getWorker()
	if w == nil {
		return Abyss_SUMMIT_ERROR
	}
	w.taskQueue <- t
	return nil
}

func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

// return the number of running goroutine.
func (p *Pool) Busy() int {
	return int(atomic.LoadInt32(&p.runningWorker))
}

func (p *Pool) Free() int {
	return p.Cap() - p.Busy()
}

func (p *Pool) closed() bool {
	return atomic.LoadInt32(&p.state) == Abyss_CLOSED
}

func (p *Pool) increaseWorker() {
	atomic.AddInt32(&p.runningWorker, 1)
}

func (p *Pool) decreaseWorker() {
	atomic.AddInt32(&p.runningWorker, -1)
}

// get a available worker from abyss to run the tasks.
func (p *Pool) getWorker() (w *worker) {
	createNewWorker := func() {
		w = p.cache.Get().(*worker)
		w.run()
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	w = p.wq.get()
	if w != nil {
		return
	} else if p.Free() > 0 {
		createNewWorker() // create a new worker, waiting for the task.
	} else {
		// p.Free() == 0, workers are all at work.
		if p.options.NonBlocking {
			return
		}
	}
	return
}

// if worker finishs its task, put it back to abyss.
func (p *Pool) putWorker(w *worker) {

}

func (p *Pool) Release() {

}

func (p *Pool) _cache() {
	w := p.cache.Get().(*worker)
	w.run()
}
