package BlueAbyss

import "time"

type worker struct {
	owner       *Pool
	taskQueue   chan func()
	recycleTime time.Time
}

func (w *worker) run() {
	w.owner.increaseWorker()
	go func() {
		defer func() {
			w.owner.decreaseWorker()
			w.owner.cache.Put(w)
		}()

		for work := range w.taskQueue {
			if work == nil {
				return
			}
			work()
			// finish task
			// let worker go back the pool
			w.owner.putWorker(w)
		}
	}()
}
