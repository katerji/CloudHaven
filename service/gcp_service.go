package service

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/katerji/UserAuthKit/gcp"
	"github.com/katerji/UserAuthKit/model"
	"google.golang.org/api/iterator"
)

type gcpService struct{}

func (service *gcpService) GetObjectSize(path string) (int64, bool) {
	object, err := gcp.GetBucketClient().Object(path).NewReader(context.Background())
	if err != nil {
		fmt.Println(err)
		return 0, false
	}
	return object.Attrs.Size, true
}

func (service *gcpService) ListUserObjects(userID int) ([]model.File, bool) {
	objectIterator := gcp.GetBucketClient().Objects(context.Background(), getUserQuery(userID))
	var files []model.File
	for {
		objectAttrs, err := objectIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return files, false
		}
		file := model.FromGCSObject(objectAttrs)
		file.OwnerID = userID
		files = append(files, file)
	}
	return files, true
}

func getUserQuery(userID int) *storage.Query {
	return &storage.Query{Prefix: fmt.Sprintf("%d/", userID)}
}

func (service *gcpService) DeleteObject(fileInput model.FileInput) bool {
	return gcp.GetBucketClient().Object(fileInput.GetPath()).Delete(context.Background()) == nil
}

func (service *gcpService) CreateObject(fileInput model.FileInput) bool {
	object := gcp.GetBucketClient().Object(fileInput.GetPath())
	writer := object.NewWriter(context.Background())
	writer.ContentType = fileInput.ContentType
	_, err := writer.Write(fileInput.Content)
	if err != nil {
		return false
	}
	return writer.Close() == nil
}

func (service *gcpService) SignObject(fileInput model.FileInput) (string, error) {
	url, err := gcp.GetStorageClient().Bucket(gcp.GetBucketName()).SignedURL(fileInput.GetPath(), gcp.GetDefaultSignOptions())
	if err != nil {
		fmt.Println(err)
		return "", errors.New("could not share file")
	}
	return url, nil
}
