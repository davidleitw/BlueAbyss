package BlueAbyss

type worker struct {
	owner     *Pool
	taskQueue chan func()
}

func (w *worker) run() {
	w.owner.runUp()
	go func() {
		defer func() {
			w.owner.runDown()
		}()

		for work := range w.taskQueue {
			if work == nil {
				return
			}
			work() // finish task

			// let worker go back the pool
		}
	}()
}
