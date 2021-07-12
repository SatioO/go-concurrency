package main

import (
	"fmt"
	"time"

	"github.com/concurrency/batch/dispatcher"
	"github.com/concurrency/batch/worker"
)

func main() {
	dd := dispatcher.NewDispatcher(16).Start()

	terms := []int{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3,
		4, 1, 2, 3, 4}

	for i := range terms {
		dd.Submit(worker.Job{
			ID:        i,
			Name:      fmt.Sprintf("JobID::%d", i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

}
