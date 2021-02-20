package BlueAbyss

import (
	"log"
	"time"
)

type worker struct {
	owner      *Pool
	taskQueue  chan func()
	expirytime time.Time
}

func (w *worker) run() {
	w.owner.increaseWorker()
	go func() {
		defer func() {
			w.owner.decreaseWorker()
			w.owner.cache.Put(w)

			// recover from panic
			if r := recover(); r != nil {
				if w.owner.opts.PanicHandler != nil {
					w.owner.opts.PanicHandler(r)
				} else {
					log.Println("worker panic: ", r)
				}
			}
		}()

		for work := range w.taskQueue {
			if work == nil {
				return
			}
			work()
			w.owner.putWorker(w)
		}
	}()
}
