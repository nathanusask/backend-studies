package main

import (
	"fmt"
	"sync"
)

var (
	mu sync.Mutex
)

func printOddAndIncrement(counter *int) {
	mu.Lock()
	defer mu.Unlock()

	if *counter < 21 && *counter%2 == 1 {
		fmt.Println("A->", *counter)
		*counter++
	}
}

func printEvenAndIncrement(counter *int) {
	mu.Lock()
	defer mu.Unlock()

	if *counter < 21 && *counter%2 == 0 {
		fmt.Println("B:", *counter)
		*counter++
	}
}

func main() {
	counter := 1

	for {
		go printOddAndIncrement(&counter)
		go printEvenAndIncrement(&counter)

		if counter > 20 {
			break
		}
	}
}
