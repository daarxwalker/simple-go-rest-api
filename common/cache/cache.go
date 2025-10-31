package cache

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrPingFailed = errors.New("ping failed")
)

func Connect(host, password string, db int) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := redis.NewClient(
		&redis.Options{
			Addr:      host,
			Password:  password,
			DB:        db,
			TLSConfig: nil,
		},
	)
	if pingErr := client.Ping(ctx).Err(); pingErr != nil {
		return nil, errors.Join(ErrPingFailed, pingErr)
	}
	log.Println("🗄️ cache(" + strconv.Itoa(db) + "): ✅")
	return client, nil
}
