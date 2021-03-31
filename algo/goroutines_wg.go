package main

import (
	"fmt"
	"sync"
)

// var mu sync.Mutex

func printOdd(i *int, wg *sync.WaitGroup) {
	defer wg.Done()
	// mu.Lock()
	// defer mu.Unlock()
	if *i%2 == 1 {
		fmt.Println(fmt.Sprintf("A=%d", *i))
		*i = *i + 1
	}
}

func printEven(i *int, wg *sync.WaitGroup) {
	defer wg.Done()
	// mu.Lock()
	// defer mu.Unlock()
	if *i%2 == 0 {
		fmt.Println(fmt.Sprintf("B:%d", *i))
		*i = *i + 1
	}
}

func main() {
	var wg sync.WaitGroup
	i := 1

	for {
		wg.Add(1)
		go printOdd(&i, &wg)
		wg.Wait()

		wg.Add(1)
		go printEven(&i, &wg)
		wg.Wait()

		if i > 20 {
			break
		}
	}

	// wg.Wait()
}
