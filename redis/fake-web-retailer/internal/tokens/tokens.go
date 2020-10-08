package tokens

import (
	"context"
	"fmt"

	"time"

	redis "github.com/go-redis/redis/v8"
)

const (
	loginKey      = "login:"
	recentViewKey = "recent:"
	viewedKey     = "viewd:%s"
)

type service struct {
	redisClient *redis.Client
}

// New initiates a new token service
func New(conn *redis.Client) Service {
	return &service{
		redisClient: conn,
	}
}

// Service is a service for token-related operations
type Service interface {
	CheckToken(ctx context.Context, token string) (string, error)
	UpdateToken(ctx context.Context, token string, user string, item string) error
}

func (s *service) CheckToken(ctx context.Context, token string) (string, error) {
	return s.redisClient.HGet(ctx, loginKey, token).Result()
}

func (s *service) UpdateToken(ctx context.Context, token string, user string, item string) error {
	timestamp := time.Now().Unix()
	_, err := s.redisClient.HSet(ctx, loginKey, token, user).Result()
	if err != nil {
		err = fmt.Errorf("failed to set token %s and user %s with error: %s", token, user, err.Error())
		return err
	}

	_, err = s.redisClient.ZAdd(ctx, recentViewKey, &redis.Z{
		Member: token,
		Score:  float64(timestamp),
	}).Result()
	if err != nil {
		err = fmt.Errorf("failed to add token %s to sorted set %s with error: %s", token, recentViewKey, err.Error())
		return err
	}

	if item != "" {
		viewedSetKey := fmt.Sprintf(viewedKey, token)
		_, err = s.redisClient.ZAdd(ctx, viewedKey, &redis.Z{Member: item, Score: float64(timestamp)}).Result()
		if err != nil {
			err = fmt.Errorf("failed to add item %s to the viewed set %s with error: %s", item, viewedSetKey, err.Error())
			return err
		}

		_, err = s.redisClient.ZRemRangeByRank(ctx, viewedSetKey, 0, -26).Result() // only keep 25 most recent views
		if err != nil {
			err = fmt.Errorf("failed to keep 25 most recent views for the viewed set %s with error: %s", viewedSetKey, err.Error())
			return err
		}
	}

	return nil
}
