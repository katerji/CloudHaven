package model

import (
	"cloud.google.com/go/storage"
	"fmt"
	"time"
)

type File struct {
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

func FromGCSObject(attrs *storage.ObjectAttrs) File {
	return File{
		Name:        attrs.Name,
		Size:        attrs.Size,
		ContentType: attrs.ContentType,
		CreatedOn:   attrs.Created,
		ModifiedOn:  attrs.Updated,
	}
}
