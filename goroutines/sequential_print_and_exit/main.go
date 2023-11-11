package main

import (
	"fmt"
)

// worker prints numbers sequentially, but each worker prints every numWorkers-th number.
// After numbers are printed sequentially, each worker then waits for an exit signal.
// i: worker id, also the starting number+1
// numWorkers: number of workers
// n: the last number to print
// r: receive channel
// w: send channel
// exitSignal: receive channel for exit signal
// nextExitSignal: send channel for exit signal
// last: send channel for last number
// done: send channel for done signal
func worker(i, numWorkers, n int, r <-chan bool, w chan<- bool, exitSignal <-chan bool, nextExitSignal chan<- bool, last chan<- bool, done chan<- bool) {
	cur := i + 1
	for {
		<-r
		if cur <= n {
			fmt.Printf("worker %d: %d\n", i, cur)
			if cur == n {
				close(last)
			}
			cur += numWorkers
		}
		if cur > n {
			close(w)
			break
		}
		w <- true
	}

	<-exitSignal
	if i != numWorkers-1 {
		nextExitSignal <- true
	} else {
		close(nextExitSignal)
	}
	fmt.Printf("worker %d has exited\n", i)
	done <- true
}

// printNumbersWithWorkers prints numbers from 1 to n with numWorkers workers.
// It prints numbers sequentially, but each worker prints every numWorkers-th number.
// After numbers are printed sequentially, each worker then prints an exit message in sequential order.
// n: the last number to print
// numWorkers: number of workers
// For example, if n = 10 and numWorkers = 3, then the output will be:
// worker 0: 1
// worker 1: 2
// worker 2: 3
// worker 0: 4
// worker 1: 5
// worker 2: 6
// worker 0: 7
// worker 1: 8
// worker 2: 9
// worker 0: 10
// worker 0 has exited
// worker 1 has exited
// worker 2 has exited
func printNumbersWithWorkers(n, numWorkers int) {
	// init channels
	ch, exitChan := make([]chan bool, numWorkers), make([]chan bool, numWorkers)
	for i := 0; i < numWorkers; i++ {
		ch[i] = make(chan bool)
		exitChan[i] = make(chan bool)
	}
	done := make(chan bool, numWorkers)
	last := make(chan bool)

	// make sure all channels are closed or ready for gc
	defer func() {
		ch, exitChan = nil, nil // for gc
		close(done)
	}()

	// start workers
	for i := 0; i < numWorkers; i++ {
		go worker(i, numWorkers, n, ch[i], ch[(i+1)%numWorkers], exitChan[i], exitChan[(i+1)%numWorkers], last, done)
	}
	ch[0] <- true       // signal the first worker to start printing
	<-last              // wait for the last number to be printed
	exitChan[0] <- true // start the first worker to exit
	// wait for all workers to exit
	count := 0
	for {
		<-done
		count++
		if count == numWorkers {
			// all workers have exited
			return
		}
	}
}

// main is the entry point of the program.
func main() {
	printNumbersWithWorkers(40, 6)
	printNumbersWithWorkers(10, 3)
	printNumbersWithWorkers(9, 3)
}
