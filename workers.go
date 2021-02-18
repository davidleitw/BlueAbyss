package BlueAbyss

type workers interface {
	get() *worker
	expiryChan() chan *worker
	insertWorker(*worker)
}

type workerQueue struct {
	ws     chan *worker
	exp    chan *worker
	backUp chan *worker
}

func newWorkQueue(size int) workers {
	return &workerQueue{
		ws:     make(chan *worker, size),
		exp:    make(chan *worker, size),
		backUp: nil,
	}
}

func (wq *workerQueue) insertWorker(w *worker) {
	wq.ws <- w
}

func (wq *workerQueue) get() *worker {
	select {
	case w := <-wq.ws:
		return w
	default:
		return nil
	}
}

func (wq *workerQueue) expiryChan() chan *worker {
	return wq.exp
}
