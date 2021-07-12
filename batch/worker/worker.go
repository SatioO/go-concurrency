package worker

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Job struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Worker is a a single processor. Typically its possible to
// start multiple workers for better throughput
type Worker struct {
	ID      int           // id of the worker
	JobChan JobChannel    // a channel to receive single unit of work
	Queue   JobQueue      // shared between all workers.
	Quit    chan struct{} // a channel to quit working
}

type JobChannel chan Job
type JobQueue chan chan Job

func NewWorker(ID int, JobChan JobChannel, Queue JobQueue, Quit chan struct{}) *Worker {
	return &Worker{
		ID:      ID,
		JobChan: JobChan,
		Queue:   Queue,
		Quit:    Quit,
	}
}

func (wr *Worker) Start() {
	c := &http.Client{Timeout: time.Millisecond * 15000}
	go func() {
		for {
			// when available, put the JobChan again on the JobPool
			// and wait to receive a job
			wr.Queue <- wr.JobChan
			select {
			case job := <-wr.JobChan:
				// when a job is received, process
				callApi(job.ID, wr.ID, c)
			case <-wr.Quit:
				// a signal on this channel means someone triggered
				// a shutdown for this worker
				close(wr.JobChan)
				return
			}
		}
	}()
}

func callApi(num, id int, c *http.Client) {
	baseURL := "https://jsonplaceholder.typicode.com/photos/%d"

	ur := fmt.Sprintf(baseURL, num)
	// fmt.Println(fmt.Sprintf(baseURL, num), id)
	req, err := http.NewRequest(http.MethodGet, ur, nil)
	if err != nil {
		//log.Printf("error creating a request for term %d :: error is %+v", num, err)
		return
	}
	res, err := c.Do(req)
	if err != nil {
		//log.Printf("error querying for term %d :: error is %+v", num, err)
		return
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		//log.Printf("error reading response body :: error is %+v", err)
		return
	}
	log.Printf("%d  :: ok", id)
}

func (wr *Worker) Stop() {
	close(wr.Quit)
}
