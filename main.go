package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
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
	group := new(errgroup.Group)
	out := make(chan int)

	for _, channel := range channels {
		// For each channel, closure over the variables
		// the `func` passed to `.Go()` cannot take parameters
		channel := channel
		group.Go(func() error {
			// Consume each value
			for value := range channel {
				// Forward the value
				out <- value
			}

			return nil
		})
	}

	go func() {
		// All channels are closed
		group.Wait()
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
