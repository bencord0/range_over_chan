package main

import (
	"fmt"
	"sync"
)

func arrangeData(from, to int) <-chan int {
	out := make(chan int)

	increment := 1
	if from > to {
		increment = -1
	}

	go func() {
		value := from
		for value != to {
			out <- value
			value += increment
		}
		close(out)
	}()

	return out
}

func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	wg.Add(len(channels))
	for _, channel := range channels {
		// For each channel
		go func(channel <-chan int) {
			// Consume each value
			for value := range channel {
				// Forward the value
				out <- value
			}
			// When channel closed, unblock.
			wg.Done()
		}(channel)
	}

	go func() {
		// All channels are closed
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	c1 := arrangeData(1, 10)
	c2 := arrangeData(7, -3)
	c3 := arrangeData(100, 110)

	for value := range merge(c1, c2, c3) {
		fmt.Println(value)
	}
}
