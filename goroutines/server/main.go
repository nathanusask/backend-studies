package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const maxPingTime = 10

var domains = []string{
	"domain1.com",
	"domain2.com",
	"domain3.com",
	"domain4.com",
	"domain5.com",
	"domain6.com",
	"domain7.com",
	"domain8.com",
}

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	fastest := make(chan string, len(domains)+1)
	defer close(fastest)

	wg := &sync.WaitGroup{}

	for _, domain := range domains {
		wg.Add(1)
		go func(ctx context.Context, domain string) {
			defer wg.Done()

			_ = ping(ctx, domain)
			select {
			case fastest <- domain:
			}
		}(ctx, domain)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case res := <-fastest:
			fmt.Println("Fastest domain is ", res)
			return
		case <-ctx.Done():
			return
		}
	}()

	wg.Wait()
}

func ping(ctx context.Context, domain string) int {
	rand.Seed(time.Now().UnixNano())
	pingTime := rand.Intn(maxPingTime)
	fmt.Println(domain, "needs", pingTime, "seconds")
	time.Sleep(time.Duration(pingTime) * time.Second)

	return pingTime
}
