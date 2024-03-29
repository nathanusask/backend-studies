package main

/*
只用channel交替打印数字并有序退出。

要求：

1. 不能使用`WaitGroup/Mutex/`,只能使用`channel;
2. 不能出现deadlock，程序没有运行时错误，即没有panic；
3. 数字交替打印后，worker有序退出，不能出现数字还没打完但已经有worker退出的情况。
*/

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
		fmt.Printf("worker %d: %d\n", i, cur)
		if cur == n {
			close(last) // signal the last number is printed
		}
		cur += numWorkers
		if cur > n {
			// when the last number of the current worker is printed, the next worker will print its last number
			// if here we instead send a signal to the next worker which is not receiving the signal
			// the program will panic
			// remember, sending to a channel that has no receivers will panic
			close(w)
			break
		}
		w <- true
	}

	<-exitSignal
	// similar to the above, we need to make sure the next worker is ready to receive the exit signal
	// otherwise, just close the channel
	if i != numWorkers-1 {
		nextExitSignal <- true
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
	fmt.Println("print to number", n, "with", numWorkers, "workers")
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
	fmt.Println()
	printNumbersWithWorkers(10, 3)
	fmt.Println()
	printNumbersWithWorkers(9, 3)
	fmt.Println()
}
