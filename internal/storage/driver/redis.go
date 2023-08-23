package driver

import (
	"context"
	"fmt"
	"time"

	driver "github.com/go-redis/redis/v8"
)

var pingTimeout = 10 * time.Second

type redis struct {
	client *driver.Client
}

func Redis(opt RedisOpts) (*redis, error) {
	client := driver.NewClient(&driver.Options{
		Addr:     opt.Address,
		Username: opt.Username,
		Password: opt.Password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping error: %w", err.Err())
	}

	return &redis{client: client}, nil
}

// Set updates the value associated with a key in Redis.
func (r *redis) Set(ctx context.Context, key, value string) error {
	err := r.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves the value associated with a key from Redis.
func (r *redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == driver.Nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return val, nil
}
