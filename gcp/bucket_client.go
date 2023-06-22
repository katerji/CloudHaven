package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"time"
)

const defaultBucket = "cloudhaven-storage"

type storageClient struct {
	client        *storage.Client
	defaultBucket *storage.BucketHandle
}

var storageInstance *storageClient

func GetBucketClient() *storage.BucketHandle {
	if storageInstance == nil {
		storageInstance = initGCS()
	}
	return storageInstance.defaultBucket
}

func GetBucketName() string {
	return defaultBucket
}

func GetDefaultSignOptions() *storage.SignedURLOptions {
	return &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(24 * time.Hour),
	}
}

func GetStorageClient() *storage.Client {
	if storageInstance == nil {
		storageInstance = initGCS()
	}
	return storageInstance.client
}

func initGCS() *storageClient {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	return &storageClient{
		client:        client,
		defaultBucket: client.Bucket(defaultBucket),
	}
}
