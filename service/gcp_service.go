package service

import (
	"context"
	"fmt"
	"github.com/katerji/UserAuthKit/gcp"
)

type GCPService struct{}

func (service *GCPService) GetObjectSize(path string) (int64, bool) {
	object, err := gcp.GetCloudClient().Object(path).NewReader(context.Background())
	if err != nil {
		fmt.Println(err)
		return 0, false
	}
	return object.Attrs.Size, true
}
