package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"time"
)

const bucket = "cloudhaven-storage"

type storageClient struct {
	client *storage.Client
	bucket *storage.BucketHandle
}

var storageInstance *storageClient

func GetBucketClient() *storage.BucketHandle {
	if storageInstance == nil {
		storageInstance = initGCS()
	}
	return storageInstance.bucket
}

func GetBucketName() string {
	return bucket
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
		client: client,
		bucket: client.Bucket(bucket),
	}
}
