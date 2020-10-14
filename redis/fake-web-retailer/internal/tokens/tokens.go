package tokens

import (
	"context"
	"fmt"
	"math"

	"time"

	redis "github.com/go-redis/redis/v8"
)

const (
	loginKey      = "login:"
	recentViewKey = "recent:"
	viewedKey     = "viewd:%s"
	// QUIT tells when to stop cleaning sessions; it's never going to stop however in this way :-<
	QUIT = false
	// LIMIT tells the limit of sessions
	LIMIT = 10000000
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
	CleanSessions(ctx context.Context)
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

// CleanSessions should be run as a daemon
func (s *service) CleanSessions(ctx context.Context) {
	var size int64
	var err error
	for !QUIT {
		size, err = s.redisClient.ZCard(ctx, recentViewKey).Result()
		if err != nil {
			panic(fmt.Sprintf("failed to find out how many tokens are known with error: %s", err.Error()))
		}
		if size <= LIMIT {
			time.Sleep(time.Second)
		} else {
			endIndex := math.Min(float64(size-LIMIT), float64(100))
			tokens, err := s.redisClient.ZRange(ctx, recentViewKey, int64(0), int64(endIndex-1)).Result()
			if err != nil {
				panic(fmt.Sprintf("failed to get tokens from %s between 0 and %d with error: %s", recentViewKey, int64(endIndex-1), err.Error()))
			}
			sessionKeys := []string{}
			for _, token := range tokens {
				sessionKeys = append(sessionKeys, fmt.Sprintf(viewedKey, token))
			}
			_, err = s.redisClient.Del(ctx, sessionKeys...).Result()
			if err != nil {
				panic(fmt.Sprintf("failed to delete session keys: %s", err.Error()))
			}
			_, err = s.redisClient.HDel(ctx, loginKey, tokens...).Result()
			if err != nil {
				panic(fmt.Sprintf("failed to delete tokens: %s", err.Error()))
			}
			_, err = s.redisClient.ZRem(ctx, recentViewKey, tokens).Result()
			if err != nil {
				panic(fmt.Sprintf("failed to remove tokens: %s", err.Error()))
			}
		}
	}
}
