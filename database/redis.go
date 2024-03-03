package database

import (
	"github.com/redis/go-redis/v9"
	"jinya-fonts/config"
)

func getRedisClient() (*redis.Client, error) {
	opts, err := redis.ParseURL(config.LoadedConfiguration.GoogleRedisUrl)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opts), nil
}

func GetCachedFontFile(path string) ([]byte, error) {
	client, err := getRedisClient()
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := getContext()
	defer cancelFunc()

	return client.Get(ctx, path).Bytes()
}

func AddCachedFontFile(name, weight, style, fileType string, data []byte, googleFont bool) error {
	client, err := getRedisClient()
	if err != nil {
		return err
	}

	ctx, cancelFunc := getContext()
	defer cancelFunc()

	return client.Set(ctx, GetFontFileName(name, weight, style, fileType, googleFont), data, 0).Err()
}

func RemoveCachedFontFile(name, weight, style, fileType string, googleFont bool) error {
	client, err := getRedisClient()
	if err != nil {
		return err
	}

	ctx, cancelFunc := getContext()
	defer cancelFunc()

	return client.Del(ctx, GetFontFileName(name, weight, style, fileType, googleFont)).Err()
}

func ClearGoogleFontsCache() {
	ctx, cancelFunc := getContext()
	defer cancelFunc()

	if client, err := getRedisClient(); err == nil {
		_ = client.FlushDB(ctx).Err()
	}
}
