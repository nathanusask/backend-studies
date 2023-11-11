package main

import (
	"fmt"
)

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
