package service

import (
	"cloud.google.com/go/storage"
	"context"
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
