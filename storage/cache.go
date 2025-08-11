package storage

import (
	"context"
	"jinya-fonts/config"
	path2 "path"
	"time"

	"github.com/redis/go-redis/v9"
)

func getRedisContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func getRedisClient() (*redis.Client, error) {
	opts, err := redis.ParseURL(config.LoadedConfiguration.RedisUrl)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opts), nil
}

func getCachedFontFile(path string) ([]byte, error) {
	client, err := getRedisClient()
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := getRedisContext()
	defer cancelFunc()

	return client.Get(ctx, path2.Base(path)).Bytes()
}

func addCachedFontFile(path string, data []byte) error {
	client, err := getRedisClient()
	if err != nil {
		return err
	}

	ctx, cancelFunc := getRedisContext()
	defer cancelFunc()

	return client.Set(ctx, path2.Base(path), data, 0).Err()
}

func removeCachedFontFile(path string) error {
	client, err := getRedisClient()
	if err != nil {
		return err
	}

	ctx, cancelFunc := getRedisContext()
	defer cancelFunc()

	return client.Del(ctx, path2.Base(path)).Err()
}

func ClearGoogleFontsCache() {
	ctx, cancelFunc := getRedisContext()
	defer cancelFunc()

	client, err := getRedisClient()
	if err != nil {
		return
	}

	cursor := uint64(0)
	iter := client.Scan(ctx, cursor, "google*", 0).Iterator()
	for iter.Next(ctx) {
		_ = client.Del(ctx, iter.Val())
	}
}

func CheckRedis() bool {
	ctx, cancelFunc := getRedisContext()
	defer cancelFunc()

	client, err := getRedisClient()
	if err != nil {
		return false
	}

	return client.Ping(ctx).Err() == nil
}
