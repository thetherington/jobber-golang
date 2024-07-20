package redis

import (
	"context"
	"fmt"

	apmgoredis "github.com/ekrucio/apm-agent-go/module/apmgoredisv9/v2"
	"github.com/redis/go-redis/v9"
	"github.com/thetherington/jobber-gateway/internal/adapter/config"
)

const (
	KEY_LOGGED_IN_USERS = "loggedInUsers"
	CATEGORY_KEY        = "selectedCategories"
)

/**
 * Redis implements port.CacheRepository and scs.Store interface
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

func (r *Redis) SaveUserSelectedCategory(ctx context.Context, username string, category string) error {
	return r.client.Set(ctx, fmt.Sprintf("%s:%s", CATEGORY_KEY, username), category, 0).Err()
}

func (r *Redis) SaveLoggedInUserToCache(ctx context.Context, value string) ([]string, error) {
	_, err := r.client.LPos(ctx, KEY_LOGGED_IN_USERS, value, redis.LPosArgs{Rank: 1}).Result()
	if err != nil {
		switch {
		case err == redis.Nil:
			r.client.LPush(ctx, KEY_LOGGED_IN_USERS, value)
		default:
			return nil, err
		}
	}

	return r.GetLoggedInUsersFromCache(ctx)
}

func (r *Redis) GetLoggedInUsersFromCache(ctx context.Context) ([]string, error) {
	users, err := r.client.LRange(ctx, KEY_LOGGED_IN_USERS, 0, -1).Result()
	if err != nil {
		if err == redis.Nil {
			return []string{}, nil
		}

		return nil, err
	}

	return users, nil
}

func (r *Redis) RemoveLoggedInUserFromCache(ctx context.Context, value string) ([]string, error) {
	r.client.LRem(ctx, KEY_LOGGED_IN_USERS, 1, value)

	return r.GetLoggedInUsersFromCache(ctx)
}

func (r *Redis) DeleteLoggedInUsers(ctx context.Context) error {
	_, err := r.client.Del(ctx, KEY_LOGGED_IN_USERS).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Close() {
	r.client.Close()
}
