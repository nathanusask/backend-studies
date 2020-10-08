package main

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	pong, err := rdb.Ping(ctx).Result()
	fmt.Println(pong, err)
}
