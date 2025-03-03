package services

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheService struct {
	client *redis.Client
}

func NewCacheService() *CacheService {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &CacheService{client: client}
}

func (s *CacheService) Get(ctx context.Context, key string, value interface{}) error {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), value)
}

func (s *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, val, expiration).Err()
}

func (s *CacheService) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

func (s *CacheService) InvalidateProjectCache(projectID uint64) {
	ctx := context.Background()
	key := "project:" + strconv.FormatUint(projectID, 10)
	err := s.Delete(ctx, key)
	if err != nil {
		log.Printf("Error invalidating project cache: %v", err)
	}
}

func (s *CacheService) InvalidateUserCache(userID uint) {
	ctx := context.Background()
	key := "user:" + strconv.FormatUint(uint64(userID), 10)
	err := s.Delete(ctx, key)
	if err != nil {
		log.Printf("Error invalidating user cache: %v", err)
	}
}
