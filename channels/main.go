package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var incr = make(chan int32)
	var counter int32 = 1

	go func() {
		incr <- counter
	}()

	for counter = <-incr; counter < 23; counter = <-incr {
		go func() {
			fmt.Println("A->", counter)
			atomic.AddInt32(&counter, 1)
			incr <- counter
		}()

		counter = <-incr
		go func() {
			fmt.Println("B: ", counter)
			atomic.AddInt32(&counter, 1)
			incr <- counter
		}()
	}
}
