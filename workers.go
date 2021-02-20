package BlueAbyss

type workers interface {
	get() *worker
	getback() *worker
	backup(int)
	hasbackup() bool
	expiredChan() <-chan *worker
	insertWorker(*worker) bool
	insertWorkerBackup(*worker) bool
}

type workerQueue struct {
	queue   chan *worker
	expired chan *worker
	back    chan *worker
}

func newWorkQueue(size int) workers {
	return &workerQueue{
		queue:   make(chan *worker, size),
		expired: make(chan *worker, size),
		back:    nil,
	}
}

func (wq *workerQueue) get() *worker {
	select {
	case w := <-wq.queue:
		return w
	default:
		return nil
	}
}

func (wq *workerQueue) getback() *worker {
	select {
	case w := <-wq.back:
		return w
	default:
		return nil
	}
}

func (wq *workerQueue) backup(size int) {
	wq.back = make(chan *worker, size)
}

func (wq *workerQueue) hasbackup() bool {
	return wq.back != nil
}

func (wq *workerQueue) expiredChan() <-chan *worker {
	return wq.expired
}

func (wq *workerQueue) insertWorker(w *worker) bool {
	select {
	case wq.queue <- w:
		return true
	default:
		return false
	}
}

func (wq *workerQueue) insertWorkerBackup(w *worker) bool {
	if wq.back == nil {
		return false
	}
	select {
	case wq.back <- w:
		return true
	default:
		return false
	}
}
