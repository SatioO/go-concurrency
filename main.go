package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	input := buffer(numbers)

	re1 := getInputChan(input)

	merged := fanIn(re1)

	for i := range merged {
		fmt.Println(i)
	}
}

func buffer(numbers []int) <-chan int {
	input := make(chan int)

	go func() {
		for num := range numbers {
			input <- num
		}
		close(input)
	}()

	return input
}

func getInputChan(input <-chan int) <-chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for num := range input {
			time.Sleep(time.Duration(time.Second))
			output <- num * num
		}
	}()

	return output
}

func fanIn(chans ...<-chan int) <-chan int {
	var wg sync.WaitGroup

	output := make(chan int)

	// increase counter to number of channels `len(outputsChan)`
	// as we will spawn number of goroutines equal to number of channels received to merge
	wg.Add(len(chans))

	for _, optChan := range chans {
		go func(sc <-chan int) {
			// run until channel (square numbers sender) closes
			for sqr := range sc {
				output <- sqr
			}
			// once channel (square numbers sender) closes,
			// call `Done` on `WaitGroup` to decrement counter
			wg.Done()
		}(optChan)
	}

	// run goroutine to close merged channel once done
	go func() {
		// wait until WaitGroup finishesh
		wg.Wait()
		close(output)
	}()

	return output
}
