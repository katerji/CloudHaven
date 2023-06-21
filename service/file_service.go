package service

import (
	"fmt"
	"github.com/katerji/UserAuthKit/model"
)

type fileService struct{}

func (service fileService) UpsertUserFiles(files []model.File) bool {
	for _, file := range files {
		fmt.Println(file.Name)
	}
	return true
}
