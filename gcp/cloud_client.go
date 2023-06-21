package gcp

import (
	"cloud.google.com/go/storage"
	"context"
)

const bucket = "cloudhaven-storage"

type CloudClient struct {
	*storage.BucketHandle
}

var instance *CloudClient

func GetCloudClient() *CloudClient {
	if instance == nil {
		instance = initClient()
	}
	return instance
}

func initClient() *CloudClient {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	return &CloudClient{client.Bucket(bucket)}
}
