package tests

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	testredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

func FlushCache(cache *redis.Client) {
	cache.FlushDB(context.Background())
}

func CreateTestCache(t testing.TB) (*redis.Client, func()) {
	container, err := testredis.Run(
		t.Context(),
		"redis:7",
	)
	require.NoError(t, err)
	t.Cleanup(
		func() {
			_ = container.Terminate(t.Context())
		},
	)
	endpoint, err := container.Endpoint(t.Context(), "")
	require.NoError(t, err)
	client := redis.NewClient(
		&redis.Options{
			Addr: endpoint,
		},
	)
	return client, func() {
		_ = client.Close()
	}
}
