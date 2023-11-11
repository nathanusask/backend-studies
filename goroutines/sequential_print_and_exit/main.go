package main

import (
	"fmt"
)

// worker prints numbers sequentially, but each worker prints every numWorkers-th number.
// After numbers are printed sequentially, each worker then waits for an exit signal.
func worker(i, numWorkers, n int, r <-chan bool, w chan<- bool, exitSignal <-chan bool, next chan<- bool, last chan<- bool, done chan<- bool) {
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
		next <- true
	} else {
		close(next)
	}
	fmt.Printf("worker %d has exited\n", i)
	done <- true
}

// printNumbersWithWorkers prints numbers from 1 to n with numWorkers workers.
// It prints numbers sequentially, but each worker prints every numWorkers-th number.
// After numbers are printed sequentially, each worker then prints an exit message in sequential order.
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
	ch, exitChan := make([]chan bool, numWorkers), make([]chan bool, numWorkers)
	for i := 0; i < numWorkers; i++ {
		ch[i] = make(chan bool)
		exitChan[i] = make(chan bool)
	}
	done := make(chan bool, numWorkers)
	last := make(chan bool)
	defer func() {
		ch, exitChan = nil, nil // for gc
		close(done)
	}()
	for i := 0; i < numWorkers; i++ {
		go worker(i, numWorkers, n, ch[i], ch[(i+1)%numWorkers], exitChan[i], exitChan[(i+1)%numWorkers], last, done)
	}
	ch[0] <- true
	<-last
	exitChan[0] <- true
	count := 0
	for {
		<-done
		count++
		if count == numWorkers {
			return
		}
	}
}

func main() {
	printNumbersWithWorkers(40, 6)
	printNumbersWithWorkers(10, 3)
	printNumbersWithWorkers(9, 3)
}
