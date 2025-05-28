package keydb

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewClient() *redis.Client {
	opts, err := redis.ParseURL(os.Getenv("KEYDB_DSN"))

	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opts)
	err = client.Ping(context.Background()).Err()

	if err != nil {
		panic(err)
	}

	return client
}
