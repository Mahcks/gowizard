package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
}

func New(ctx context.Context, host, port, password string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", host, port),
		Password: password,
		DB:       0,
	})

	// Ping client to make sure we have a good connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{
		Client: client,
	}, nil
}

func (r *Redis) Close() error {
	if r.Client != nil {
		err := r.Client.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
