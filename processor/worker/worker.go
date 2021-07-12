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

type JobChannel chan Job

type Worker struct {
	ID      int
	JobChan JobChannel
}

func NewWorker(ID int, JobChan chan Job) *Worker {
	return &Worker{ID, JobChan}
}

func (wr *Worker) Start() {
	c := &http.Client{Timeout: time.Millisecond * 15000}
	go func() {
		for job := range wr.JobChan {
			callApi(job.ID, wr.ID, c)
		}
	}()
}

func callApi(num, id int, c *http.Client) {
	baseURL := "https://jsonplaceholder.typicode.com/albums/%d"

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
	re, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//log.Printf("error reading response body :: error is %+v", err)
		return
	}
	log.Printf("%d  :: ok", string(re))
}

func (w1 *Worker) Send(message Job) {
	w1.JobChan <- message
}
