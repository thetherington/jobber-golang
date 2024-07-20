package redis

import (
	"context"
	"fmt"

	apmgoredis "github.com/ekrucio/apm-agent-go/module/apmgoredisv9/v2"
	"github.com/redis/go-redis/v9"
	"github.com/thetherington/jobber-gig/internal/adapters/config"
)

const KEY = "selectedCategories"

/**
 * Redis implements port.CacheRepository interface
 * and provides an access to the go-redis library
 */
type Redis struct {
	client *redis.Client
	url    string
}

// New creates a new instance of Redis
func New(ctx context.Context, config *config.Redis) (*Redis, error) {
	opt, err := redis.ParseURL(config.Host)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	client.AddHook(apmgoredis.NewHook())

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{client, config.Host}, nil
}

func (r *Redis) GetUserSelectedGigCategory(ctx context.Context, username string) (string, error) {
	resp, err := r.client.Get(ctx, fmt.Sprintf("%s:%s", KEY, username)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}

		return "", err
	}

	return resp, nil
}

func (r *Redis) Close() {
	r.client.Close()
}
