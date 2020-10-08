package main

import (
	"context"
	"time"

	"github.com/nathanusask/backend-studies/fake-web-retailer/internal/tokens"
	redis "github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	tokenService := tokens.New(redisClient)

	token := "some random token"
	tokenService.CheckToken(ctx, token)
}