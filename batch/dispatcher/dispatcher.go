package dispatcher

import (
	"github.com/concurrency/batch/worker"
)

type Dispatcher struct {
	Workers  []*worker.Worker  // this is the list of workers that dispatcher tracks
	WorkChan worker.JobChannel // client submits job to this channel
	Queue    worker.JobQueue   // this is the shared JobPool between the workers
}

// NewDispatcher returns a new dispatcher. A Dispatcher communicates between the client
// and the worker. Its main job is to receive a job and share it on the WorkPool
// WorkPool is the link between the dispatcher and all the workers as
// the WorkPool of the dispatcher is common JobPool for all the workers
func NewDispatcher(num int) *Dispatcher {
	return &Dispatcher{
		Workers:  make([]*worker.Worker, num),
		WorkChan: make(worker.JobChannel),
		Queue:    make(worker.JobQueue),
	}
}

// Start creates pool of num count of workers.
func (d *Dispatcher) Start() *Dispatcher {
	l := len(d.Workers)

	for i := 1; i <= l; i++ {
		wrk := worker.NewWorker(i, make(worker.JobChannel), d.Queue, make(chan struct{}))
		wrk.Start()
		d.Workers = append(d.Workers, wrk)
	}

	go d.process()

	return d
}

// process listens to a job submitted on WorkChan and
// relays it to the WorkPool. The WorkPool is shared between
// the workers.
func (d *Dispatcher) process() {
	for job := range d.WorkChan {
		// wait for a worker to submit JobChan to Queue
		// note that this Queue is shared among all workers.
		// Whenever there is an available JobChan on Queue pull it

		jobChan := <-d.Queue

		// Once a jobChan is available, send the submitted Job on this JobChan
		jobChan <- job
	}
}

func (d *Dispatcher) Submit(job worker.Job) {
	d.WorkChan <- job
}
