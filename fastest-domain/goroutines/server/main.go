package main

import (
	"context"
	"fmt"
	"math/rand"
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

	resp := findFastest(ctx)
	fmt.Println(resp)
}

func findFastest(ctx context.Context) string {
	fastest := make(chan string)

	for _, domain := range domains {
		go func(ctx context.Context, domain string) {
			_ = ping(ctx, domain)
			fastest <- domain
		}(ctx, domain)
	}

	select {
	case res := <-fastest:
		resp := fmt.Sprintf("Fastest domain is %s", res)
		close(fastest)
		return resp
	case <-ctx.Done():
		resp := "No fastest domain was found!"
		close(fastest)
		return resp
	}
}

func ping(ctx context.Context, domain string) int {
	rand.Seed(time.Now().UnixNano())

	pingTime := rand.Intn(maxPingTimeInMilliseconds)
	fmt.Println(domain, "needs", pingTime, "milliseconds to ping")

	time.Sleep(time.Duration(pingTime) * time.Millisecond)

	return pingTime
}
