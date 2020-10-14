package main

import (
	"context"

	redis "github.com/go-redis/redis/v8"
	"github.com/nathanusask/backend-studies/fake-web-retailer/internal/tokens"
)

func main() {
	ctx := context.Background()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	tokenService := tokens.New(redisClient)
}
