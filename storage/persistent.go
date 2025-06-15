package storage

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"jinya-fonts/config"
	"strings"
	"time"
)

func getMinioContext(timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Minute)
}

func getMinioClient() (*minio.Client, error) {
	endpoint := config.LoadedConfiguration.S3ServerUrl
	accessKeyID := config.LoadedConfiguration.S3AccessKey
	secretAccessKey := config.LoadedConfiguration.S3SecretKey

	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
}

func savePersistentFontFile(path string, data []byte, fileType string) error {
	client, err := getMinioClient()
	if err != nil {
		return err
	}

	size := len(data)
	reader := bytes.NewReader(data)

	ctx, cancelFunc := getMinioContext(15)
	defer cancelFunc()

	_, err = client.PutObject(ctx, config.LoadedConfiguration.S3Bucket, path, reader, int64(size), minio.PutObjectOptions{
		ContentType: fileType,
	})

	return err
}

func removePersistentFontFile(path string) error {
	client, err := getMinioClient()
	if err != nil {
		return err
	}

	ctx, cancelFunc := getMinioContext(15)
	defer cancelFunc()

	return client.RemoveObject(ctx, config.LoadedConfiguration.S3Bucket, path, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
}

func getPersistentFontFile(path string) ([]byte, error) {
	client, err := getMinioClient()
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(path, "/fonts/") {
		path = "/fonts/" + path
	}

	ctx, cancelFunc := getMinioContext(1)
	defer cancelFunc()

	res, err := client.GetObject(ctx, config.LoadedConfiguration.S3Bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	defer res.Close()

	return io.ReadAll(res)
}
