package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const maxPingTimeInMilliseconds = 5000

var domains = []string{
	"domain1.com",
	"domain2.com",
	"domain3.com",
	"domain4.com",
	"domain5.com",
	"domain6.com",
	"domain7.com",
	"domain8.com",
	"domain9.com",
	"domain10.com",
	"domain11.com",
	"domain12.com",
	"domain13.com",
	"domain14.com",
	"domain15.com",
}

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	fastest := make(chan string)

	wg := &sync.WaitGroup{}

	for _, domain := range domains {
		wg.Add(1)
		go func(ctx context.Context, domain string) {
			defer wg.Done()

			_ = ping(ctx, domain)
			fastest <- domain
		}(ctx, domain)
	}

	defer wg.Done()
	select {
	case res := <-fastest:
		fmt.Println("Fastest domain is", res)
		close(fastest)
		return
	case <-ctx.Done():
		close(fastest)
		return
	}
}

func ping(ctx context.Context, domain string) int {
	rand.Seed(time.Now().UnixNano())

	pingTime := rand.Intn(maxPingTimeInMilliseconds)
	fmt.Println(domain, "needs", pingTime, "milliseconds to ping")

	time.Sleep(time.Duration(pingTime) * time.Millisecond)

	return pingTime
}
