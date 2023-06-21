package gcp

import (
	"cloud.google.com/go/storage"
	"context"
)

const bucket = "cloudhaven-storage"

type bucketClient struct {
	*storage.BucketHandle
}

var instance *bucketClient

func GetBucketClient() *bucketClient {
	if instance == nil {
		instance = initClient()
	}
	return instance
}

func initClient() *bucketClient {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	return &bucketClient{client.Bucket(bucket)}
}
