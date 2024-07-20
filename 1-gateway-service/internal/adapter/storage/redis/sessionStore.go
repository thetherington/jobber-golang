package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const PREFIX = "scs:session:"

// helper function to make a key with the prefix and token
func key(token string) string {
	return fmt.Sprintf("%s%s", PREFIX, token)
}

// Find returns the data for a given session token from the RedisStore instance.
// If the session token is not found or is expired, the returned exists flag
// will be set to false.
func (r *Redis) Find(token string) (b []byte, exists bool, err error) {
	b, err = r.client.Get(context.Background(), key(token)).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	return b, true, nil
}

// Commit adds a session token and data to the RedisStore instance with the
// given expiry time. If the session token already exists then the data and
// expiry time are updated.
func (r *Redis) Commit(token string, b []byte, expiry time.Time) error {
	pipe := r.client.Pipeline()

	pipe.Do(context.Background(), "SET", key(token), b)

	pipe.Do(context.Background(), "PEXPIREAT", key(token), makeMillisecondTimestamp(expiry))

	_, err := pipe.Exec(context.Background())

	return err
}

// Delete removes a session token and corresponding data from the RedisStore
// instance.
func (r *Redis) Delete(token string) error {
	_, err := r.client.Del(context.Background(), key(token)).Result()
	return err
}

// All returns a map containing the token and data for all active (i.e.
// not expired) sessions in the RedisStore instance.
func (r *Redis) All() (map[string][]byte, error) {
	keys, err := r.client.Keys(context.Background(), PREFIX+"*").Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	sessions := make(map[string][]byte)

	for _, key := range keys {
		token := key[len(PREFIX):]

		data, exists, err := r.Find(token)
		if err == redis.Nil {
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		if exists {
			sessions[token] = data
		}
	}

	return sessions, nil
}

func makeMillisecondTimestamp(t time.Time) int64 {
	return t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
