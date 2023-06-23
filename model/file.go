package model

import (
	"cloud.google.com/go/storage"
	"fmt"
	"github.com/katerji/UserAuthKit/utils"
	"time"
)

type File struct {
	ID          int
	Name        string
	Size        int64
	ContentType string
	CreatedOn   time.Time
	ModifiedOn  time.Time
	OwnerID     int
}

func (file *File) GetPath() string {
	return fmt.Sprintf("%d/%s", file.OwnerID, file.Name)
}

func (file *File) ToOutput() FileOutput {
	return FileOutput{
		ID:          file.ID,
		Name:        file.Name,
		Size:        file.Size,
		ContentType: file.ContentType,
		CreatedOn:   utils.TimeToString(file.CreatedOn),
		ModifiedOn:  utils.TimeToString(file.ModifiedOn),
		OwnerID:     file.OwnerID,
	}
}

func FromGCSObject(attrs *storage.ObjectAttrs) File {
	return File{
		Name:        attrs.Name,
		Size:        attrs.Size,
		ContentType: attrs.ContentType,
		CreatedOn:   attrs.Created,
		ModifiedOn:  attrs.Updated,
	}
}
