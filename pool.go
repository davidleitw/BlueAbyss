package BlueAbyss

import (
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	capacity int32
	running  int32
	blocking int
	wq       workers
	state    int32
	lock     sync.Locker
	cache    sync.Pool
	opts     *Options
}

func NewPool(size int, options ...option) (*Pool, error) {
	opts := loadOptions(options...)
	if size <= 0 {
		size = DefaultAbyssSize
	}
	if opts.BackupWqSize <= 0 {
		opts.BackupWqSize = int32(size)
	}

	p := &Pool{
		capacity: int32(size),
		running:  int32(0),
		blocking: 0,
		wq:       newWorkQueue(size),
		state:    Abyss_OPENED,
		lock:     &sync.Mutex{},
		opts:     opts,
	}

	p.cache = sync.Pool{
		New: func() interface{} {
			return &worker{
				owner:     p,
				taskQueue: make(chan func(), 1),
			}
		},
	}

	expiryT := time.Now().Add(p.opts.ExpiryDuration)
	if opts.PreAlloc {
		for i := 0; i < p.Cap(); i++ {
			w := p.cache.Get().(*worker)
			w.expirytime = expiryT
			w.run()
			p.wq.insertWorker(w)
		}
	}

	go p.expiredCollection()
	go p.backupWorkerQueue()
	return p, nil
}

func (p *Pool) expiredCollection() {
	exp := p.opts.ExpiryDuration
	heartBeat := time.NewTicker(exp)
	defer heartBeat.Stop()

	for range heartBeat.C {
		if p.closed() {
			return
		}

		expirychan := p.wq.expiredChan()
		expirytime := time.Now().Add(p.opts.ExpiryDuration)
		for w := range expirychan {
			w.expirytime = expirytime
			if ok := p.wq.insertWorker(w); !ok {
				if p.wq.hasbackup() {
					p.wq.insertWorkerBackup(w)
				} else {
					break
				}
			}
		}
	}
}

func (p *Pool) backupWorkerQueue() {
	heartBeat := time.NewTicker(20 * time.Second)

	for range heartBeat.C {
		if p.closed() {
			return
		}

		p.lock.Lock()
		if !p.wq.hasbackup() && p.usageRate() >= p.opts.BackupRate {
			// ready to set backup worker queue.
			p.wq.backup(int(p.opts.BackupWqSize))
			for i := 0; i < int(p.opts.BackupWqSize); i++ {
				w := p.cache.Get().(*worker)
				w.run()
				p.wq.insertWorkerBackup(w)
			}
		}
		p.lock.Unlock()
	}
}

func (p *Pool) usageRate() float32 {
	return float32(p.Busy()) / float32(p.Cap())
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
	return int(atomic.LoadInt32(&p.running))
}

func (p *Pool) Full() bool {
	return p.Cap() == p.Busy()
}

func (p *Pool) Free() int {
	return p.Cap() - p.Busy()
}

func (p *Pool) closed() bool {
	return atomic.LoadInt32(&p.state) == Abyss_CLOSED
}

func (p *Pool) increaseWorker() {
	atomic.AddInt32(&p.running, 1)
}

func (p *Pool) decreaseWorker() {
	atomic.AddInt32(&p.running, -1)
}

// get a available worker from abyss to run the tasks.
func (p *Pool) getWorker() (w *worker) {
	createNewWorker := func() {
		w = p.cache.Get().(*worker)
		w.run()
	}

	p.lock.Lock()

	w = p.wq.get()
	if w != nil {
		p.lock.Unlock()
		return
	} else if !p.Full() {
		createNewWorker() // create a new worker, waiting for the task.
		p.lock.Unlock()
	} else {
		// p.Free() == 0, workers are all at work.

	}
	return
}

// if worker finishs its task, put it back to abyss.
func (p *Pool) putWorker(w *worker) {

}

func (p *Pool) Release() {

}
