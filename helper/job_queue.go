package helper

import "sync"

// Job interface for job processing
type Job interface {
	Process()
}

// Worker threads that process jobs
type Worker struct {
	done             sync.WaitGroup
	readyPool        chan chan Job
	assignedJobQueue chan Job
	quit             chan bool
}

// JobQueue for enqueueing jobs
type JobQueue struct {
	internalQueue     chan Job
	readyPool         chan chan Job
	workers           []*Worker
	dispatcherStopped sync.WaitGroup
	workersStopped    sync.WaitGroup
	quit              chan bool
}

// NewJobQueue creates a new job queue
func NewJobQueue(maxWorkers int) *JobQueue {
	workersStopped := sync.WaitGroup{}
	readyPool := make(chan chan Job, maxWorkers)
	workers := make([]*Worker, maxWorkers, maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		workers[i] = NewWorker(readyPool, workersStopped)
	}

	return &JobQueue{
		internalQueue:     make(chan Job),
		readyPool:         readyPool,
		workers:           workers,
		dispatcherStopped: sync.WaitGroup{},
		workersStopped:    workersStopped,
		quit:              make(chan bool),
	}
}

// Start workers and dispatcher
func (q *JobQueue) Start() {
	for i := 0; i < len(q.workers); i++ {
		q.workers[i].Start()
	}
	go q.dispatch()
}

// Stop workers and dispatcher
func (q *JobQueue) Stop() {
	q.quit <- true
	q.dispatcherStopped.Wait()
}

func (q *JobQueue) dispatch() {
	q.dispatcherStopped.Add(1)
	for {
		select {
		case job := <-q.internalQueue:
			workerChannel := <-q.readyPool
			workerChannel <- job
		case <-q.quit:
			for i := 0; i < len(q.workers); i++ {
				q.workers[i].Stop()
			}
			q.workersStopped.Wait()
			q.dispatcherStopped.Done()
			return
		}
	}
}

// Submit adds a new job
func (q *JobQueue) Submit(job Job) {
	q.internalQueue <- job
}

// NewWorker creates a new worker
func NewWorker(readyPool chan chan Job, done sync.WaitGroup) *Worker {
	return &Worker{
		done:             done,
		readyPool:        readyPool,
		assignedJobQueue: make(chan Job),
		quit:             make(chan bool),
	}
}

// Start begins job processing loop
func (w *Worker) Start() {
	go func() {
		w.done.Add(1)
		for {
			w.readyPool <- w.assignedJobQueue
			select {
			case job := <-w.assignedJobQueue:
				job.Process()
			case <-w.quit:
				w.done.Done()
				return
			}
		}
	}()
}

// Stop the worker
func (w *Worker) Stop() {
	w.quit <- true
}
