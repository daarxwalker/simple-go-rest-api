package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Client *redis.Client
}

func New(client *redis.Client) *Storage {
	return &Storage{
		Client: client,
	}
}

func (s *Storage) Exists(c context.Context, key ...string) bool {
	return s.Client.Exists(c, key...).Val() == 1
}

func (s *Storage) Get(c context.Context, key string, target any) error {
	value := s.Client.Get(c, key).Val()
	if len(value) == 0 {
		return nil
	}
	if err := json.Unmarshal([]byte(value), target); err != nil {
		return fmt.Errorf("deserialize cached value failed: %w", err)
	}
	return nil
}

func (s *Storage) Set(c context.Context, key string, value any, expiration time.Duration) error {
	valueBytes, marshalErr := json.Marshal(value)
	if marshalErr != nil {
		return fmt.Errorf("serialize value to cache failed: %w", marshalErr)
	}
	if err := s.Client.Set(c, key, string(valueBytes), expiration).Err(); err != nil {
		return fmt.Errorf("unable to set cache value: %w", err)
	}
	return nil
}

func (s *Storage) Destroy(c context.Context, key string) error {
	return s.Set(c, key, "", time.Millisecond)
}
